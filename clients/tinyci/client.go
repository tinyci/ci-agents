package tinyci

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/clients/jsonbuffer"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"golang.org/x/net/websocket"
)

type roundTripper struct {
	token string
	under http.RoundTripper
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", rt.token)
	return rt.under.RoundTrip(req)
}

// Client is the official wrapper around the swagger uisvc client.
type Client struct {
	client  *uisvc.Client
	baseURL *url.URL
	tls     bool
}

type httpClientWrapper struct{ *http.Client }

func (hcw *httpClientWrapper) Do(req *http.Request) (*http.Response, error) {
	resp, err := hcw.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 500 {
		var uerr uisvc.Error
		buf := bytes.NewBuffer(nil)

		if err := json.NewDecoder(io.TeeReader(resp.Body, buf)).Decode(&uerr); err != nil {
			return nil, errors.New(buf.String())
		}

		var err error

		if uerr.Errors != nil {
			err = errors.New(strings.Join(*uerr.Errors, "; "))
		}

		return nil, err
	}

	return resp, nil
}

// New constructs a new *Client
func New(u, token string, cert *transport.Cert) (*Client, error) {
	baseURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	gt, err := transport.NewHTTP(cert)
	if err != nil {
		return nil, err
	}

	client := gt.Client(nil)
	client.Transport = &roundTripper{token: token, under: client.Transport}

	c, err := uisvc.NewClient(baseURL.String(), uisvc.WithHTTPClient(&httpClientWrapper{client}))
	if err != nil {
		return nil, err
	}

	return &Client{tls: cert != nil, baseURL: baseURL, client: c}, nil
}

