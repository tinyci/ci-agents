package db

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestRunValidate(t *testing.T) {
	m := testInit(t)
	parent, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	fork, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	ref := &models.Ref{
		RepositoryID: fork.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	baseref := &models.Ref{
		RepositoryID: parent.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, ref.Insert(ctx, m.db, boil.Infer()))
	assert.NilError(t, baseref.Insert(ctx, m.db, boil.Infer()))

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Queue:   "default",
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	sub := &models.Submission{
		HeadRefID: null.Int64From(ref.ID),
		BaseRefID: baseref.ID,
		CreatedAt: time.Now(),
		TicketID:  null.Int64From(10),
	}

	assert.NilError(t, sub.Insert(ctx, m.db, boil.Infer()))

	content, err := json.Marshal(ts)
	assert.NilError(t, err)

	task := &models.Task{
		TaskSettings: content,
		SubmissionID: sub.ID,
	}

	assert.NilError(t, task.Insert(ctx, m.db, boil.Infer()))

	failures := []struct {
		name        string
		RunSettings *types.RunSettings
		task        *models.Task
	}{
		{"", ts.Runs["foobar"], task},
		{"foobar", nil, task},
		{"foobar", ts.Runs["foobar"], nil},
	}

	for i, failure := range failures {
		var content []byte
		if failure.RunSettings != nil {
			var err error
			content, err = json.Marshal(failure.RunSettings)
			assert.NilError(t, err)
		}

		var taskID int64
		if failure.task != nil {
			taskID = failure.task.ID
		}

		r := &models.Run{
			Name:        failure.name,
			RunSettings: content,
			TaskID:      taskID,
		}

		assert.Assert(t, r.Insert(ctx, m.db, boil.Infer()) != nil, "iteration %d", i)
	}

	content, err = json.Marshal(ts.Runs["foobar"])
	assert.NilError(t, err)
	r := &models.Run{
		Name:        "foobar",
		RunSettings: content,
		TaskID:      task.ID,
	}

	assert.NilError(t, r.Insert(ctx, m.db, boil.Infer()))
	assert.Assert(t, r.ID != 0)

	r2, err := models.FindRun(ctx, m.db, r.ID)
	assert.NilError(t, err)

	assert.Assert(t, cmp.Equal(r.ID, r2.ID))
	assert.Assert(t, r2.RunSettings != nil)

	runs, err := m.RunsForTicket(ctx, parent.Name, 10)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs), 1))

	runs, err = m.RunsForTicket(ctx, parent.Name, 11)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs), 0))
}

func TestRunList(t *testing.T) {
	m := testInit(t)
	parent, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	fork, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	ref := &models.Ref{
		RepositoryID: fork.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, ref.Insert(ctx, m.db, boil.Infer()))

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Queue:   "default",
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	sub := &models.Submission{
		BaseRefID: ref.ID,
	}

	assert.NilError(t, sub.Insert(ctx, m.db, boil.Infer()))

	content, err := json.Marshal(ts)
	assert.NilError(t, err)

	task := &models.Task{
		TaskSettings: content,
		SubmissionID: sub.ID,
	}

	assert.NilError(t, task.Insert(ctx, m.db, boil.Infer()))

	content, err = json.Marshal(ts.Runs["foobar"])
	assert.NilError(t, err)

	for i := 0; i < 1000; i++ {
		r := &models.Run{
			Name:        "foobar",
			RunSettings: content,
			TaskID:      task.ID,
		}

		assert.NilError(t, r.Insert(ctx, m.db, boil.Infer()))
	}

	count, err := m.RunTotalCount(ctx)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(count, int64(1000)))

	runs, err := m.RunList(ctx, 0, 0, "", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs), 10))

	runs, err = m.RunList(ctx, 0, 10000, "", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs), 100))

	// check for off-by-one in pagination
	runs0, err := m.RunList(ctx, 0, 100, "", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs0), 100))

	runs1, err := m.RunList(ctx, 1, 100, "", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs1), 100))

	// check for no overlap
	for _, run0 := range runs0 {
		for _, run1 := range runs1 {
			assert.Assert(t, run0.ID != run1.ID)
		}
	}

	runs, err = m.RunList(ctx, 0, 0, fork.Name, "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs), 10))

	runs, err = m.RunList(ctx, 0, 0, fork.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") // baseref sha
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs), 10))

	runs, err = m.RunList(ctx, 0, 0, parent.Name, "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(runs), 0))

	runs, err = m.RunList(ctx, 0, 0, "quux/foobar", "")
	assert.Error(t, err, sql.ErrNoRows.Error())
	assert.Assert(t, cmp.Equal(len(runs), 0))
}
