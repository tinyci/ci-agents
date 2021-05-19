package github

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/github"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	topTypes "github.com/tinyci/ci-agents/types"
	"golang.org/x/oauth2"
)

// RepositoryServer is the external handle for reposvc.
type RepositoryServer struct {
	H *grpcHandler.H
}

func (rs *RepositoryServer) getClientForRepo(ctx context.Context, repoName string) (*github.Client, error) {
	repo, err := rs.H.Clients.Data.GetRepository(ctx, repoName)
	if err != nil {
		return nil, err
	}

	return rs.getClientForUser(ctx, repo.Owner)
}

func (rs *RepositoryServer) getClientForUser(ctx context.Context, u *types.User) (*github.Client, error) {
	var token topTypes.OAuthToken

	if err := json.Unmarshal(u.TokenJSON, &token); err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.Token},
	)

	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}
