package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	gtypes "github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/clients/github"
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

	TasksCount int64 `json:"tasks_count" gorm:"-"`

	Status *bool `json:"status" gorm:"-"`

	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at" gorm:"-"`
	FinishedAt *time.Time `json:"finished_at" gorm:"-"`

	Canceled bool `json:"canceled" gorm:"-"`

	TicketID int64 `json:"ticket_id"`
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

	var status bool
	if s.Status != nil {
		status = *s.Status
	}

	return &gtypes.Submission{
		Id:         s.ID,
		BaseRef:    s.BaseRef.ToProto(),
		HeadRef:    hr,
		User:       pu,
		TasksCount: s.TasksCount,
		CreatedAt:  MakeTimestamp(&s.CreatedAt),
		StartedAt:  MakeTimestamp(s.StartedAt),
		FinishedAt: MakeTimestamp(s.FinishedAt),
		StatusSet:  s.Status != nil,
		Status:     status,
		Canceled:   s.Canceled,
		TicketID:   s.TicketID,
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

	var status *bool
	if gt.StatusSet {
		status = &gt.Status
	}

	created := MakeTime(gt.CreatedAt, false)

	if created.IsZero() {
		// this is a new record and hasn't been updated. Bump the created_at time.
		t := time.Now()
		created = &t
	}

	finished := MakeTime(gt.FinishedAt, true)
	started := MakeTime(gt.StartedAt, true)

	return &Submission{
		ID:         gt.Id,
		User:       u,
		BaseRef:    baseref,
		HeadRef:    headref,
		TasksCount: gt.TasksCount,
		CreatedAt:  *created,
		FinishedAt: finished,
		StartedAt:  started,
		Status:     status,
		Canceled:   gt.Canceled,
		TicketID:   gt.TicketID,
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
	if err := m.WrapError(obj.Find(&subs), "listing submissions"); err != nil {
		return nil, err
	}

	return subs, m.assignSubmissionPostFetch(subs)
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
	obj := m.Group("submissions.id").Order("submissions.id DESC").
		Joins("inner join refs on refs.id = submissions.base_ref_id or refs.id = submissions.head_ref_id")

	if repository != "" && sha != "" {
		ref, err := m.GetRefByNameAndSHA(repository, sha)
		if err != nil {
			return nil, err
		}
		obj = obj.Where("submissions.base_ref_id = ? or submissions.head_ref_id = ?", ref.ID, ref.ID)
	} else if repository != "" {
		var err *errors.Error
		repo, err := m.GetRepositoryByName(repository)
		if err != nil {
			return nil, err
		}
		obj = obj.Where("refs.repository_id = ?", repo.ID)
	}

	return obj, nil
}

func (m *Model) assignSubmissionPostFetch(subs []*Submission) *errors.Error {
	idmap := map[int64]*Submission{}
	ids := []int64{}
	for _, sub := range subs {
		idmap[sub.ID] = sub
		ids = append(ids, sub.ID)
	}

	rows, err := m.Raw("select distinct submission_id, count(*) from tasks where submission_id in (?) group by submission_id", ids).Rows()
	if err != nil {
		return errors.New(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, count int64
		if err := rows.Scan(&id, &count); err != nil {
			return errors.New(err)
		}

		idmap[id].TasksCount = count
	}

	return m.populateStates(ids, idmap)
}

func (m *Model) selectOne(query string, id int64, out interface{}) *errors.Error {
	rows, err := m.Raw(query, id).Rows()
	if err != nil {
		return errors.New(err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(out); err != nil {
			return errors.New(err)
		}
	}

	return nil
}

func (m *Model) gatherStates(ids []int64) (map[int64][]*bool, map[int64]bool, *errors.Error) {
	rows, err := m.Raw("select submission_id, status, canceled from tasks where submission_id in (?)", ids).Rows()
	if err != nil {
		return nil, nil, errors.New(err)
	}
	defer rows.Close()

	overallStatus := map[int64][]*bool{}
	submissionCanceled := map[int64]bool{}

	for rows.Next() {
		var (
			id       int64
			status   *bool
			canceled bool
		)

		if err := rows.Scan(&id, &status, &canceled); err != nil {
			return nil, nil, errors.New(err)
		}

		if _, ok := overallStatus[id]; !ok {
			overallStatus[id] = []*bool{}
		}

		if _, ok := submissionCanceled[id]; canceled && !ok {
			submissionCanceled[id] = true
		} else if !canceled && ok {
			submissionCanceled[id] = false
		}

		overallStatus[id] = append(overallStatus[id], status)
	}

	return overallStatus, submissionCanceled, nil
}

func (m *Model) populateStates(ids []int64, idmap map[int64]*Submission) *errors.Error {
	overallStatus, submissionCanceled, err := m.gatherStates(ids)
	if err != nil {
		return err
	}

	for id, states := range overallStatus {
		failed := false
		unfinished := false
		for _, status := range states {
			if status == nil {
				unfinished = true
				idmap[id].Status = nil
				break
			}

			if !*status {
				failed = true
			}
		}

		if !unfinished {
			f := !failed
			idmap[id].Status = &f

			var t *time.Time
			if err := m.selectOne("select max(finished_at) from tasks where submission_id = ?", id, &t); err != nil {
				return errors.New(err)
			}

			idmap[id].FinishedAt = t
		}

		var t *time.Time
		if err := m.selectOne("select min(started_at) from tasks where submission_id = ? and started_at is not null", id, &t); err != nil {
			return errors.New(err)
		}

		idmap[id].StartedAt = t
		idmap[id].Canceled = submissionCanceled[id]
	}

	return nil
}

// SubmissionListForRepository returns a list of submissions with pagination. If ref
// is non-nil, it will isolate to the ref only and ignore the repo.
func (m *Model) SubmissionListForRepository(repo, sha string, page, perPage int64) ([]*Submission, *errors.Error) {
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
	if err := m.WrapError(obj.Find(&subs), "listing submissions for repository"); err != nil {
		return nil, err
	}

	return subs, m.assignSubmissionPostFetch(subs)
}

func (m *Model) submissionTasksQuery(sub *Submission) *gorm.DB {
	return m.Order("tasks.id DESC").
		Joins("inner join submissions on tasks.submission_id = submissions.id").
		Where("submissions.id = ?", sub.ID)
}

// TasksForSubmission returns all the tasks for a given submission.
func (m *Model) TasksForSubmission(sub *Submission, page, perPage int64) ([]*Task, *errors.Error) {
	tasks := []*Task{}

	obj := m.submissionTasksQuery(sub).Offset(page * perPage).Limit(perPage)
	if err := m.WrapError(obj.Find(&tasks), "listing tasks for a submission"); err != nil {
		return nil, err
	}

	return tasks, m.assignRunCountsToTask(tasks)
}

// GetSubmissionByID returns a submission by internal identifier
func (m *Model) GetSubmissionByID(id int64) (*Submission, *errors.Error) {
	sub := &Submission{}
	if err := m.WrapError(m.Where("submissions.id = ?", id).First(sub), "getting submission by id"); err != nil {
		return nil, err
	}
	return sub, m.assignSubmissionPostFetch([]*Submission{sub})
}

// CancelSubmissionByID cancels all related tasks (and thusly, runs) in a submission that are outstanding.
func (m *Model) CancelSubmissionByID(id int64, baseURL string, client github.Client) *errors.Error {
	sub, err := m.GetSubmissionByID(id)
	if err != nil {
		return err
	}

	var tasks []*Task

	if err := m.WrapError(m.submissionTasksQuery(sub).Find(&tasks), "canceling tasks for a submission"); err != nil {
		return err
	}

	for _, task := range tasks {
		if err := m.CancelTaskByID(task.ID, baseURL, client); err != nil {
			// we want to skip returning this error so all tasks have the chance to be canceled.
			fmt.Println(err) // FIXME can't log here. need to fix that.
		}
	}

	return nil
}
