package data

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

func makeRepoList(list *types.RepositoryList) (model.RepositoryList, *errors.Error) {
	rl := model.RepositoryList{}

	for _, repo := range list.List {
		pr, err := model.NewRepositoryFromProto(repo)
		if err != nil {
			return nil, errors.New(err)
		}

		rl = append(rl, pr)
	}

	return rl, nil
}

// GetRepository retrieves a repository by name.
func (c *Client) GetRepository(name string) (*model.Repository, *errors.Error) {
	repo, err := c.client.GetRepository(context.Background(), &data.Name{Name: name})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewRepositoryFromProto(repo)
}

// PutRepositories takes a list of github repositories and adds them to the database for the user as owner.
func (c *Client) PutRepositories(name string, github []*github.Repository, autoCreated bool) *errors.Error {
	content, err := json.Marshal(github)
	if err != nil {
		return errors.New(err)
	}

	_, err = c.client.SaveRepositories(context.Background(), &data.GithubJSON{JSON: content, Username: name, AutoCreated: autoCreated})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// EnableRepository enables a repository in CI for a user as owner.
func (c *Client) EnableRepository(user, name string) *errors.Error {
	_, err := c.client.EnableRepository(context.Background(), &data.RepoUserSelection{Username: user, RepoName: name})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DisableRepository disabls a repository in CI for a user as owner.
func (c *Client) DisableRepository(user, name string) *errors.Error {
	_, err := c.client.DisableRepository(context.Background(), &data.RepoUserSelection{Username: user, RepoName: name})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// OwnedRepositories lists the owned repositories by the user.
func (c *Client) OwnedRepositories(name, search string) (model.RepositoryList, *errors.Error) {
	list, err := c.client.OwnedRepositories(context.Background(), &data.NameSearch{Name: name, Search: search})
	if err != nil {
		return nil, errors.New(err)
	}

	return makeRepoList(list)
}

// AllRepositories lists all visible repositories by the user.
func (c *Client) AllRepositories(name, search string) (model.RepositoryList, *errors.Error) {
	list, err := c.client.AllRepositories(context.Background(), &data.NameSearch{Name: name, Search: search})
	if err != nil {
		return nil, errors.New(err)
	}

	return makeRepoList(list)
}

// PrivateRepositories lists all visible private repositories by the user.
func (c *Client) PrivateRepositories(name, search string) (model.RepositoryList, *errors.Error) {
	list, err := c.client.PrivateRepositories(context.Background(), &data.NameSearch{Name: name, Search: search})
	if err != nil {
		return nil, errors.New(err)
	}

	return makeRepoList(list)
}

// PublicRepositories lists all owned public repositories by the user.
func (c *Client) PublicRepositories(search string) (model.RepositoryList, *errors.Error) {
	list, err := c.client.PublicRepositories(context.Background(), &data.Search{Search: search})
	if err != nil {
		return nil, errors.New(err)
	}

	return makeRepoList(list)
}
