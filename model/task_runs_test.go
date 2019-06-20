package model

import (
	"fmt"
	"math/rand"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

func (ms *modelSuite) TestRunsForTask(c *check.C) {
	parent, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	fork, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	ref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(ref).Error, check.IsNil)

	tasks := map[int64]map[string]*types.RunSettings{}

	now := time.Now()
	fmt.Print("generating tasks... ")
	for i := 0; i < 1000; i++ {
		count := rand.Intn(100) + 1
		runs := map[string]*types.RunSettings{}

		writeRuns := []*Run{}

		for x := count - 1; x >= 0; x-- {
			name := testutil.RandString(8)
			runs[name] = &types.RunSettings{
				Image:   testutil.RandString(8),
				Command: []string{testutil.RandString(30)},
			}
			writeRuns = append(writeRuns, &Run{
				Name:        name,
				CreatedAt:   time.Now(),
				RunSettings: runs[name],
			})
		}

		ts := &types.TaskSettings{
			Mountpoint: "/tmp",
			Runs:       runs,
		}

		t2 := &Task{
			Ref:          ref,
			Parent:       parent,
			BaseSHA:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			TaskSettings: ts,
		}

		c.Assert(ms.model.Create(t2).Error, check.IsNil)

		for _, run := range writeRuns {
			run.Task = t2
			c.Assert(ms.model.Create(run).Error, check.IsNil)
		}

		tasks[t2.ID] = runs
	}

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing runs for task validation... ")
	var tasklen int
	for i := 0; i < 1000; i++ {
		t, err := ms.model.ListTasks(fork.Name, "", int64(i), 100)
		c.Assert(err, check.IsNil)
		if len(t) > 0 {
			tasklen += len(t)
		} else {
			break
		}
	}

	for task, runs := range tasks {
		x, err := ms.model.CountRunsForTask(task)
		c.Assert(err, check.IsNil)
		c.Assert(x, check.Equals, int64(len(runs)))
		newRuns, err := ms.model.GetRunsForTask(task, 0, 100)
		c.Assert(err, check.IsNil)
		c.Assert(len(newRuns), check.Equals, int(x))
		var lastRunID int64
		for _, run := range newRuns {
			c.Assert(runs[run.Name], check.NotNil)
			c.Assert(run.Task.ID, check.Equals, task)
			if lastRunID > 0 {
				c.Assert(run.ID < lastRunID, check.Equals, true)
				lastRunID = run.ID
			}
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))
}
