package queue

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/queue"
	"github.com/tinyci/ci-agents/types"
	"google.golang.org/grpc"
)

// Submit submits a push or pull request to the queue.
func (c *Client) Submit(sub *types.Submission) *errors.Error {
	_, err := c.client.Submit(context.Background(), &queue.Submission{
		Headsha:     sub.HeadSHA,
		Basesha:     sub.BaseSHA,
		Parent:      sub.Parent,
		Fork:        sub.Fork,
		All:         sub.All,
		SubmittedBy: sub.SubmittedBy,
		Manual:      sub.Manual,
		PullRequest: sub.PullRequest,
	}, grpc.WaitForReady(true))
	return errors.New(err)
}
