package datasvc

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
	topTypes "github.com/tinyci/ci-agents/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserByName retrieves the user by name and returns it.
func (ds *DataServer) UserByName(ctx context.Context, name *data.Name) (*types.User, error) {
	user, err := ds.H.Model.FindUserByName(name.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return user.ToProto(), nil
}

// PatchUser loads a user record, overlays the changes and pushes it back to
// the db.
func (ds *DataServer) PatchUser(ctx context.Context, u *types.User) (*empty.Empty, error) {
	origUser, err := ds.H.Model.FindUserByName(u.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	newUser, err := model.NewUserFromProto(u)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	// for now this is the only edit possible. :)
	origUser.Token = newUser.Token
	if err := ds.H.Model.Save(origUser).Error; err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// PutUser creates a new user.
func (ds *DataServer) PutUser(ctx context.Context, u *types.User) (*types.User, error) {
	ot := &topTypes.OAuthToken{}
	if err := json.Unmarshal(u.TokenJSON, ot); err != nil {
		return &types.User{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	um, err := ds.H.Model.CreateUser(u.Username, ot)
	if err != nil {
		ds.H.Clients.Log.Errorf(ctx, "Could not create user %q: %v", u.Username, err)
		return &types.User{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	ds.H.Clients.Log.Infof(ctx, "Created user %q", u.Username)

	return um.ToProto(), nil
}

// ListUsers returns a list of users registered with the system.
func (ds *DataServer) ListUsers(ctx context.Context, e *empty.Empty) (*types.UserList, error) {
	list := []*model.User{}

	if err := ds.H.Model.Find(&list).Error; err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	tu := &types.UserList{}

	for _, u := range list {
		tu.Users = append(tu.Users, u.ToProto())
	}

	return tu, nil
}

// HasCapability returns true if the capability requested exists for the user provided.
func (ds *DataServer) HasCapability(ctx context.Context, cr *data.CapabilityRequest) (*types.Bool, error) {
	u, err := ds.H.Model.FindUserByID(cr.Id)
	if err != nil {
		return &types.Bool{Result: false}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	res, err := ds.H.Model.HasCapability(u, model.Capability(cr.Capability), ds.H.Auth.FixedCapabilities)
	if err != nil {
		return &types.Bool{Result: false}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &types.Bool{Result: res}, nil
}

// AddCapability adds a capability for a user.
func (ds *DataServer) AddCapability(ctx context.Context, cr *data.CapabilityRequest) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByID(cr.Id)
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.AddCapabilityToUser(u, model.Capability(cr.Capability)); err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// RemoveCapability removes a capability from a user.
func (ds *DataServer) RemoveCapability(ctx context.Context, cr *data.CapabilityRequest) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByID(cr.Id)
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.RemoveCapabilityFromUser(u, model.Capability(cr.Capability)); err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// GetCapabilities retrieves the capabilities for a user.
func (ds *DataServer) GetCapabilities(ctx context.Context, u *types.User) (*data.Capabilities, error) {
	mu, err := model.NewUserFromProto(u)
	if err != nil {
		return &data.Capabilities{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	caps, err := ds.H.Model.GetCapabilities(mu, ds.H.Auth.FixedCapabilities)
	if err != nil {
		return &data.Capabilities{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	strCaps := []string{}
	for _, cap := range caps {
		strCaps = append(strCaps, string(cap))
	}

	return &data.Capabilities{Capabilities: strCaps}, nil
}
