package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"golang.org/x/oauth2"
)

// RepositoryServer is the external handle for reposvc.
type RepositoryServer struct {
	H *handler.H
}

func (rs *RepositoryServer) getClientForRepo(ctx context.Context, repoName string) (*github.Client, error) {
	repo, err := rs.H.Clients.Data.GetRepository(ctx, repoName)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: repo.Owner.Token.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}

func (rs *RepositoryServer) getClientForUser(ctx context.Context, u *types.User) (*github.Client, error) {
	user, err := model.NewUserFromProto(u)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.Token.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}
