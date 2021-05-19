package data

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"google.golang.org/grpc"
)

// GetRepository retrieves a repository by name.
func (c *Client) GetRepository(ctx context.Context, name string) (*types.Repository, error) {
	repo, err := c.client.GetRepository(ctx, &data.Name{Name: name}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// PutRepositories takes a list of github repositories and adds them to the database for the user as owner.
func (c *Client) PutRepositories(ctx context.Context, name string, github []*github.Repository, autoCreated bool) error {
	content, err := json.Marshal(github)
	if err != nil {
		return err
	}

	_, err = c.client.SaveRepositories(ctx, &data.GithubJSON{JSON: content, Username: name, AutoCreated: autoCreated}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	return nil
}

// EnableRepository enables a repository in CI for a user as owner.
func (c *Client) EnableRepository(ctx context.Context, user, name string) error {
	_, err := c.client.EnableRepository(ctx, &data.RepoUserSelection{Username: user, RepoName: name}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	return nil
}

// DisableRepository disabls a repository in CI for a user as owner.
func (c *Client) DisableRepository(ctx context.Context, user, name string) error {
	_, err := c.client.DisableRepository(ctx, &data.RepoUserSelection{Username: user, RepoName: name}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	return nil
}

// OwnedRepositories lists the owned repositories by the user.
func (c *Client) OwnedRepositories(ctx context.Context, name string, search *string) (*types.RepositoryList, error) {
	s := ""
	if search != nil {
		s = *search
	}

	return c.client.OwnedRepositories(ctx, &data.NameSearch{Name: name, Search: s}, grpc.WaitForReady(true))
}

// AllRepositories lists all visible repositories by the user.
func (c *Client) AllRepositories(ctx context.Context, name string, search *string) (*types.RepositoryList, error) {
	s := ""
	if search != nil {
		s = *search
	}

	return c.client.AllRepositories(ctx, &data.NameSearch{Name: name, Search: s}, grpc.WaitForReady(true))
}

// PrivateRepositories lists all visible private repositories by the user.
func (c *Client) PrivateRepositories(ctx context.Context, name, search string) (*types.RepositoryList, error) {
	return c.client.PrivateRepositories(ctx, &data.NameSearch{Name: name, Search: search}, grpc.WaitForReady(true))
}

// PublicRepositories lists all owned public repositories by the user.
func (c *Client) PublicRepositories(ctx context.Context, search string) (*types.RepositoryList, error) {
	return c.client.PublicRepositories(ctx, &data.Search{Search: search}, grpc.WaitForReady(true))
}
