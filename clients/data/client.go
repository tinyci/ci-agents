package data

import (
	"io"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc"
)

// Client is a datasvc client.
type Client struct {
	client data.DataClient
	closer io.Closer
}

// New creates a new *Client.
func New(addr string, cert *transport.Cert, trace bool) (*Client, error) {
	var (
		closer  io.Closer
		options []grpc.DialOption
		eErr    error
	)

	if trace {
		closer, options, eErr = utils.SetUpGRPCTracing("data")
		if eErr != nil {
			return nil, eErr
		}
	}

	c, err := transport.GRPCDial(cert, addr, options...)
	if err != nil {
		return nil, err
	}

	return &Client{closer: closer, client: data.NewDataClient(c)}, nil
}

// Close closes the client's tracing functionality
func (c *Client) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}

	return nil
}
