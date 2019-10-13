package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// GetSession retrieves a session from the database by id.
func (c *Client) GetSession(ctx context.Context, id string) (*model.Session, *errors.Error) {
	s, err := c.client.LoadSession(ctx, &types.StringID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewSessionFromProto(s), nil
}

// PutSession adds a session to the database.
func (c *Client) PutSession(ctx context.Context, s *model.Session) *errors.Error {
	_, err := c.client.PutSession(ctx, s.ToProto(), grpc.WaitForReady(true))
	return errors.New(err)
}
