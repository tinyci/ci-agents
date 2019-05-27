package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// GetToken returns a newly minted access token to tinyCI or error otherwise.
// To get a new token with this method, call the DeleteToken method first if
// one exists already.
func (c *Client) GetToken(username string) (string, *errors.Error) {
	token, err := c.client.GetToken(context.Background(), &data.Name{Name: username}, grpc.WaitForReady(true))
	if err != nil {
		return "", errors.New(err)
	}

	return token.ID, nil
}

// DeleteToken removes the existing access token and makes it available to be
// regenerated.
func (c *Client) DeleteToken(username string) *errors.Error {
	_, err := c.client.DeleteToken(context.Background(), &data.Name{Name: username}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}
	return nil
}

// ValidateToken validates the token and returns error if it is not valid somehow.
func (c *Client) ValidateToken(token string) (*model.User, *errors.Error) {
	user, err := c.client.ValidateToken(context.Background(), &types.StringID{ID: token}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(user)
}
