package tinyci

import (
	"context"
	"io"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/gen/client/uisvc/client/operations"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
)

// Client is the official wrapper around the swagger uisvc client.
type Client struct {
	client *operations.Client
}

// New constructs a new *Client
func New(url, token string, cert *transport.Cert) (*Client, error) {
	c, err := operations.New(url, token, cert)
	if err != nil {
		return nil, err
	}

	return &Client{client: c}, nil
}

// DeleteToken removes your token. You won't be able to request anything after making this call.
func (c *Client) DeleteToken(ctx context.Context) error {
	return c.client.DeleteToken(ctx)
}

// Errors gets the user errors logged into the system.
func (c *Client) Errors(ctx context.Context) ([]*model.UserError, error) {
	errs, err := c.client.GetErrors(ctx)
	if err != nil {
		return nil, err
	}

	ue := []*model.UserError{}

	return ue, utils.JSONIO(errs, &ue)
}

// Submit submits a request to test a repository to tinyCI.
func (c *Client) Submit(ctx context.Context, repository, sha string, all bool) error {
	return c.client.GetSubmit(ctx, all, repository, sha)
}

// LogAttach attaches to a and retrieves it's output. Attach will block the
// stream assuming that that job is not completed.
func (c *Client) LogAttach(ctx context.Context, id int64, w io.WriteCloser) error {
	return c.client.GetLogAttachID(ctx, id, w)
}

// LoadRepositories loads your repos from github and returns the objects tinyci recorded.
func (c *Client) LoadRepositories(ctx context.Context, search string) ([]*model.Repository, error) {
	if err := c.client.GetRepositoriesScan(ctx); err != nil {
		return nil, err
	}

	repos, err := c.client.GetRepositoriesMy(ctx, search)
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
func (c *Client) AddToCI(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}

	return c.client.GetRepositoriesCiAddOwnerRepo(ctx, owner, reponame)
}

// DeleteFromCI deletes a repository from CI.
func (c *Client) DeleteFromCI(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	return c.client.GetRepositoriesCiDelOwnerRepo(ctx, owner, reponame)
}

// Subscribed lists all subscribed repositories.
func (c *Client) Subscribed(ctx context.Context, search string) ([]*model.Repository, error) {
	repos, err := c.client.GetRepositoriesSubscribed(ctx, search)
	if err != nil {
		return nil, err
	}

	ret := []*model.Repository{}
	return ret, utils.JSONIO(repos, &ret)
}

// Visible lists all visible repositories.
func (c *Client) Visible(ctx context.Context, search string) ([]*model.Repository, error) {
	repos, err := c.client.GetRepositoriesVisible(ctx, search)
	if err != nil {
		return nil, err
	}

	ret := []*model.Repository{}
	return ret, utils.JSONIO(repos, &ret)
}

// Subscribe to a repository.
func (c *Client) Subscribe(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	return c.client.GetRepositoriesSubAddOwnerRepo(ctx, owner, reponame)
}

// Unsubscribe from a repository.
func (c *Client) Unsubscribe(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	return c.client.GetRepositoriesSubDelOwnerRepo(ctx, owner, reponame)
}

// Tasks returns the tasks with pagination and optional filtering. (Just pass empty values for no filters)
func (c *Client) Tasks(ctx context.Context, repository, sha string, page, perPage int64) ([]*model.Task, error) {
	tasks, err := c.client.GetTasks(ctx, page, perPage, repository, sha)
	if err != nil {
		return nil, err
	}

	ret := []*model.Task{}

	return ret, utils.JSONIO(tasks, &ret)
}

// TaskCount returns the total number of tasks matching the filter.
func (c *Client) TaskCount(ctx context.Context, repository, sha string) (int64, error) {
	count, err := c.client.GetTasksCount(ctx, repository, sha)
	return count, err
}

// RunsForTask returns the runs for the provided task id.
func (c *Client) RunsForTask(ctx context.Context, taskID, page, perPage int64) ([]*model.Run, error) {
	runs, err := c.client.GetTasksRunsID(ctx, taskID, page, perPage)
	if err != nil {
		return nil, err
	}

	ret := []*model.Run{}
	return ret, utils.JSONIO(runs, &ret)
}

