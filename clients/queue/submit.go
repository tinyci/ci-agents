package queue

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/queue"
	"github.com/tinyci/ci-agents/types"
	"google.golang.org/grpc"
)

// Submit submits a push or pull request to the queue.
func (c *Client) Submit(ctx context.Context, sub *types.Submission) error {
	_, err := c.client.Submit(ctx, &queue.Submission{
		Headsha:     sub.HeadSHA,
		Basesha:     sub.BaseSHA,
		Parent:      sub.Parent,
		Fork:        sub.Fork,
		All:         sub.All,
		SubmittedBy: sub.SubmittedBy,
		Manual:      sub.Manual,
		TicketID:    sub.TicketID,
	}, grpc.WaitForReady(true))
	return err
}
