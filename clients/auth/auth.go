package auth

import (
	"context"
	"io"

	transport "github.com/erikh/go-transport"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc"
)

// Client is a handle into the auth client.
type Client struct {
	ac     auth.AuthClient
	closer io.Closer
}

// NewClient creates a new *Client for use.
func NewClient(addr string, cert *transport.Cert, trace bool) (*Client, *errors.Error) {
	var (
		closer  io.Closer
		options []grpc.DialOption
		eErr    *errors.Error
	)

	if trace {
		closer, options, eErr = utils.SetUpGRPCTracing("auth")
		if eErr != nil {
			return nil, eErr
		}
	}

	t, err := transport.GRPCDial(cert, addr, options...)
	if err != nil {
		return nil, errors.New(err)
	}

	return &Client{closer: closer, ac: auth.NewAuthClient(t)}, nil
}

// Close closes the client's tracing functionality
func (c *Client) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}

	return nil
}

// Capabilities notes what types of auth this server supports.
func (c *Client) Capabilities() ([]string, *errors.Error) {
	caps, err := c.ac.Capabilities(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, errors.New(err)
	}

	return caps.List, nil
}
