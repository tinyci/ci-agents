package queue

import (
	"context"
	"io"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/queue"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc"
)

// Client is the queue client.
type Client struct {
	client queue.QueueClient
	closer io.Closer
}

// New constructs a new *Client.
func New(addr string, cert *transport.Cert, trace bool) (*Client, *errors.Error) {
	var (
		closer  io.Closer
		options []grpc.DialOption
		eErr    *errors.Error
	)

	if trace {
		closer, options, eErr = utils.SetUpGRPCTracing("data")
		if eErr != nil {
			return nil, eErr
		}
	}

	cc, err := transport.GRPCDial(cert, addr, options...)
	if err != nil {
		return nil, errors.New(err)
	}

	return &Client{closer: closer, client: queue.NewQueueClient(cc)}, nil
}

// Close closes the client's tracing functionality
func (c *Client) Close() *errors.Error {
	if c.closer != nil {
		return errors.New(c.closer.Close())
	}

	return nil
}

// GetCancel retrieves the cancel state of the run.
func (c *Client) GetCancel(ctx context.Context, id int64) (bool, *errors.Error) {
	status, err := c.client.GetCancel(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return false, errors.New(err)
	}

	return status.Status, nil
}

// SetCancel sets the cancel state for a given run id.
func (c *Client) SetCancel(ctx context.Context, id int64) *errors.Error {
	_, err := c.client.SetCancel(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// NextQueueItem returns the next item in the queue.
func (c *Client) NextQueueItem(ctx context.Context, queueName, hostname string) (*model.QueueItem, *errors.Error) {
	qi, err := c.client.NextQueueItem(ctx, &types.QueueRequest{QueueName: queueName, RunningOn: hostname}, grpc.WaitForReady(false))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewQueueItemFromProto(qi)
}

// SetStatus completes the run by returning its status back to the system.
func (c *Client) SetStatus(ctx context.Context, id int64, status bool) *errors.Error {
	_, err := c.client.PutStatus(ctx, &types.Status{Id: id, Status: status}, grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	return nil
}
