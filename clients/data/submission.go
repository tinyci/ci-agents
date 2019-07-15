package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
)

// PutSubmission puts a submission into the datasvc
func (c *Client) PutSubmission(sub *model.Submission) (*model.Submission, *errors.Error) {
	s, err := c.client.PutSubmission(context.Background(), sub.ToProto())
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewSubmissionFromProto(s)
}
