package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
)

// MyLogin returns the login username for the token provided.
func (rs *RepositoryServer) MyLogin(ctx context.Context, token *repository.String) (*repository.String, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.Name},
	)
	gh := github.NewClient(oauth2.NewClient(ctx, ts))

	u, _, err := gh.Users.Get(ctx, "")
	if err != nil {
		return nil, errors.New(err).(errors.Error).Wrap("trying to get login username for token").ToGRPC(codes.FailedPrecondition)
	}

	return &repository.String{Name: u.GetLogin()}, nil
}
