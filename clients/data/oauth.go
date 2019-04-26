package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"google.golang.org/grpc"
)

// OAuthValidateState validates the state in the database.
func (c *Client) OAuthValidateState(state string) ([]string, *errors.Error) {
	oas, err := c.client.OAuthValidateState(context.Background(), &data.OAuthState{State: state}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return oas.Scopes, nil
}

// OAuthRegisterState registers the oauth state in the database.
func (c *Client) OAuthRegisterState(state string, scopes []string) *errors.Error {
	_, err := c.client.OAuthRegisterState(context.Background(), &data.OAuthState{State: state, Scopes: scopes}, grpc.WaitForReady(true))
	return errors.New(err)
}
