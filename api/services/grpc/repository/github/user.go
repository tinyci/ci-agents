package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MyLogin returns the login username for the token provided.
func (rs *RepositoryServer) MyLogin(ctx context.Context, token *repository.String) (*repository.String, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.Name},
	)
	gh := github.NewClient(oauth2.NewClient(ctx, ts))

	u, _, err := gh.Users.Get(ctx, "")
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "trying to get login username for token: %v", err)
	}

	return &repository.String{Name: u.GetLogin()}, nil
}
