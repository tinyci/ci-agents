package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
)

// MyRepositories returns all the writable repositories accessible to user
// owning the access key
func (rs *RepositoryServer) MyRepositories(ctx context.Context, user *types.User) (*repository.RepositoryList, error) {
	var i int
	ret := map[string]*github.Repository{}
	order := []string{}

	gh, err := rs.getClientForUser(ctx, user)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	for {
		repos, _, err := gh.Repositories.List(
			ctx,
			"",
			&github.RepositoryListOptions{
				Visibility: "owner",
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 100,
				},
			},
		)
		if err != nil {
			return nil, errors.New(err)
		}

		for _, repo := range repos {
			if repo.GetPermissions()["admin"] {
				if _, ok := ret[repo.GetFullName()]; !ok {
					ret[repo.GetFullName()] = repo
					order = append(order, repo.GetFullName())
				}
			}
		}

		if len(repos) < 100 {
			break
		}

		i++
	}

	vals := &repository.RepositoryList{}

	for _, value := range order {
		repo := &repository.RepositoryData{
			Name:         ret[value].GetFullName(),
			MasterBranch: ret[value].GetMasterBranch(),
		}

		vals.Repositories = append(vals.Repositories, repo)
	}

	return vals, nil
}

// GetRepository retrieves a repository from github and filters and returns it.
func (rs *RepositoryServer) GetRepository(ctx context.Context, uwn *repository.UserWithRepo) (*repository.RepositoryData, error) {
	owner, repo, err := utils.OwnerRepo(uwn.RepoName)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	gh, err := rs.getClientForUser(ctx, uwn.User)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	r, _, eErr := gh.Repositories.Get(ctx, owner, repo)
	if eErr != nil {
		return nil, errors.New(eErr).Wrapf("Could not fetch repository %v/%v", owner, repo).ToGRPC(codes.FailedPrecondition)
	}

	outRepo := &repository.RepositoryData{
		Name:         r.GetFullName(),
		MasterBranch: r.GetMasterBranch(),
	}

	return outRepo, nil
}
