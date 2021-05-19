package db

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// CancelRefByName is used in auto cancellation on new queue item arrivals. It finds
// all the runs associated with the ref and cancels them if they are in queue
// still.
// Do note that it does not match the SHA; more often than not this is caused
// by an --amend + force push which updates the SHA, or a new commit which also
// changes the SHA. The name is the only reliable data in this case.
func (m *Model) CancelRefByName(ctx context.Context, repoID int64, refName string) error {
	repo, err := models.FindRepository(ctx, m.db, repoID)
	if err != nil {
		return err
	}

	gh := &github.Repository{}
	if err := json.Unmarshal(repo.Github, gh); err != nil {
		return err
	}

	mb := gh.GetDefaultBranch()
	if refName == mb || (mb == "" && refName == "heads/master") { // FIXME constantize this reference.
		return nil
	}

	tasks, err := models.Tasks(
		qm.Where("refs.ref = ? and refs.repository_id = ? and tasks.status is null and tasks.finished_at is null", refName, repoID),
		qm.InnerJoin("submissions on submissions.id = tasks.submission_id"),
		qm.InnerJoin("refs on refs.id = submissions.head_ref_id"),
	).All(ctx, m.db)
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

// CancelRun is a thin facade for the CancelTask functionality, loading the record from the ID provided.
func (m *Model) CancelRun(ctx context.Context, runID int64) error {
	run, err := models.FindRun(ctx, m.db, runID)
	if err != nil {
		return err
	}

	task, err := run.Task().One(ctx, m.db)
	if err != nil {
		return err
	}

	return m.CancelTask(ctx, task.ID)
}
