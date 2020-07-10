package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/utils"
	"golang.org/x/oauth2"
)

var (
	// DefaultUsername controls the default username in the event NoAuth is in
	// effect; if set it will be used, otherwise an API call will be made.
	DefaultUsername string

	// Readonly disables certain actions -- most of which write content back to
	// github -- to ensure testing can be done with a real github client, without
	// affecting too much ON github. This is set through config/service.go.
	Readonly bool
)

// Client is the generic client to github operations.
type Client interface {
	CommentError(context.Context, string, int64, error) *errors.Error
	MyRepositories(context.Context) ([]*github.Repository, *errors.Error)
	GetRepository(context.Context, string) (*github.Repository, *errors.Error)
	MyLogin(context.Context) (string, *errors.Error)
	GetFileList(context.Context, string, string) ([]string, *errors.Error)
	GetSHA(context.Context, string, string) (string, *errors.Error)
	GetRefs(context.Context, string, string) ([]string, *errors.Error)
	GetFile(context.Context, string, string, string) ([]byte, *errors.Error)
	GetDiffFiles(context.Context, string, string, string) ([]string, *errors.Error)
	SetupHook(context.Context, string, string, string, string) *errors.Error
	TeardownHook(context.Context, string, string, string) *errors.Error
	PendingStatus(context.Context, string, string, string, string, string) *errors.Error
	StartedStatus(context.Context, string, string, string, string, string) *errors.Error
	ErrorStatus(context.Context, string, string, string, string, string, error) *errors.Error
	FinishedStatus(context.Context, string, string, string, string, string, bool, string) *errors.Error
	ClearStates(context.Context, string, string) *errors.Error
}

// HTTPClient encapsulates the "real world", or http client.
type HTTPClient struct {
	github *github.Client
}

// NewClientFromAccessToken turns an accessToken into a new Client.
func NewClientFromAccessToken(accessToken string) Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &HTTPClient{github: github.NewClient(tc)}
}

// PendingStatus updates the status for the sha for the given repo on github.
func (c *HTTPClient) PendingStatus(ctx context.Context, owner, repo, name, sha, url string) *errors.Error {
	if Readonly {
		return nil
	}

	_, _, err := c.github.Repositories.CreateStatus(ctx, owner, repo, sha, &github.RepoStatus{
		TargetURL:   github.String(url),
		State:       github.String("pending"),
		Description: github.String("The run will be starting soon."),
		Context:     github.String(name),
	})

	return errors.New(err)
}

// StartedStatus updates the status for the sha for the given repo on github.
func (c *HTTPClient) StartedStatus(ctx context.Context, owner, repo, name, sha, url string) *errors.Error {
	if Readonly {
		return nil
	}

	_, _, err := c.github.Repositories.CreateStatus(ctx, owner, repo, sha, &github.RepoStatus{
		TargetURL:   github.String(url),
		State:       github.String("pending"),
		Description: github.String("The run has started!"),
		Context:     github.String(name),
	})

	return errors.New(err)
}

func capStatus(str string) *string {
	if len(str) > 140 {
		return github.String(str[:140])
	}

	return github.String(str)
}

// ErrorStatus updates the status for the sha for the given repo on github.
func (c *HTTPClient) ErrorStatus(ctx context.Context, owner, repo, name, sha, url string, outErr error) *errors.Error {
	if Readonly {
		return nil
	}

	_, _, err := c.github.Repositories.CreateStatus(ctx, owner, repo, sha, &github.RepoStatus{
		TargetURL: github.String(url),
		State:     github.String("error"),
		// github statuses cap at 140c
		Description: capStatus(errors.New(outErr).Wrap("The run encountered an error").Error()),
		Context:     github.String(name),
	})

	return errors.New(err)
}

// FinishedStatus updates the status for the sha for the given repo on github.
func (c *HTTPClient) FinishedStatus(ctx context.Context, owner, repo, name, sha, url string, status bool, addlMessage string) *errors.Error {
	if Readonly {
		return nil
	}

	statusString := "failure"
	if status {
		statusString = "success"
	}

	_, _, err := c.github.Repositories.CreateStatus(ctx, owner, repo, sha, &github.RepoStatus{
		TargetURL: github.String(url),
		State:     github.String(statusString),
		// github statuses cap at 140c
		Description: capStatus(fmt.Sprintf("The run finished: %s! %v", statusString, addlMessage)),
		Context:     github.String(name),
	})

	return errors.New(err)
}

