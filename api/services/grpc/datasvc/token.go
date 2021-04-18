package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetToken retrieves a token, creating it if necessary, for the user supplied.
// This token can be presented by the user as a part of the authentication
// process to login to tinyci and perform operations.
//
// If the token is already set, then this function will return an error and
// refuse to yield the existing token.
func (ds *DataServer) GetToken(ctx context.Context, name *data.Name) (*types.StringID, error) {
	token, err := ds.H.Model.GetToken(name.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &types.StringID{ID: token}, nil
}

// DeleteToken removes the existing token; a new GetToken request will generate a new one.
func (ds *DataServer) DeleteToken(ctx context.Context, name *data.Name) (*empty.Empty, error) {
	if err := ds.H.Model.DeleteToken(name.Name); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// ValidateToken accepts the token and returns an error if the user cannot login.
func (ds *DataServer) ValidateToken(ctx context.Context, id *types.StringID) (*types.User, error) {
	u, err := ds.H.Model.ValidateToken(id.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return u.ToProto(), nil
}
