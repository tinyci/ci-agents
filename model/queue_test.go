package model

import (
	"context"
	"fmt"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/types"
)

func (ms *modelSuite) TestQueueValidate(c *check.C) {
	for iter := 0; iter < 100; iter++ {
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

		ref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		}

		c.Assert(ms.model.Save(ref).Error, check.IsNil)

		sub := &Submission{
			BaseRef: baseref,
			HeadRef: ref,
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

		task := &Task{
			TaskSettings: ts,
			Submission:   sub,
		}

		c.Assert(ms.model.Save(task).Error, check.IsNil)
		run := &Run{
			Name:        "foobar",
			RunSettings: ts.Runs["foobar"],
			Task:        task,
		}

		c.Assert(ms.model.Save(run).Error, check.IsNil)

		failures := []struct {
			queueName string
			run       *Run
		}{
			{"", run},
			{"default", nil},
		}

		for i, failure := range failures {
			qi := &QueueItem{
				Run:       failure.run,
				QueueName: failure.queueName,
			}
			c.Assert(ms.model.Create(qi).Error, check.NotNil, check.Commentf("iteration %d", i))
			c.Assert(ms.model.Save(qi).Error, check.NotNil, check.Commentf("iteration %d", i))
		}

		qi := &QueueItem{
			QueueName: "default",
			Run:       run,
		}

		c.Assert(ms.model.Save(qi).Error, check.IsNil)

		qi2 := &QueueItem{}
		c.Assert(ms.model.Where("id = ?", qi.ID).First(qi2).Error, check.IsNil)
		c.Assert(qi.ID, check.Not(check.Equals), 0)
		c.Assert(qi.ID, check.Equals, qi2.ID)
		c.Assert(qi2.QueueName, check.Equals, "default")
		c.Assert(qi2.Run, check.NotNil)

		qis := []*QueueItem{}

		for i := 0; i < 10; i++ {
			qi.ID = 0
			qi.Run, err = ms.CreateRun()
			c.Assert(err, check.IsNil)
			c.Assert(ms.model.Save(qi).Error, check.IsNil)

			qi.ID = 0
			qi.Run, err = ms.CreateRun()
			c.Assert(err, check.IsNil)
			c.Assert(ms.model.Save(qi).Error, check.IsNil)

			qis = append(qis, qi)
		}

		i, err := ms.model.QueueTotalCount()
		c.Assert(err, check.IsNil)
		c.Assert(i, check.Equals, int64(21*(iter+1))) // relative to test iteration

		for _, qi := range qis {
			ro := "test"
			qi.RunningOn = &ro
			qi2, err := NewQueueItemFromProto(qi.ToProto())
			c.Assert(err, check.IsNil)
			c.Assert(qi2.Run.ID, check.Equals, qi.Run.ID)
			c.Assert(qi2.ID, check.Equals, qi.ID)
			c.Assert(qi2.QueueName, check.Equals, qi.QueueName)
			c.Assert(qi2.Running, check.Equals, qi.Running)
			c.Assert(*qi2.RunningOn, check.Equals, "test")
			c.Assert(qi2.StartedAt, check.Equals, qi.StartedAt)
		}
	}
}

