package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// NextQueueItem return the next queue item. The runningOn is a hostname which
// is provided for tracking purposes. It should be unique (but, is ultimately not necessary).
func (c *Client) NextQueueItem(queueName, runningOn string) (*model.QueueItem, *errors.Error) {
	item, err := c.client.QueueNext(context.Background(), &types.QueueRequest{QueueName: queueName, RunningOn: runningOn})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewQueueItemFromProto(item)
}

// PutStatus returns the status of the run.
func (c *Client) PutStatus(runID int64, status bool, msg string) *errors.Error {
	_, err := c.client.PutStatus(context.Background(), &types.Status{AdditionalMessage: msg, Id: runID, Status: status})
	return errors.New(err)
}

// PutQueue adds many QueueItems to the queue.
func (c *Client) PutQueue(qis []*model.QueueItem) ([]*model.QueueItem, *errors.Error) {
	ql := &data.QueueList{}

	for _, qi := range qis {
		ql.Items = append(ql.Items, qi.ToProto())
	}

	ql, err := c.client.QueueAdd(context.Background(), ql)
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
func (c *Client) SetCancel(id int64) *errors.Error {
	_, err := c.client.SetCancel(context.Background(), &types.IntID{ID: id})
	return errors.New(err)
}

// GetCancel returns the state for the run.
func (c *Client) GetCancel(id int64) (bool, *errors.Error) {
	status, err := c.client.GetCancel(context.Background(), &types.IntID{ID: id})
	if err != nil {
		return false, errors.New(err)
	}

	return status.Status, nil
}
