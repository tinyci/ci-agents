package model

import (
	"encoding/json"
	"fmt"
	"time"

	"errors"

	"github.com/jinzhu/gorm"
	gtypes "github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// Task is the organizational unit of a single source-controlled directory. It
// contains many runs; it the settings are kept in a file named `task.yml` and
// lives in the directory it is testing.
type Task struct {
	ID int64 `gorm:"primary_key" json:"id"`

	Path string `json:"path"`

	TaskSettingsJSON []byte `gorm:"column:task_settings" json:"-"`

	Canceled   bool       `json:"canceled"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	Status     *bool      `json:"status,omitempty"`

	TaskSettings *types.TaskSettings `json:"settings"`

	Runs int64 `json:"runs" gorm:"-"`

	Submission   *Submission `gorm:"association_autoupdate:false" json:"submission"`
	SubmissionID int64       `json:"-"`
}

// NewTaskFromProto converts the proto representation to the task type.
func NewTaskFromProto(gt *gtypes.Task) (*Task, error) {
	var sub *Submission
	if gt.Submission != nil {
		var err error
		sub, err = NewSubmissionFromProto(gt.Submission)
		if err != nil {
			return nil, err
		}
	}

	return &Task{
		ID:           gt.Id,
		Path:         gt.Path,
		Canceled:     gt.Canceled,
		FinishedAt:   MakeTime(gt.FinishedAt, true),
		StartedAt:    MakeTime(gt.StartedAt, true),
		CreatedAt:    *MakeTime(gt.CreatedAt, false),
		Status:       MakeStatus(gt.Status, gt.StatusSet),
		TaskSettings: types.NewTaskSettingsFromProto(gt.Settings),
		Runs:         gt.Runs,
		Submission:   sub,
	}, nil
}

// ToProto converts the task to a protobuf representation
func (t *Task) ToProto() *gtypes.Task {
	var status, set bool
	if t.Status != nil {
		status = *t.Status
		set = true
	}

	var sub *gtypes.Submission

	if t.Submission != nil {
		sub = t.Submission.ToProto()
	}

	return &gtypes.Task{
		Id:         t.ID,
		Path:       t.Path,
		Canceled:   t.Canceled,
		FinishedAt: MakeTimestamp(t.FinishedAt),
		StartedAt:  MakeTimestamp(t.StartedAt),
		CreatedAt:  MakeTimestamp(&t.CreatedAt),
		Status:     status,
		StatusSet:  set,
		Settings:   t.TaskSettings.ToProto(),
		Runs:       t.Runs,
		Submission: sub,
	}
}

// Validate ensures all parameters are set properly.
func (t *Task) Validate() error {
	if t.Submission == nil {
		return errors.New("invalid submission in task")
	}

	if t.TaskSettings == nil {
		return errors.New("task settings are missing")
	}

	return t.TaskSettings.Validate(true)
}

// AfterFind validates the output from the database before releasing it to the
// hook chain
func (t *Task) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal(t.TaskSettingsJSON, &t.TaskSettings); err != nil {
		return utils.WrapError(err, "unpacking task settings for task %d", t.ID)
	}

	if err := t.Validate(); err != nil {
		return utils.WrapError(err, "reading task id %d", t.ID)
	}

	return nil
}

// BeforeCreate just calls BeforeSave.
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	return t.BeforeSave(tx)
}

// BeforeSave is a gorm hook to marshal the token JSON before saving the record
func (t *Task) BeforeSave(tx *gorm.DB) error {
	if err := t.Validate(); err != nil {
		return utils.WrapError(err, "saving task id %d", t.ID)
	}

	var err error
	t.TaskSettingsJSON, err = json.Marshal(t.TaskSettings)
	if err != nil {
		return utils.WrapError(err, "marshaling settings for task id %d", t.ID)
	}

	return nil
}

// CancelTasksForPR cancels all tasks for a PR.
func (m *Model) CancelTasksForPR(repository string, prID int64, baseURL string) error {
	tasks := []*Task{}

	err := m.WrapError(m.
		Joins("inner join submissions on submissions.id = tasks.submission_id").
		Joins("inner join refs on refs.id = submissions.base_ref_id").
		Joins("inner join repositories on repositories.id = refs.repository_id").
		Where("repositories.name = ? and tasks.pull_request_id = ?", repository, prID).Find(&tasks), "locating pull request tasks")
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Submission.BaseRef.Repository.Owner != nil {
			client := github.NewClientFromAccessToken(task.Submission.BaseRef.Repository.Owner.Token.Token)
			if task.FinishedAt != nil {
				continue
			}

			if err := m.CancelTask(task, baseURL, client); err != nil {
				return err
			}
		}
	}

	return nil
}

// CancelTaskByID cancels a task by ID
func (m *Model) CancelTaskByID(id int64, baseURL string, gh github.Client) error {
	var task Task
	if err := m.WrapError(m.Where("id = ?", id).First(&task), "locating task for cancellation"); err != nil {
		return err
	}

	return m.CancelTask(&task, baseURL, gh)
}

// CancelTask finds the queue items and runs for the task, removes them,
// cancels the associated runs for the task, and finally, saves the task itself. It will
// fail to do all of this if the task is already finished.
func (m *Model) CancelTask(task *Task, baseURL string, gh github.Client) error {
	if task.FinishedAt != nil {
		return fmt.Errorf("task %d was already finished; cannot cancel", task.ID)
	}

	runs := []*Run{}

	if err := m.WrapError(m.Where("task_id = ?", task.ID).Find(&runs), "finding runs by task ID"); err != nil {
		return utils.WrapError(err, "locating runs to be canceled for task %d", task.ID)
	}

	for _, thisRun := range runs {
		if thisRun.Status == nil {
			if err := m.SetRunStatus(thisRun.ID, gh, false, true, baseURL, ""); err != nil {
				return utils.WrapError(err, "setting run state for ID %d", thisRun.ID)
			}
		}
	}

	task.Canceled = true
	now := time.Now()
	var b bool
	task.Status = &b
	task.FinishedAt = &now
	return m.WrapError(m.Save(task), "saving task")
}

// UpdateTaskStatus is triggered when a run state change happens that is *not* a cancellation.
func (m *Model) UpdateTaskStatus(task *Task) error {
	runs := []*Run{}

	if task.FinishedAt != nil && task.Status != nil {
		return nil
	}

	err := m.WrapError(m.Where(`task_id = ?`, task.ID).Order("id DESC").Find(&runs), "looking up runs by task id")
	if err != nil {
		if err == utils.ErrNotFound {
			return nil
		}
		return err
	}

	for _, run := range runs {
		if run.Status == nil || run.FinishedAt == nil {
			return nil
		}
	}

	status := true
	now := time.Now()
	task.FinishedAt = &now
	// i hate myself right now
	for _, run := range runs {
		if !*run.Status {
			status = false
		}
	}
	task.Status = &status

	return m.WrapError(m.Save(task), "saving task")
}

func (m *Model) prepTaskListQuery(repository, sha string) (*gorm.DB, error) {
	db := m.Model(&Task{})

	if repository != "" {
		db = db.
			Joins("inner join submissions on tasks.submission_id = submissions.id").
			Joins("inner join refs on submissions.head_ref_id = refs.id or submissions.base_ref_id = refs.id").
			Joins("inner join repositories on refs.repository_id = repositories.id")

		if sha != "" {
			db = db.Where("refs.sha = ? or refs.ref = ?", sha, sha)
		}

		db = db.Where("repositories.name = ?", repository)
	}

	return db.Order("id DESC"), nil
}

// CountTasks counts all the tasks, optionally based on the repository and sha.
func (m *Model) CountTasks(repository, sha string) (int64, error) {
	var count int64

	db, err := m.prepTaskListQuery(repository, sha)
	if err != nil {
		return 0, err
	}

	dbErr := db.Count(&count).Error
	if dbErr != nil {
		return 0, dbErr
	}

	return count, nil
}

func (m *Model) assignRunCountsToTask(tasks []*Task) error {
	idmap := map[int64]*Task{}
	ids := []int64{}

	for _, task := range tasks {
		idmap[task.ID] = task
		ids = append(ids, task.ID)
	}

	rows, eErr := m.Raw("select distinct task_id, count(*) from runs where task_id in (?) group by task_id", ids).Rows()
	if eErr != nil {
		return eErr
	}
	defer rows.Close()

	for rows.Next() {
		var id, count int64
		if err := rows.Scan(&id, &count); err != nil {
			return err
		}

		idmap[id].Runs = count
	}

	return nil
}

// ListTasks gathers all the tasks based on the page and perPage values. It can optionally filter by repository and SHA.
func (m *Model) ListTasks(repository, sha string, page, perPage int64) ([]*Task, error) {
	page, perPage, err := utils.ScopePaginationInt(page, perPage)
	if err != nil {
		return nil, err
	}

	db, err := m.prepTaskListQuery(repository, sha)
	if err != nil {
		return nil, err
	}

	tasks := []*Task{}

	dbErr := db.Limit(perPage).Offset(page * perPage).Find(&tasks).Error
	if dbErr != nil {
		return nil, dbErr
	}

	return tasks, m.assignRunCountsToTask(tasks)
}
