package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
)

// PutSubmission puts a submission into the datasvc. Updates the created_at time.
func (c *Client) PutSubmission(ctx context.Context, sub *types.Submission) (*types.Submission, error) {
	return c.client.PutSubmission(ctx, sub)
}

// GetSubmissionByID returns the submission for the given ID.
func (c *Client) GetSubmissionByID(ctx context.Context, id int64) (*types.Submission, error) {
	return c.client.GetSubmission(ctx, &types.IntID{ID: id})
}

// GetRunsForSubmission returns the runs for the given submission; with pagination
func (c *Client) GetRunsForSubmission(ctx context.Context, sub *types.Submission, page, perPage int64) (*types.RunList, error) {
	return c.client.GetSubmissionRuns(ctx, &data.SubmissionQuery{Submission: sub, Page: page, PerPage: perPage})
}

// GetTasksForSubmission returns the tasks for the given submission; with pagination
func (c *Client) GetTasksForSubmission(ctx context.Context, sub *types.Submission, page, perPage int64) (*types.TaskList, error) {
	return c.client.GetSubmissionTasks(ctx, &data.SubmissionQuery{Submission: sub, Page: page, PerPage: perPage})
}

// ListSubmissions lists the submissions with pagination, and an optional (just
// pass empty strings if undesired) repository and sha filter.
func (c *Client) ListSubmissions(ctx context.Context, page, perPage int64, repository, sha string) (*types.SubmissionList, error) {
	return c.client.ListSubmissions(ctx, &data.RepositoryFilterRequestWithPagination{Page: page, PerPage: perPage, Repository: repository, Sha: sha})
}

// CountSubmissions returns the count of all submissions that meet the optional
// filtering requirements.
func (c *Client) CountSubmissions(ctx context.Context, repository, sha string) (int64, error) {
	count, err := c.client.CountSubmissions(ctx, &data.RepositoryFilterRequest{Repository: repository, Sha: sha})
	if err != nil {
		return 0, err
	}

	return count.Count, nil
}

// CancelSubmission cancels a submission by ID.
func (c *Client) CancelSubmission(ctx context.Context, id int64) error {
	if _, err := c.client.CancelSubmission(ctx, &types.IntID{ID: id}); err != nil {
		return err
	}

	return nil
}
