package processors

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc/codes"
)

// GetRefByNameAndSHA retrieves the ref from repository and sha data.
func (ds *DataServer) GetRefByNameAndSHA(ctx context.Context, rp *data.RefPair) (*types.Ref, error) {
	ref, err := ds.H.Model.GetRefByNameAndSHA(rp.RepoName, rp.Sha)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}
	return ref.ToProto(), nil
}

// PutRef adds a ref to the database.
func (ds *DataServer) PutRef(ctx context.Context, ref *types.Ref) (*types.Ref, error) {
	ret, err := model.NewRefFromProto(ref)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	if err := ds.H.Model.PutRef(ret); err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	return ret.ToProto(), nil
}

// CancelRefByName cancels by the ref for all runs related to it. It is looked
// up by ref and repository information. It is used by the queuesvc to auto
// cancel runs as new ones are being submitted.
func (ds *DataServer) CancelRefByName(ctx context.Context, rr *data.RepoRef) (*empty.Empty, error) {
	if err := ds.H.Model.CancelRefByName(rr.Repository, rr.RefName, ds.H.URL, config.DefaultGithubClient); err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}
