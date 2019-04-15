package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// GetToken returns a newly minted access token to tinyCI or error otherwise.
// To get a new token with this method, call the DeleteToken method first if
// one exists already.
func (c *Client) GetToken(username string) (string, *errors.Error) {
	token, err := c.client.GetToken(context.Background(), &data.Name{Name: username})
	if err != nil {
		return "", errors.New(err)
	}

	return token.ID, nil
}

// DeleteToken removes the existing access token and makes it available to be
// regenerated.
func (c *Client) DeleteToken(username string) *errors.Error {
	_, err := c.client.DeleteToken(context.Background(), &data.Name{Name: username})
	if err != nil {
		return errors.New(err)
	}
	return nil
}

// ValidateToken validates the token and returns error if it is not valid somehow.
func (c *Client) ValidateToken(token string) (*model.User, *errors.Error) {
	user, err := c.client.ValidateToken(context.Background(), &types.StringID{ID: token})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(user)
}
