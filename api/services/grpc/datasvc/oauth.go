package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OAuthRegisterState registers the state code with the datasvc; it will be
// used to validate the other side of the handshake when github redirects back
// to us.
func (ds *DataServer) OAuthRegisterState(ctx context.Context, oas *data.OAuthState) (*empty.Empty, error) {
	if err := ds.H.Model.OAuthRegisterState(ctx, oas.State, oas.Scopes); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// OAuthValidateState registers the state code with the datasvc; it will be
// used to validate the other side of the handshake when github redirects back
// to us.
func (ds *DataServer) OAuthValidateState(ctx context.Context, oas *data.OAuthState) (*data.OAuthState, error) {
	o, err := ds.H.Model.OAuthValidateState(ctx, oas.State)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	oas.Scopes = o

	return oas, err
}
