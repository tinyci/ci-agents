package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
)

// PutSubmission puts a submission into the datasvc. Updates the created_at time.
func (c *Client) PutSubmission(ctx context.Context, sub *model.Submission) (*model.Submission, error) {
	s, err := c.client.PutSubmission(ctx, sub.ToProto())
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewSubmissionFromProto(s)
}

// GetSubmissionByID returns the submission for the given ID.
func (c *Client) GetSubmissionByID(ctx context.Context, id int64) (*model.Submission, error) {
	s, err := c.client.GetSubmission(ctx, &types.IntID{ID: id})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewSubmissionFromProto(s)
}

// GetRunsForSubmission returns the runs for the given submission; with pagination
func (c *Client) GetRunsForSubmission(ctx context.Context, sub *model.Submission, page, perPage int64) ([]*model.Run, error) {
	runs, err := c.client.GetSubmissionRuns(ctx, &data.SubmissionQuery{Submission: sub.ToProto(), Page: page, PerPage: perPage})
	if err != nil {
		return nil, errors.New(err)
	}

	list := []*model.Run{}

	for _, run := range runs.List {
		r, err := model.NewRunFromProto(run)
		if err != nil {
			return nil, err
		}
		list = append(list, r)
	}

	return list, nil
}

// GetTasksForSubmission returns the tasks for the given submission; with pagination
func (c *Client) GetTasksForSubmission(ctx context.Context, sub *model.Submission, page, perPage int64) ([]*model.Task, error) {
	tasks, err := c.client.GetSubmissionTasks(ctx, &data.SubmissionQuery{Submission: sub.ToProto(), Page: page, PerPage: perPage})
	if err != nil {
		return nil, errors.New(err)
	}

	list := []*model.Task{}

	for _, task := range tasks.Tasks {
		t, err := model.NewTaskFromProto(task)
		if err != nil {
			return nil, err
		}
		list = append(list, t)
	}

	return list, nil
}

// ListSubmissions lists the submissions with pagination, and an optional (just
// pass empty strings if undesired) repository and sha filter.
func (c *Client) ListSubmissions(ctx context.Context, page, perPage int64, repository, sha string) ([]*model.Submission, error) {
	list, err := c.client.ListSubmissions(ctx, &data.RepositoryFilterRequestWithPagination{Page: page, PerPage: perPage, Repository: repository, Sha: sha})
	if err != nil {
		return nil, errors.New(err)
	}

	subs := []*model.Submission{}

	for _, sub := range list.Submissions {
		newSub, err := model.NewSubmissionFromProto(sub)
		if err != nil {
			return nil, err
		}

		subs = append(subs, newSub)
	}

	return subs, nil
}

// CountSubmissions returns the count of all submissions that meet the optional
// filtering requirements.
func (c *Client) CountSubmissions(ctx context.Context, repository, sha string) (int64, error) {
	count, err := c.client.CountSubmissions(ctx, &data.RepositoryFilterRequest{Repository: repository, Sha: sha})
	if err != nil {
		return 0, errors.New(err)
	}

	return count.Count, nil
}

// CancelSubmission cancels a submission by ID.
func (c *Client) CancelSubmission(ctx context.Context, id int64) error {
	if _, err := c.client.CancelSubmission(ctx, &types.IntID{ID: id}); err != nil {
		return errors.New(err)
	}

	return nil
}
