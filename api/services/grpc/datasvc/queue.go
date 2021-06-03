package datasvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db"
	"github.com/tinyci/ci-agents/db/models"
	topTypes "github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueueCount is the count of items in the queue
func (ds *DataServer) QueueCount(ctx context.Context, empty *empty.Empty) (*data.Count, error) {
	res, err := ds.H.Model.QueueTotalCount(ctx)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &data.Count{Count: res}, nil
}

// QueueCountForRepository is the count of items in the queue for the given repository
func (ds *DataServer) QueueCountForRepository(ctx context.Context, repo *data.Name) (*data.Count, error) {
	r, err := ds.H.Model.GetRepositoryByName(ctx, repo.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	res, err := ds.H.Model.QueueTotalCountForRepository(ctx, r.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &data.Count{Count: res}, nil
}

// QueueListForRepository lists the queue with pagination
func (ds *DataServer) QueueListForRepository(ctx context.Context, qlr *data.QueueListRequest) (*data.QueueList, error) {
	r, err := ds.H.Model.GetRepositoryByName(ctx, qlr.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	list, err := ds.H.Model.QueueListForRepository(ctx, r.ID, qlr.Page, qlr.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	retList := &data.QueueList{}

	for _, item := range list {
		qi, err := ds.C.ToProto(ctx, item)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		retList.Items = append(retList.Items, qi.(*types.QueueItem))
	}

	return retList, nil
}

// QueueAdd adds an item to the queue
func (ds *DataServer) QueueAdd(ctx context.Context, list *data.QueueList) (*data.QueueList, error) {
	modelItems := []*models.QueueItem{}

	for _, item := range list.Items {
		r, err := ds.C.FromProto(ctx, item.Run)
		if err != nil {
			return nil, err
		}

		if err := ds.H.Model.PutRun(ctx, r.(*models.Run)); err != nil {
			return nil, err
		}

		run, err := ds.C.ToProto(ctx, r)
		if err != nil {
			return nil, err
		}

		item.Run = run.(*types.Run)

		it, err := ds.C.FromProto(ctx, item)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		}

		modelItems = append(modelItems, it.(*models.QueueItem))
	}

	if err := ds.H.Model.QueuePipelineAdd(ctx, modelItems); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	retList := &data.QueueList{}

	for _, item := range modelItems {
		qi, err := ds.C.ToProto(ctx, item)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		retList.Items = append(retList.Items, qi.(*types.QueueItem))
	}

	return retList, nil
}

// QueueNext returns the next item for the named queue.
func (ds *DataServer) QueueNext(ctx context.Context, r *types.QueueRequest) (*types.QueueItem, error) {
	qi, err := ds.H.Model.NextQueueItem(ctx, r.RunningOn, r.QueueName)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		if stat, ok := status.FromError(err); ok {
			return nil, stat.Err()
		}

		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	ret, err := ds.C.ToProto(ctx, qi)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return ret.(*types.QueueItem), nil
}

// PutStatus sets the status for the given run_id
func (ds *DataServer) PutStatus(ctx context.Context, s *types.Status) (*empty.Empty, error) {
	u, err := ds.H.Model.GetOwnerForRun(ctx, s.Id)
	if err != nil {
		return nil, err
	}

	if err := ds.H.Model.SetRunStatus(ctx, s.Id, s.Status); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	var token topTypes.OAuthToken

	if err := u.Token.Unmarshal(&token); err != nil {
		return nil, err
	}

	bits, err := ds.H.Model.GetRunDetail(ctx, s.Id)
	if err != nil {
		return nil, err
	}

	go func(ds *DataServer, u *models.User, bits *db.RunDetail) {
		client := ds.H.OAuth.GithubClient(u.Username, token.Token)

		if err := client.FinishedStatus(context.Background(), bits.Owner, bits.Repo, bits.Run.Name, bits.HeadSHA, fmt.Sprintf("%s/log/%d", ds.H.URL, s.Id), s.Status, "The run completed!"); err != nil {
			ds.H.Clients.Log.Error(context.Background(), err)
		}
	}(ds, u, bits)

	return &empty.Empty{}, nil
}

// SetCancel flags the run (which will flag the rest of the task's runs) as
// canceled. Will fail on finished tasks.
func (ds *DataServer) SetCancel(ctx context.Context, id *types.IntID) (*empty.Empty, error) {
	u, err := ds.H.Model.GetOwnerForRun(ctx, id.ID)
	if err != nil {
		return nil, err
	}

	var token topTypes.OAuthToken

	if err := u.Token.Unmarshal(&token); err != nil {
		return nil, err
	}

	if err := ds.H.Model.CancelRun(ctx, id.ID); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	bits, err := ds.H.Model.GetRunDetail(ctx, id.ID)
	if err != nil {
		return nil, err
	}

	go func(ds *DataServer, u *models.User, bits *db.RunDetail) {
		client := ds.H.UserConfig.OAuth.GithubClient(u.Username, token.Token)

		if err := client.ErrorStatus(context.Background(), bits.Owner, bits.Repo, bits.Run.Name, bits.HeadSHA, fmt.Sprintf("%s/log/%d", ds.H.URL, id.ID), utils.ErrRunCanceled); err != nil {
			ds.H.Clients.Log.Error(context.Background(), err)
		}
	}(ds, u, bits)

	return &empty.Empty{}, nil
}

// GetCancel returns the canceled state for the run.
func (ds *DataServer) GetCancel(ctx context.Context, id *types.IntID) (*types.Status, error) {
	s := &types.Status{Id: id.ID}
	task, err := ds.H.Model.GetTaskForRun(ctx, id.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	s.Status = task.Canceled
	return s, nil
}
