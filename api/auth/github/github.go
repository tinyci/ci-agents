package github

import (
	"context"
	"encoding/base32"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	authconsts "github.com/tinyci/ci-agents/api/auth"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	topTypes "github.com/tinyci/ci-agents/types"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
)

// ErrRedirect indicates that the error intends to redirect the user to the proper spot.
var ErrRedirect = errors.New("redirection")

// AuthServer is the handle/entrypoint for all github-based authsvc implementation.
type AuthServer struct {
	H *handler.H
}

// Capabilities denotes what capabilities the auth server can manage.
func (as *AuthServer) Capabilities(ctx context.Context, e *empty.Empty) (*auth.StringList, error) {
	return &auth.StringList{List: []string{authconsts.CapabilityOAuth}}, nil
}

// OAuthChallenge is a remote endpoint for performing the final steps of the oauth handshake.
func (as *AuthServer) OAuthChallenge(ctx context.Context, ocr *auth.OAuthChallengeRequest) (*auth.OAuthInfo, error) {
	scopes, eErr := as.H.Clients.Data.OAuthValidateState(ctx, ocr.State)
	if eErr != nil {
		return nil, eErr.Wrap("Locating state").ToGRPC(codes.FailedPrecondition)
	}

	conf := as.H.OAuth.Config(scopes)

	tok, err := conf.Exchange(ctx, ocr.Code)
	if err != nil {
		switch err.(type) {
		case *oauth2.RetrieveError:
			return nil, errors.New(err).Wrap("exchanging code for a token").ToGRPC(codes.FailedPrecondition)
		default:
			as.H.Clients.Log.Error(ctx, err)
			url, err := as.makeOAuthURL(ctx, scopes)
			if err != nil {
				return nil, err.ToGRPC(codes.FailedPrecondition)
			}

			return &auth.OAuthInfo{Url: url, Redirect: true}, nil
		}
	}

	client := conf.Client(ctx, tok)
	c := github.NewClient(client)
	u, _, err := c.Users.Get(ctx, "")
	if err != nil {
		return nil, errors.New(err).Wrap("Looking up token user").ToGRPC(codes.FailedPrecondition)
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
			return nil, eErr.Wrapf("Could not create user %v", u.GetLogin()).ToGRPC(codes.FailedPrecondition)
		}
	} else {
		if err := as.H.Clients.Data.PatchUser(ctx, user); err != nil {
			return nil, eErr.Wrapf("Could not patch user %v", u.GetLogin()).ToGRPC(codes.FailedPrecondition)
		}
	}

	return &auth.OAuthInfo{Username: user.Username}, nil
}

func (as *AuthServer) makeOAuthURL(ctx context.Context, scopes []string) (string, *errors.Error) {
	conf := as.H.OAuth.Config(scopes)
	state := strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)), "=")

	if err := as.H.Clients.Data.OAuthRegisterState(ctx, state, scopes); err != nil {
		return "", err.Wrap("registering state")
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
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	return &auth.String{Str: url}, nil
}
