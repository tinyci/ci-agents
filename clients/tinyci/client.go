package tinyci

import (
	"context"
	"io"

	"github.com/tinyci/ci-agents/ci-gen/gen/client/uisvc/client/operations"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
)

// Client is the official wrapper around the swagger uisvc client.
type Client struct {
	client *operations.Client
}

// New constructs a new *Client
func New(url, token string) (*Client, *errors.Error) {
	c, err := operations.New(url, token)
	if err != nil {
		return nil, err
	}

	return &Client{client: c}, nil
}

// DeleteToken removes your token. You won't be able to request anything after making this call.
func (c *Client) DeleteToken() *errors.Error {
	return c.client.DeleteToken(context.Background())
}

// Errors gets the user errors logged into the system.
func (c *Client) Errors() ([]*model.UserError, *errors.Error) {
	errs, err := c.client.GetErrors(context.Background())
	if err != nil {
		return nil, err
	}

	ue := []*model.UserError{}

	return ue, utils.JSONIO(errs, &ue)
}

// Submit submits a request to test a repository to tinyCI.
func (c *Client) Submit(repository, sha string, all bool) *errors.Error {
	return c.client.GetSubmit(context.Background(), all, repository, sha)
}

// LogAttach attaches to a and retrieves it's output. Attach will block the
// stream assuming that that job is not completed.
func (c *Client) LogAttach(id int64, w io.WriteCloser) *errors.Error {
	return c.client.GetLogAttachID(context.Background(), id, w)
}

// LoadRepositories loads your repos from github and returns the objects tinyci recorded.
func (c *Client) LoadRepositories(search string) ([]*model.Repository, *errors.Error) {
	if err := c.client.GetRepositoriesScan(context.Background()); err != nil {
		return nil, err
	}

	repos, err := c.client.GetRepositoriesMy(context.Background(), search)
	if err != nil {
		return nil, err
	}

	ret := []*model.Repository{}

	if err := utils.JSONIO(repos, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// AddToCI adds a repository to CI.
func (c *Client) AddToCI(repository string) *errors.Error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}

	return c.client.GetRepositoriesCiAddOwnerRepo(context.Background(), owner, reponame)
}

// DeleteFromCI deletes a repository from CI.
func (c *Client) DeleteFromCI(repository string) *errors.Error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	return c.client.GetRepositoriesCiDelOwnerRepo(context.Background(), owner, reponame)
}

// Subscribed lists all subscribed repositories.
func (c *Client) Subscribed(search string) ([]*model.Repository, *errors.Error) {
	repos, err := c.client.GetRepositoriesSubscribed(context.Background(), search)
	if err != nil {
		return nil, err
	}

	ret := []*model.Repository{}
	return ret, utils.JSONIO(repos, &ret)
}

// Visible lists all visible repositories.
func (c *Client) Visible(search string) ([]*model.Repository, *errors.Error) {
	repos, err := c.client.GetRepositoriesVisible(context.Background(), search)
	if err != nil {
		return nil, err
	}

	ret := []*model.Repository{}
	return ret, utils.JSONIO(repos, &ret)
}

// Subscribe to a repository.
func (c *Client) Subscribe(repository string) *errors.Error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	return c.client.GetRepositoriesSubAddOwnerRepo(context.Background(), owner, reponame)
}

// Unsubscribe from a repository.
func (c *Client) Unsubscribe(repository string) *errors.Error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	return c.client.GetRepositoriesSubDelOwnerRepo(context.Background(), owner, reponame)
}

// Tasks returns the tasks with pagination and optional filtering. (Just pass empty values for no filters)
func (c *Client) Tasks(repository, sha string, page, perPage int64) ([]*model.Task, *errors.Error) {
	tasks, err := c.client.GetTasks(context.Background(), page, perPage, repository, sha)
	if err != nil {
		return nil, err
	}

	ret := []*model.Task{}

	return ret, utils.JSONIO(tasks, &ret)
}

// TaskCount returns the total number of tasks matching the filter.
func (c *Client) TaskCount(repository, sha string) (int64, *errors.Error) {
	return c.client.GetTasksCount(context.Background(), repository, sha)
}

// RunsForTask returns the runs for the provided task id.
func (c *Client) RunsForTask(taskID, page, perPage int64) ([]*model.Run, *errors.Error) {
	runs, err := c.client.GetTasksRunsID(context.Background(), taskID, page, perPage)
	if err != nil {
		return nil, err
	}

	ret := []*model.Run{}
	return ret, utils.JSONIO(runs, &ret)
}

// RunsForTaskCount returns the count of the runs for the provided task id.
func (c *Client) RunsForTaskCount(taskID int64) (int64, *errors.Error) {
	return c.client.GetTasksRunsIDCount(context.Background(), taskID)
}

// Runs returns all the runs matching the filter set with pagination.
func (c *Client) Runs(repository, sha string, page, perPage int64) ([]*model.Run, *errors.Error) {
	runs, err := c.client.GetRuns(context.Background(), page, perPage, repository, sha)
	if err != nil {
		return nil, err
	}

	ret := []*model.Run{}
	return ret, utils.JSONIO(runs, &ret)
}

// RunsCount returns the count of the runs.
func (c *Client) RunsCount(repository, sha string) (int64, *errors.Error) {
	return c.client.GetRunsCount(context.Background(), repository, sha)
}

// GetRun retrieves a run by id.
func (c *Client) GetRun(id int64) (*model.Run, *errors.Error) {
	run, err := c.client.GetRunRunID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	ret := &model.Run{}
	return ret, utils.JSONIO(run, ret)
}

// CancelRun cancels the run by id. It may also cancel other runs based on rules
// around queue management.
func (c *Client) CancelRun(id int64) *errors.Error {
	return c.client.PostCancelRunID(context.Background(), id)
}

// AddCapability adds a capability for a user. Must have the modify:user capability to interact.
func (c *Client) AddCapability(username string, capability model.Capability) *errors.Error {
	return c.client.PostCapabilitiesUsernameCapability(context.Background(), string(capability), username)
}

// RemoveCapability removes a capability from a user. Must have the modify:user capability to interact.
func (c *Client) RemoveCapability(username string, capability model.Capability) *errors.Error {
	return c.client.DeleteCapabilitiesUsernameCapability(context.Background(), string(capability), username)
}

// GetUserProperties returns some properties about the requesting account; like username and capabilities.
func (c *Client) GetUserProperties() (map[string]interface{}, *errors.Error) {
	res, err := c.client.GetUserProperties(context.Background())
	if err != nil {
		return nil, err
	}

	return res.(map[string]interface{}), nil
}

// VisibleRepos retrieves the visible repositories to the user, a search may also be provided to limit scope.
func (c *Client) VisibleRepos(search string) ([]*model.Repository, *errors.Error) {
	repos, err := c.client.GetRepositoriesVisible(context.Background(), search)
	if err != nil {
		return nil, err
	}

	newRepos := model.RepositoryList{}
	return newRepos, utils.JSONIO(repos, newRepos)
}
