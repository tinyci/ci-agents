package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
)

// CommentError is for commenting on PRs when there is no better means of bubbling up an error.
func (c *Client) CommentError(ctx context.Context, repoName string, prID int64, errstr string) error {
	_, err := c.client.CommentError(ctx, &repository.CommentErrorRequest{RepoName: repoName, PrID: prID, Error: errstr})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// PendingStatus updates the status for the sha for the given repo.
func (c *Client) PendingStatus(ctx context.Context, repoName, sha, runName, url string) error {
	_, err := c.client.PendingStatus(ctx, &repository.StatusRequest{RepoName: repoName, Sha: sha, RunName: runName, Url: url})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// StartedStatus updates the status for the sha for the given repo.
func (c *Client) StartedStatus(ctx context.Context, repoName, sha, runName, url string) error {
	_, err := c.client.StartedStatus(ctx, &repository.StatusRequest{RepoName: repoName, Sha: sha, RunName: runName, Url: url})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// ErrorStatus updates the status for the sha for the given repo.
func (c *Client) ErrorStatus(ctx context.Context, repoName, sha, runName, url string, eErr error) error {
	_, err := c.client.ErrorStatus(ctx, &repository.ErrorStatusRequest{
		RepoName: repoName,
		Sha:      sha,
		RunName:  runName,
		Url:      url,
		Error:    eErr.Error(),
	})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// FinishedStatus updates the status for the sha for the given repo.
func (c *Client) FinishedStatus(ctx context.Context, repoName, sha, runName, url string, status bool, msg string) error {
	_, err := c.client.FinishedStatus(ctx, &repository.FinishedStatusRequest{
		RepoName: repoName,
		Sha:      sha,
		RunName:  runName,
		Url:      url,
		Status:   status,
		Msg:      msg,
	})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// ClearStates removes all status reports from a SHA in an attempt to restart
// the process.
func (c *Client) ClearStates(ctx context.Context, repoName, sha string) error {
	_, err := c.client.ClearStates(ctx, &repository.RepoSHAPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return errors.New(err)
	}

	return nil
}
