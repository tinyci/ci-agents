package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CancelTask cancels a task by ID.
func (ds *DataServer) CancelTask(ctx context.Context, id *types.IntID) (*empty.Empty, error) {
	if err := ds.H.Model.CancelTask(ctx, id.ID); err != nil {
		return nil, utils.WrapError(err, "could not cancel runs for for task_id %d", id.ID)
	}

	return &empty.Empty{}, nil
}

// CancelTasksByPR cancels multiple tasks by Pull Request ID.
func (ds *DataServer) CancelTasksByPR(ctx context.Context, prq *types.CancelPRRequest) (*empty.Empty, error) {
	if err := ds.H.Model.CancelTaskForPR(ctx, prq.Repository, prq.Id); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Could not cancel tasks for repo %q, PR #%d: %v", prq.Repository, prq.Id, err)
	}
	return &empty.Empty{}, nil
}

// PutTask creates a new task for use. Returns the ID of the created task.
func (ds *DataServer) PutTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	t, err := ds.C.FromProto(ctx, task)
	if err != nil {
		return nil, err
	}

	if err := ds.H.Model.PutTask(ctx, t.(*models.Task)); err != nil {
		return nil, err
	}

	task2, err := ds.C.ToProto(ctx, t)
	if err != nil {
		return nil, err
	}

	return task2.(*types.Task), nil
}

// ListTasks returns a list of tasks based on the query.
func (ds *DataServer) ListTasks(ctx context.Context, req *data.TaskListRequest) (*types.TaskList, error) {
	tasks, err := ds.H.Model.ListTasks(ctx, req.Repository, req.Sha, req.Page, req.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	retTasks := &types.TaskList{}

	for _, task := range tasks {
		t, err := ds.C.ToProto(ctx, task)
		if err != nil {
			return nil, err
		}
		retTasks.Tasks = append(retTasks.Tasks, t.(*types.Task))
	}

	return retTasks, nil
}

// CountTasks counts the number of tasks that would be found by the query
func (ds *DataServer) CountTasks(ctx context.Context, req *data.TaskListRequest) (*data.Count, error) {
	count, err := ds.H.Model.CountTasks(ctx, req.Repository, req.Sha)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &data.Count{Count: count}, nil
}

// RunsForTask retrieves all the runs for the task, with optional pagination.
func (ds *DataServer) RunsForTask(ctx context.Context, req *data.RunsForTaskRequest) (*types.RunList, error) {
	runs, err := ds.H.Model.GetRunsForTask(ctx, req.Id, req.Page, req.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	rl := &types.RunList{}

	for _, run := range runs {
		r, err := ds.C.ToProto(ctx, run)
		if err != nil {
			return nil, err
		}
		rl.List = append(rl.List, r.(*types.Run))
	}

	return rl, nil
}

// CountRunsForTask counts all the runs for a given task.
func (ds *DataServer) CountRunsForTask(ctx context.Context, id *types.IntID) (*data.Count, error) {
	count, err := ds.H.Model.CountRunsForTask(ctx, id.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &data.Count{Count: count}, nil
}
