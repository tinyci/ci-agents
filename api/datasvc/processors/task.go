package processors

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc/codes"
)

// CancelTasksByPR cancels multiple tasks by Pull Request ID.
func (ds *DataServer) CancelTasksByPR(ctx context.Context, prq *types.CancelPRRequest) (*empty.Empty, error) {
	if err := ds.H.Model.CancelTasksForPR(prq.Repository, prq.Id, ds.H.URL); err != nil {
		return nil, err.Wrapf("Could not cancel tasks for repo %q, PR #%d", prq.Repository, prq.Id).ToGRPC(codes.FailedPrecondition)
	}
	return &empty.Empty{}, nil
}

// PutTask creates a new task for use. Returns the ID of the created task.
func (ds *DataServer) PutTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	t, err := model.NewTaskFromProto(task)
	if err != nil {
		return nil, err
	}

	if err := ds.H.Model.Create(t).Error; err != nil {
		return nil, errors.New(err).ToGRPC(codes.FailedPrecondition)
	}

	return t.ToProto(), nil
}

// ListTasks returns a list of tasks based on the query.
func (ds *DataServer) ListTasks(ctx context.Context, req *data.TaskListRequest) (*types.TaskList, error) {
	tasks, err := ds.H.Model.ListTasks(req.Repository, req.Sha, req.Page, req.PerPage)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	retTasks := &types.TaskList{}

	for _, task := range tasks {
		retTasks.Tasks = append(retTasks.Tasks, task.ToProto())
	}

	return retTasks, nil
}

// CountTasks counts the number of tasks that would be found by the query
func (ds *DataServer) CountTasks(ctx context.Context, req *data.TaskListRequest) (*data.Count, error) {
	count, err := ds.H.Model.CountTasks(req.Repository, req.Sha)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	return &data.Count{Count: count}, nil
}

// RunsForTask retrieves all the runs for the task, with optional pagination.
func (ds *DataServer) RunsForTask(ctx context.Context, req *data.RunsForTaskRequest) (*types.RunList, error) {
	runs, err := ds.H.Model.GetRunsForTask(req.Id, req.Page, req.PerPage)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	rl := &types.RunList{}

	for _, run := range runs {
		rl.List = append(rl.List, run.ToProto())
	}

	return rl, nil
}

// CountRunsForTask counts all the runs for a given task.
func (ds *DataServer) CountRunsForTask(ctx context.Context, id *types.IntID) (*data.Count, error) {
	count, err := ds.H.Model.CountRunsForTask(id.ID)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	return &data.Count{Count: count}, nil
}

// ListSubscribedTasksForUser mirrors the model call by the same name.
func (ds *DataServer) ListSubscribedTasksForUser(ctx context.Context, lstr *data.ListSubscribedTasksRequest) (*types.TaskList, error) {
	tasks, err := ds.H.Model.ListSubscribedTasksForUser(lstr.Id, lstr.Page, lstr.PerPage)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}
	grpcTask := &types.TaskList{}

	for _, task := range tasks {
		grpcTask.Tasks = append(grpcTask.Tasks, task.ToProto())
	}

	return grpcTask, nil
}
