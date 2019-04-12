package processors

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"golang.org/x/oauth2"
)

// UserByName retrieves the user by name and returns it.
func (ds *DataServer) UserByName(ctx context.Context, name *data.Name) (*types.User, error) {
	user, err := ds.H.Model.FindUserByName(name.Name)
	if err != nil {
		return nil, err
	}

	return user.ToProto(), nil
}

// PatchUser loads a user record, overlays the changes and pushes it back to
// the db.
func (ds *DataServer) PatchUser(ctx context.Context, u *types.User) (*empty.Empty, error) {
	origUser, err := ds.H.Model.FindUserByName(u.Username)
	if err != nil {
		return nil, err
	}

	newUser, err := model.NewUserFromProto(u)
	if err != nil {
		return nil, err
	}

	// for now this is the only edit possible. :)
	origUser.Token = newUser.Token
	if err := ds.H.Model.Save(origUser).Error; err != nil {
		return nil, errors.New(err)
	}

	return &empty.Empty{}, nil
}

// PutUser creates a new user.
func (ds *DataServer) PutUser(ctx context.Context, u *types.User) (*types.User, error) {
	ot := &oauth2.Token{}
	if err := json.Unmarshal(u.TokenJSON, ot); err != nil {
		return nil, err
	}

	um, err := ds.H.Model.CreateUser(u.Username, ot)
	if err != nil {
		ds.H.Clients.Log.Errorf("Could not create user %q: %v", u.Username, err)
		return nil, err
	}

	ds.H.Clients.Log.Infof("Created user %q", u.Username)

	return um.ToProto(), nil
}

// ListUsers returns a list of users registered with the system.
func (ds *DataServer) ListUsers(ctx context.Context, e *empty.Empty) (*types.UserList, error) {
	list := []*model.User{}

	if err := ds.H.Model.Find(&list).Error; err != nil {
		return nil, errors.New(err)
	}

	tu := &types.UserList{}

	for _, u := range list {
		tu.Users = append(tu.Users, u.ToProto())
	}

	return tu, nil
}
