package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// GetErrors retrieves all the errors for the user.
func (c *Client) GetErrors(ctx context.Context, name string) ([]*model.UserError, error) {
	errs, err := c.client.GetErrors(ctx, &data.Name{Name: name}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}

	errList := []*model.UserError{}

	for _, e := range errs.Errors {
		errList = append(errList, model.NewUserErrorFromProto(e))
	}

	return errList, nil
}

// AddError adds an error.
func (c *Client) AddError(ctx context.Context, msg, username string) error {
	u, err := c.client.UserByName(ctx, &data.Name{Name: username}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	_, err = c.client.AddError(ctx, &types.UserError{Error: msg, UserID: u.Id}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}

	return nil
}

// DeleteError removes an error.
func (c *Client) DeleteError(ctx context.Context, id, userID int64) error {
	_, err := c.client.DeleteError(ctx, &types.UserError{Id: id, UserID: userID}, grpc.WaitForReady(true))
	return err
}
