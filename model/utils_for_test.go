package model

import (
	"math/rand"
	"path"

	"github.com/google/go-github/github"
	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

// CreateUsers creates `count` random users.
func (ms *modelSuite) CreateUsers(count int) ([]*User, error) {
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
func (ms *modelSuite) CreateRepository() (*Repository, error) {
	owners, err := ms.CreateUsers(rand.Intn(8) + 1)
	if err != nil {
		return nil, err
	}
	name := path.Join(testutil.RandString(8), testutil.RandString(8))
	r := &Repository{
		Name:   name,
		Github: &gh.Repository{FullName: github.String(name)},
		Owners: owners,
	}

	return r, ms.model.Save(r).Error
}

func (ms *modelSuite) CreateRun() (*Run, error) {
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
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	if err := ms.model.Save(ref).Error; err != nil {
		return nil, err
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
		BaseSHA:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	run := &Run{
		Name:        runName,
		RunSettings: ts.Runs[runName],
		Task:        task,
	}

	return run, ms.model.Save(run).Error
}
