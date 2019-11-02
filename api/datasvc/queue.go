package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
)

// QueueCount is the count of items in the queue
func (ds *DataServer) QueueCount(ctx context.Context, empty *empty.Empty) (*data.Count, error) {
	res, err := ds.H.Model.QueueTotalCount()
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	return &data.Count{Count: res}, nil
}

// QueueCountForRepository is the count of items in the queue for the given repository
func (ds *DataServer) QueueCountForRepository(ctx context.Context, repo *data.Name) (*data.Count, error) {
	r, err := ds.H.Model.GetRepositoryByName(repo.Name)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	res, err := ds.H.Model.QueueTotalCountForRepository(r)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	return &data.Count{Count: res}, nil
}

// QueueListForRepository lists the queue with pagination
func (ds *DataServer) QueueListForRepository(ctx context.Context, qlr *data.QueueListRequest) (*data.QueueList, error) {
	r, err := ds.H.Model.GetRepositoryByName(qlr.Name)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	page, perPage, err := utils.ScopePaginationInt(qlr.Page, qlr.PerPage)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	list, err := ds.H.Model.QueueListForRepository(r, page, perPage)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	retList := &data.QueueList{}

	for _, item := range list {
		retList.Items = append(retList.Items, item.ToProto())
	}

	return retList, nil
}

// QueueAdd adds an item to the queue
func (ds *DataServer) QueueAdd(ctx context.Context, list *data.QueueList) (*data.QueueList, error) {
	modelItems := []*model.QueueItem{}

	for _, item := range list.Items {
		it, err := model.NewQueueItemFromProto(item)
		if err != nil {
			return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
		}

		modelItems = append(modelItems, it)
	}

	var err error
	if modelItems, err = ds.H.Model.QueuePipelineAdd(modelItems); err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	retList := &data.QueueList{}

	for _, item := range modelItems {
		retList.Items = append(retList.Items, item.ToProto())
	}

	return retList, nil
}

// QueueNext returns the next item for the named queue.
func (ds *DataServer) QueueNext(ctx context.Context, r *types.QueueRequest) (*types.QueueItem, error) {
	qi, err := ds.H.Model.NextQueueItem(r.RunningOn, r.QueueName)
	if err != nil {
		switch err.(type) {
		case errors.Error, *errors.Error:
			var e2 *errors.Error
			if e, ok := err.(errors.Error); ok {
				e2 = &e
			}
			if e2.Contains(errors.ErrNotFound) {
				e2.SetLog(false)
			}
		}
		return nil, err
	}

	return qi.ToProto(), nil
}

// PutStatus sets the status for the given run_id
func (ds *DataServer) PutStatus(ctx context.Context, s *types.Status) (*empty.Empty, error) {
	if err := ds.H.Model.SetRunStatus(s.Id, config.DefaultGithubClient(), s.Status, false, ds.H.URL, s.AdditionalMessage); err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

// SetCancel flags the run (which will flag the rest of the task's runs) as
// canceled. Will fail on finished tasks.
func (ds *DataServer) SetCancel(ctx context.Context, id *types.IntID) (*empty.Empty, error) {
	if err := ds.H.Model.CancelRun(id.ID, ds.H.URL, config.DefaultGithubClient()); err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

// GetCancel returns the canceled state for the run.
func (ds *DataServer) GetCancel(ctx context.Context, id *types.IntID) (*types.Status, error) {
	s := &types.Status{Id: id.ID}
	res, err := ds.H.Model.GetCancelForRun(id.ID)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	s.Status = res
	return s, nil
}
