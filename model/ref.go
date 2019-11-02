package model

import (
	"github.com/jinzhu/gorm"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/errors"
)

// Ref encapsulates git ref -- sha or branch name -- which is tied to a task
// (and multiple runs, potentially)
type Ref struct {
	ID           int64       `gorm:"primary_key" json:"id"`
	Repository   *Repository `json:"repository" gorm:"unique;association_autoupdate:false"`
	RepositoryID int64       `json:"-"`
	RefName      string      `gorm:"column:ref" json:"ref_name"`
	SHA          string      `json:"sha"`
}

// NewRefFromProto converts a proto ref to a real ref.
func NewRefFromProto(r *types.Ref) (*Ref, error) {
	repo, err := NewRepositoryFromProto(r.Repository)
	if err != nil {
		return nil, err
	}

	return &Ref{
		ID:         r.Id,
		Repository: repo,
		RefName:    r.RefName,
		SHA:        r.Sha,
	}, nil
}

// ToProto returns a ref in the protobuf representation
func (r *Ref) ToProto() *types.Ref {
	return &types.Ref{
		Id:         r.ID,
		Repository: r.Repository.ToProto(),
		RefName:    r.RefName,
		Sha:        r.SHA,
	}
}

// Validate validates the ref before saving it and after fetching it.
func (r *Ref) Validate() error {
	if r.Repository == nil {
		return errors.New("invalid repository")
	}

	if r.RefName == "" {
		return errors.New("empty ref")
	}

	if r.SHA == "" {
		return errors.New("empty SHA")
	}

	if len(r.SHA) != 40 {
		return errors.New("invalid SHA")
	}

	return nil
}

// AfterFind validates the output from the database before releasing it to the
// hook chain
func (r *Ref) AfterFind(tx *gorm.DB) error {
	if err := r.Validate(); err != nil {
		return errors.New(err).(errors.Error).Wrapf("reading ref id %d (%q)", r.ID, r.SHA)
	}

	return nil
}

// BeforeCreate just calls BeforeSave.
func (r *Ref) BeforeCreate(tx *gorm.DB) error {
	return r.BeforeSave(tx)
}

// BeforeSave is a gorm hook to marshal the token JSON before saving the record
func (r *Ref) BeforeSave(tx *gorm.DB) error {
	if err := r.Validate(); err != nil {
		return errors.New(err).(errors.Error).Wrapf("saving ref %q (%q)", r.RefName, r.SHA)
	}

	return nil
}

// GetRefByNameAndSHA returns the ref matching the name/sha combination.
func (m *Model) GetRefByNameAndSHA(repoName string, sha string) (*Ref, error) {
	ref := &Ref{}
	err := m.WrapError(
		m.Joins("inner join repositories on refs.repository_id = repositories.id").
			Where("repositories.name = ? and refs.sha = ?", repoName, sha).
			First(ref),
		"getting ref by name and sha",
	)
	return ref, errors.New(err)
}

// PutRef adds the ref to the database.
func (m *Model) PutRef(ref *Ref) error {
	return m.WrapError(m.Create(ref), "creating ref")
}

// CancelRefByName is used in auto cancellation on new queue item arrivals. It finds
// all the runs associated with the ref and cancels them if they are in queue
// still.
// Do note that it does not match the SHA; more often than not this is caused
// by an --amend + force push which updates the SHA, or a new commit which also
// changes the SHA. The name is the only reliable data in this case.
func (m *Model) CancelRefByName(repoID int64, refName, baseURL string, gh github.Client) error {
	tasks := []*Task{}

	repo := &Repository{}
	err := m.WrapError(m.Where("id = ?", repoID).First(repo), "finding repository during cancel ref operation")
	if err != nil {
		return errors.New(err)
	}

	mb := repo.Github.GetMasterBranch()
	if refName == mb || (mb == "" && refName == "heads/master") {
		return nil
	}

	err = m.WrapError(
		m.Where("refs.ref = ? and refs.repository_id = ? and tasks.status is null and tasks.finished_at is null", refName, repoID).
			Joins("inner join submissions on submissions.id = tasks.submission_id").
			Joins("inner join refs on refs.id = submissions.head_ref_id").
			Find(&tasks),
		"finding tasks during cancel ref operation",
	)
	if err != nil {
		return errors.New(err)
	}

	for _, task := range tasks {
		if err := m.CancelTask(task, baseURL, gh); err != nil {
			return err
		}
	}

	return nil
}
