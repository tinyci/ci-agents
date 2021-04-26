package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/utils"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Validate the item. if passed true, will validate for creation scenarios
func queueItemValidateHook(ctx context.Context, db boil.ContextExecutor, qi *models.QueueItem) error {
	if qi.QueueName == "" {
		return errors.New("queue name was empty")
	}

	return nil
}

// QueueTotalCount returns the number of items in the queue
func (m *Model) QueueTotalCount(ctx context.Context) (int64, error) {
	return models.QueueItems().Count(ctx, m.db)
}

// QueueList returns a list of queue items with pagination.
func (m *Model) QueueList(ctx context.Context, page, perPage int64) ([]*models.QueueItem, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	return models.QueueItems(qm.OrderBy("id DESC"), qm.Limit(ppg), qm.Offset(pg*ppg)).All(ctx, m.db)
}

func getQueueRepoQueryMods(repoID int64) []qm.QueryMod {
	return []qm.QueryMod{
		qm.InnerJoin("runs on run_id = runs.id"),
		qm.InnerJoin("tasks on runs.task_id = tasks.id"),
		qm.InnerJoin("submissions on submissions.id = tasks.submission_id"),
		qm.InnerJoin("refs on refs.id = submissions.base_ref_id"),
		qm.InnerJoin("repositories on refs.repository_id = repositories.id"),
		qm.Where("repositories.id = ?", repoID),
	}
}

// QueueListForRepository returns a list of queue items with pagination.
func (m *Model) QueueListForRepository(ctx context.Context, repoID, page, perPage int64) ([]*models.QueueItem, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	mods := getQueueRepoQueryMods(repoID)
	return models.QueueItems(
		append(mods, qm.Offset(pg*ppg), qm.Limit(ppg), qm.OrderBy("id DESC"))...,
	).All(ctx, m.db)
}

// QueueTotalCountForRepository returns the number of items in the queue where
// the parent fork matches the repository name given
func (m *Model) QueueTotalCountForRepository(ctx context.Context, repoID int64) (int64, error) {
	return models.QueueItems(getQueueRepoQueryMods(repoID)...).Count(ctx, m.db)
}

// NextQueueItem returns the next item in the named queue. If for some reason the
// queueName is an empty string, the string `default` will be used instead.
func (m *Model) NextQueueItem(ctx context.Context, runningOn string, queueName string) (qi *models.QueueItem, retErr error) {
	if queueName == "" {
		queueName = "default"
	}

	if runningOn == "" {
		return nil, errors.New("no runner hostname provided")
	}

	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer func() {
		if retErr != nil {
			tx.Rollback()
		}
	}()

	if _, err := tx.Exec("lock table queue_items"); err != nil {
		return nil, err
	}

	qi, err = models.QueueItems(qm.Limit(1), qm.OrderBy("id"), qm.Where("queue_name = ? and not running", queueName)).One(ctx, tx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	run, err := qi.Run().One(ctx, tx)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	run.StartedAt = null.TimeFrom(t)
	run.RanOn = null.StringFrom(runningOn)
	if _, err := run.Update(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}

	task, err := run.Task().One(ctx, tx)
	if err != nil {
		return nil, err
	}

	if !task.StartedAt.Valid {
		task.StartedAt = null.TimeFrom(t)
		if _, err := task.Update(ctx, tx, boil.Infer()); err != nil {
			return nil, err
		}
	}

	qi.StartedAt = null.TimeFrom(t)
	qi.Running = true
	qi.RunningOn = null.StringFrom(runningOn)

	if _, err := qi.Update(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}

	return qi, tx.Commit()
}

// QueuePipelineAdd adds a group of queue items in a transaction.
func (m *Model) QueuePipelineAdd(ctx context.Context, qis []*models.QueueItem) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, qi := range qis {
		if err := qi.Insert(ctx, tx, boil.Infer()); err != nil {
			return err
		}
	}

	return tx.Commit()
}
