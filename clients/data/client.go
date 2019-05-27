package data

import (
	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/errors"
)

// Client is a datasvc client.
type Client struct {
	client data.DataClient
}

// New creates a new *Client.
func New(addr string, cert *transport.Cert) (*Client, *errors.Error) {
	c, err := transport.GRPCDial(cert, addr)
	if err != nil {
		return nil, errors.New(err)
	}

	return &Client{client: data.NewDataClient(c)}, nil
}
