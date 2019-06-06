package testclients

import (
	"fmt"
	"io/ioutil"
	"strings"

	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/clients/queue"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// QueueClient is the queuesvc client
type QueueClient struct {
	client     *queue.Client
	dataClient *DataClient
}

// NewQueueClient returns a new queuesvc client with window dressings for tests.
func NewQueueClient(dc *DataClient) (*QueueClient, error) {
	ops, err := queue.New("localhost:6001", nil, false)
	return &QueueClient{client: ops, dataClient: dc}, err
}

// Client returns the underlying client.
func (qc *QueueClient) Client() *queue.Client {
	return qc.client
}

// SetUpSubmissionRepo takes a name of a repo; and configures the submission
// repo and a user belonging to it. Returns the name of the owner and any error.
func (qc *QueueClient) SetUpSubmissionRepo(name string) *errors.Error {
	parentUser, _, err := utils.OwnerRepo(name)
	if err != nil {
		return err
	}

	if _, err := qc.dataClient.MakeUser(parentUser); err != nil {
		return err
	}

	if err := qc.dataClient.MakeRepo(name, parentUser, false); err != nil {
		return err
	}

	if err := qc.dataClient.Client().EnableRepository(parentUser, name); err != nil {
		return err
	}

	return nil
}

// SetMockSubmissionSuccess creates all the mock tooling necessary to set up a submission
func (qc *QueueClient) SetMockSubmissionSuccess(mock *github.MockClientMockRecorder, sub *types.Submission) error {
	repoConfigBytes, err := ioutil.ReadFile("../../testdata/standard_repoconfig.yml")
	if err != nil {
		return err
	}

	taskBytes, err := ioutil.ReadFile("../../testdata/standard_task.yml")
	if err != nil {
		return err
	}

	mock.GetRepository(sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	mock.GetRepository(sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	mock.GetRepository(sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	mock.GetRefs(sub.Fork, sub.HeadSHA).Return([]string{"heads/fork-branch"}, nil)
	mock.GetRefs(sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	mock.GetDiffFiles(sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar"}, nil)
	mock.GetFileList(sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	mock.GetRepository(sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	mock.GetFile(sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)

	mock.GetFile(sub.Fork, sub.HeadSHA, "bar/task.yml").Return(taskBytes, nil)
	mock.GetFile(sub.Fork, sub.HeadSHA, "foo/task.yml").Return(taskBytes, nil)
	mock.GetFile(sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)

	parts := strings.SplitN(sub.Parent, "/", 2)

	for _, name := range []string{"*root*", "foo", "bar"} {
		for x := 1; x <= 5; x++ {
			mock.PendingStatus(parts[0], parts[1], fmt.Sprintf("%s:%d", name, x), sub.HeadSHA, "url")
		}
	}

	return nil
}
