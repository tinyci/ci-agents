package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
)

// SetupHook sets up the pr webhook in the backing repository service.
func (c *Client) SetupHook(repoName, url, secret string) *errors.Error {
	_, err := c.client.SetupHook(context.Background(), &repository.HookSetupRequest{RepoName: repoName, HookURL: url, HookSecret: secret})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// TeardownHook destroys the pr webhook in the backing repository service.
func (c *Client) TeardownHook(repoName, url string) *errors.Error {
	_, err := c.client.TeardownHook(context.Background(), &repository.HookTeardownRequest{RepoName: repoName, HookURL: url})
	if err != nil {
		return errors.New(err)
	}

	return nil
}
