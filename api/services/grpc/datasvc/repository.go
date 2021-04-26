package datasvc

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ds *DataServer) toTypesRepository(ctx context.Context, r *models.Repository) (*types.Repository, error) {
	repo, err := ds.C.ToProto(ctx, r)
	if err != nil {
		return nil, err
	}

	return repo.(*types.Repository), nil
}

func (ds *DataServer) toTypesRepositoryList(ctx context.Context, repos []*models.Repository) (*types.RepositoryList, error) {
	list := &types.RepositoryList{}

	for _, repo := range repos {
		r, err := ds.C.ToProto(ctx, repo)
		if err != nil {
			return nil, err
		}

		list.List = append(list.List, r.(*types.Repository))
	}

	return list, nil
}

// EnableRepository adds the repository to the CI system. It must already exist.
func (ds *DataServer) EnableRepository(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	user, err := ds.H.Model.FindUserByName(ctx, rus.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	repo, err := ds.H.Model.GetRepositoryByNameForUser(ctx, rus.RepoName, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.AddSubscriptionsForUser(ctx, user.ID, []*models.Repository{repo}); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.EnableRepository(ctx, repo, user.ID); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// DisableRepository removes the repository from the CI system. It must already be enabled.
func (ds *DataServer) DisableRepository(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	user, err := ds.H.Model.FindUserByName(ctx, rus.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	repo, err := ds.H.Model.GetRepositoryByNameForUser(ctx, rus.RepoName, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.DisableRepository(ctx, repo); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// SaveRepositories saves the provided user's repositories with the username as the owner.
func (ds *DataServer) SaveRepositories(ctx context.Context, gh *data.GithubJSON) (*empty.Empty, error) {
	repos := []*github.Repository{}

	if err := json.Unmarshal(gh.JSON, &repos); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.SaveRepositories(ctx, repos, gh.Username, gh.AutoCreated); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// PrivateRepositories returns the list of private repositories the user can see.
func (ds *DataServer) PrivateRepositories(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(ctx, nameSearch.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	var search *string
	if nameSearch.Search != "" {
		search = &nameSearch.Search
	}
	repos, err := ds.H.Model.GetPrivateReposForUser(ctx, u.ID, search)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return ds.toTypesRepositoryList(ctx, repos)
}

// OwnedRepositories returns the list of owned repositories by the user
func (ds *DataServer) OwnedRepositories(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(ctx, nameSearch.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	var search *string
	if nameSearch.Search != "" {
		search = &nameSearch.Search
	}

	repos, err := ds.H.Model.GetOwnedRepos(ctx, u.ID, search)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return ds.toTypesRepositoryList(ctx, repos)
}

// AllRepositories returns the list of repositories the user can see.
func (ds *DataServer) AllRepositories(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(ctx, nameSearch.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	var search *string
	if nameSearch.Search != "" {
		search = &nameSearch.Search
	}

	repos, err := ds.H.Model.GetVisibleReposForUser(ctx, u.ID, search)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return ds.toTypesRepositoryList(ctx, repos)
}

// PublicRepositories returns the list of repositories the user can see.
func (ds *DataServer) PublicRepositories(ctx context.Context, search *data.Search) (*types.RepositoryList, error) {
	var s *string
	if search.Search != "" {
		s = &search.Search
	}

	repos, err := ds.H.Model.GetAllPublicRepos(ctx, s)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return ds.toTypesRepositoryList(ctx, repos)
}

// GetRepository returns the repository information for the provided name.
func (ds *DataServer) GetRepository(ctx context.Context, name *data.Name) (*types.Repository, error) {
	repo, err := ds.H.Model.GetRepositoryByName(ctx, name.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return ds.toTypesRepository(ctx, repo)
}
