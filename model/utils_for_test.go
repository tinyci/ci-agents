package model

import (
	"fmt"
	"path"
	"time"

	"github.com/google/go-github/github"
	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

// CreateUsers creates `count` random users.
func (ms *modelSuite) CreateUsers(count int) ([]*User, *errors.Error) {
	users := []*User{}

	for i := 0; i < count; i++ {
		u, err := ms.model.CreateUser(testutil.RandString(8), testToken)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

// CreateRepository creates a random repository.
func (ms *modelSuite) CreateRepository() (*Repository, *errors.Error) {
	owners, err := ms.CreateUsers(1)
	if err != nil {
		return nil, err
	}
	name := path.Join(testutil.RandString(8), testutil.RandString(8))
	r := &Repository{
		Name:   name,
		Github: &gh.Repository{FullName: github.String(name)},
		Owner:  owners[0],
	}

	return r, errors.New(ms.model.Save(r).Error)
}

func (ms *modelSuite) CreateRun() (*Run, *errors.Error) {
	parent, err := ms.CreateRepository()
	if err != nil {
		return nil, err
	}

	fork, err := ms.CreateRepository()
	if err != nil {
		return nil, err
	}

	ref := &Ref{
		Repository: fork,
		RefName:    testutil.RandString(8),
		SHA:        testutil.RandHexString(40),
	}

	if err := ms.model.Save(ref).Error; err != nil {
		return nil, errors.New(err)
	}

	runName := testutil.RandString(8)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			runName: {
				Image:   "foo",
				Command: []string{"run", "me"},
				Queue:   "default",
			},
		},
	}

	task := &Task{
		TaskSettings: ts,
		Parent:       parent,
		Ref:          ref,
		BaseSHA:      testutil.RandHexString(40),
	}

	run := &Run{
		Name:        runName,
		RunSettings: ts.Runs[runName],
		Task:        task,
	}

	return run, errors.New(ms.model.Save(run).Error)
}

func (ms *modelSuite) FillQueue(count int64) ([]*QueueItem, *errors.Error) {
	fillstart := time.Now()
	qis := []*QueueItem{}

	for i := int64(1); i <= count; i++ {
		run, err := ms.CreateRun()
		if err != nil {
			return nil, err
		}

		qi := &QueueItem{
			Run:       run,
			QueueName: "default",
		}

		qis = append(qis, qi)
	}

	var err *errors.Error
	qis, err = ms.model.QueuePipelineAdd(qis)
	if err != nil {
		return nil, err
	}

	fmt.Println("Filling queue took", time.Since(fillstart))

	return qis, nil
}
