package data

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/model"
)

// PatchUser adjusts the token for the user.
func (c *Client) PatchUser(u *model.User) *errors.Error {
	_, err := c.client.PatchUser(context.Background(), u.ToProto())
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// PutUser inserts the user provided.
func (c *Client) PutUser(u *model.User) (*model.User, *errors.Error) {
	u2, err := c.client.PutUser(context.Background(), u.ToProto())
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(u2)
}

// GetUser obtains a user record by name
func (c *Client) GetUser(name string) (*model.User, *errors.Error) {
	u, err := c.client.UserByName(context.Background(), &data.Name{Name: name})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewUserFromProto(u)
}

// ListUsers lists the users in the system.
func (c *Client) ListUsers() ([]*model.User, *errors.Error) {
	users, err := c.client.ListUsers(context.Background(), &empty.Empty{})
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
