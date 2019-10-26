package model

import (
	"fmt"
	"math/rand"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

func (ms *modelSuite) TestTaskValidate(c *check.C) {
	parent, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	fork, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)
	baseref := &Ref{
		Repository: parent,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(baseref).Error, check.IsNil)

	headref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(headref).Error, check.IsNil)

	sub := &Submission{
		BaseRef: baseref,
		HeadRef: headref,
	}
	c.Assert(ms.model.Save(sub).Error, check.IsNil)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	failures := []struct {
		sub          *Submission
		TaskSettings *types.TaskSettings
	}{
		{nil, ts},
		{sub, nil},
	}

	for i, failure := range failures {
		t := &Task{
			Submission:   failure.sub,
			TaskSettings: failure.TaskSettings,
		}

		c.Assert(ms.model.Create(t).Error, check.NotNil, check.Commentf("iteration %d", i))
		c.Assert(ms.model.Save(t).Error, check.NotNil, check.Commentf("iteration %d", i))
	}

	t := &Task{
		Submission:   sub,
		TaskSettings: ts,
	}

	c.Assert(ms.model.Save(t).Error, check.IsNil)

	t2 := &Task{}
	c.Assert(ms.model.Where("id = ?", t.ID).First(t2).Error, check.IsNil)
	c.Assert(t2.ID, check.Equals, t.ID)
	c.Assert(len(t2.TaskSettings.Runs), check.Equals, 1)
	c.Assert(t2.Submission.BaseRef.Repository.Name, check.Equals, t.Submission.BaseRef.Repository.Name)
}

func (ms *modelSuite) TestTaskList(c *check.C) {
	parent, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	fork, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	baseref := &Ref{
		Repository: parent,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(baseref).Error, check.IsNil)

	headref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(headref).Error, check.IsNil)

	sub := &Submission{
		BaseRef: baseref,
		HeadRef: headref,
	}

	c.Assert(ms.model.Save(sub).Error, check.IsNil)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	t := &Task{
		Submission:   sub,
		TaskSettings: ts,
	}

	c.Assert(ms.model.Save(t).Error, check.IsNil)

	tasks, err := ms.model.ListTasks(fork.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(tasks[0].Submission.BaseRef.SHA, check.Not(check.Equals), "")

	count, err := ms.model.CountTasks("", "")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1))

	count, err = ms.model.CountTasks(fork.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1))

	us, err := ms.CreateUsers(1)
	c.Assert(err, check.IsNil)
	tasks, err = ms.model.ListSubscribedTasksForUser(us[0].ID, 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(tasks), check.Equals, 0)

	c.Assert(ms.model.AddSubscriptionsForUser(us[0], []*Repository{parent}), check.IsNil)

	tasks, err = ms.model.ListSubscribedTasksForUser(us[0].ID, 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(tasks), check.Not(check.Equals), 0)
}

func (ms *modelSuite) TestTaskListSHAList(c *check.C) {
	parent, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	fork, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	baseref := &Ref{
		Repository: parent,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(baseref).Error, check.IsNil)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	shas := map[string]int{}

	now := time.Now()
	fmt.Print("generating tasks... ")
	for i := 0; i < 1000; i++ {
		count := rand.Intn(50)
		sha := testutil.RandString(40)
		shas[sha] = count
		ref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		c.Assert(ms.model.Save(ref).Error, check.IsNil)

		sub := &Submission{
			BaseRef: baseref,
			HeadRef: ref,
		}
		c.Assert(ms.model.Save(sub).Error, check.IsNil)

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				Submission:   sub,
				TaskSettings: ts,
			}
			c.Assert(ms.model.Save(t2).Error, check.IsNil)
		}
	}

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing single repo multi-SHA... ")
	var tasklen int
	for i := 0; i < 1000; i++ {
		t, err := ms.model.ListTasks(fork.Name, "", int64(i), 100)
		c.Assert(err, check.IsNil)
		if len(t) > 0 {
			tasklen += len(t)
		} else {
			break
		}

		var lastID int64
		for _, tsk := range t {
			if lastID != 0 {
				c.Assert(tsk.ID < lastID, check.Equals, true)
				lastID = tsk.ID
			}
		}
	}

	var totalcount int

	for _, count := range shas {
		totalcount += count
	}

	c.Assert(tasklen, check.Equals, totalcount)

	for sha, count := range shas {
		x, err := ms.model.CountTasks(fork.Name, sha)
		c.Assert(err, check.IsNil)
		c.Assert(x, check.Equals, int64(count))
		tasks, err := ms.model.ListTasks(fork.Name, sha, 0, 100)
		c.Assert(err, check.IsNil)
		c.Assert(len(tasks), check.Equals, count)
		for _, task := range tasks {
			c.Assert(task.Submission.HeadRef.SHA, check.Equals, sha)
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))

	// totalcount is already calculated in this test, so re-use it.

	count, err := ms.model.CountTasks("", "")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(totalcount))
}

