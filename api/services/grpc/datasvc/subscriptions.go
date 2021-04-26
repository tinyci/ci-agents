package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RemoveSubscription removes a subscription from the user's subscriptions.
func (ds *DataServer) RemoveSubscription(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByName(ctx, rus.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	r, err := ds.H.Model.GetRepositoryByNameForUser(ctx, rus.RepoName, u.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.RemoveSubscriptionForUser(ctx, u.ID, []*models.Repository{r}); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// AddSubscription subscribes a repository to a user's account.
func (ds *DataServer) AddSubscription(ctx context.Context, rus *data.RepoUserSelection) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByName(ctx, rus.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	r, err := ds.H.Model.GetRepositoryByNameForUser(ctx, rus.RepoName, u.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.AddSubscriptionsForUser(ctx, u.ID, []*models.Repository{r}); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// ListSubscriptions lists all subscriptions for a user
func (ds *DataServer) ListSubscriptions(ctx context.Context, nameSearch *data.NameSearch) (*types.RepositoryList, error) {
	u, err := ds.H.Model.FindUserByName(ctx, nameSearch.Name)
	if err != nil {
		return nil, err
	}

	subs, err := ds.H.Model.GetSubscriptionsForUser(ctx, u.ID, &nameSearch.Search, 0, 20)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	rl := &types.RepositoryList{}

	for _, sub := range subs {
		r, err := ds.C.ToProto(ctx, sub)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		rl.List = append(rl.List, r.(*types.Repository))
	}

	return rl, nil
}
