package processors

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// PutSession saves a session created for a user.
func (ds *DataServer) PutSession(ctx context.Context, s *types.Session) (*empty.Empty, error) {
	if err := ds.H.Model.SaveSession(model.NewSessionFromProto(s)); err != nil {
		return nil, errors.New(err)
	}

	return &empty.Empty{}, nil
}

// LoadSession retrieves a session by ID.
func (ds *DataServer) LoadSession(ctx context.Context, id *types.StringID) (*types.Session, error) {
	s, err := ds.H.Model.LoadSession(id.ID)
	if err != nil {
		return nil, err
	}

	return s.ToProto(), nil
}
