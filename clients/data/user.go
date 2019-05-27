package data

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// PatchUser adjusts the token for the user.
func (c *Client) PatchUser(u *model.User) *errors.Error {
	_, err := c.client.PatchUser(context.Background(), u.ToProto(), grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// PutUser inserts the user provided.
func (c *Client) PutUser(u *model.User) (*model.User, *errors.Error) {
	u2, err := c.client.PutUser(context.Background(), u.ToProto(), grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(u2)
}

// GetUser obtains a user record by name
func (c *Client) GetUser(name string) (*model.User, *errors.Error) {
	u, err := c.client.UserByName(context.Background(), &data.Name{Name: name}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(u)
}

// ListUsers lists the users in the system.
func (c *Client) ListUsers() ([]*model.User, *errors.Error) {
	users, err := c.client.ListUsers(context.Background(), &empty.Empty{}, grpc.WaitForReady(true))
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

// HasCapability returns true if the user has the specified capability.
func (c *Client) HasCapability(u *model.User, cap model.Capability) (bool, *errors.Error) {
	res, err := c.client.HasCapability(context.Background(), &data.CapabilityRequest{Id: u.ID, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return false, errors.New(err)
	}

	return res.Result, nil
}

// AddCapability adds a capability for a user.
func (c *Client) AddCapability(u *model.User, cap model.Capability) *errors.Error {
	_, err := c.client.AddCapability(context.Background(), &data.CapabilityRequest{Id: u.ID, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// RemoveCapability removes a capability from a user.
func (c *Client) RemoveCapability(u *model.User, cap model.Capability) *errors.Error {
	_, err := c.client.RemoveCapability(context.Background(), &data.CapabilityRequest{Id: u.ID, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}