// SetupHook sets up the pr webhook in github.
func (c *HTTPClient) SetupHook(ctx context.Context, owner, repo, configAddress, hookSecret string) *errors.Error {
	if Readonly {
		return nil
	}

	_, _, err := c.github.Repositories.CreateHook(ctx, owner, repo, &github.Hook{
		URL:    github.String(configAddress),
		Events: []string{"push", "pull_request"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          configAddress,
			"content_type": "json",
			"secret":       hookSecret,
		},
	})

	return errors.New(err)
}

// TeardownHook removes the pr webhook in github.
func (c *HTTPClient) TeardownHook(ctx context.Context, owner, repo, hookURL string) *errors.Error {
	if Readonly {
		return nil
	}

	var id int64
	var i int

	for {
		hooks, _, err := c.github.Repositories.ListHooks(ctx, owner, repo, &github.ListOptions{Page: i, PerPage: 20})
		if err != nil || len(hooks) == 0 {
			break
		}
		for _, hook := range hooks {
			if hook.Config["url"] == hookURL {
				id = hook.GetID()
				goto finish
			}
		}
		i++
	}

finish:
	if id != 0 {
		_, err := c.github.Repositories.DeleteHook(ctx, owner, repo, id)
		return errors.New(err)
	}

	return nil
}

// MyRepositories returns all the writable repositories accessible to user
// owning the access key
func (c *HTTPClient) MyRepositories(ctx context.Context) ([]*github.Repository, *errors.Error) {
	var i int
	ret := map[string]*github.Repository{}
	order := []string{}

	for {
		repos, _, err := c.github.Repositories.List(
			ctx,
			"",
			&github.RepositoryListOptions{
				Visibility: "owner",
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 100,
				},
			},
		)
		if err != nil {
			return nil, errors.New(err)
		}

		for _, repo := range repos {
			if repo.GetPermissions()["admin"] {
				if _, ok := ret[repo.GetFullName()]; !ok {
					ret[repo.GetFullName()] = repo
					order = append(order, repo.GetFullName())
				}
			}
		}

		if len(repos) < 100 {
			break
		}

		i++
	}

	vals := []*github.Repository{}

	for _, value := range order {
		vals = append(vals, ret[value])
	}

	return vals, nil
}

// MyLogin returns the username calling out to the API with its key. Can either
// be seeded by OAuth or Personal Token.
func (c *HTTPClient) MyLogin(ctx context.Context) (string, *errors.Error) {
	if DefaultUsername != "" {
		return DefaultUsername, nil
	}

	u, _, err := c.github.Users.Get(ctx, "")
	if err != nil {
		return "", errors.New(err)
	}

	return u.GetLogin(), nil
}

// GetRepository retrieves the github response for a given repository.
func (c *HTTPClient) GetRepository(ctx context.Context, name string) (*github.Repository, *errors.Error) {
	owner, repo, eErr := utils.OwnerRepo(name)
	if eErr != nil {
		return nil, eErr
	}

	r, _, err := c.github.Repositories.Get(ctx, owner, repo)
	return r, errors.New(err)
}

// GetFileList finds all the files in the tree for the given repository
func (c *HTTPClient) GetFileList(ctx context.Context, repoName, sha string) ([]string, *errors.Error) {
	owner, repo, eErr := utils.OwnerRepo(repoName)
	if eErr != nil {
		return nil, eErr
	}

	tree, _, err := c.github.Git.GetTree(ctx, owner, repo, sha, true)
	if err != nil {
		return nil, errors.New(err)
	}

	files := []string{}

	for _, entry := range tree.Entries {
		files = append(files, entry.GetPath())
	}

	return files, nil
}

// GetSHA retrieves the SHA for the branch in the given repository
func (c *HTTPClient) GetSHA(ctx context.Context, repoName, refName string) (string, *errors.Error) {
	owner, repo, eErr := utils.OwnerRepo(repoName)
	if eErr != nil {
		return "", eErr
	}

	ref, _, err := c.github.Git.GetRef(ctx, owner, repo, refName)
	if err != nil {
		return "", errors.New(err)
	}

	return ref.GetObject().GetSHA(), nil
}

