package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"google.golang.org/grpc"
)

// GetSession retrieves a session from the database by id.
func (c *Client) GetSession(ctx context.Context, id string) (*types.Session, error) {
	s, err := c.client.LoadSession(ctx, &types.StringID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return s, nil
}

// PutSession adds a session to the database.
func (c *Client) PutSession(ctx context.Context, s *types.Session) error {
	_, err := c.client.PutSession(ctx, s, grpc.WaitForReady(true))
	return err
}
