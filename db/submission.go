package db

import (
	"context"
	"fmt"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// NewSubmissionFromMessage creates a blank record populated by the appropriate data
// reaped from the message type in types/submission.go.
func (m *Model) NewSubmissionFromMessage(ctx context.Context, sub *types.Submission) (*models.Submission, error) {
	if err := sub.Validate(); err != nil {
		return nil, utils.WrapError(err, "did not pass validation")
	}

	var (
		u    *models.User
		head *models.Ref
		err  error
	)

	if sub.SubmittedBy != "" {
		u, err = m.FindUserByName(ctx, sub.SubmittedBy)
		if err != nil {
			return nil, utils.WrapError(err, "for use in submission")
		}
	}

	if sub.HeadSHA != "" {
		head, err = m.GetRefByNameAndSHA(ctx, sub.Fork, sub.HeadSHA)
		if err != nil {
			return nil, utils.WrapError(err, "head")
		}
	}

	base, err := m.GetRefByNameAndSHA(ctx, sub.Parent, sub.BaseSHA)
	if err != nil {
		return nil, utils.WrapError(err, "head")
	}

	var id *int64
	if u != nil {
		id = &u.ID
	}

	var headrefID *int64
	if head != nil {
		headrefID = &head.ID
	}

	return &models.Submission{
		UserID:    null.Int64FromPtr(id),
		HeadRefID: null.Int64FromPtr(headrefID),
		BaseRefID: base.ID,
	}, nil
}

// PutSubmission creates the submission.
func (m *Model) PutSubmission(ctx context.Context, sub *models.Submission) error {
	return sub.Insert(ctx, m.db, boil.Infer())
}

// SubmissionList returns a list of submissions with pagination and repo/sha filtering.
func (m *Model) SubmissionList(ctx context.Context, page, perPage int64, repository, sha string) ([]*models.Submission, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	qms, err := m.submissionListQuery(ctx, repository, sha, false)
	if err != nil {
		return nil, err
	}

	qms = append(qms,
		qm.Offset(pg*ppg),
		qm.Limit(ppg),
		qm.GroupBy("submissions.id"),
		qm.OrderBy("submissions.id DESC"),
	)

	subs, err := models.Submissions(qms...).All(ctx, m.db)
	if err != nil {
		return nil, err
	}

	return subs, nil //m.assignSubmissionPostFetch(subs)
}

func (m *Model) submissionListQuery(ctx context.Context, repository, sha string, count bool) ([]qm.QueryMod, error) {
	mods := []qm.QueryMod{}

	if repository != "" && sha != "" {
		ref, err := m.GetRefByNameAndSHA(ctx, repository, sha)
		if err != nil {
			return nil, err
		}

		mods = append(mods, qm.Where("submissions.base_ref_id = ? or submissions.head_ref_id = ?", ref.ID, ref.ID))
	} else if repository != "" {
		repo, err := m.GetRepositoryByName(ctx, repository)
		if err != nil {
			return nil, err
		}

		mods = append(mods,
			qm.Where("refs.repository_id = ?", repo.ID),
			qm.InnerJoin("refs on refs.id = submissions.base_ref_id or refs.id = submissions.head_ref_id"),
		)
	}

	return mods, nil
}

// SubmissionCount counts the number of submissions with an optional repo/sha filter
func (m *Model) SubmissionCount(ctx context.Context, repository, sha string) (int64, error) {
	qms, err := m.submissionListQuery(ctx, repository, sha, true)
	if err != nil {
		return 0, err
	}

	return models.Submissions(qms...).Count(ctx, m.db)
}

// GetSubmissionByID returns a submission by internal identifier
func (m *Model) GetSubmissionByID(ctx context.Context, id int64) (*models.Submission, error) {
	return models.FindSubmission(ctx, m.db, id)
}

func (m *Model) submissionTasksQuery(id int64) []qm.QueryMod {
	return []qm.QueryMod{
		qm.OrderBy("tasks.id"),
		qm.InnerJoin("submissions on tasks.submission_id = submissions.id"),
		qm.Where("submissions.id = ?", id),
	}
}

// TasksForSubmission returns all the tasks for a given submission.
func (m *Model) TasksForSubmission(ctx context.Context, id int64, page, perPage int64) ([]*models.Task, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	return models.Tasks(append(m.submissionTasksQuery(id), qm.Offset(pg*ppg), qm.Limit(ppg))...).All(ctx, m.db)
}

func (m *Model) runsForSubmission(id int64) []qm.QueryMod {
	return []qm.QueryMod{
		qm.InnerJoin("tasks on tasks.id = runs.task_id"),
		qm.InnerJoin("submissions on submissions.id = tasks.submission_id"),
		qm.Where("submissions.id = ?", id),
	}
}

// RunsForSubmission returns all the runs for the provided submission id.
func (m *Model) RunsForSubmission(ctx context.Context, id int64, page, perPage int64) ([]*models.Run, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	return models.Runs(append(m.runsForSubmission(id),
		qm.GroupBy("runs.id"),
		qm.OrderBy("runs.id"),
		qm.Offset(pg*ppg),
		qm.Limit(ppg),
	)...).All(ctx, m.db)
}

// CountRunsForSubmission counts all the runs for the provided submission id.
func (m *Model) CountRunsForSubmission(ctx context.Context, id int64) (int64, error) {
	return models.Runs(m.runsForSubmission(id)...).Count(ctx, m.db)
}

// CancelSubmissionByID cancels all related tasks (and thusly, runs) in a submission that are outstanding.
func (m *Model) CancelSubmissionByID(ctx context.Context, id int64) error {
	sub, err := m.GetSubmissionByID(ctx, id)
	if err != nil {
		return err
	}

	tasks, err := sub.Tasks().All(ctx, m.db)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := m.CancelTask(ctx, task.ID); err != nil {
			// we want to skip returning this error so all tasks have the chance to be canceled.
			fmt.Println(err) // FIXME can't log here. need to fix that.
		}
	}

	return nil
}

