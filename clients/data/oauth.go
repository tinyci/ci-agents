package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"google.golang.org/grpc"
)

// OAuthValidateState validates the state in the database.
func (c *Client) OAuthValidateState(ctx context.Context, state string) ([]string, error) {
	oas, err := c.client.OAuthValidateState(ctx, &data.OAuthState{State: state}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	return oas.Scopes, nil
}

// OAuthRegisterState registers the oauth state in the database.
func (c *Client) OAuthRegisterState(ctx context.Context, state string, scopes []string) error {
	_, err := c.client.OAuthRegisterState(ctx, &data.OAuthState{State: state, Scopes: scopes}, grpc.WaitForReady(true))
	return err
}
