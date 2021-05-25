package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

func taskValidateHook(ctx context.Context, db boil.ContextExecutor, task *models.Task) error {
	if task.TaskSettings == nil {
		return errors.New("task settings are missing")
	}

	ts := &types.TaskSettings{}

	if err := task.TaskSettings.Unmarshal(ts); err != nil {
		return err
	}

	return ts.Validate(true)
}

// PutTask commits the task to the database.
func (m *Model) PutTask(ctx context.Context, task *models.Task) error {
	return task.Insert(ctx, m.db, boil.Infer())
}

// ListTasks gathers all the tasks based on the page and perPage values. It can optionally filter by repository and SHA.
func (m *Model) ListTasks(ctx context.Context, repository, sha string, page, perPage int64) ([]*models.Task, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	mods, err := m.prepTaskListQuery(repository, sha)
	if err != nil {
		return nil, err
	}

	return models.Tasks(append(mods, qm.Limit(ppg), qm.Offset(pg*ppg))...).All(ctx, m.db)
}

func (m *Model) prepTaskListQuery(repository, sha string) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	if repository != "" {
		mods = append(mods,
			qm.InnerJoin("submissions on tasks.submission_id = submissions.id"),
			qm.InnerJoin("refs on submissions.head_ref_id = refs.id or submissions.base_ref_id = refs.id"),
			qm.InnerJoin("repositories on refs.repository_id = repositories.id"),
		)

		if sha != "" {
			mods = append(mods, qm.Where("refs.sha = ? or refs.ref = ?", sha, sha))
		}

		mods = append(mods, qm.Where("repositories.name = ?", repository))
	}

	return append(mods, qm.OrderBy("id DESC")), nil
}

// GetRunsForTask is just a join of all runs that belong to a task.
func (m *Model) GetRunsForTask(ctx context.Context, id, page, perPage int64) ([]*models.Run, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	return models.Runs(qm.OrderBy("runs.id DESC"), qm.Limit(ppg), qm.Offset(pg*ppg), models.RunWhere.TaskID.EQ(id)).All(ctx, m.db)
}

// CountRunsForTask retrieves the total count of runs for the given task.
func (m *Model) CountRunsForTask(ctx context.Context, id int64) (int64, error) {
	task, err := models.Tasks(models.TaskWhere.ID.EQ(id)).One(ctx, m.db)
	if err != nil {
		return 0, err
	}

	return task.Runs().Count(ctx, m.db)
}

// CancelTask finds the queue items and runs for the task, removes them,
// cancels the associated runs for the task, and finally, saves the task itself. It will
// fail to do all of this if the task is already finished.
func (m *Model) CancelTask(ctx context.Context, taskID int64) error {
	task, err := models.FindTask(ctx, m.db, taskID)
	if err != nil {
		return err
	}

	if task.FinishedAt.Valid {
		return fmt.Errorf("task %d was already finished; cannot cancel", task.ID)
	}

	runs, err := task.Runs().All(ctx, m.db)
	if err != nil {
		return utils.WrapError(err, "locating runs to be canceled for task %d", task.ID)
	}

	for _, thisRun := range runs {
		if !thisRun.Status.Valid {
			if err := m.SetRunStatus(ctx, thisRun.ID, false); err != nil {
				fmt.Println(err)
			}
		}
	}

	task.Canceled = true
	task.Status = null.BoolFrom(false)
	task.FinishedAt = null.TimeFrom(time.Now())

	_, err = task.Update(ctx, m.db, boil.Infer())
	return err
}

// CancelTaskForPR cancels a task for a given pull request ID, by repository.
func (m *Model) CancelTaskForPR(ctx context.Context, repoName string, prID int64) error {
	tasks, err := models.Tasks(
		qm.InnerJoin("submissions on submissions.id = tasks.submission_id"),
		qm.InnerJoin("refs on refs.id = submissions.base_ref_id"),
		qm.InnerJoin("repositories on repositories.id = refs.repository_id"),
		qm.Where("repositories.name = ? and tasks.pull_request_id = ?", repoName, prID)).All(ctx, m.db)

	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := m.CancelTask(ctx, task.ID); err != nil {
			return err
		}
	}

	return nil
}

// CountTasks counts the task for a given repository and sha.
func (m *Model) CountTasks(ctx context.Context, repoName, sha string) (int64, error) {
	return models.Tasks(
		qm.InnerJoin("submissions on tasks.submission_id = submissions.id"),
		qm.InnerJoin("refs on submissions.head_ref_id = refs.id or submissions.base_ref_id = refs.id"),
		qm.InnerJoin("repositories on refs.repository_id = repositories.id"),
		qm.Where("refs.sha = ? or refs.ref = ?", sha, sha),
		qm.Where("repositories.name = ?", repoName),
	).Count(ctx, m.db)
}

// UpdateTaskStatus is triggered when a run state change happens that is *not* a cancellation.
func (m *Model) UpdateTaskStatus(ctx context.Context, taskID int64, status null.Bool, finishedAt null.Time) error {
	task, err := models.Tasks(models.TaskWhere.ID.EQ(taskID)).One(ctx, m.db)
	if err != nil {
		return err
	}

	if task.FinishedAt.Valid && task.Status.Valid {
		return nil
	}

	runs, err := task.Runs(qm.OrderBy("id DESC")).All(ctx, m.db)
	if err != nil {
		return err
	}

	for _, run := range runs {
		if !run.Status.Valid || !run.FinishedAt.Valid {
			return nil
		}
	}

	task.FinishedAt = finishedAt
	task.Status = status

	_, err = task.Update(ctx, m.db, boil.Infer())
	return err
}