//
// func (m *Model) assignSubmissionPostFetch(subs []*models.Submission) error {
// 	idmap := map[int64]*models.Submission{}
// 	ids := []int64{}
// 	for _, sub := range subs {
// 		idmap[sub.ID] = sub
// 		ids = append(ids, sub.ID)
// 	}
//
// 	idmap2, err := m.assignCountFromQuery(`
// 		select
// 		distinct submission_id, count(*)
// 		from tasks
// 		where submission_id in (?)
// 		group by submission_id
// 	`, ids)
// 	if err != nil {
// 		return err
// 	}
//
// 	for id, count := range idmap2 {
// 		idmap[id].TasksCount = count
// 	}
//
// 	idmap2, err = m.assignCountFromQuery(`
// 		select
// 		distinct tasks.submission_id, count(*)
// 		from runs
// 			inner join tasks on tasks.id = runs.task_id
// 		where tasks.submission_id in (?)
// 		group by tasks.submission_id
// 	`, ids)
// 	if err != nil {
// 		return err
// 	}
//
// 	for id, count := range idmap2 {
// 		idmap[id].RunsCount = count
// 	}
//
// 	return m.populateStates(ids, idmap)
// }
//
// func (m *Model) assignCountFromQuery(query string, ids []int64) (map[int64]int64, error) {
// 	idmap := map[int64]int64{}
//
// 	rows, err := m.Raw(query, ids).Rows()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	for rows.Next() {
// 		var id, count int64
// 		if err := rows.Scan(&id, &count); err != nil {
// 			return nil, err
// 		}
//
// 		idmap[id] = count
// 	}
//
// 	return idmap, nil
// }
//
// func (m *Model) populateStates(ids []int64, idmap map[int64]*Submission) error {
// 	overallStatus, submissionCanceled, err := m.gatherStates(ids)
// 	if err != nil {
// 		return err
// 	}
//
// 	for id, states := range overallStatus {
// 		failed := false
// 		unfinished := false
// 		for _, status := range states {
// 			if status == nil {
// 				unfinished = true
// 				idmap[id].Status = nil
// 				break
// 			}
//
// 			if !*status {
// 				failed = true
// 			}
// 		}
//
// 		if !unfinished {
// 			f := !failed
// 			idmap[id].Status = &f
//
// 			var t *time.Time
// 			if err := m.selectOne("select max(finished_at) from tasks where submission_id = ?", id, &t); err != nil {
// 				return err
// 			}
//
// 			idmap[id].FinishedAt = t
// 		}
//
// 		var t *time.Time
// 		if err := m.selectOne("select min(started_at) from tasks where submission_id = ? and started_at is not null", id, &t); err != nil {
// 			return err
// 		}
//
// 		idmap[id].StartedAt = t
// 		idmap[id].Canceled = submissionCanceled[id]
// 	}
//
// 	return nil
// }