func (ms *modelSuite) TestQueueManipulation(c *check.C) {
	var (
		firstID, lastRunID int64
	)

	oldModel := ms.model
	db := ms.model.Begin()
	ms.model = &Model{DB: db}

	fillstart := time.Now()
	for i := 1; i <= 1000; i++ {
		run, err := ms.CreateRun()
		c.Assert(err, check.IsNil)

		qi := &QueueItem{
			Run:       run,
			QueueName: "default",
		}

		c.Assert(ms.model.Save(qi).Error, check.IsNil)
		if firstID == 0 {
			firstID = qi.ID
		}
		lastRunID = run.ID
	}
	fmt.Println("Filling queue took", time.Since(fillstart))

	c.Assert(ms.model.Commit().Error, check.IsNil)
	ms.model = oldModel

	count, err := ms.model.QueueTotalCount()
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1000))

	_, err = ms.model.QueueList(-1, 100)
	c.Assert(err, check.NotNil)

	for i := 0; i < 10; i++ {
		list, err := ms.model.QueueList(int64(i), 100)
		c.Assert(err, check.IsNil)
		c.Assert(len(list), check.Equals, 100)
		for _, qi := range list {
			c.Assert(qi.ID, check.Not(check.Equals), 0)
			c.Assert(qi.Running, check.Equals, false)
			c.Assert(qi.RunningOn, check.IsNil)

			count, err := ms.model.QueueTotalCountForRepository(qi.Run.Task.Submission.BaseRef.Repository)
			c.Assert(err, check.IsNil)
			c.Assert(count, check.Equals, int64(1)) // repo names are uniq'd

			_, err = ms.model.QueueListForRepository(qi.Run.Task.Submission.BaseRef.Repository, -1, 100)
			c.Assert(err, check.NotNil)

			tmp, err := ms.model.QueueListForRepository(qi.Run.Task.Submission.BaseRef.Repository, 0, 100)
			c.Assert(err, check.IsNil)
			c.Assert(tmp[0], check.DeepEquals, qi) // repo names are uniq'd
		}
	}

	list, err := ms.model.QueueList(0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(list), check.Equals, 100)

	list2, err := ms.model.QueueList(1, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(list), check.Equals, 100)

	// validate no overlap.
	for _, qi := range list {
		for _, qi2 := range list2 {
			c.Assert(qi.ID, check.Not(check.Equals), qi2.ID)
		}
	}

	start := time.Now()

	for i := lastRunID - 999; i < lastRunID; i++ {
		qi, err := ms.model.NextQueueItem("hostname", "") // testing empty string handling
		c.Assert(err, check.IsNil)
		c.Assert(qi.ID, check.Equals, firstID, check.Commentf("%d", lastRunID-i))
		c.Assert(qi.Run.ID, check.Equals, int64(i)) // ensures same order
		c.Assert(qi.Run.Name, check.Not(check.Equals), "")
		c.Assert(qi.Running, check.Equals, true)
		c.Assert(*qi.RunningOn, check.Equals, "hostname")
		c.Assert(qi.StartedAt, check.NotNil)
		c.Assert(qi.Run.Task.Submission.BaseRef.Repository, check.NotNil) // checking the ORM works
		c.Assert(*qi.Run.RanOn, check.Equals, "hostname")
		firstID++
	}

	fmt.Println("Iterating queue took", time.Since(start))
}

func (ms *modelSuite) TestQueueNamed(c *check.C) {
	names := []string{
		"default", // testing default functionality
		"foo",
		"bar",
		"quux",
	}

	var firstID int64

	oldModel := ms.model
	db := ms.model.Begin()
	ms.model = &Model{DB: db}

	fillstart := time.Now()
	for i := 1; i <= 1000; i++ {
		for _, name := range names {
			run, err := ms.CreateRun()
			c.Assert(err, check.IsNil)

			qi := &QueueItem{
				Run:       run,
				QueueName: name,
			}

			c.Assert(ms.model.Save(qi).Error, check.IsNil)
			if firstID == 0 {
				firstID = run.ID
			}
		}
	}
	fmt.Println("Filling queue took", time.Since(fillstart))

	c.Assert(ms.model.Commit().Error, check.IsNil)
	ms.model = oldModel

	start := time.Now()

	for i := firstID; i < firstID+int64(1000*len(names)); i++ {
		qi, err := ms.model.NextQueueItem("hostname", names[(int64(i)-firstID)%int64(len(names))])
		c.Assert(err, check.IsNil)
		c.Assert(qi.Run.ID, check.Equals, i)
		c.Assert(qi.Run.Name, check.Not(check.Equals), "")
		c.Assert(qi.Running, check.Equals, true)
		c.Assert(*qi.RunningOn, check.Equals, "hostname")
		c.Assert(qi.StartedAt, check.NotNil)
		c.Assert(qi.Run.Task.Submission.BaseRef.Repository, check.NotNil) // checking the ORM works
		c.Assert(*qi.Run.RanOn, check.Equals, "hostname")
	}

	fmt.Println("Iterating queue took", time.Since(start))
}

