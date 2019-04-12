package queue

import (
	"context"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/queue"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// Client is the queue client.
type Client struct {
	client queue.QueueClient
}

// New constructs a new *Client.
func New(addr string, cert *transport.Cert) (*Client, *errors.Error) {
	cc, err := transport.GRPCDial(cert, addr)
	if err != nil {
		return nil, errors.New(err)
	}

	return &Client{client: queue.NewQueueClient(cc)}, nil
}

// GetCancel retrieves the cancel state of the run.
func (c *Client) GetCancel(id int64) (bool, *errors.Error) {
	status, err := c.client.GetCancel(context.Background(), &types.IntID{ID: id})
	if err != nil {
		return false, errors.New(err)
	}

	return status.Status, nil
}

// SetCancel sets the cancel state for a given run id.
func (c *Client) SetCancel(id int64) *errors.Error {
	_, err := c.client.SetCancel(context.Background(), &types.IntID{ID: id})
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// NextQueueItem returns the next item in the queue.
func (c *Client) NextQueueItem(queueName, hostname string) (*model.QueueItem, *errors.Error) {
	qi, err := c.client.NextQueueItem(context.Background(), &types.QueueRequest{QueueName: queueName, RunningOn: hostname})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewQueueItemFromProto(qi)
}

// SetStatus completes the run by returning its status back to the system.
func (c *Client) SetStatus(id int64, status bool) *errors.Error {
	_, err := c.client.PutStatus(context.Background(), &types.Status{Id: id, Status: status})
	if err != nil {
		return errors.New(err)
	}

	return nil
}
