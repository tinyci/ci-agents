package github

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetupHook sets up the pr webhook in github.
func (rs *RepositoryServer) SetupHook(ctx context.Context, hsr *repository.HookSetupRequest) (*empty.Empty, error) {
	owner, repo, err := utils.OwnerRepo(hsr.RepoName)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	gh, err := rs.getClientForRepo(ctx, hsr.RepoName)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	_, _, eErr := gh.Repositories.CreateHook(ctx, owner, repo, &github.Hook{
		URL:    github.String(hsr.HookURL),
		Events: []string{"push", "pull_request"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          hsr.HookURL,
			"content_type": "json",
			"secret":       hsr.HookSecret,
		},
	})

	if eErr != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "configuring hook on repo %v/%v: %v", owner, repo, err)
	}

	return &empty.Empty{}, nil
}

// TeardownHook removes the pr webhook in github.
func (rs *RepositoryServer) TeardownHook(ctx context.Context, htr *repository.HookTeardownRequest) (*empty.Empty, error) {
	owner, repo, err := utils.OwnerRepo(htr.RepoName)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	gh, err := rs.getClientForRepo(ctx, htr.RepoName)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	var id int64
	var i int

	for {
		hooks, _, err := gh.Repositories.ListHooks(ctx, owner, repo, &github.ListOptions{Page: i, PerPage: 20})
		if err != nil || len(hooks) == 0 {
			break
		}
		for _, hook := range hooks {
			if hook.Config["url"] == htr.HookURL {
				id = hook.GetID()
				goto finish
			}
		}
		i++
	}

finish:
	if id != 0 {
		_, err := gh.Repositories.DeleteHook(context.Background(), owner, repo, id)
		if err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
		}
	}

	return &empty.Empty{}, nil
}
