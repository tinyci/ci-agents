package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
)

// MyRepositories returns all the writable repositories accessible to user
// owning the access key
func (c *Client) MyRepositories(u *model.User) ([]*repository.RepositoryData, *errors.Error) {
	list, err := c.client.MyRepositories(context.Background(), u.ToProto())
	if err != nil {
		return nil, errors.New(err)
	}

	return list.Repositories, nil
}

// GetRepository retrieves the github response for a given repository.
func (c *Client) GetRepository(u *model.User, repoName string) (*repository.RepositoryData, *errors.Error) {
	data, err := c.client.GetRepository(context.Background(), &repository.UserWithRepo{User: u.ToProto(), RepoName: repoName})
	if err != nil {
		return nil, errors.New(err)
	}

	return data, nil
}
