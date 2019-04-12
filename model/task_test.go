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

	ref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(ref).Error, check.IsNil)

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
		ref          *Ref
		parent       *Repository
		baseSHA      string
		TaskSettings *types.TaskSettings
	}{
		{nil, parent, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ts},
		{ref, nil, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ts},
		{ref, parent, "", ts},
		{ref, parent, "123", ts},
		{ref, parent, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil},
	}

	for i, failure := range failures {
		t := &Task{
			Ref:          failure.ref,
			Parent:       failure.parent,
			BaseSHA:      failure.baseSHA,
			TaskSettings: failure.TaskSettings,
		}

		c.Assert(ms.model.Create(t).Error, check.NotNil, check.Commentf("iteration %d", i))
		c.Assert(ms.model.Save(t).Error, check.NotNil, check.Commentf("iteration %d", i))
	}

	t := &Task{
		Ref:          ref,
		Parent:       parent,
		BaseSHA:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		TaskSettings: ts,
	}

	c.Assert(ms.model.Save(t).Error, check.IsNil)

	t2 := &Task{}
	c.Assert(ms.model.Where("id = ?", t.ID).First(t2).Error, check.IsNil)
	c.Assert(t2.ID, check.Equals, t.ID)
	c.Assert(len(t2.TaskSettings.Runs), check.Equals, 1)
	c.Assert(t2.Parent.Name, check.Equals, t.Parent.Name)
	c.Assert(t2.Ref.Repository.Name, check.Equals, t.Ref.Repository.Name)
	c.Assert(t2.Parent.Name, check.Not(check.Equals), t2.Ref.Repository.Name)
}

func (ms *modelSuite) TestTaskList(c *check.C) {
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
		Ref:          ref,
		Parent:       parent,
		BaseSHA:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		TaskSettings: ts,
	}

	c.Assert(ms.model.Save(t).Error, check.IsNil)

	tasks, err := ms.model.ListTasks(fork.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(tasks[0].BaseSHA, check.Not(check.Equals), "")

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

	ref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(ref).Error, check.IsNil)

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
		count := rand.Intn(100)
		sha := testutil.RandString(40)
		shas[sha] = count
		ref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		c.Assert(ms.model.Save(ref).Error, check.IsNil)

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				Ref:          ref,
				Parent:       parent,
				BaseSHA:      sha,
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
			c.Assert(task.BaseSHA, check.Equals, sha)
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

	ref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(ref).Error, check.IsNil)

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
		ref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		c.Assert(ms.model.Save(ref).Error, check.IsNil)

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				Ref:          ref,
				Parent:       parent,
				BaseSHA:      sha,
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
			c.Assert(task.Parent.Name, check.Equals, parentName)
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

	fork, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	ref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(ref).Error, check.IsNil)

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
		ref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		c.Assert(ms.model.Save(ref).Error, check.IsNil)

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				Ref:          ref,
				Parent:       parent,
				BaseSHA:      sha,
				TaskSettings: ts,
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
			c.Assert(task.Ref.Repository.Name, check.Equals, forkName)
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
