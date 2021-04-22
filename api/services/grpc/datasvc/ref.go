package datasvc

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRefByNameAndSHA retrieves the ref from repository and sha data.
func (ds *DataServer) GetRefByNameAndSHA(ctx context.Context, rp *data.RefPair) (*types.Ref, error) {
	ref, err := ds.H.Model.GetRefByNameAndSHA(rp.RepoName, rp.Sha)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "%v", err)
		}

		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}
	return ref.ToProto(), nil
}

// PutRef adds a ref to the database.
func (ds *DataServer) PutRef(ctx context.Context, ref *types.Ref) (*types.Ref, error) {
	ret, err := model.NewRefFromProto(ref)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.PutRef(ret); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return ret.ToProto(), nil
}

// CancelRefByName cancels by the ref for all runs related to it. It is looked
// up by ref and repository information. It is used by the queuesvc to auto
// cancel runs as new ones are being submitted.
func (ds *DataServer) CancelRefByName(ctx context.Context, rr *data.RepoRef) (*empty.Empty, error) {
	if err := ds.H.Model.CancelRefByName(rr.Repository, rr.RefName, ds.H.URL, config.DefaultGithubClient("")); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}
