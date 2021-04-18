package github

import (
	"context"
	"encoding/base32"
	"strings"

	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	"github.com/tinyci/ci-agents/api/handlers/grpc"
	authconsts "github.com/tinyci/ci-agents/api/services/grpc/auth"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/model"
	topTypes "github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrRedirect indicates that the error intends to redirect the user to the proper spot.
var ErrRedirect = errors.New("redirection")

// AuthServer is the handle/entrypoint for all github-based authsvc implementation.
type AuthServer struct {
	H *grpc.H
}

// Capabilities denotes what capabilities the auth server can manage.
func (as *AuthServer) Capabilities(ctx context.Context, e *empty.Empty) (*auth.StringList, error) {
	return &auth.StringList{List: []string{authconsts.CapabilityOAuth}}, nil
}

// OAuthChallenge is a remote endpoint for performing the final steps of the oauth handshake.
func (as *AuthServer) OAuthChallenge(ctx context.Context, ocr *auth.OAuthChallengeRequest) (*auth.OAuthInfo, error) {
	scopes, eErr := as.H.Clients.Data.OAuthValidateState(ctx, ocr.State)
	if eErr != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "locating state")
	}

	conf := as.H.OAuth.Config(scopes)

	tok, err := conf.Exchange(ctx, ocr.Code)
	if err != nil {
		switch err.(type) {
		case *oauth2.RetrieveError:
			return nil, status.Errorf(codes.FailedPrecondition, "exchanging code for a token")
		default:
			as.H.Clients.Log.Error(ctx, err)
			url, err := as.makeOAuthURL(ctx, scopes)
			if err != nil {
				return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
			}

			return &auth.OAuthInfo{Url: url, Redirect: true}, nil
		}
	}

	client := conf.Client(ctx, tok)
	c := github.NewClient(client)
	u, _, err := c.Users.Get(ctx, "")
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Looking up token user")
	}

	user, eErr := as.H.Clients.Data.GetUser(ctx, u.GetLogin())
	if eErr != nil {
		user = &model.User{
			Username: u.GetLogin(),
		}
	}

	user.Token = &topTypes.OAuthToken{Token: tok.AccessToken, Scopes: scopes, Username: u.GetLogin()}
	if eErr != nil { // same check as above; to determine whether to add or patch
		user, eErr = as.H.Clients.Data.PutUser(ctx, user)
		if eErr != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "Could not create user %v: %v", u.GetLogin(), eErr)
		}
	} else {
		if err := as.H.Clients.Data.PatchUser(ctx, user); err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "Could not patch user %v: %v", u.GetLogin(), eErr)
		}
	}

	return &auth.OAuthInfo{Username: user.Username}, nil
}

func (as *AuthServer) makeOAuthURL(ctx context.Context, scopes []string) (string, error) {
	conf := as.H.OAuth.Config(scopes)
	state := strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)), "=")

	if err := as.H.Clients.Data.OAuthRegisterState(ctx, state, scopes); err != nil {
		return "", utils.WrapError(err, "registering state")
	}

	return conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	), nil
}

// GetOAuthURL returns the url to redirect the user to.
func (as *AuthServer) GetOAuthURL(ctx context.Context, scopes *auth.Scopes) (*auth.String, error) {
	url, err := as.makeOAuthURL(ctx, scopes.List)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &auth.String{Str: url}, nil
}
