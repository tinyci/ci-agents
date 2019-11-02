package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc/codes"
)

// RemoveSubscription removes a subscription from the user's subscriptions.
func (ds *DataServer) RemoveSubscription(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByName(rus.Username)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	r, err := ds.H.Model.GetRepositoryByNameForUser(rus.RepoName, u)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	if err := ds.H.Model.RemoveSubscriptionForUser(u, r); err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

// AddSubscription subscribes a repository to a user's account.
func (ds *DataServer) AddSubscription(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByName(rus.Username)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	r, err := ds.H.Model.GetRepositoryByNameForUser(rus.RepoName, u)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	if err := ds.H.Model.AddSubscriptionsForUser(u, []*model.Repository{r}); err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

// ListSubscriptions lists all subscriptions for a user
func (ds *DataServer) ListSubscriptions(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByNameWithSubscriptions(nameSearch.Name, nameSearch.Search)
	if err != nil {
		return nil, err.(errors.Error).ToGRPC(codes.FailedPrecondition)
	}

	return model.RepositoryList(u.Subscribed).ToProto(), nil
}
