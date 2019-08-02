package repository

import (
	"io"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc"
)

// Client is a repository client.
type Client struct {
	client repository.RepositoryClient
	closer io.Closer
}

// New creates a new *Client.
func New(addr string, cert *transport.Cert, trace bool) (*Client, *errors.Error) {
	var (
		closer  io.Closer
		options []grpc.DialOption
		eErr    *errors.Error
	)

	if trace {
		closer, options, eErr = utils.SetUpGRPCTracing("repository")
		if eErr != nil {
			return nil, eErr
		}
	}

	c, err := transport.GRPCDial(cert, addr, options...)
	if err != nil {
		return nil, errors.New(err)
	}

	return &Client{closer: closer, client: repository.NewRepositoryClient(c)}, nil
}

// Close closes the client's tracing functionality
func (c *Client) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}

	return nil
}
