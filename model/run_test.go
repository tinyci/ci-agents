package model

import (
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/types"
)

func (ms *modelSuite) TestRunValidate(c *check.C) {
	parent, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	fork, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	ref := &Ref{
		Repository: fork,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	baseref := &Ref{
		Repository: parent,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ms.model.Save(ref).Error, check.IsNil)
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

	sub := &Submission{
		HeadRef:   ref,
		BaseRef:   baseref,
		CreatedAt: time.Now(),
		TicketID:  10,
	}
	c.Assert(ms.model.Save(sub).Error, check.IsNil)

	task := &Task{
		TaskSettings: ts,
		Submission:   sub,
	}

	c.Assert(ms.model.Save(task).Error, check.IsNil)

	failures := []struct {
		name        string
		RunSettings *types.RunSettings
		task        *Task
	}{
		{"", ts.Runs["foobar"], task},
		{"foobar", nil, task},
		{"foobar", ts.Runs["foobar"], nil},
	}

	for i, failure := range failures {
		r := &Run{
			Name:        failure.name,
			RunSettings: failure.RunSettings,
			Task:        failure.task,
		}
		c.Assert(ms.model.Create(r).Error, check.NotNil, check.Commentf("iteration %d", i))
		c.Assert(ms.model.Save(r).Error, check.NotNil, check.Commentf("iteration %d", i))
	}

	r := &Run{
		Name:        "foobar",
		RunSettings: ts.Runs["foobar"],
		Task:        task,
	}

	c.Assert(ms.model.Save(r).Error, check.IsNil)

	r2 := &Run{}
	c.Assert(ms.model.Where("id = ?", r.ID).First(r2).Error, check.IsNil)
	c.Assert(r.ID, check.Not(check.Equals), 0)
	c.Assert(r.ID, check.Equals, r2.ID)
	c.Assert(r2.RunSettings, check.NotNil)

	runs, err := ms.model.RunsForTicket(parent.Name, 10)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 1)

	runs, err = ms.model.RunsForTicket(parent.Name, 11)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 0)
}

func (ms *modelSuite) TestRunList(c *check.C) {
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

	sub := &Submission{
		BaseRef: ref,
	}

	c.Assert(ms.model.Save(sub).Error, check.IsNil)
	task := &Task{
		TaskSettings: ts,
		Submission:   sub,
	}

	c.Assert(ms.model.Save(task).Error, check.IsNil)

	for i := 0; i < 1000; i++ {
		r := &Run{
			Name:        "foobar",
			RunSettings: ts.Runs["foobar"],
			Task:        task,
		}
		c.Assert(ms.model.Create(r).Error, check.IsNil, check.Commentf("iteration %d", i))
	}

	count, err := ms.model.RunTotalCount()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1000))

	runs, err := ms.model.RunList(0, 0, "", "")
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 10)

	runs, err = ms.model.RunList(0, 10000, "", "")
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 100)

	// check for off-by-one in pagination
	runs0, err := ms.model.RunList(0, 100, "", "")
	c.Assert(err, check.IsNil)
	c.Assert(len(runs0), check.Equals, 100)

	runs1, err := ms.model.RunList(1, 100, "", "")
	c.Assert(err, check.IsNil)
	c.Assert(len(runs1), check.Equals, 100)

	// check for no overlap
	for _, run0 := range runs0 {
		for _, run1 := range runs1 {
			c.Assert(run0.ID, check.Not(check.Equals), run1.ID)
		}
	}

	runs, err = ms.model.RunList(0, 0, fork.Name, "")
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 10)

	runs, err = ms.model.RunList(0, 0, fork.Name, task.Submission.BaseRef.SHA)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 10)

	runs, err = ms.model.RunList(0, 0, parent.Name, "")
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 0)

	runs, err = ms.model.RunList(0, 0, "quux/foobar", "")
	c.Assert(err, check.NotNil)
	c.Assert(len(runs), check.Equals, 0)
}
