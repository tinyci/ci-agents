package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
)

// CommentError is for commenting on PRs when there is no better means of bubbling up an error.
func (c *Client) CommentError(repoName string, prID int64, errstr string) *errors.Error {
	_, err := c.client.CommentError(context.Background(), &repository.CommentErrorRequest{RepoName: repoName, PrID: prID, Error: errstr})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// PendingStatus updates the status for the sha for the given repo.
func (c *Client) PendingStatus(repoName, sha, runName, url string) *errors.Error {
	_, err := c.client.PendingStatus(context.Background(), &repository.StatusRequest{RepoName: repoName, Sha: sha, RunName: runName, Url: url})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// StartedStatus updates the status for the sha for the given repo.
func (c *Client) StartedStatus(repoName, sha, runName, url string) *errors.Error {
	_, err := c.client.StartedStatus(context.Background(), &repository.StatusRequest{RepoName: repoName, Sha: sha, RunName: runName, Url: url})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// ErrorStatus updates the status for the sha for the given repo.
func (c *Client) ErrorStatus(repoName, sha, runName, url string, eErr error) *errors.Error {
	_, err := c.client.ErrorStatus(context.Background(), &repository.ErrorStatusRequest{
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
func (c *Client) FinishedStatus(repoName, sha, runName, url string, status bool, msg string) *errors.Error {
	_, err := c.client.FinishedStatus(context.Background(), &repository.FinishedStatusRequest{
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
func (c *Client) ClearStates(repoName, sha string) *errors.Error {
	_, err := c.client.ClearStates(context.Background(), &repository.RepoSHAPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return errors.New(err)
	}

	return nil
}
