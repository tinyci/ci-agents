package auth

import (
	"context"
	"io"

	transport "github.com/erikh/go-transport"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc"
)

// Client is a handle into the auth client.
type Client struct {
	ac     auth.AuthClient
	closer io.Closer
}

// NewClient creates a new *Client for use.
func NewClient(addr string, cert *transport.Cert, trace bool) (*Client, error) {
	var (
		closer  io.Closer
		options []grpc.DialOption
		eErr    error
	)

	if trace {
		closer, options, eErr = utils.SetUpGRPCTracing("auth")
		if eErr != nil {
			return nil, eErr
		}
	}

	t, err := transport.GRPCDial(cert, addr, options...)
	if err != nil {
		return nil, err
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
func (c *Client) Capabilities(ctx context.Context) ([]string, error) {
	caps, err := c.ac.Capabilities(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	return caps.List, nil
}
