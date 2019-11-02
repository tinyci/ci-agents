package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
)

// MyRepositories returns all the writable repositories accessible to user
// owning the access key
func (c *Client) MyRepositories(ctx context.Context, u *model.User) ([]*repository.RepositoryData, error) {
	list, err := c.client.MyRepositories(ctx, u.ToProto())
	if err != nil {
		return nil, errors.New(err)
	}

	return list.Repositories, nil
}

// GetRepository retrieves the github response for a given repository.
func (c *Client) GetRepository(ctx context.Context, u *model.User, repoName string) (*repository.RepositoryData, error) {
	data, err := c.client.GetRepository(ctx, &repository.UserWithRepo{User: u.ToProto(), RepoName: repoName})
	if err != nil {
		return nil, errors.New(err)
	}

	return data, nil
}
