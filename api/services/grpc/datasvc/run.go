package datasvc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RunCount is the count of items in the queue
func (ds *DataServer) RunCount(ctx context.Context, rp *data.RefPair) (*data.Count, error) {
	var (
		res int64
	)

	// we don't want to return an error upon finding zero runs. We want to return 0.
	// this comes up in a few places in this routine.
	dead := &data.Count{Count: 0}

	if rp.RepoName != "" {
		repo, err := ds.H.Model.GetRepositoryByName(ctx, rp.RepoName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return dead, nil
			}
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		}

		if rp.Sha != "" {
			res, err = ds.H.Model.RunTotalCountForRepositoryAndSHA(ctx, repo, rp.Sha)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return dead, nil
				}
				return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
			}
		} else {
			res, err = ds.H.Model.RunTotalCountForRepository(ctx, repo)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return dead, nil
				}
				return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
			}
		}
	} else {
		var err error
		res, err = ds.H.Model.RunTotalCount(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return dead, nil
			}
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		}
	}

	return &data.Count{Count: res}, nil
}

// RunList lists the runs with pagination
func (ds *DataServer) RunList(ctx context.Context, rq *data.RunListRequest) (*types.RunList, error) {
	pg, ppg, err := utils.ScopePaginationInt(&rq.Page, &rq.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	list, err := ds.H.Model.RunList(ctx, int64(pg), int64(ppg), rq.Repository, rq.Sha)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	ret := &types.RunList{}

	for _, run := range list {
		r, err := ds.C.ToProto(ctx, run)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		ret.List = append(ret.List, r.(*types.Run))
	}

	return ret, nil
}

// GetRun retrieves a run by id.
func (ds *DataServer) GetRun(ctx context.Context, id *types.IntID) (*types.Run, error) {
	run, err := ds.H.Model.GetRun(ctx, id.ID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	ret, err := ds.C.ToProto(ctx, run)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return ret.(*types.Run), nil
}

// GetRunUI retrieves a run by id.
func (ds *DataServer) GetRunUI(ctx context.Context, id *types.IntID) (*types.Run, error) {
	return ds.GetRun(ctx, id)
}
