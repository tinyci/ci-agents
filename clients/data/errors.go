package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// GetErrors retrieves all the errors for the user.
func (c *Client) GetErrors(name string) ([]*model.UserError, *errors.Error) {
	errs, err := c.client.GetErrors(context.Background(), &data.Name{Name: name})
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
	u, err := c.client.UserByName(context.Background(), &data.Name{Name: username})
	if err != nil {
		return errors.New(err)
	}

	_, err = c.client.AddError(context.Background(), &types.UserError{Error: msg, UserID: u.Id})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteError removes an error.
func (c *Client) DeleteError(id, userID int64) *errors.Error {
	_, err := c.client.DeleteError(context.Background(), &types.UserError{Id: id, UserID: userID})
	return errors.New(err)
}
