package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// ListSubscriptions lists the subscriptions that the user has selected.
func (c *Client) ListSubscriptions(name, search string) (model.RepositoryList, *errors.Error) {
	rl, err := c.client.ListSubscriptions(context.Background(), &data.NameSearch{Name: name, Search: search}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return makeRepoList(rl)
}

// AddSubscription adds a subscription for the user.
func (c *Client) AddSubscription(name, repo string) *errors.Error {
	_, err := c.client.AddSubscription(context.Background(), &data.RepoUserSelection{RepoName: repo, Username: name}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteSubscription removes a subscription for the user.
func (c *Client) DeleteSubscription(name, repo string) *errors.Error {
	// sigh.. these names.
	_, err := c.client.RemoveSubscription(context.Background(), &data.RepoUserSelection{RepoName: repo, Username: name}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}
