package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// GetErrors retrieves all the errors for the user.
func (c *Client) GetErrors(name string) ([]*model.UserError, *errors.Error) {
	errs, err := c.client.GetErrors(context.Background(), &data.Name{Name: name}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	errList := []*model.UserError{}

	for _, e := range errs.Errors {
		errList = append(errList, model.NewUserErrorFromProto(e))
	}

	return errList, nil
}

// AddError adds an error.
func (c *Client) AddError(msg, username string) *errors.Error {
	u, err := c.client.UserByName(context.Background(), &data.Name{Name: username}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	_, err = c.client.AddError(context.Background(), &types.UserError{Error: msg, UserID: u.Id}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteError removes an error.
func (c *Client) DeleteError(id, userID int64) *errors.Error {
	_, err := c.client.DeleteError(context.Background(), &types.UserError{Id: id, UserID: userID}, grpc.WaitForReady(true))
	return errors.New(err)
}
