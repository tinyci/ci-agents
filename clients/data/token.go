package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// GetToken returns a newly minted access token to tinyCI or error otherwise.
// To get a new token with this method, call the DeleteToken method first if
// one exists already.
func (c *Client) GetToken(ctx context.Context, username string) (string, error) {
	token, err := c.client.GetToken(ctx, &data.Name{Name: username}, grpc.WaitForReady(true))
	if err != nil {
		return "", err
	}

	return token.ID, nil
}

// DeleteToken removes the existing access token and makes it available to be
// regenerated.
func (c *Client) DeleteToken(ctx context.Context, username string) error {
	_, err := c.client.DeleteToken(ctx, &data.Name{Name: username}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	return nil
}

// ValidateToken validates the token and returns error if it is not valid somehow.
func (c *Client) ValidateToken(ctx context.Context, token string) (*model.User, error) {
	user, err := c.client.ValidateToken(ctx, &types.StringID{ID: token}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return model.NewUserFromProto(user)
}
