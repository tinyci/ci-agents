package auth

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/errors"
)

// OAuthChallenge handles oauth codes and on success, returns the user (Created or patched with latest token)
func (c *Client) OAuthChallenge(ctx context.Context, state, code string) (*auth.OAuthInfo, *errors.Error) {
	userinfo, err := c.ac.OAuthChallenge(ctx, &auth.OAuthChallengeRequest{Code: code, State: state})
	if err != nil {
		return nil, errors.New(err)
	}

	return userinfo, nil
}

// GetOAuthURL retrieves the OAuth redirection URL based on the provided requirements.
func (c *Client) GetOAuthURL(ctx context.Context, scopes []string) (string, *errors.Error) {
	str, err := c.ac.GetOAuthURL(ctx, &auth.Scopes{List: scopes})
	if err != nil {
		return "", errors.New(err)
	}

	return str.Str, nil
}
