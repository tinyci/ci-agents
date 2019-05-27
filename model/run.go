package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/errors"
	gtypes "github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// orchestration code for certain status operations.
type runBits struct {
	run    *Run
	github github.Client
	parts  []string
}

// Run is the individual parallel run. Runs are organized into the queue
// service for handing out to runners in order. Multiple runs are make up a
// task.
type Run struct {
	ID int64 `gorm:"primary_key" json:"id"`

	RunSettingsJSON []byte `gorm:"column:run_settings" json:"-"`

	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Status     *bool      `json:"status,omitempty"`
	Task       *Task      `gorm:"association_autoupdate:false" json:"task"`
	TaskID     int64      `json:"-"`

	RunSettings *types.RunSettings `json:"settings"`
}

// NewRunFromProto yields a new run from a protobuf message.
func NewRunFromProto(r *gtypes.Run) (*Run, *errors.Error) {
	task, err := NewTaskFromProto(r.Task)
	if err != nil {
		return nil, err
	}

	return &Run{
		ID:          r.Id,
		Name:        r.Name,
		CreatedAt:   *MakeTime(r.CreatedAt, false),
		StartedAt:   MakeTime(r.StartedAt, true),
		FinishedAt:  MakeTime(r.FinishedAt, true),
		Status:      MakeStatus(r.Status, r.StatusSet),
		Task:        task,
		RunSettings: types.NewRunSettingsFromProto(r.Settings),
	}, nil
}

// ToProto converts the run to its protobuf representation
func (r *Run) ToProto() *gtypes.Run {
	var status, set bool
	if r.Status != nil {
		status = *r.Status
		set = true
	}

	return &gtypes.Run{
		Id:         r.ID,
		Name:       r.Name,
		CreatedAt:  MakeTimestamp(&r.CreatedAt),
		StartedAt:  MakeTimestamp(r.StartedAt),
		FinishedAt: MakeTimestamp(r.FinishedAt),
		Status:     status,
		StatusSet:  set,
		Task:       r.Task.ToProto(),
		Settings:   r.RunSettings.ToProto(),
	}
}

// Validate validates the run record, accounting for creation or modification.
// Returns *errors.Error on any found.
func (r *Run) Validate() *errors.Error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	if r.Task == nil {
		return errors.New("no task associated")
	}

	if r.RunSettings == nil {
		return errors.New("no settings provided")
	}

	return r.RunSettings.Validate()
}

// AfterFind validates the output from the database before releasing it to the
// hook chain
func (r *Run) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal(r.RunSettingsJSON, &r.RunSettings); err != nil {
		return errors.New(err).Wrapf("unpacking task settings for task %d", r.ID)
	}

	if err := r.Validate(); err != nil {
		return errors.New(err).Wrapf("reading task id %d", r.ID)
	}

	return nil
}

// BeforeCreate just calls BeforeSave.
func (r *Run) BeforeCreate(tx *gorm.DB) error {
	return r.BeforeSave(tx)
}

// BeforeSave is a gorm hook to marshal the token JSON before saving the record
func (r *Run) BeforeSave(tx *gorm.DB) error {
	if err := r.Validate(); err != nil {
		return errors.New(err).Wrapf("saving task id %d", r.ID)
	}

	var err error
	r.RunSettingsJSON, err = json.Marshal(r.RunSettings)
	if err != nil {
		return errors.New(err).Wrapf("marshaling settings for task id %d", r.ID)
	}

	return nil
}

// SetRunStatus sets the status for the run; fails if it is already set.
func (m *Model) SetRunStatus(runID int64, gh github.Client, status, canceled bool, url, addlMessage string) *errors.Error {
	run := &Run{}

	if err := m.WrapError(m.Where("id = ?", runID).First(run), "finding run to set status"); err != nil {
		return errors.New(err)
	}

	if run.Status != nil {
		return errors.Errorf("status already set for run %d", runID)
	}

	qi := &QueueItem{}

	if err := m.WrapError(m.Where("run_id = ?", runID).First(qi), "finding queue item to clear during status update"); err != nil {
		return errors.New(err).Wrapf("locating queue item for run %d", runID)
	}

	if err := m.WrapError(m.Delete(qi), "deleting queue item during status update"); err != nil {
		return errors.New(err).Wrapf("while deleting queue item %d", qi.ID)
	}

	bits, err := m.getRunBits(runID, gh)
	if err != nil {
		return err
	}

	if canceled {
		go func() {
			if err := bits.github.ErrorStatus(bits.parts[0], bits.parts[1], bits.run.Name, bits.run.Task.Ref.SHA, fmt.Sprintf("%s/log/%d", url, runID), errors.ErrRunCanceled); err != nil {
				fmt.Println(err) // FIXME log
			}
		}()
	} else {
		go func() {
			if err := bits.github.FinishedStatus(bits.parts[0], bits.parts[1], bits.run.Name, bits.run.Task.Ref.SHA, fmt.Sprintf("%s/log/%d", url, runID), status, addlMessage); err != nil {
				fmt.Println(err)
			}
		}()
	}

	run.Status = &status
	now := time.Now()
	run.FinishedAt = &now
	if err := m.WrapError(m.Save(run), "saving updated run times"); err != nil {
		return errors.New(err)
	}

	return m.UpdateTaskStatus(run.Task)
}

