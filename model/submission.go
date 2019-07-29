package model

import (
	"time"

	"github.com/jinzhu/gorm"
	gtypes "github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// Submission is the concrete type for a test submission, unlike
// types/submission.go which is the ephemeral type for messages and so on.
// Each submission has it own ID and list of tasks so that we can further give
// the user the opportunity to aggregate their tasks into a single unit.
type Submission struct {
	ID int64 `gorm:"priamry_key" json:"id"`

	User   *User `gorm:"association_autoupdate:false,nullable:true" json:"user"`
	UserID int64 `json:"-"`

	HeadRef   *Ref  `gorm:"association_autoupdate:false,column:head_ref_id,nullable:true" json:"head_ref"`
	HeadRefID int64 `json:"-"`

	BaseRef   *Ref  `gorm:"association_autoupdate:false,column:base_ref_id" json:"base_ref"`
	BaseRefID int64 `json:"-"`

	CreatedAt time.Time
}

// ToProto converts the submissions to the protobuf version
func (s *Submission) ToProto() *gtypes.Submission {
	var (
		pu *gtypes.User
		hr *gtypes.Ref
	)

	if s.User != nil {
		pu = s.User.ToProto()
	}

	if s.HeadRef != nil {
		hr = s.HeadRef.ToProto()
	}

	return &gtypes.Submission{
		Id:        s.ID,
		BaseRef:   s.BaseRef.ToProto(),
		HeadRef:   hr,
		User:      pu,
		CreatedAt: MakeTimestamp(&s.CreatedAt),
	}
}

// NewSubmissionFromProto converts the proto representation to the task type.
func NewSubmissionFromProto(gt *gtypes.Submission) (*Submission, *errors.Error) {
	var (
		u       *User
		headref *Ref
		err     *errors.Error
	)

	if gt.User != nil {
		u, err = NewUserFromProto(gt.User)
		if err != nil {
			return nil, err.Wrap("converting for use in submission")
		}
	}

	if gt.HeadRef != nil {
		headref, err = NewRefFromProto(gt.HeadRef)
		if err != nil {
			return nil, err.Wrap("converting for use in submission")
		}
	}

	baseref, err := NewRefFromProto(gt.BaseRef)
	if err != nil {
		return nil, err.Wrap("converting for use in submission")
	}

	return &Submission{
		ID:      gt.Id,
		User:    u,
		BaseRef: baseref,
		HeadRef: headref,
	}, nil
}

// NewSubmissionFromMessage creates a blank record populated by the appropriate data
// reaped from the message type in types/submission.go.
func (m *Model) NewSubmissionFromMessage(sub *types.Submission) (*Submission, *errors.Error) {
	if err := sub.Validate(); err != nil {
		return nil, err.Wrap("did not pass validation")
	}

	var (
		u    *User
		base *Ref
		err  *errors.Error
	)

	if sub.SubmittedBy != "" {
		u, err = m.FindUserByName(sub.SubmittedBy)
		if err != nil {
			return nil, err.Wrap("for use in submission")
		}
	}

	if sub.BaseSHA != "" {
		base, err = m.GetRefByNameAndSHA(sub.Parent, sub.BaseSHA)
		if err != nil {
			return nil, err.Wrap("base")
		}
	}

	head, err := m.GetRefByNameAndSHA(sub.Fork, sub.HeadSHA)
	if err != nil {
		return nil, err.Wrap("head")
	}

	return &Submission{
		User:    u,
		HeadRef: head,
		BaseRef: base,
	}, nil
}

// SubmissionList returns a list of submissions with pagination and repo/sha filtering.
func (m *Model) SubmissionList(page, perPage int64, repository, sha string) ([]*Submission, *errors.Error) {
	subs := []*Submission{}

	page, perPage, err := utils.ScopePaginationInt(page, perPage)
	if err != nil {
		return nil, err
	}

	obj, err := m.submissionListQuery(repository, sha)
	if err != nil {
		return nil, err
	}

	obj = obj.Offset(page * perPage).Limit(perPage)
	return subs, m.WrapError(obj.Find(&subs), "listing submissions")
}

// SubmissionCount counts the number of submissions with an optional repo/sha filter
func (m *Model) SubmissionCount(repository, sha string) (int64, *errors.Error) {
	var count int64

	obj, err := m.submissionListQuery(repository, sha)
	if err != nil {
		return 0, err
	}

	return count, m.WrapError(obj.Model(&Submission{}).Count(&count), "listing submissions")
}

func (m *Model) submissionListQuery(repository, sha string) (*gorm.DB, *errors.Error) {
	obj := m.Order("submissions.id DESC").
		Joins("inner join refs on refs.id = submissions.base_ref_id").
		Joins("inner join repositories on repositories.id = refs.repository_id")

	var (
		repo *Repository
		ref  *Ref
	)

	if repository != "" {
		var err *errors.Error
		repo, err = m.GetRepositoryByName(repository)
		if err != nil {
			return nil, err
		}

		if sha != "" {
			ref, err = m.GetRefByNameAndSHA(repository, sha)
			if err != nil {
				return nil, err
			}
		}
	}

	if ref != nil {
		obj = obj.Where("submissions.base_ref_id = ?", ref.ID)
	} else if repo != nil {
		obj = obj.Where("repositories.id = ?", repo.ID)
	}

	return obj, nil
}

// SubmissionListForRepository returns a list of submissions with pagination. If ref
// is non-nil, it will isolate to the ref only and ignore the repo.
func (m *Model) SubmissionListForRepository(repo, sha string, page, perPage int64) ([]*Submission, *errors.Error) {
	// FIXME this should probably be two independent calls with a query builder or something
	subs := []*Submission{}

	page, perPage, err := utils.ScopePaginationInt(page, perPage)
	if err != nil {
		return nil, err
	}

	obj, err := m.submissionListQuery(repo, sha)
	if err != nil {
		return nil, err
	}

	obj = obj.Offset(page * perPage).Limit(perPage)
	return subs, m.WrapError(obj.Find(&subs), "listing submissions for repository")
}

// TasksForSubmission returns all the tasks for a given submission.
func (m *Model) TasksForSubmission(sub *Submission, page, perPage int64) ([]*Task, *errors.Error) {
	tasks := []*Task{}

	obj := m.Offset(page*perPage).
		Limit(perPage).
		Order("tasks.id DESC").
		Joins("inner join submissions on tasks.submission_id = submissions.id").
		Where("submissions.id = ?", sub.ID)

	return tasks, m.WrapError(obj.Find(&tasks), "listing tasks for a submission")
}

// GetSubmissionByID returns a submission by internal identifier
func (m *Model) GetSubmissionByID(id int64) (*Submission, *errors.Error) {
	sub := &Submission{}
	return sub, m.WrapError(m.Where("submissions.id = ?", id).First(sub), "getting submission by id")
}
