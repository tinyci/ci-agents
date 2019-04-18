package processors

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/grpc/services/data"
)

// OAuthRegisterState registers the state code with the datasvc; it will be
// used to validate the other side of the handshake when github redirects back
// to us.
func (ds *DataServer) OAuthRegisterState(ctx context.Context, oas *data.OAuthState) (*empty.Empty, error) {
	if err := ds.H.Model.OAuthRegisterState(oas.State, oas.Scopes); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// OAuthValidateState registers the state code with the datasvc; it will be
// used to validate the other side of the handshake when github redirects back
// to us.
func (ds *DataServer) OAuthValidateState(ctx context.Context, oas *data.OAuthState) (*data.OAuthState, error) {
	o, err := ds.H.Model.OAuthValidateState(oas.State)
	if err != nil {
		return nil, err
	}

	return o.ToProto(), nil
}
