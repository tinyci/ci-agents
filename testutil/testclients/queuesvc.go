package testclients

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/golang/mock/gomock"
	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/clients/queue"
	"github.com/tinyci/ci-agents/config"
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
	ops, err := queue.New(config.DefaultServices.Queue.String(), nil, false)
	return &QueueClient{client: ops, dataClient: dc}, err
}

// Client returns the underlying client.
func (qc *QueueClient) Client() *queue.Client {
	return qc.client
}

// SetUpSubmissionRepo takes a name of a repo; and configures the submission
// repo and a user belonging to it. Returns the name of the owner and any error.
func (qc *QueueClient) SetUpSubmissionRepo(name string, forkOf string) error {
	parentUser, _, err := utils.OwnerRepo(name)
	if err != nil {
		return err
	}

	if _, err := qc.dataClient.MakeUser(parentUser); err != nil {
		return err
	}

	if err := qc.dataClient.MakeRepo(name, parentUser, false, forkOf); err != nil {
		return err
	}

	if err := qc.dataClient.Client().EnableRepository(context.Background(), parentUser, name); err != nil {
		return err
	}

	return nil
}

// SetMockSubmissionOnFork sets mock submissions for fork-only repositories. Used in a few tests.
func (qc *QueueClient) SetMockSubmissionOnFork(mock *github.MockClientMockRecorder, sub *types.Submission, parent, resolvedSHA, pathadd string) error {
	repoConfigBytes, err := ioutil.ReadFile(pathadd + "../../../testdata/standard_repoconfig.yml")
	if err != nil {
		return err
	}

	taskBytes, err := ioutil.ReadFile(pathadd + "../../../testdata/standard_task.yml")
	if err != nil {
		return err
	}

	mock.GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent), Fork: gh.Bool(true)}, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(parent)}}, nil)
	mock.GetSHA(gomock.Any(), sub.Parent, "heads/master").Return(sub.BaseSHA, nil)
	mock.GetSHA(gomock.Any(), sub.Fork, "heads/master").Return(resolvedSHA, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(parent)}}, nil)
	mock.GetSHA(gomock.Any(), sub.Fork, "heads/master").Return(resolvedSHA, nil)
	mock.GetSHA(gomock.Any(), sub.Fork, "heads/master").Return(resolvedSHA, nil)
	mock.ClearStates(gomock.Any(), sub.Fork, resolvedSHA).Return(nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(parent)}}, nil)
	mock.GetRefs(gomock.Any(), sub.Fork, resolvedSHA).Return([]string{"heads/master"}, nil)
	mock.GetRefs(gomock.Any(), sub.Fork, resolvedSHA).Return([]string{"heads/master"}, nil)
	mock.GetDiffFiles(gomock.Any(), sub.Fork, resolvedSHA, resolvedSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar"}, nil)
	mock.GetFileList(gomock.Any(), sub.Fork, resolvedSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(parent)}}, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(parent)}}, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(parent)}}, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(parent)}}, nil)
	mock.GetFile(gomock.Any(), sub.Fork, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)

	mock.GetFile(gomock.Any(), sub.Fork, resolvedSHA, "bar/task.yml").Return(taskBytes, nil)
	mock.GetFile(gomock.Any(), sub.Fork, resolvedSHA, "foo/task.yml").Return(taskBytes, nil)
	mock.GetFile(gomock.Any(), sub.Fork, resolvedSHA, "task.yml").Return(taskBytes, nil)

	parts := strings.SplitN(sub.Fork, "/", 2)

	for _, name := range []string{"*root*", "foo", "bar"} {
		for x := 1; x <= 5; x++ {
			mock.PendingStatus(gomock.Any(), parts[0], parts[1], fmt.Sprintf("%s:%d", name, x), resolvedSHA, "url")
		}
	}

	return nil
}

// GetYAMLs finds the test yamls needed to run the tests.
func (qc *QueueClient) GetYAMLs(pathadd string) ([]byte, []byte, error) {
	repoConfigBytes, err := ioutil.ReadFile(pathadd + "../../../testdata/standard_repoconfig.yml")
	if err != nil {
		return nil, nil, err
	}

	taskBytes, err := ioutil.ReadFile(pathadd + "../../../testdata/standard_task.yml")
	if err != nil {
		return nil, nil, err
	}

	return repoConfigBytes, taskBytes, nil
}

// SetMockSubmissionSuccess creates all the mock tooling necessary to set up a submission
func (qc *QueueClient) SetMockSubmissionSuccess(mock *github.MockClientMockRecorder, sub *types.Submission, forkBranch string, pathadd string) error {
	repoConfigBytes, taskBytes, err := qc.GetYAMLs(pathadd)
	if err != nil {
		return err
	}

	mock.GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	mock.GetSHA(gomock.Any(), sub.Fork, forkBranch).Return(sub.HeadSHA, nil) // also here
	mock.GetSHA(gomock.Any(), sub.Parent, "heads/master").Return(sub.BaseSHA, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	mock.GetRefs(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{forkBranch}, nil)
	mock.GetRefs(gomock.Any(), sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	mock.GetRefs(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{forkBranch}, nil)
	mock.GetDiffFiles(gomock.Any(), sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar"}, nil)
	mock.GetFileList(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	mock.GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	mock.GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	mock.GetFile(gomock.Any(), sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)

	mock.GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "bar/task.yml").Return(taskBytes, nil)
	mock.GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "foo/task.yml").Return(taskBytes, nil)
	mock.GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)

	parts := strings.SplitN(sub.Parent, "/", 2)

	for _, name := range []string{"*root*", "foo", "bar"} {
		for x := 1; x <= 5; x++ {
			mock.PendingStatus(gomock.Any(), parts[0], parts[1], fmt.Sprintf("%s:%d", name, x), sub.HeadSHA, "url")
		}
	}

	return nil
}
