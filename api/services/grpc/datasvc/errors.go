package datasvc

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetErrors retrieves the errors for the provided user.
func (ds *DataServer) GetErrors(ctx context.Context, name *data.Name) (*types.UserErrors, error) {
	u, err := ds.H.Model.FindUserByName(ctx, name.Name)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	errs, err := ds.H.Model.GetErrors(ctx, u)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	errors := &types.UserErrors{}

	for _, err := range errs {
		e, err := ds.C.ToProto(ctx, err)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		errors.Errors = append(errors.Errors, e.(*types.UserError))
	}

	return errors, nil
}

// AddError adds an error to the user's error stack.
func (ds *DataServer) AddError(ctx context.Context, ue *types.UserError) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByID(ctx, ue.UserID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.AddError(ctx, u, errors.New(ue.Error)); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &empty.Empty{}, nil
}

// DeleteError removes an error from errors list. The error string does not need to be provided.
func (ds *DataServer) DeleteError(ctx context.Context, ue *types.UserError) (*empty.Empty, error) {
	u, err := ds.H.Model.FindUserByID(ctx, ue.UserID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.DeleteError(ctx, u, ue.Id); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}
