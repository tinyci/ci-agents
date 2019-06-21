package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// CancelTasksByPR cancels tasks by PR ID.
func (c *Client) CancelTasksByPR(repository string, prID int64) *errors.Error {
	if _, err := c.client.CancelTasksByPR(context.Background(), &types.CancelPRRequest{Repository: repository, Id: prID}, grpc.WaitForReady(true)); err != nil {
		return errors.New(err)
	}

	return nil
}

// PutTask adds a task to the database.
func (c *Client) PutTask(task *model.Task) (*model.Task, *errors.Error) {
	t, err := c.client.PutTask(context.Background(), task.ToProto(), grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewTaskFromProto(t)
}

// ListTasks returns the items in the task list that match the repository and
// sha parameters; they may also be blank to select all items. page and perPage
// are limiters to define pagination rules.
func (c *Client) ListTasks(repository, sha string, page, perPage int64) ([]*model.Task, *errors.Error) {
	tasks, err := c.client.ListTasks(context.Background(), &data.TaskListRequest{
		Repository: repository,
		Sha:        sha,
		Page:       page,
		PerPage:    perPage,
	}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	retTask := []*model.Task{}

	for _, task := range tasks.Tasks {
		t, err := model.NewTaskFromProto(task)
		if err != nil {
			return nil, err
		}

		retTask = append(retTask, t)
	}

	return retTask, nil
}

// CountTasks counts the tasks with the filters applied.
func (c *Client) CountTasks(repository, sha string) (int64, *errors.Error) {
	count, err := c.client.CountTasks(context.Background(), &data.TaskListRequest{Repository: repository, Sha: sha}, grpc.WaitForReady(true))
	if err != nil {
		return 0, errors.New(err)
	}

	return count.Count, nil
}

// GetRunsForTask retrieves all the runs by task ID.
func (c *Client) GetRunsForTask(taskID, page, perPage int64) ([]*model.Run, *errors.Error) {
	runs, err := c.client.RunsForTask(context.Background(), &data.RunsForTaskRequest{Id: taskID, Page: page, PerPage: perPage}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	modelRuns := []*model.Run{}

	for _, run := range runs.List {
		r, err := model.NewRunFromProto(run)
		if err != nil {
			return nil, errors.New(err)
		}

		modelRuns = append(modelRuns, r)
	}

	return modelRuns, nil
}

// CountRunsForTask counts all the runs associated with the task.
func (c *Client) CountRunsForTask(taskID int64) (int64, *errors.Error) {
	count, err := c.client.CountRunsForTask(context.Background(), &types.IntID{ID: taskID}, grpc.WaitForReady(true))
	if err != nil {
		return 0, errors.New(err)
	}

	return count.Count, nil
}

// ListSubscribedTasksForUser lists all the tasks for the repos the user is subscribed to.
func (c *Client) ListSubscribedTasksForUser(userID, page, perPage int64) ([]*model.Task, *errors.Error) {
	modelTasks := []*model.Task{}

	tasks, err := c.client.ListSubscribedTasksForUser(context.Background(), &data.ListSubscribedTasksRequest{Id: userID, Page: page, PerPage: perPage}, grpc.WaitForReady(true))
	if err != nil {
		return modelTasks, errors.New(err)
	}

	for _, task := range tasks.Tasks {
		t, err := model.NewTaskFromProto(task)
		if err != nil {
			return modelTasks, err
		}

		modelTasks = append(modelTasks, t)
	}

	return modelTasks, nil
}

// CancelTask cancels a task by id.
func (c *Client) CancelTask(id int64) *errors.Error {
	_, err := c.client.CancelTask(context.Background(), &types.IntID{ID: id})
	return errors.New(err)
}