// DeleteToken removes your token. You won't be able to request anything after making this call.
func (c *Client) DeleteToken(ctx context.Context) error {
	resp, err := c.client.DeleteToken(ctx)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// Errors gets the user errors logged into the system.
func (c *Client) Errors(ctx context.Context) ([]*model.UserError, error) {
	resp, err := c.client.GetErrors(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ue := []*model.UserError{}

	return ue, json.NewDecoder(io.TeeReader(resp.Body, os.Stdout)).Decode(&ue)
}

// Submit submits a request to test a repository to tinyCI.
func (c *Client) Submit(ctx context.Context, repository, sha string, all bool) error {
	resp, err := c.client.GetSubmit(ctx, &uisvc.GetSubmitParams{All: &all, Repository: repository, Sha: sha})
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

// LogAttach attaches to a and retrieves it's output. Attach will block the
// stream assuming that that job is not completed.
func (c *Client) LogAttach(ctx context.Context, id int64, w io.WriteCloser) error {
	baseURL := *c.baseURL

	baseURL.Scheme = "ws"
	if c.tls {
		baseURL.Scheme = "wss"
	}

	baseURL.Path += fmt.Sprintf("/log/attach/%d", id)
	conn, err := websocket.Dial(baseURL.String(), "", baseURL.String())
	if err != nil {
		return err
	}

	_, err = io.Copy(w, jsonbuffer.NewReadWrapper(conn))
	return err
}

// LoadRepositories loads your repos from github and returns the objects tinyci recorded.
func (c *Client) LoadRepositories(ctx context.Context, search *string) ([]*model.Repository, error) {
	resp, err := c.client.GetRepositoriesScan(ctx)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	resp, err = c.client.GetRepositoriesMy(ctx, &uisvc.GetRepositoriesMyParams{Search: search})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []*model.Repository{}

	return ret, json.NewDecoder(resp.Body).Decode(&ret)
}

// AddToCI adds a repository to CI.
func (c *Client) AddToCI(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}

	resp, err := c.client.GetRepositoriesCiAddOwnerRepo(ctx, owner, reponame)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

// DeleteFromCI deletes a repository from CI.
func (c *Client) DeleteFromCI(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	resp, err := c.client.GetRepositoriesCiDelOwnerRepo(ctx, owner, reponame)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

// Subscribed lists all subscribed repositories.
func (c *Client) Subscribed(ctx context.Context, search *string) ([]*model.Repository, error) {
	resp, err := c.client.GetRepositoriesSubscribed(ctx, &uisvc.GetRepositoriesSubscribedParams{Search: search})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []*model.Repository{}
	return ret, json.NewDecoder(resp.Body).Decode(&ret)
}

// Visible lists all visible repositories.
func (c *Client) Visible(ctx context.Context, search *string) ([]*model.Repository, error) {
	resp, err := c.client.GetRepositoriesVisible(ctx, &uisvc.GetRepositoriesVisibleParams{Search: search})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []*model.Repository{}
	return ret, json.NewDecoder(resp.Body).Decode(&ret)
}

// Subscribe to a repository.
func (c *Client) Subscribe(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	resp, err := c.client.GetRepositoriesSubAddOwnerRepo(ctx, owner, reponame)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

// Unsubscribe from a repository.
func (c *Client) Unsubscribe(ctx context.Context, repository string) error {
	owner, reponame, err := utils.OwnerRepo(repository)
	if err != nil {
		return err
	}
	resp, err := c.client.GetRepositoriesSubDelOwnerRepo(ctx, owner, reponame)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

// Tasks returns the tasks with pagination and optional filtering. (Just pass empty values for no filters)
func (c *Client) Tasks(ctx context.Context, repository, sha *string, page, perPage *int64) ([]*model.Task, error) {
	resp, err := c.client.GetTasks(ctx, &uisvc.GetTasksParams{Page: page, PerPage: perPage, Repository: repository, Sha: sha})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []*model.Task{}

	return ret, json.NewDecoder(resp.Body).Decode(&ret)
}

// TaskCount returns the total number of tasks matching the filter.
func (c *Client) TaskCount(ctx context.Context, repository, sha *string) (int64, error) {
	resp, err := c.client.GetTasksCount(ctx, &uisvc.GetTasksCountParams{Repository: repository, Sha: sha})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var count int64
	return count, json.NewDecoder(resp.Body).Decode(&count)
}

// RunsForTask returns the runs for the provided task id.
func (c *Client) RunsForTask(ctx context.Context, taskID int64, page, perPage *int64) ([]*model.Run, error) {
	resp, err := c.client.GetTasksRunsId(ctx, taskID, &uisvc.GetTasksRunsIdParams{Page: page, PerPage: perPage})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []*model.Run{}
	return ret, json.NewDecoder(resp.Body).Decode(&ret)
}

// RunsForTaskCount returns the count of the runs for the provided task id.
func (c *Client) RunsForTaskCount(ctx context.Context, taskID int64) (int64, error) {
	resp, err := c.client.GetTasksRunsIdCount(ctx, taskID)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var count int64
	return count, json.NewDecoder(resp.Body).Decode(&count)
}

// Runs returns all the runs matching the filter set with pagination.
func (c *Client) Runs(ctx context.Context, repository, sha *string, page, perPage *int64) ([]*model.Run, error) {
	resp, err := c.client.GetRuns(ctx, &uisvc.GetRunsParams{Page: page, PerPage: perPage, Repository: repository, Sha: sha})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []*model.Run{}
	return ret, json.NewDecoder(resp.Body).Decode(&ret)
}

// RunsCount returns the count of the runs.
func (c *Client) RunsCount(ctx context.Context, repository, sha *string) (int64, error) {
	resp, err := c.client.GetRunsCount(ctx, &uisvc.GetRunsCountParams{Repository: repository, Sha: sha})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var count int64
	return count, json.NewDecoder(resp.Body).Decode(&count)
}

// GetRun retrieves a run by id.
func (c *Client) GetRun(ctx context.Context, id int64) (*model.Run, error) {
	resp, err := c.client.GetRunRunId(ctx, id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := &model.Run{}
	return ret, json.NewDecoder(resp.Body).Decode(&ret)
}

// CancelRun cancels the run by id. It may also cancel other runs based on rules
// around queue management.
func (c *Client) CancelRun(ctx context.Context, id int64) error {
	resp, err := c.client.PostCancelRunId(ctx, id)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// AddCapability adds a capability for a user. Must have the modify:user capability to interact.
func (c *Client) AddCapability(ctx context.Context, username string, capability model.Capability) error {
	resp, err := c.client.PostCapabilitiesUsernameCapability(ctx, string(capability), username)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// RemoveCapability removes a capability from a user. Must have the modify:user capability to interact.
func (c *Client) RemoveCapability(ctx context.Context, username string, capability model.Capability) error {
	resp, err := c.client.DeleteCapabilitiesUsernameCapability(ctx, string(capability), username)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// GetUserProperties returns some properties about the requesting account; like username and capabilities.
func (c *Client) GetUserProperties(ctx context.Context) (map[string]interface{}, error) {
	resp, err := c.client.GetUserProperties(ctx)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := map[string]interface{}{}

	return res, json.NewDecoder(resp.Body).Decode(&res)
}

// VisibleRepos retrieves the visible repositories to the user, a search may also be provided to limit scope.
func (c *Client) VisibleRepos(ctx context.Context, search *string) ([]*model.Repository, error) {
	resp, err := c.client.GetRepositoriesVisible(ctx, &uisvc.GetRepositoriesVisibleParams{Search: search})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	newRepos := model.RepositoryList{}
	return newRepos, json.NewDecoder(resp.Body).Decode(&newRepos)
}

// Submissions returns a list of submissions, paginated and optionally filtered by repository and SHA.
func (c *Client) Submissions(ctx context.Context, repository, sha *string, page, perPage *int64) ([]*model.Submission, error) {
	resp, err := c.client.GetSubmissions(ctx, &uisvc.GetSubmissionsParams{Page: page, PerPage: perPage, Repository: repository, Sha: sha})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	newSubs := []*model.Submission{}
	return newSubs, json.NewDecoder(resp.Body).Decode(&newSubs)
}

// TasksForSubmission returns the tasks for the given submission.
func (c *Client) TasksForSubmission(ctx context.Context, sub *model.Submission) ([]*model.Task, error) {
	perPage := int64(20)
	page := int64(0)

	totalTasks := []*model.Task{}

	for {
		resp, err := c.client.GetSubmissionIdTasks(ctx, sub.ID, &uisvc.GetSubmissionIdTasksParams{Page: &page, PerPage: &perPage})
		if err != nil {
			return nil, err
		}

		jsonTasks := []*model.Task{}
		if err := json.NewDecoder(resp.Body).Decode(&jsonTasks); err != nil {
			resp.Body.Close()
			return nil, err
		}

		resp.Body.Close()

		totalTasks = append(totalTasks, jsonTasks...)

		if len(jsonTasks) == 0 {
			break
		}

		page++
	}

	return totalTasks, nil
}
