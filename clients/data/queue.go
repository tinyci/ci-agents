package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"google.golang.org/grpc"
)

// NextQueueItem return the next queue item. The runningOn is a hostname which
// is provided for tracking purposes. It should be unique (but, is ultimately not necessary).
func (c *Client) NextQueueItem(ctx context.Context, queueName, runningOn string) (*types.QueueItem, error) {
	item, err := c.client.QueueNext(ctx, &types.QueueRequest{QueueName: queueName, RunningOn: runningOn}, grpc.WaitForReady(false))
	if err != nil {
		return nil, err
	}

	return item, nil
}

// PutStatus returns the status of the run.
func (c *Client) PutStatus(ctx context.Context, runID int64, status bool, msg string) error {
	_, err := c.client.PutStatus(ctx, &types.Status{AdditionalMessage: msg, Id: runID, Status: status}, grpc.WaitForReady(true))
	return err
}

// PutQueue adds many QueueItems to the queue.
func (c *Client) PutQueue(ctx context.Context, qis []*types.QueueItem) ([]*types.QueueItem, error) {
	ret, err := c.client.QueueAdd(ctx, &data.QueueList{Items: qis}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}
	return []*types.QueueItem(ret.Items), nil
}

// SetCancel cancels a run, and any other task-level runs.
func (c *Client) SetCancel(ctx context.Context, id int64) error {
	_, err := c.client.SetCancel(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	return err
}

// GetCancel returns the state for the run.
func (c *Client) GetCancel(ctx context.Context, id int64) (bool, error) {
	status, err := c.client.GetCancel(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return false, err
	}

	return status.Status, nil
}
