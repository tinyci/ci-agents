package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func runValidateHook(ctx context.Context, db boil.ContextExecutor, r *models.Run) error {
	if r.Name == "" {
		return errors.New("missing name")
	}

	if r.RunSettings == nil {
		return errors.New("no settings provided")
	}

	rs := &types.RunSettings{}

	if err := json.Unmarshal(r.RunSettings, rs); err != nil {
		return err
	}

	return rs.Validate()
}

// GetRun retrieves a run by ID
func (m *Model) GetRun(ctx context.Context, runID int64) (*models.Run, error) {
	return models.FindRun(ctx, m.db, runID)
}

// PutRun stores a new run
func (m *Model) PutRun(ctx context.Context, run *models.Run) error {
	return run.Insert(ctx, m.db, boil.Infer())
}

// GetTaskForRun returns the task for the provided runID
func (m *Model) GetTaskForRun(ctx context.Context, runID int64) (*models.Task, error) {
	run, err := models.FindRun(ctx, m.db, runID)
	if err != nil {
		return nil, err
	}

	return run.Task().One(ctx, m.db)
}

// GetOwnerForRun retrieves the owner of a run's repository.
func (m *Model) GetOwnerForRun(ctx context.Context, runID int64) (*models.User, error) {
	return models.Users(
		qm.InnerJoin("repositories on repositories.owner_id = users.id"),
		qm.InnerJoin("refs on refs.repository_id = repositories.id"),
		qm.InnerJoin("submissions on submissions.base_ref_id = refs.id"),
		qm.InnerJoin("tasks on tasks.submission_id = submissions.id"),
		qm.InnerJoin("runs on runs.task_id = tasks.id"),
	).One(ctx, m.db)
}

// SetRunStatus sets the status for the run; fails if it is already set.
func (m *Model) SetRunStatus(ctx context.Context, runID int64, status bool) error {
	run, err := models.Runs(models.RunWhere.ID.EQ(runID)).One(ctx, m.db)
	if err != nil {
		return err
	}

	qi, err := models.QueueItems(models.QueueItemWhere.RunID.EQ(runID)).One(ctx, m.db)
	if err != nil {
		return err
	}

	if _, err := qi.Delete(ctx, m.db); err != nil {
		return err
	}

	if run.Status.Valid && run.FinishedAt.Valid {
		return fmt.Errorf("status already set for run %d", runID)
	}

	run.Status = null.BoolFrom(status)
	now := time.Now()
	run.FinishedAt = null.TimeFrom(now)

	if _, err := run.Update(ctx, m.db, boil.Infer()); err != nil {
		return err
	}

	return m.UpdateTaskStatus(ctx, run.TaskID, run.Status, run.FinishedAt)
}

// RunDetail contains a number of parameters from inner joins in the run that
// would be hard to get out otherwise.
type RunDetail struct {
	Run     *models.Run
	Owner   string
	Repo    string
	HeadSHA string
}

// GetRunDetail returns a RunDetail based on the runID passed. Used for a
// number of focused operations in the datasvc.
func (m *Model) GetRunDetail(ctx context.Context, runID int64) (*RunDetail, error) {
	run, err := models.FindRun(ctx, m.db, runID)
	if err != nil {
		return nil, err
	}

	task, err := run.Task().One(ctx, m.db)
	if err != nil {
		return nil, err
	}

	sub, err := task.Submission().One(ctx, m.db)
	if err != nil {
		return nil, err
	}

	baseref, err := sub.BaseRef().One(ctx, m.db)
	if err != nil {
		return nil, err
	}

	headref, err := sub.HeadRef().One(ctx, m.db)
	if err != nil {
		baseref = headref
	}

	repo, err := baseref.Repository().One(ctx, m.db)
	if err != nil {
		return nil, err
	}

	owner, repoName, err := utils.OwnerRepo(repo.Name)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid repository for run %d: %v", err, run.ID, repo.Name)
	}

	return &RunDetail{
		Run:     run,
		Owner:   owner,
		Repo:    repoName,
		HeadSHA: headref.Sha,
	}, nil
}