func (ms *modelSuite) TestQueueConcurrent(c *check.C) {
	count := int64(1000)
	goRoutines := 10

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queueChan := make(chan *QueueItem, goRoutines)
	errChan := make(chan error, goRoutines)

	var firstID int64

	fillstart := time.Now()
	qis := []*QueueItem{}
	for i := int64(1); i <= count; i++ {
		run, err := ms.CreateRun()
		c.Assert(err, check.IsNil)

		qi := &QueueItem{
			Run:       run,
			QueueName: "default",
		}

		qis = append(qis, qi)

		if firstID == 0 {
			firstID = run.ID
		}
	}

	_, err := ms.model.QueuePipelineAdd(qis)
	c.Assert(err, check.IsNil)

	fmt.Println("Filling queue took", time.Since(fillstart))

	start := time.Now()

	for i := 0; i < goRoutines; i++ {
		go func(i int) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				qi, err := ms.model.NextQueueItem("hostname", "default")
				if err != nil {
					errChan <- err
					return
				}

				queueChan <- qi
			}
		}(i)
	}

	for i := firstID; i <= firstID+count-1; i++ {
		select {
		case err := <-errChan:
			c.Assert(err, check.IsNil)
		case qi := <-queueChan:
			c.Assert(qi.Run.ID, check.Equals, i, check.Commentf("%d", i-firstID)) // ensures same order
			c.Assert(qi.Run.Name, check.Not(check.Equals), "")
			c.Assert(qi.Running, check.Equals, true)
			c.Assert(*qi.RunningOn, check.Equals, "hostname")
			c.Assert(qi.StartedAt, check.NotNil)
			c.Assert(qi.Run.Task.Submission.BaseRef.Repository, check.NotNil) // checking the ORM works
			c.Assert(*qi.Run.RanOn, check.Equals, "hostname")
		}
	}

	fmt.Println("Iterating queue took", time.Since(start))
}

func (ms *modelSuite) TestQueueNamedConcurrent(c *check.C) {
	names := []string{
		"default",
		"foo",
		"bar",
		"quux",
	}

	multiplier := 2
	count := 1000

	queueChan := make(chan *QueueItem, len(names)*multiplier)
	errChan := make(chan error, len(names)*multiplier)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var y int

	fillstart := time.Now()
	qis := []*QueueItem{}

	for i := 1; i <= count; i++ {
		for _, name := range names {
			y++
			run, err := ms.CreateRun()
			c.Assert(err, check.IsNil)

			qi := &QueueItem{
				Run:       run,
				QueueName: name,
			}

			qis = append(qis, qi)
		}
	}

	var err error
	_, err = ms.model.QueuePipelineAdd(qis)
	c.Assert(err, check.IsNil)

	fmt.Println("Filling queue took", time.Since(fillstart))

	start := time.Now()

	for _, name := range names {
		for i := 0; i < multiplier; i++ {
			go func(name string) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}

					qi, err := ms.model.NextQueueItem("hostname", name)
					if err != nil {
						errChan <- err
						return
					}

					queueChan <- qi
				}
			}(name)
		}
	}

	for i := 1; i <= count; i++ {
		select {
		case err := <-errChan:
			c.Assert(err, check.IsNil)
		case qi := <-queueChan:
			c.Assert(qi.Run.ID, check.Not(check.Equals), 0) // order isn't as deterministic here
			c.Assert(qi.Run.Name, check.Not(check.Equals), "")
			c.Assert(qi.Running, check.Equals, true)
			c.Assert(*qi.RunningOn, check.Equals, "hostname")
			c.Assert(qi.StartedAt, check.NotNil)
			c.Assert(qi.Run.Task.Submission.BaseRef.Repository, check.NotNil) // checking the ORM works
			c.Assert(*qi.Run.RanOn, check.Equals, "hostname")
		}
	}

	fmt.Println("Iterating queue took", time.Since(start))
}
