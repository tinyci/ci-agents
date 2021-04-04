package datasvc

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EnableRepository adds the repository to the CI system. It must already exist.
func (ds *DataServer) EnableRepository(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	user, err := ds.H.Model.FindUserByName(rus.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	repo, err := ds.H.Model.GetRepositoryByNameForUser(rus.RepoName, user)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.AddSubscriptionsForUser(user, model.RepositoryList{repo}); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.EnableRepository(repo, user); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// DisableRepository removes the repository from the CI system. It must already be enabled.
func (ds *DataServer) DisableRepository(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	user, err := ds.H.Model.FindUserByName(rus.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	repo, err := ds.H.Model.GetRepositoryByNameForUser(rus.RepoName, user)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.DisableRepository(repo); err != nil {
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

	if err := ds.H.Model.SaveRepositories(repos, gh.Username, gh.AutoCreated); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// PrivateRepositories returns the list of private repositories the user can see.
func (ds *DataServer) PrivateRepositories(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(nameSearch.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	repos, err := ds.H.Model.GetPrivateReposForUser(u, nameSearch.Search)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// OwnedRepositories returns the list of owned repositories by the user
func (ds *DataServer) OwnedRepositories(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(nameSearch.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	repos, err := ds.H.Model.GetOwnedRepos(u, nameSearch.Search)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// AllRepositories returns the list of repositories the user can see.
func (ds *DataServer) AllRepositories(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(nameSearch.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	repos, err := ds.H.Model.GetVisibleReposForUser(u, nameSearch.Search)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// PublicRepositories returns the list of repositories the user can see.
func (ds *DataServer) PublicRepositories(ctx context.Context, search *data.Search) (*types.RepositoryList, error) {
	repos, err := ds.H.Model.GetAllPublicRepos(search.Search)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// GetRepository returns the repository information for the provided name.
func (ds *DataServer) GetRepository(ctx context.Context, name *data.Name) (*types.Repository, error) {
	repo, err := ds.H.Model.GetRepositoryByName(name.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return repo.ToProto(), nil
}
