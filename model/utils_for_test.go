package model

import (
	"fmt"
	"path"
	"time"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

func (ms *modelSuite) CreateRef(r *Repository, name, sha string) (*Ref, *errors.Error) {
	ref := &Ref{
		Repository: r,
		SHA:        sha,
		RefName:    name,
	}

	return ref, ms.model.WrapError(ms.model.Create(ref), "creating ref for test")
}

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
	return ms.CreateRepositoryWithName(path.Join(testutil.RandString(8), testutil.RandString(8)))
}

func (ms *modelSuite) CreateRepositoryWithName(name string) (*Repository, *errors.Error) {
	owners, err := ms.CreateUsers(1)
	if err != nil {
		return nil, errors.New(err)
	}
	r := &Repository{
		Name:   name,
		Github: &github.Repository{FullName: github.String(name)},
		Owner:  owners[0],
	}

	if err := ms.model.Save(r).Error; err != nil {
		return nil, errors.New(err)
	}

	return r, nil
}

func (ms *modelSuite) CreateTaskForSubmission(sub *Submission) (*Task, error) {
	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"default": {
				Image:   "foo",
				Command: []string{"run", "me"},
				Queue:   "default",
			},
		},
	}

	task := &Task{
		TaskSettings: ts,
		Submission:   sub,
	}

	run := &Run{
		Name:        "default",
		RunSettings: ts.Runs["default"],
		Task:        task,
	}

	if err := ms.model.Save(run).Error; err != nil {
		return nil, errors.New(err)
	}

	qi := &QueueItem{
		Run:       run,
		QueueName: "default",
	}

	if err := ms.model.Save(qi).Error; err != nil {
		return nil, errors.New(err)
	}

	return run.Task, nil
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
	baseref := &Ref{
		Repository: parent,
		RefName:    testutil.RandString(8),
		SHA:        testutil.RandHexString(40),
	}

	if err := ms.model.Save(baseref).Error; err != nil {
		return nil, errors.New(err)
	}

	headref := &Ref{
		Repository: fork,
		RefName:    testutil.RandString(8),
		SHA:        testutil.RandHexString(40),
	}

	if err := ms.model.Save(headref).Error; err != nil {
		return nil, errors.New(err)
	}

	sub := &Submission{HeadRef: headref, BaseRef: baseref}

	if err := ms.model.Save(sub).Error; err != nil {
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
		Submission:   sub,
	}

	run := &Run{
		Name:        runName,
		RunSettings: ts.Runs[runName],
		Task:        task,
	}

	return run, errors.New(ms.model.Save(run).Error)
}

func (ms *modelSuite) FillQueue(count int64) ([]*QueueItem, error) {
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

	var err error
	qis, err = ms.model.QueuePipelineAdd(qis)
	if err != nil {
		return nil, err
	}

	fmt.Println("Filling queue took", time.Since(fillstart))

	return qis, nil
}

func (ms *modelSuite) CreateSubmission(sub *types.Submission) (*Submission, *errors.Error) {
	if sub.SubmittedBy != "" {
		if _, err := ms.model.CreateUser(sub.SubmittedBy, testutil.DummyToken); err != nil {
			return nil, err
		}
	}

	if sub.Fork != "" {
		r, err := ms.CreateRepositoryWithName(sub.Fork)
		if err != nil {
			r, err = ms.model.GetRepositoryByName(sub.Fork)
			if err != nil {
				return nil, err
			}
		}
		if _, err := ms.CreateRef(r, "dummy", sub.HeadSHA); err != nil {
			return nil, err
		}
	}

	r, err := ms.CreateRepositoryWithName(sub.Parent)
	if err != nil {
		r, err = ms.model.GetRepositoryByName(sub.Parent)
		if err != nil {
			return nil, err
		}
	}

	if _, err := ms.CreateRef(r, "dummy", sub.BaseSHA); err != nil {
		return nil, err
	}

	s, err := ms.model.NewSubmissionFromMessage(sub)
	if err != nil {
		return nil, errors.New(err)
	}

	if err := ms.model.Save(s).Error; err != nil {
		return nil, errors.New(err)
	}

	return s, nil
}
