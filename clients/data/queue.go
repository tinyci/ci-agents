package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// NextQueueItem return the next queue item. The runningOn is a hostname which
// is provided for tracking purposes. It should be unique (but, is ultimately not necessary).
func (c *Client) NextQueueItem(ctx context.Context, queueName, runningOn string) (*model.QueueItem, error) {
	item, err := c.client.QueueNext(ctx, &types.QueueRequest{QueueName: queueName, RunningOn: runningOn}, grpc.WaitForReady(false))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewQueueItemFromProto(item)
}

// PutStatus returns the status of the run.
func (c *Client) PutStatus(ctx context.Context, runID int64, status bool, msg string) error {
	_, err := c.client.PutStatus(ctx, &types.Status{AdditionalMessage: msg, Id: runID, Status: status}, grpc.WaitForReady(true))
	return errors.New(err)
}

// PutQueue adds many QueueItems to the queue.
func (c *Client) PutQueue(ctx context.Context, qis []*model.QueueItem) ([]*model.QueueItem, error) {
	ql := &data.QueueList{}

	for _, qi := range qis {
		ql.Items = append(ql.Items, qi.ToProto())
	}

	ql, err := c.client.QueueAdd(ctx, ql, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	qis2 := []*model.QueueItem{}

	for _, qi := range ql.Items {
		pqi, err := model.NewQueueItemFromProto(qi)
		if err != nil {
			return nil, err
		}

		qis2 = append(qis2, pqi)
	}

	return qis2, nil
}

// SetCancel cancels a run, and any other task-level runs.
func (c *Client) SetCancel(ctx context.Context, id int64) error {
	_, err := c.client.SetCancel(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	return errors.New(err)
}

// GetCancel returns the state for the run.
func (c *Client) GetCancel(ctx context.Context, id int64) (bool, error) {
	status, err := c.client.GetCancel(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return false, errors.New(err)
	}

	return status.Status, nil
}
