package data

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// PatchUser adjusts the token for the user.
func (c *Client) PatchUser(ctx context.Context, u *model.User) error {
	_, err := c.client.PatchUser(ctx, u.ToProto(), grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// PutUser inserts the user provided.
func (c *Client) PutUser(ctx context.Context, u *model.User) (*model.User, error) {
	u2, err := c.client.PutUser(ctx, u.ToProto(), grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(u2)
}

// GetUser obtains a user record by name
func (c *Client) GetUser(ctx context.Context, name string) (*model.User, error) {
	u, err := c.client.UserByName(ctx, &data.Name{Name: name}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(u)
}

// ListUsers lists the users in the system.
func (c *Client) ListUsers(ctx context.Context) ([]*model.User, error) {
	users, err := c.client.ListUsers(ctx, &empty.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	u := []*model.User{}

	for _, user := range users.Users {
		u2, err := model.NewUserFromProto(user)
		if err != nil {
			return nil, err
		}

		u = append(u, u2)
	}

	return u, nil
}

// GetCapabilities yields the capabilities that belong to the user.
func (c *Client) GetCapabilities(ctx context.Context, u *model.User) ([]model.Capability, error) {
	caps, err := c.client.GetCapabilities(ctx, u.ToProto())
	if err != nil {
		return nil, errors.New(err)
	}

	realCaps := []model.Capability{}
	for _, cap := range caps.Capabilities {
		realCaps = append(realCaps, model.Capability(cap))
	}

	return realCaps, nil
}

// HasCapability returns true if the user has the specified capability.
func (c *Client) HasCapability(ctx context.Context, u *model.User, cap model.Capability) (bool, error) {
	res, err := c.client.HasCapability(ctx, &data.CapabilityRequest{Id: u.ID, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return false, errors.New(err)
	}

	return res.Result, nil
}

// AddCapability adds a capability for a user.
func (c *Client) AddCapability(ctx context.Context, u *model.User, cap model.Capability) error {
	_, err := c.client.AddCapability(ctx, &data.CapabilityRequest{Id: u.ID, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// RemoveCapability removes a capability from a user.
func (c *Client) RemoveCapability(ctx context.Context, u *model.User, cap model.Capability) error {
	_, err := c.client.RemoveCapability(ctx, &data.CapabilityRequest{Id: u.ID, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}
