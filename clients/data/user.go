package data

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	topTypes "github.com/tinyci/ci-agents/types"
	"google.golang.org/grpc"
)

// PatchUser adjusts the token for the user.
func (c *Client) PatchUser(ctx context.Context, u *types.User) error {
	_, err := c.client.PatchUser(ctx, u, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	return nil
}

// PutUser inserts the user provided.
func (c *Client) PutUser(ctx context.Context, u *types.User) (*types.User, error) {
	return c.client.PutUser(ctx, u, grpc.WaitForReady(true))
}

// GetUser obtains a user record by name
func (c *Client) GetUser(ctx context.Context, name string) (*types.User, error) {
	return c.client.UserByName(ctx, &data.Name{Name: name}, grpc.WaitForReady(true))
}

// ListUsers lists the users in the system.
func (c *Client) ListUsers(ctx context.Context) (*types.UserList, error) {
	return c.client.ListUsers(ctx, &empty.Empty{}, grpc.WaitForReady(true))
}

// GetCapabilities yields the capabilities that belong to the user.
func (c *Client) GetCapabilities(ctx context.Context, u *types.User) ([]topTypes.Capability, error) {
	caps, err := c.client.GetCapabilities(ctx, u)
	if err != nil {
		return nil, err
	}

	realCaps := []topTypes.Capability{}
	for _, cap := range caps.Capabilities {
		realCaps = append(realCaps, topTypes.Capability(cap))
	}

	return realCaps, nil
}

// HasCapability returns true if the user has the specified capability.
func (c *Client) HasCapability(ctx context.Context, u *types.User, cap topTypes.Capability) (bool, error) {
	res, err := c.client.HasCapability(ctx, &data.CapabilityRequest{Id: u.Id, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return false, err
	}

	return res.Result, nil
}

// AddCapability adds a capability for a user.
func (c *Client) AddCapability(ctx context.Context, u *types.User, cap topTypes.Capability) error {
	_, err := c.client.AddCapability(ctx, &data.CapabilityRequest{Id: u.Id, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	return nil
}

// RemoveCapability removes a capability from a user.
func (c *Client) RemoveCapability(ctx context.Context, u *types.User, cap topTypes.Capability) error {
	_, err := c.client.RemoveCapability(ctx, &data.CapabilityRequest{Id: u.Id, Capability: string(cap)}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	return nil
}