func (ms *modelSuite) TestTaskListParents(c *check.C) {
	fork, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	now := time.Now()
	fmt.Print("generating tasks... ")
	parents := map[string]int{}
	for i := 0; i < 1000; i++ {
		count := rand.Intn(100)
		parent, err := ms.CreateRepository()
		c.Assert(err, check.IsNil)
		parents[parent.Name] = count
		sha := testutil.RandString(40)
		baseref := &Ref{
			Repository: parent,
			RefName:    "refs/heads/master",
			SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		}
		c.Assert(ms.model.Save(baseref).Error, check.IsNil)

		headref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		c.Assert(ms.model.Save(headref).Error, check.IsNil)

		sub := &Submission{BaseRef: baseref, HeadRef: headref}
		c.Assert(ms.model.Save(sub).Error, check.IsNil)

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				Submission:   sub,
				TaskSettings: ts,
			}
			c.Assert(ms.model.Save(t2).Error, check.IsNil)
		}
	}

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing multi parent any-SHA... ")
	for parentName, count := range parents {
		x, err := ms.model.CountTasks(parentName, "")
		c.Assert(err, check.IsNil)
		c.Assert(x, check.Equals, int64(count))
		tasks, err := ms.model.ListTasks(parentName, "", 0, 100)
		c.Assert(err, check.IsNil)
		c.Assert(len(tasks), check.Equals, count)
		var lastID int64
		for _, task := range tasks {
			c.Assert(task.Submission.BaseRef.Repository.Name, check.Equals, parentName)
			if lastID != 0 {
				c.Assert(task.ID < lastID, check.Equals, true)
				lastID = task.ID
			}
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))

	var totalcount int64
	for _, count := range parents {
		totalcount += int64(count)
	}

	count, err := ms.model.CountTasks("", "")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, totalcount)
}

func (ms *modelSuite) TestTaskListForks(c *check.C) {
	parent, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	baseref := &Ref{
		Repository: parent,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(baseref).Error, check.IsNil)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	now := time.Now()
	fmt.Print("generating tasks... ")
	forks := map[string]int{}
	for i := 0; i < 1000; i++ {
		count := rand.Intn(100)
		fork, err := ms.CreateRepository()
		c.Assert(err, check.IsNil)
		forks[fork.Name] = count
		sha := testutil.RandString(40)
		headref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		c.Assert(ms.model.Save(headref).Error, check.IsNil)

		sub := &Submission{HeadRef: headref, BaseRef: baseref}

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				TaskSettings: ts,
				Submission:   sub,
			}
			c.Assert(ms.model.Save(t2).Error, check.IsNil)
		}
	}

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing multi fork any-SHA... ")
	for forkName, count := range forks {
		x, err := ms.model.CountTasks(forkName, "")
		c.Assert(err, check.IsNil)
		c.Assert(x, check.Equals, int64(count))
		tasks, err := ms.model.ListTasks(forkName, "", 0, 100)
		c.Assert(err, check.IsNil)
		c.Assert(len(tasks), check.Equals, count)
		var lastID int64
		for _, task := range tasks {
			c.Assert(task.Submission.HeadRef.Repository.Name, check.Equals, forkName)
			if lastID != 0 {
				c.Assert(task.ID < lastID, check.Equals, true)
				lastID = task.ID
			}
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))

	var totalcount int64
	for _, count := range forks {
		totalcount += int64(count)
	}

	count, err := ms.model.CountTasks("", "")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, totalcount)
}
