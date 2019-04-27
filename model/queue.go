package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/utils"
)

// QueueItem represents an item in the queue table in the database.
type QueueItem struct {
	ID int64 `gorm:"primary_key" json:"id"`

	Run       *Run       `json:"run"`
	RunID     int64      `json:"-"`
	Running   bool       `json:"running"`
	RunningOn *string    `json:"running_on,omitempty"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	QueueName string     `json:"queue_name"`
}

// NewQueueItemFromProto converts in the opposite direction of ToProto.
func NewQueueItemFromProto(tqi *types.QueueItem) (*QueueItem, *errors.Error) {
	run, err := NewRunFromProto(tqi.Run)
	if err != nil {
		return nil, err
	}

	qi := &QueueItem{
		ID:        tqi.Id,
		Run:       run,
		Running:   tqi.Running,
		StartedAt: MakeTime(tqi.StartedAt, true),
		QueueName: tqi.QueueName,
	}

	if tqi.RunningOn != "" {
		qi.RunningOn = &tqi.RunningOn
	}

	return qi, nil
}

// ToProto converts the queueitem to the protobuf version
func (qi *QueueItem) ToProto() *types.QueueItem {
	ro := ""

	if qi.RunningOn != nil {
		ro = *qi.RunningOn
	}

	return &types.QueueItem{
		Id:        qi.ID,
		Running:   qi.Running,
		RunningOn: ro,
		StartedAt: MakeTimestamp(qi.StartedAt),
		QueueName: qi.QueueName,
		Run:       qi.Run.ToProto(),
	}
}

// AfterFind validates the output from the database before releasing it to the
// hook chain
func (qi *QueueItem) AfterFind(tx *gorm.DB) error {
	if err := qi.Validate(); err != nil {
		return errors.New(err).Wrapf("reading queue item %d", qi.ID)
	}

	return nil
}

// BeforeCreate just calls BeforeSave.
func (qi *QueueItem) BeforeCreate(tx *gorm.DB) error {
	return qi.BeforeSave(tx)
}

// BeforeSave is a gorm hook to marshal the token JSON before saving the record
func (qi *QueueItem) BeforeSave(tx *gorm.DB) error {
	if err := qi.Validate(); err != nil {
		return errors.New(err).Wrapf("saving queue item %d", qi.ID)
	}

	return nil
}

// Validate the item. if passed true, will validate for creation scenarios
func (qi *QueueItem) Validate() *errors.Error {
	if qi.Run == nil {
		return errors.New("run was nil")
	}

	if qi.QueueName == "" {
		return errors.New("queue name was empty")
	}

	return nil
}

// ValidateRunning ensures the state is running.
func (qi *QueueItem) ValidateRunning() *errors.Error {
	if qi.RunningOn == nil || *qi.RunningOn == "" {
		return errors.New("missing run target hostname")
	}

	if !qi.Running {
		return errors.New("was not flagged running, yet should be")
	}

	return nil
}

// QueueTotalCount returns the number of items in the queue
func (m *Model) QueueTotalCount() (int64, *errors.Error) {
	var ret int64
	return ret, m.WrapError(m.Table("queue_items").Count(&ret), "computing total queue count")
}

// QueueTotalCountForRepository returns the number of items in the queue where
// the parent fork matches the repository name given
func (m *Model) QueueTotalCountForRepository(repo *Repository) (int64, *errors.Error) {
	var ret int64
	return ret, m.WrapError(
		m.Table("queue_items").
			Joins("inner join runs on runs.id = queue_items.run_id").
			Joins("inner join tasks on runs.task_id = tasks.id").
			Joins("inner join repositories on tasks.parent_id = repositories.id").
			Where("repositories.id = ?", repo.ID).
			Count(&ret), "computing repository queue count",
	)
}

// NextQueueItem returns the next item in the named queue. If for some reason the
// queueName is an empty string, the string `default` will be used instead.
func (m *Model) NextQueueItem(runningOn string, queueName string) (*QueueItem, *errors.Error) {
	if queueName == "" {
		queueName = "default"
	}

	if runningOn == "" {
		return nil, errors.New("no runner hostname provided")
	}

	db := m.Begin()
	defer db.Rollback()

	if err := m.WrapError(db.Exec("lock table queue_items"), "locking queue table"); err != nil {
		return nil, err
	}

	qi := &QueueItem{}

	err := m.WrapError(
		db.Preload("Run.Task.Parent").
			Order("id").
			Where("queue_name = ? and not running", queueName).
			First(qi),
		"getting task owners during queue next",
	)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	qi.StartedAt = &t
	qi.Run.StartedAt = &t
	if qi.Run.Task.StartedAt == nil {
		qi.Run.Task.StartedAt = &t
		if err := m.WrapError(db.Save(qi.Run.Task), "updating task started_at"); err != nil {
			return nil, err
		}
	}

	qi.Running = true
	qi.RunningOn = &runningOn

	if err := m.WrapError(db.Save(qi), "updating newly shifted queue item"); err != nil {
		return nil, errors.New(err)
	}

	return qi, m.WrapError(db.Commit(), "committing queue update")
}

// QueueList returns a list of queue items with pagination.
func (m *Model) QueueList(page, perPage int64) ([]*QueueItem, *errors.Error) {
	qis := []*QueueItem{}

	page, perPage, err := utils.ScopePaginationInt(page, perPage)
	if err != nil {
		return nil, err
	}

	return qis, m.WrapError(m.Offset(page*perPage).Limit(perPage).Order("id DESC").Find(&qis), "listing queue")
}

// QueueListForRepository returns a list of queue items with pagination.
func (m *Model) QueueListForRepository(repo *Repository, page, perPage int64) ([]*QueueItem, *errors.Error) {
	qis := []*QueueItem{}

	page, perPage, err := utils.ScopePaginationInt(page, perPage)
	if err != nil {
		return nil, err
	}

	return qis, m.WrapError(
		m.Offset(page*perPage).
			Limit(perPage).
			Order("id DESC").
			Joins("inner join runs on run_id = runs.id").
			Joins("inner join tasks on runs.task_id = tasks.id").
			Joins("inner join repositories on tasks.parent_id = repositories.id").
			Where("repositories.id = ?", repo.ID).
			Find(&qis),
		"listing queue for repository",
	)
}

// QueuePipelineAdd adds a group of queue items in a transaction.
func (m *Model) QueuePipelineAdd(qis []*QueueItem) ([]*QueueItem, *errors.Error) {
	db := m.Begin()
	defer db.Rollback()

	for _, qi := range qis {
		if err := m.WrapError(db.Create(qi), "creating queue item in pipeline add"); err != nil {
			return nil, err
		}
	}

	return qis, m.WrapError(db.Commit(), "committing pipeline queue add")
}