// RunsForTaskCount returns the count of the runs for the provided task id.
func (c *Client) RunsForTaskCount(ctx context.Context, taskID int64) (int64, error) {
	count, err := c.client.GetTasksRunsIDCount(ctx, taskID)
	return count, err
}

// Runs returns all the runs matching the filter set with pagination.
func (c *Client) Runs(ctx context.Context, repository, sha string, page, perPage int64) ([]*model.Run, error) {
	runs, err := c.client.GetRuns(ctx, page, perPage, repository, sha)
	if err != nil {
		return nil, err
	}

	ret := []*model.Run{}
	return ret, utils.JSONIO(runs, &ret)
}

// RunsCount returns the count of the runs.
func (c *Client) RunsCount(ctx context.Context, repository, sha string) (int64, error) {
	count, err := c.client.GetRunsCount(ctx, repository, sha)

	return count, err
}

// GetRun retrieves a run by id.
func (c *Client) GetRun(ctx context.Context, id int64) (*model.Run, error) {
	run, err := c.client.GetRunRunID(ctx, id)
	if err != nil {
		return nil, err
	}

	ret := &model.Run{}
	return ret, utils.JSONIO(run, ret)
}

// CancelRun cancels the run by id. It may also cancel other runs based on rules
// around queue management.
func (c *Client) CancelRun(ctx context.Context, id int64) error {
	return c.client.PostCancelRunID(ctx, id)
}

// AddCapability adds a capability for a user. Must have the modify:user capability to interact.
func (c *Client) AddCapability(ctx context.Context, username string, capability model.Capability) error {
	return c.client.PostCapabilitiesUsernameCapability(ctx, string(capability), username)
}

// RemoveCapability removes a capability from a user. Must have the modify:user capability to interact.
func (c *Client) RemoveCapability(ctx context.Context, username string, capability model.Capability) error {
	return c.client.DeleteCapabilitiesUsernameCapability(ctx, string(capability), username)
}

// GetUserProperties returns some properties about the requesting account; like username and capabilities.
func (c *Client) GetUserProperties(ctx context.Context) (map[string]interface{}, error) {
	res, err := c.client.GetUserProperties(ctx)
	if err != nil {
		return nil, err
	}

	return res.(map[string]interface{}), nil
}

// VisibleRepos retrieves the visible repositories to the user, a search may also be provided to limit scope.
func (c *Client) VisibleRepos(ctx context.Context, search string) ([]*model.Repository, error) {
	repos, err := c.client.GetRepositoriesVisible(ctx, search)
	if err != nil {
		return nil, err
	}

	newRepos := model.RepositoryList{}
	return newRepos, utils.JSONIO(repos, newRepos)
}

// Submissions returns a list of submissions, paginated and optionally filtered by repository and SHA.
func (c *Client) Submissions(ctx context.Context, repository, sha string, page, perPage int64) ([]*model.Submission, error) {
	subs, err := c.client.GetSubmissions(ctx, page, perPage, repository, sha)
	if err != nil {
		return nil, err
	}

	newSubs := []*model.Submission{}
	return newSubs, utils.JSONIO(subs, &newSubs)
}

// TasksForSubmission returns the tasks for the given submission.
func (c *Client) TasksForSubmission(ctx context.Context, sub *model.Submission) ([]*model.Task, error) {
	perPage := int64(20)
	page := int64(0)

	totalTasks := []*model.Task{}

	for {
		tasks, err := c.client.GetSubmissionIDTasks(ctx, sub.ID, page, perPage)
		if err != nil {
			return nil, err
		}

		jsonTasks := []*model.Task{}

		if err := utils.JSONIO(tasks, &jsonTasks); err != nil {
			return nil, err
		}

		totalTasks = append(totalTasks, jsonTasks...)

		if len(tasks) == 0 {
			break
		}

		page++
	}

	return totalTasks, nil
}
