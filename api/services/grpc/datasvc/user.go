package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	topTypes "github.com/tinyci/ci-agents/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ds *DataServer) toTypesUser(ctx context.Context, u *models.User) (*types.User, error) {
	user, err := ds.C.ToProto(ctx, u)
	if err != nil {
		return nil, err
	}

	return user.(*types.User), nil
}

func (ds *DataServer) fromTypesUser(ctx context.Context, u *types.User) (*models.User, error) {
	user, err := ds.C.FromProto(ctx, u)
	if err != nil {
		return nil, err
	}

	return user.(*models.User), nil
}

// UserByName retrieves the user by name and returns it.
func (ds *DataServer) UserByName(ctx context.Context, name *data.Name) (*types.User, error) {
	user, err := ds.H.Model.FindUserByName(ctx, name.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return ds.toTypesUser(ctx, user)
}

// PatchUser loads a user record, overlays the changes and pushes it back to
// the db.
func (ds *DataServer) PatchUser(ctx context.Context, u *types.User) (*empty.Empty, error) {
	origUser, err := ds.H.Model.FindUserByName(ctx, u.Username)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	newUser, err := ds.fromTypesUser(ctx, u)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	// for now this is the only edit possible. :)
	origUser.Token = newUser.Token
	if err := ds.H.Model.UpdateUser(ctx, origUser); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// PutUser creates a new user.
func (ds *DataServer) PutUser(ctx context.Context, u *types.User) (*types.User, error) {
	um, err := ds.H.Model.CreateUser(ctx, u.Username, u.TokenJSON)
	if err != nil {
		ds.H.Clients.Log.Errorf(ctx, "Could not create user %q: %v", u.Username, err)
		return &types.User{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	ds.H.Clients.Log.Infof(ctx, "Created user %q", u.Username)

	return ds.toTypesUser(ctx, um)
}

// ListUsers returns a list of users registered with the system.
func (ds *DataServer) ListUsers(ctx context.Context, e *empty.Empty) (*types.UserList, error) {
	list, err := ds.H.Model.ListUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	tu := &types.UserList{}

	for _, u := range list {
		user, err := ds.toTypesUser(ctx, u)
		if err != nil {
			return nil, err
		}

		tu.Users = append(tu.Users, user)
	}

	return tu, nil
}

// HasCapability returns true if the capability requested exists for the user provided.
func (ds *DataServer) HasCapability(ctx context.Context, cr *data.CapabilityRequest) (*types.Bool, error) {
	u, err := ds.H.Model.FindUserByID(ctx, cr.Id)
	if err != nil {
		return &types.Bool{Result: false}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	res, err := ds.H.Model.HasCapability(ctx, u, topTypes.Capability(cr.Capability), ds.H.Auth.FixedCapabilities)
	if err != nil {
		return &types.Bool{Result: false}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &types.Bool{Result: res}, nil
}

// AddCapability adds a capability for a user.
func (ds *DataServer) AddCapability(ctx context.Context, cr *data.CapabilityRequest) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByID(ctx, cr.Id)
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.AddCapabilityToUser(ctx, u, topTypes.Capability(cr.Capability)); err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// RemoveCapability removes a capability from a user.
func (ds *DataServer) RemoveCapability(ctx context.Context, cr *data.CapabilityRequest) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByID(ctx, cr.Id)
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.RemoveCapabilityFromUser(ctx, u, topTypes.Capability(cr.Capability)); err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// GetCapabilities retrieves the capabilities for a user.
func (ds *DataServer) GetCapabilities(ctx context.Context, u *types.User) (*data.Capabilities, error) {
	mu, err := ds.fromTypesUser(ctx, u)
	if err != nil {
		return &data.Capabilities{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	caps, err := ds.H.Model.GetCapabilities(ctx, mu, ds.H.Auth.FixedCapabilities)
	if err != nil {
		return &data.Capabilities{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	strCaps := []string{}
	for _, cap := range caps {
		strCaps = append(strCaps, string(cap))
	}

	return &data.Capabilities{Capabilities: strCaps}, nil
}
