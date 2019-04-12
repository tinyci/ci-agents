package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
)

// OAuthValidateState validates the state in the database.
func (c *Client) OAuthValidateState(state string) *errors.Error {
	_, err := c.client.OAuthValidateState(context.Background(), &data.OAuthState{State: state})
	return errors.New(err)
}

// OAuthRegisterState registers the oauth state in the database.
func (c *Client) OAuthRegisterState(state string) *errors.Error {
	_, err := c.client.OAuthRegisterState(context.Background(), &data.OAuthState{State: state})
	return errors.New(err)
}
