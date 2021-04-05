package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PutSession saves a session created for a user.
func (ds *DataServer) PutSession(ctx context.Context, s *types.Session) (*empty.Empty, error) {
	if err := ds.H.Model.SaveSession(model.NewSessionFromProto(s)); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// LoadSession retrieves a session by ID.
func (ds *DataServer) LoadSession(ctx context.Context, id *types.StringID) (*types.Session, error) {
	s, err := ds.H.Model.LoadSession(id.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return s.ToProto(), nil
}
