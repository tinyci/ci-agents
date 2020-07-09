package datasvc

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
)

// RunCount is the count of items in the queue
func (ds *DataServer) RunCount(ctx context.Context, rp *data.RefPair) (*data.Count, error) {
	var (
		res int64
	)

	if rp.RepoName != "" {
		repo, err := ds.H.Model.GetRepositoryByName(rp.RepoName)
		if err != nil {
			return nil, err.ToGRPC(codes.FailedPrecondition)
		}

		if rp.Sha != "" {
			res, err = ds.H.Model.RunTotalCountForRepositoryAndSHA(repo, rp.Sha)
			if err != nil {
				return nil, err.ToGRPC(codes.FailedPrecondition)
			}
		} else {
			res, err = ds.H.Model.RunTotalCountForRepository(repo)
			if err != nil {
				return nil, err.ToGRPC(codes.FailedPrecondition)
			}
		}
	} else {
		var err *errors.Error
		res, err = ds.H.Model.RunTotalCount()
		if err != nil {
			return nil, err.ToGRPC(codes.FailedPrecondition)
		}
	}

	return &data.Count{Count: res}, nil
}

// RunList lists the runs with pagination
func (ds *DataServer) RunList(ctx context.Context, rq *data.RunListRequest) (*types.RunList, error) {
	page, perPage, err := utils.ScopePaginationInt(rq.Page, rq.PerPage)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	list, err := ds.H.Model.RunList(page, perPage, rq.Repository, rq.Sha)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	ret := &types.RunList{}

	for _, run := range list {
		ret.List = append(ret.List, run.ToProto())
	}

	return ret, nil
}

// GetRun retrieves a run by id.
func (ds *DataServer) GetRun(ctx context.Context, id *types.IntID) (*types.Run, error) {
	run := &model.Run{}

	if err := ds.H.Model.Preload("Task.Parent").Where("id = ?", id.ID).First(run).Error; err != nil {
		return nil, errors.New(err).ToGRPC(codes.FailedPrecondition)
	}

	return run.ToProto(), nil
}

// GetRunUI retrieves a run by id.
func (ds *DataServer) GetRunUI(ctx context.Context, id *types.IntID) (*types.Run, error) {
	run := &model.Run{}

	if err := ds.H.Model.Where("id = ?", id.ID).First(run).Error; err != nil {
		return nil, errors.New(err).ToGRPC(codes.FailedPrecondition)
	}

	return run.ToProto(), nil
}