// RunsForTicket all the runs that belong to a repository's PR.
func (m *Model) RunsForTicket(ctx context.Context, repoName string, ticketID int) ([]*models.Run, error) {
	return models.Runs(
		qm.InnerJoin("tasks on runs.task_id = tasks.id"),
		qm.InnerJoin("submissions on submissions.id = tasks.submission_id"),
		qm.InnerJoin("refs on refs.id = submissions.base_ref_id"),
		qm.InnerJoin("repositories on refs.repository_id = repositories.id"),
		qm.Where("submissions.ticket_id = ? and repositories.name = ?", ticketID, repoName),
	).All(ctx, m.db)
}

// RunTotalCount returns the number of items in the runs table
func (m *Model) RunTotalCount(ctx context.Context) (int64, error) {
	return models.Runs().Count(ctx, m.db)
}

// RunList returns a list of runs with pagination.
func (m *Model) RunList(ctx context.Context, page, perPage int64, repository, sha string) ([]*models.Run, error) {
	if repository != "" {
		repo, err := m.GetRepositoryByName(ctx, repository)
		if err != nil {
			return nil, err
		}

		var r *models.Ref

		if sha != "" {
			r, err = m.GetRefByNameAndSHA(ctx, repository, sha)
			if err != nil {
				return nil, err
			}
		}

		return m.RunListForRepository(ctx, repo, r, page, perPage)
	}

	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	mods := []qm.QueryMod{qm.Offset(pg * ppg), qm.Limit(ppg), qm.OrderBy("runs.id DESC"), qm.GroupBy("runs.id")}
	return models.Runs(mods...).All(ctx, m.db)
}

var forRepositoryMods = []qm.QueryMod{
	qm.InnerJoin("tasks on tasks.id = runs.task_id"),
	qm.InnerJoin("submissions on submissions.id = tasks.submission_id"),
	qm.InnerJoin("refs on refs.id = submissions.head_ref_id or refs.id = submissions.base_ref_id"),
	qm.InnerJoin("repositories on repositories.id = refs.repository_id"),
}

// RunListForRepository returns a list of queue items with pagination. If ref
// is non-nil, it will isolate to the ref only.
func (m *Model) RunListForRepository(ctx context.Context, repo *models.Repository, ref *models.Ref, page, perPage int64) ([]*models.Run, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	mods := append([]qm.QueryMod{
		qm.Offset(pg * ppg),
		qm.Limit(ppg),
		qm.OrderBy("runs.id DESC"),
	}, forRepositoryMods...)

	if ref != nil {
		mods = append(mods, qm.Where("repositories.id = ? and refs.id = ?", repo.ID, ref.ID))
	} else {
		mods = append(mods, qm.Where("repositories.id = ?", repo.ID))
	}

	return models.Runs(mods...).All(ctx, m.db)
}

// RunTotalCountForRepository returns the total count of all runs for a given repository
func (m *Model) RunTotalCountForRepository(ctx context.Context, repo *models.Repository) (int64, error) {
	return models.Runs(append(forRepositoryMods, qm.Where("repositories.id = ?", repo.ID))...).Count(ctx, m.db)
}

// RunTotalCountForRepositoryAndSHA counts runs by repository and sha
func (m *Model) RunTotalCountForRepositoryAndSHA(ctx context.Context, repo *models.Repository, sha string) (int64, error) {
	ref, err := m.GetRefByNameAndSHA(ctx, repo.Name, sha)
	if err != nil {
		return 0, err
	}

	row := models.Runs(append(forRepositoryMods, qm.Where("repositories.id = ? and refs.id = ?", repo.ID, ref.ID), qm.Select("count(runs.id)"))...).QueryRowContext(ctx, m.db)
	var feh int64

	if err := row.Scan(&feh); err != nil {
		return 0, err
	}
	fmt.Println(feh)
	return feh, nil
}
