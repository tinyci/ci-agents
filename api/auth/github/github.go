package github

import (
	"context"
	"encoding/base32"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
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

// OAuthChallenge is a remote endpoint for performing the final steps of the oauth handshake.
func (as *AuthServer) OAuthChallenge(ctx context.Context, ocr *auth.OAuthChallengeRequest) (*types.User, error) {
	conf := as.H.OAuth.Config(ocr.Scopes)

	tok, err := conf.Exchange(ctx, ocr.Code)
	if err != nil {
		switch err.(type) {
		case *oauth2.RetrieveError:
			return nil, errors.New(err).ToGRPC(codes.FailedPrecondition)
		default:
			as.H.Clients.Log.Error(ctx, err)
			return nil, ErrRedirect.ToGRPC(codes.FailedPrecondition)
		}
	}

	client := conf.Client(ctx, tok)
	c := github.NewClient(client)
	u, _, err := c.Users.Get(ctx, "")
	if err != nil {
		return nil, errors.New(err).ToGRPC(codes.FailedPrecondition)
	}

	user, eErr := as.H.Clients.Data.GetUser(u.GetLogin())
	if eErr != nil {
		user = &model.User{
			Username: u.GetLogin(),
		}
	}

	user.Token = &topTypes.OAuthToken{Token: tok.AccessToken, Scopes: ocr.Scopes, Username: u.GetLogin()}
	if eErr != nil { // same check as above; to determine whether to add or patch
		user, eErr = as.H.Clients.Data.PutUser(user)
		if eErr != nil {
			return nil, eErr.Wrapf("Could not create user %v", u.GetLogin()).ToGRPC(codes.FailedPrecondition)
		}
	} else {
		if err := as.H.Clients.Data.PatchUser(user); err != nil {
			return nil, eErr.Wrapf("Could not patch user %v", u.GetLogin()).ToGRPC(codes.FailedPrecondition)
		}
	}

	return user.ToProto(), nil
}

// GetOAuthURL returns the url to redirect the user to.
func (as *AuthServer) GetOAuthURL(ctx context.Context, scopes *auth.Scopes) (*auth.String, error) {
	conf := as.H.OAuth.Config(scopes.List)

	state := strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)), "=")
	if err := as.H.Clients.Data.OAuthRegisterState(state, scopes.List); err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	return &auth.String{Str: conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)}, nil
}
