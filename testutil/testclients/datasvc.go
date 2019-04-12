package testclients

import (
	"path"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// DataClient is the datasvc client
type DataClient struct {
	client *data.Client
}

// NewDataClient returns a new datasvc client with window dressings for tests.
func NewDataClient() (*DataClient, *errors.Error) {
	ops, err := data.New("localhost:6000", nil)
	return &DataClient{client: ops}, err
}

// Client returns the underlying client.
func (dc *DataClient) Client() *data.Client {
	return dc.client
}

// MakeUser makes a new user with the name provided. It is given a dummy access token.
func (dc *DataClient) MakeUser(username string) (*model.User, *errors.Error) {
	return dc.client.PutUser(&model.User{
		Username: username,
		Token:    testutil.DummyToken,
	})
}

// MakeRepo saves a repo with name, owner, and private state.
func (dc *DataClient) MakeRepo(fullRepo, owner string, private bool) *errors.Error {
	repos := []interface{}{
		map[string]interface{}{"full_name": fullRepo, "private": private},
	}

	ghRepos := []*github.Repository{}

	if err := utils.JSONIO(repos, &ghRepos); err != nil {
		return err
	}

	return dc.client.PutRepositories(owner, ghRepos, false)
}

// MakeQueueItem returns a queueitem that has already been stored
func (dc *DataClient) MakeQueueItem() (*model.QueueItem, *errors.Error) {
	username := testutil.RandString(8)
	_, err := dc.MakeUser(username)
	if err != nil {
		return nil, err
	}

	parentRepoOwner, parentRepoName := testutil.RandString(8), testutil.RandString(8)
	repoName := path.Join(parentRepoOwner, parentRepoName)
	if err := dc.MakeRepo(repoName, username, false); err != nil {
		return nil, err
	}

	parent, err := dc.client.GetRepository(repoName)
	if err != nil {
		return nil, err
	}

	forkRepoOwner, forkRepoName := testutil.RandString(8), testutil.RandString(8)
	forkName := path.Join(forkRepoOwner, forkRepoName)
	if err := dc.MakeRepo(forkName, username, false); err != nil {
		return nil, err
	}

	fork, err := dc.client.GetRepository(forkName)
	if err != nil {
		return nil, err
	}

	ref := &model.Ref{
		Repository: fork,
		RefName:    testutil.RandString(8),
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	id, err := dc.client.PutRef(ref)
	if err != nil {
		return nil, err
	}

	ref.ID = id

	runName := testutil.RandString(8)

	ts := &types.TaskSettings{
		WorkDir:    "/tmp",
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			runName: {
				Image:   "foo",
				Command: []string{"run", "me"},
				Queue:   "default",
			},
		},
	}

	task := &model.Task{
		TaskSettings: ts,
		Parent:       parent,
		Ref:          ref,
		BaseSHA:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	t, err := dc.client.PutTask(task)
	if err != nil {
		return nil, err
	}

	task.ID = t.ID

	qi := &model.QueueItem{
		QueueName: "default",
		Run: &model.Run{
			Name:        runName,
			RunSettings: ts.Runs[runName],
			Task:        t,
		},
	}

	qis, err := dc.client.PutQueue([]*model.QueueItem{qi})
	if err != nil {
		return nil, err
	}

	return qis[0], nil
}
