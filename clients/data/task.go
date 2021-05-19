package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"google.golang.org/grpc"
)

// CancelTasksByPR cancels tasks by PR ID.
func (c *Client) CancelTasksByPR(ctx context.Context, repository string, prID int64) error {
	if _, err := c.client.CancelTasksByPR(ctx, &types.CancelPRRequest{Repository: repository, Id: prID}, grpc.WaitForReady(true)); err != nil {
		return err
	}

	return nil
}

// PutTask adds a task to the database.
func (c *Client) PutTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	return c.client.PutTask(ctx, task, grpc.WaitForReady(true))
}

// ListTasks returns the items in the task list that match the repository and
// sha parameters; they may also be blank to select all items. page and perPage
// are limiters to define pagination rules.
func (c *Client) ListTasks(ctx context.Context, repository, sha string, page, perPage int64) (*types.TaskList, error) {
	return c.client.ListTasks(ctx, &data.TaskListRequest{
		Repository: repository,
		Sha:        sha,
		Page:       page,
		PerPage:    perPage,
	}, grpc.WaitForReady(true))
}

// CountTasks counts the tasks with the filters applied.
func (c *Client) CountTasks(ctx context.Context, repository, sha string) (int64, error) {
	count, err := c.client.CountTasks(ctx, &data.TaskListRequest{Repository: repository, Sha: sha}, grpc.WaitForReady(true))
	if err != nil {
		return 0, err
	}

	return count.Count, nil
}

// GetRunsForTask retrieves all the runs by task ID.
func (c *Client) GetRunsForTask(ctx context.Context, taskID, page, perPage int64) (*types.RunList, error) {
	return c.client.RunsForTask(ctx, &data.RunsForTaskRequest{Id: taskID, Page: page, PerPage: perPage}, grpc.WaitForReady(true))
}

// CountRunsForTask counts all the runs associated with the task.
func (c *Client) CountRunsForTask(ctx context.Context, taskID int64) (int64, error) {
	count, err := c.client.CountRunsForTask(ctx, &types.IntID{ID: taskID}, grpc.WaitForReady(true))
	if err != nil {
		return 0, err
	}

	return count.Count, nil
}

// CancelTask cancels a task by id.
func (c *Client) CancelTask(ctx context.Context, id int64) error {
	_, err := c.client.CancelTask(ctx, &types.IntID{ID: id})
	return err
}
