package processors

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// EnableRepository adds the repository to the CI system. It must already exist.
func (ds *DataServer) EnableRepository(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	user, err := ds.H.Model.FindUserByName(rus.Username)
	if err != nil {
		return nil, err
	}

	repo, err := ds.H.Model.GetRepositoryByNameForUser(rus.RepoName, user)
	if err != nil {
		return nil, err
	}

	if err := ds.H.Model.AddSubscriptionsForUser(user, model.RepositoryList{repo}); err != nil {
		return nil, err
	}

	if err := ds.H.Model.EnableRepository(repo); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// DisableRepository removes the repository from the CI system. It must already be enabled.
func (ds *DataServer) DisableRepository(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	user, err := ds.H.Model.FindUserByName(rus.Username)
	if err != nil {
		return nil, err
	}

	repo, err := ds.H.Model.GetRepositoryByNameForUser(rus.RepoName, user)
	if err != nil {
		return nil, err
	}

	if err := ds.H.Model.DisableRepository(repo); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// SaveRepositories saves the provided user's repositories with the username as the owner.
func (ds *DataServer) SaveRepositories(ctx context.Context, gh *data.GithubJSON) (*empty.Empty, error) {
	repos := []*github.Repository{}

	if err := json.Unmarshal(gh.JSON, &repos); err != nil {
		return nil, err
	}

	if err := ds.H.Model.SaveRepositories(repos, gh.Username, gh.AutoCreated); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// PrivateRepositories returns the list of private repositories the user can see.
func (ds *DataServer) PrivateRepositories(ctx context.Context, name *data.Name) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(name.Name)
	if err != nil {
		return nil, err
	}

	repos, err := ds.H.Model.GetPrivateReposForUser(u)
	if err != nil {
		return nil, err
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// OwnedRepositories returns the list of owned repositories by the user
func (ds *DataServer) OwnedRepositories(ctx context.Context, name *data.Name) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(name.Name)
	if err != nil {
		return nil, err
	}

	repos, err := ds.H.Model.GetOwnedRepos(u)
	if err != nil {
		return nil, err
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// AllRepositories returns the list of repositories the user can see.
func (ds *DataServer) AllRepositories(ctx context.Context, name *data.Name) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(name.Name)
	if err != nil {
		return nil, err
	}

	repos, err := ds.H.Model.GetVisibleReposForUser(u)
	if err != nil {
		return nil, err
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// PublicRepositories returns the list of repositories the user can see.
func (ds *DataServer) PublicRepositories(ctx context.Context, empty *empty.Empty) (*types.RepositoryList, error) {
	repos, err := ds.H.Model.GetAllPublicRepos()
	if err != nil {
		return nil, err
	}

	return model.RepositoryList(repos).ToProto(), nil
}

// GetRepository returns the repository information for the provided name.
func (ds *DataServer) GetRepository(ctx context.Context, name *data.Name) (*types.Repository, error) {
	repo, err := ds.H.Model.GetRepositoryByName(name.Name)
	if err != nil {
		return nil, err
	}

	return repo.ToProto(), nil
}
