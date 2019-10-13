package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// ListSubscriptions lists the subscriptions that the user has selected.
func (c *Client) ListSubscriptions(ctx context.Context, name, search string) (model.RepositoryList, *errors.Error) {
	rl, err := c.client.ListSubscriptions(ctx, &data.NameSearch{Name: name, Search: search}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return makeRepoList(rl)
}

// AddSubscription adds a subscription for the user.
func (c *Client) AddSubscription(ctx context.Context, name, repo string) *errors.Error {
	_, err := c.client.AddSubscription(ctx, &data.RepoUserSelection{RepoName: repo, Username: name}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteSubscription removes a subscription for the user.
func (c *Client) DeleteSubscription(ctx context.Context, name, repo string) *errors.Error {
	// sigh.. these names.
	_, err := c.client.RemoveSubscription(ctx, &data.RepoUserSelection{RepoName: repo, Username: name}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}
