package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
)

// OAuthValidateState validates the state in the database.
func (c *Client) OAuthValidateState(state string) ([]string, *errors.Error) {
	oas, err := c.client.OAuthValidateState(context.Background(), &data.OAuthState{State: state})
	if err != nil {
		return nil, errors.New(err)
	}

	return oas.Scopes, nil
}

// OAuthRegisterState registers the oauth state in the database.
func (c *Client) OAuthRegisterState(state string, scopes []string) *errors.Error {
	_, err := c.client.OAuthRegisterState(context.Background(), &data.OAuthState{State: state, Scopes: scopes})
	return errors.New(err)
}