// CancelRun is a thin facade for the CancelTask functionality, loading the record from the ID provided.
func (m *Model) CancelRun(runID int64, baseURL string, gh github.Client) *errors.Error {
	run := &Run{}

	if err := m.WrapError(m.Where("id = ?", runID).First(run), "locating run for cancel"); err != nil {
		return errors.New(err)
	}

	return m.CancelTask(run.Task, baseURL, gh)
}

// GetCancelForRun satisfies the datasvc interface by finding the underlying
// task's canceled state.
func (m *Model) GetCancelForRun(runID int64) (bool, *errors.Error) {
	run := &Run{}

	if err := m.WrapError(m.Where("id = ?", runID).First(run), "locating run to report cancel state"); err != nil {
		return false, errors.New(err)
	}

	return run.Task.Canceled, nil
}

func (m *Model) getRunBits(runID int64, gh github.Client) (*runBits, *errors.Error) {
	run := &Run{}

	load := m.DB

	if gh == nil {
		load = m.Preload("Task.Parent")
	}

	if err := m.WrapError(load.Where("id = ?", runID).First(run), "locating run"); err != nil {
		return nil, errors.New(err)
	}

	parts := strings.SplitN(run.Task.Parent.Name, "/", 2)
	if len(parts) != 2 {
		return nil, errors.Errorf("invalid repository for run %d: %v", run.ID, run.Task.Ref.Repository.Name)
	}

	if gh == nil {
		if run.Task.Parent.Owner == nil {
			return nil, errors.Errorf("No owner for repository %q corresponding to run %d", run.Task.Parent.Name, run.ID)
		}
		gh = github.NewClientFromAccessToken(run.Task.Parent.Owner.Token.Token)
	}

	return &runBits{
		run:    run,
		github: gh,
		parts:  parts,
	}, nil
}

// RunList returns a list of runs with pagination.
func (m *Model) RunList(page, perPage int64, repository, sha string) ([]*Run, *errors.Error) {
	runs := []*Run{}

	page, perPage, err := utils.ScopePaginationInt(page, perPage)
	if err != nil {
		return nil, err
	}

	obj := m.Offset(page * perPage).Limit(perPage).Order("runs.id DESC")

	if repository != "" {
		repo, err := m.GetRepositoryByName(repository)
		if err != nil {
			return nil, err
		}

		var r *Ref

		if sha != "" {
			r, err = m.GetRefByNameAndSHA(repository, sha)
			if err != nil {
				return nil, err
			}
		}

		return m.RunListForRepository(repo, r, page, perPage)
	}

	return runs, m.WrapError(obj.Find(&runs), "listing runs")
}

// RunListForRepository returns a list of queue items with pagination. If ref
// is non-nil, it will isolate to the ref only.
func (m *Model) RunListForRepository(repo *Repository, ref *Ref, page, perPage int64) ([]*Run, *errors.Error) {
	runs := []*Run{}

	page, perPage, err := utils.ScopePaginationInt(page, perPage)
	if err != nil {
		return nil, err
	}

	obj := m.Offset(page * perPage).
		Limit(perPage).
		Order("runs.id DESC").
		Joins("inner join tasks on tasks.id = runs.task_id").
		Joins("inner join refs on refs.id = tasks.ref_id").
		Joins("inner join repositories on repositories.id = refs.repository_id")

	if ref != nil {
		obj = obj.Where("repositories.id = ? and refs.id = ?", repo.ID, ref.ID)
	} else {
		obj = obj.Where("repositories.id = ?", repo.ID)
	}

	return runs, m.WrapError(obj.Find(&runs), "listing runs for repository")
}

// RunTotalCount returns the number of items in the runs table
func (m *Model) RunTotalCount() (int64, *errors.Error) {
	var ret int64
	return ret, m.WrapError(m.Table("runs").Count(&ret), "counting runs")
}

// RunTotalCountForRepository returns the number of items in the queue where
// the parent fork matches the repository name given
func (m *Model) RunTotalCountForRepository(repo *Repository) (int64, *errors.Error) {
	var ret int64
	return ret, m.WrapError(
		m.Table("runs").
			Joins("inner join tasks on runs.task_id = tasks.id").
			Joins("inner join refs on tasks.ref_id = refs.id").
			Joins("inner join repositories on refs.repository_id = repositories.id").
			Where("repositories.id = ?", repo.ID).
			Count(&ret),
		"counting runs for repository",
	)
}

// RunTotalCountForRepositoryAndSHA returns the number of items in the queue where
// the parent fork matches the repository name given
func (m *Model) RunTotalCountForRepositoryAndSHA(repo *Repository, sha string) (int64, *errors.Error) {
	var ret int64
	return ret, m.WrapError(
		m.Table("runs").
			Joins("inner join tasks on runs.task_id = tasks.id").
			Joins("inner join refs on tasks.ref_id = refs.id").
			Joins("inner join repositories on refs.repository_id = repositories.id").
			Where("repositories.id = ? and refs.sha = ?", repo.ID, sha).
			Count(&ret),
		"counting runs for repository and sha",
	)
}

// RunsForPR retrieves all the runs that belong to a repository's PR.
func (m *Model) RunsForPR(repoName string, prID int) ([]*Run, *errors.Error) {
	ret := []*Run{}

	return ret, m.WrapError(
		m.Table("runs").
			Joins("inner join tasks on runs.task_id = tasks.id").
			Joins("inner join repositories on tasks.parent_id = repositories.id").
			Where("tasks.pull_request_id = ? and repositories.name = ?", prID, repoName).Find(&ret),
		"retrieving runs for a pull request id",
	)
}
