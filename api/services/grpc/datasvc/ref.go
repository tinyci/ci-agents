package datasvc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefByNameAndSHA retrieves the ref from repository and sha data.
func (ds *DataServer) GetRefByNameAndSHA(ctx context.Context, rp *data.RefPair) (*types.Ref, error) {
	ref, err := ds.H.Model.GetRefByNameAndSHA(ctx, rp.RepoName, rp.Sha)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "%v", err)
		}

		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	ret, err := ds.C.ToProto(ctx, ref)
	if err != nil {
		return nil, err
	}

	return ret.(*types.Ref), nil
}

// PutRef adds a ref to the database.
func (ds *DataServer) PutRef(ctx context.Context, ref *types.Ref) (*types.Ref, error) {
	ret, err := ds.C.FromProto(ctx, ref)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.PutRef(ctx, ret.(*models.Ref)); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	ret, err = ds.C.ToProto(ctx, ret)
	if err != nil {
		return nil, err
	}

	return ret.(*types.Ref), nil
}

// CancelRefByName cancels by the ref for all runs related to it. It is looked
// up by ref and repository information. It is used by the queuesvc to auto
// cancel runs as new ones are being submitted.
func (ds *DataServer) CancelRefByName(ctx context.Context, rr *data.RepoRef) (*empty.Empty, error) {
	if err := ds.H.Model.CancelRefByName(ctx, rr.Repository, rr.RefName); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}
