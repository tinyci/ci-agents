package db

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCancellationByRef(t *testing.T) {
	m := testInit(t)

	qis, err := m.FillTestQueue(ctx, 1000)
	assert.NilError(t, err)

	for _, qi := range qis {
		run, err := qi.Run().One(ctx, m.db)
		assert.NilError(t, err)

		task, err := run.Task().One(ctx, m.db)
		assert.NilError(t, err)

		sub, err := task.Submission().One(ctx, m.db)
		assert.NilError(t, err)

		headRef, err := sub.HeadRef().One(ctx, m.db)
		assert.NilError(t, err)

		repo, err := headRef.Repository().One(ctx, m.db)
		assert.NilError(t, err)

		assert.NilError(t, m.CancelRefByName(ctx, repo.ID, headRef.Ref))

		runs, err := m.GetRunsForTask(ctx, task.ID, 0, 100)
		assert.NilError(t, err)

		for _, run := range runs {
			task, err := run.Task().One(ctx, m.db)
			assert.NilError(t, err)

			assert.Assert(t, task.Canceled)
		}
	}
}

func TestCancellationByTask(t *testing.T) {
	m := testInit(t)
	qis, err := m.FillTestQueue(ctx, 1000)
	assert.NilError(t, err)

	for _, qi := range qis {
		run, err := qi.Run().One(ctx, m.db)
		assert.NilError(t, err)

		task, err := run.Task().One(ctx, m.db)
		assert.NilError(t, err)

		assert.NilError(t, m.CancelTask(ctx, task.ID))

		runs, err := m.GetRunsForTask(ctx, task.ID, 0, 100)
		assert.NilError(t, err)

		for _, run := range runs {
			task, err := run.Task().One(ctx, m.db)
			assert.NilError(t, err)
			assert.Assert(t, task.Canceled)
		}
	}
}

func TestCancellationByRun(t *testing.T) {
	m := testInit(t)

	qis, err := m.FillTestQueue(ctx, 1000)
	assert.NilError(t, err)

	for _, qi := range qis {

		run, err := qi.Run().One(ctx, m.db)
		assert.NilError(t, err)

		task, err := run.Task().One(ctx, m.db)
		assert.NilError(t, err)

		assert.NilError(t, m.CancelRun(ctx, run.ID))

		runs, err := m.GetRunsForTask(ctx, task.ID, 0, 100)
		assert.NilError(t, err)

		for _, run := range runs {
			task, err := run.Task().One(ctx, m.db)
			assert.NilError(t, err)
			assert.Assert(t, task.Canceled)
		}
	}
}