// GetRefs gets the refs that match the given SHA. Only heads and tags are considered.
func (c *HTTPClient) GetRefs(ctx context.Context, repoName, sha string) ([]string, *errors.Error) {
	owner, repo, eErr := utils.OwnerRepo(repoName)
	if eErr != nil {
		return nil, eErr
	}

	// FIXME pagination (sigh)
	refs, _, err := c.github.Git.ListRefs(ctx, owner, repo, nil)
	if err != nil {
		return nil, errors.New(err)
	}

	list := []string{}

	for _, ref := range refs {
		if ref.GetObject().GetSHA() == sha {
			list = append(list, strings.TrimPrefix(ref.GetRef(), "refs/"))
		}
	}

	return list, nil
}

// GetFile retrieves a file from github directly through the api. Used for
// retrieving our configuration yamls and other stuff.
func (c *HTTPClient) GetFile(ctx context.Context, repoName, sha, filename string) ([]byte, *errors.Error) {
	owner, repo, eErr := utils.OwnerRepo(repoName)
	if eErr != nil {
		return nil, eErr
	}

	tree, _, err := c.github.Git.GetTree(ctx, owner, repo, sha, true)
	if err != nil {
		return nil, errors.New(err)
	}

	for _, entry := range tree.Entries {
		if entry.GetPath() == filename {
			content, _, err := c.github.Git.GetBlobRaw(ctx, owner, repo, entry.GetSHA())
			if err != nil {
				return nil, errors.New(err)
			}

			return content, errors.New(err)
		}
	}

	return nil, errors.New("file not found")
}

// GetDiffFiles retrieves the files present in the diff between the base and the head.
func (c *HTTPClient) GetDiffFiles(ctx context.Context, repoName, base, head string) ([]string, *errors.Error) {
	owner, repo, eErr := utils.OwnerRepo(repoName)
	if eErr != nil {
		return nil, eErr
	}

	if base == strings.Repeat("0", 40) {
		return c.GetFileList(ctx, repoName, head)
	}

	if head == strings.Repeat("0", 40) {
		return []string{}, errors.New("branch deleted")
	}

	commits, _, err := c.github.Repositories.CompareCommits(ctx, owner, repo, base, head)
	if err != nil {
		return nil, errors.New(err)
	}

	files := []string{}

	for _, file := range commits.Files {
		files = append(files, file.GetFilename())
	}

	return files, nil
}

// ClearStates removes all status reports from a SHA in an attempt to restart
// the process.
func (c *HTTPClient) ClearStates(ctx context.Context, repoName, sha string) *errors.Error {
	if Readonly {
		return nil
	}

	owner, repo, err := utils.OwnerRepo(repoName)
	if err != nil {
		return err
	}

	statuses := []*github.RepoStatus{}
	contexts := map[string]struct{}{}

	var i int
	for {
		states, _, err := c.github.Repositories.ListStatuses(ctx, owner, repo, sha, &github.ListOptions{Page: i, PerPage: 200})
		if err != nil {
			return errors.New(err)
		}

		if len(states) == 0 {
			break
		}

		statuses = append(statuses, states...)
		i++
	}

	for _, status := range statuses {
		if _, ok := contexts[status.GetContext()]; ok {
			continue
		}

		contexts[status.GetContext()] = struct{}{}

		// XXX the context MUST be preserved for this to be overwritten. Do not
		// change it here.
		status.State = github.String("error")
		status.Description = github.String("The run that this test was a part of has been overridden by a new run. Pushing a new change will remove this error.")
		_, _, err := c.github.Repositories.CreateStatus(ctx, owner, repo, sha, status)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}

// CommentError is for commenting on PRs when there is no better means of bubbling up an error.
func (c *HTTPClient) CommentError(ctx context.Context, repoName string, prID int64, err error) *errors.Error {
	if Readonly {
		return nil
	}

	owner, repo, retErr := utils.OwnerRepo(repoName)
	if retErr != nil {
		return retErr
	}

	_, _, eerr := c.github.Issues.CreateComment(ctx, owner, repo, int(prID), &github.IssueComment{
		Body: github.String(fmt.Sprintf("%v", err)),
	})

	if eerr != nil {
		return errors.New(retErr)
	}

	return nil
}
