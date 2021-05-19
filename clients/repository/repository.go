package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
)

// MyRepositories returns all the writable repositories accessible to user
// owning the access key
func (c *Client) MyRepositories(ctx context.Context, u *types.User) ([]*repository.RepositoryData, error) {
	list, err := c.client.MyRepositories(ctx, u)
	if err != nil {
		return nil, err
	}

	return list.Repositories, nil
}

// GetRepository retrieves the github response for a given repository.
func (c *Client) GetRepository(ctx context.Context, u *types.User, repoName string) (*repository.RepositoryData, error) {
	data, err := c.client.GetRepository(ctx, &repository.UserWithRepo{User: u, RepoName: repoName})
	if err != nil {
		return nil, err
	}

	return data, nil
}
