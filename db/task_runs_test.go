package db

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestRunsForTask(t *testing.T) {
	m := testInit(t)
	parent, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	fork, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	baseref := &models.Ref{
		RepositoryID: parent.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, baseref.Insert(ctx, m.db, boil.Infer()))

	headref := &models.Ref{
		RepositoryID: fork.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, headref.Insert(ctx, m.db, boil.Infer()))

	sub := &models.Submission{BaseRefID: baseref.ID, HeadRefID: null.Int64From(headref.ID)}
	assert.NilError(t, sub.Insert(ctx, m.db, boil.Infer()))

	tasks := map[int64]map[string]*types.RunSettings{}

	tx, err := m.db.Begin()
	assert.NilError(t, err)

	now := time.Now()
	fmt.Print("generating tasks... ")
	for i := 0; i < 1000; i++ {
		count := rand.Intn(100) + 1
		runs := map[string]*types.RunSettings{}

		writeRuns := []*models.Run{}

		for x := count - 1; x >= 0; x-- {
			name := testutil.RandString(8)
			runs[name] = &types.RunSettings{
				Queue:   "default",
				Image:   testutil.RandString(8),
				Command: []string{testutil.RandString(30)},
			}

			content, err := json.Marshal(runs[name])
			assert.NilError(t, err)

			writeRuns = append(writeRuns, &models.Run{
				Name:        name,
				CreatedAt:   time.Now(),
				RunSettings: content,
			})
		}

		ts := &types.TaskSettings{
			Mountpoint: "/tmp",
			Runs:       runs,
		}

		content, err := json.Marshal(ts)
		assert.NilError(t, err)

		t2 := &models.Task{
			TaskSettings: content,
			SubmissionID: sub.ID,
		}

		assert.NilError(t, t2.Insert(ctx, tx, boil.Infer()))

		for _, run := range writeRuns {
			run.TaskID = t2.ID
			assert.NilError(t, run.Insert(ctx, tx, boil.Infer()))
		}

		tasks[t2.ID] = runs
	}

	assert.NilError(t, tx.Commit())

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing runs for task validation... ")
	var tasklen int
	for i := 0; i < 1000; i++ {
		tasks, err := m.ListTasks(ctx, fork.Name, "", int64(i), 100)
		assert.NilError(t, err)
		if len(tasks) > 0 {
			tasklen += len(tasks)
		} else {
			break
		}
	}

	for task, runs := range tasks {
		x, err := m.CountRunsForTask(ctx, task)
		assert.NilError(t, err)
		assert.Assert(t, cmp.Equal(x, int64(len(runs))))
		newRuns, err := m.GetRunsForTask(ctx, task, 0, 100)
		assert.NilError(t, err)
		assert.Assert(t, cmp.Equal(len(newRuns), int(x)))
		var lastRunID int64
		for _, run := range newRuns {
			assert.Assert(t, runs[run.Name] != nil)
			assert.Assert(t, cmp.Equal(run.TaskID, task))
			if lastRunID > 0 {
				assert.Assert(t, run.ID < lastRunID)
				lastRunID = run.ID
			}
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))
}
