package queuesvc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/queue"
	gtypes "github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueueServer encapsulates a GRPC server for the queuesvc.
type QueueServer struct {
	H *grpcHandler.H
}

// SetCancel mirrors the cancel in datasvc -- just easier to access by runners.
func (qs *QueueServer) SetCancel(ctx context.Context, id *gtypes.IntID) (*empty.Empty, error) {
	if err := qs.H.Clients.Data.SetCancel(ctx, id.ID); err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// GetCancel mirrors the GetCancel in datasvc -- just easier to access by runners.
func (qs *QueueServer) GetCancel(ctx context.Context, id *gtypes.IntID) (*gtypes.Status, error) {
	state, err := qs.H.Clients.Data.GetCancel(ctx, id.ID)
	if err != nil {
		return &gtypes.Status{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &gtypes.Status{Status: state}, nil
}

// PutStatus pushes the finished run's status out to github and back into the
// datasvc.
func (qs *QueueServer) PutStatus(ctx context.Context, s *gtypes.Status) (*empty.Empty, error) {
	if err := qs.H.Clients.Data.PutStatus(ctx, s.Id, s.Status, s.AdditionalMessage); err != nil {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}

// NextQueueItem gathers the next available item from the queue, if any, and
// returns it. If there is any failure, the queue could not be read and there
// is a need to retry after a wait.
func (qs *QueueServer) NextQueueItem(ctx context.Context, qr *gtypes.QueueRequest) (*gtypes.QueueItem, error) {
	qi, err := qs.H.Clients.Data.NextQueueItem(ctx, qr.QueueName, qr.RunningOn)
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return &gtypes.QueueItem{}, stat.Err()
		}

		return &gtypes.QueueItem{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if qi.Run.Task.Submission.BaseRef.Repository.Owner == nil {
		err := errors.New("No owner for repository for queued run; skipping")
		qs.H.Clients.Log.WithFields(log.FieldMap{
			"repository": qi.Run.Task.Submission.BaseRef.Repository.Name,
			"run_id":     fmt.Sprintf("%d", qi.Run.Id),
			"ran_on":     qr.RunningOn,
		}).Error(ctx, err)

		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	owner := qi.Run.Task.Submission.BaseRef.Repository.Owner

	github, err := qs.H.OAuth.GithubClient(owner.Username, owner.TokenJSON)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "error crafting token")
	}

	parts := strings.SplitN(qi.Run.Task.Submission.BaseRef.Repository.Name, "/", 2)
	if len(parts) != 2 {
		return &gtypes.QueueItem{}, status.Errorf(codes.FailedPrecondition, "invalid repository")
	}

	go func() {
		if err := github.StartedStatus(ctx, parts[0], parts[1], qi.Run.Name, qi.Run.Task.Submission.HeadRef.Sha, fmt.Sprintf("%s/typeslog/%d", qs.H.URL, qi.Run.Id)); err != nil {
			fmt.Println(err)
		}
	}()

	return qi, nil
}

func doSubmit(ctx context.Context, h *grpcHandler.H, qis []*gtypes.QueueItem) (retErr error) {
	since := time.Now()
	defer func() {
		if retErr == nil {
			h.Clients.Log.Infof(ctx, "Successful submission took %v", time.Since(since))
		} else {
			h.Clients.Log.Errorf(ctx, "Submission failed with errors: %v", retErr)
		}
	}()

	if _, err := h.Clients.Data.PutQueue(ctx, qis); err != nil {
		return err
	}

	return nil
}

// Submit is the submission endpoint for the queue; all items gathered from the
// submission are automatically injected into the queue.
func (qs *QueueServer) Submit(ctx context.Context, sub *queue.Submission) (*empty.Empty, error) {
	submission := &types.Submission{
		Parent:      sub.Parent,
		Fork:        sub.Fork,
		HeadSHA:     sub.Headsha,
		BaseSHA:     sub.Basesha,
		TicketID:    sub.TicketID,
		SubmittedBy: sub.SubmittedBy,
		All:         sub.All,
		Manual:      sub.Manual,
	}

	submissionLogger := qs.H.Clients.Log.WithFields(
		log.FieldMap{
			"parent": sub.Parent,
			"fork":   sub.Fork,
			"head":   sub.Headsha,
			"base":   sub.Basesha,
		})

	processCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// XXX this allows the cancel function to be inept while still not leaking
	go func() { time.Sleep(5 * time.Minute); cancel() }()

	sp := qs.newSubmissionProcessor()
	qis, err := sp.process(processCtx, submission)
	if err != nil {
		submissionLogger.Errorf(ctx, "Post-processing error: %v", err)
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	submissionLogger.Infof(ctx, "Putting %d queue items from submissions", len(qis))
	if err := doSubmit(ctx, qs.H, qis); err != nil {
		for _, qi := range qis {
			if err := qs.H.Clients.Data.PutStatus(ctx, qi.Run.Id, false, fmt.Sprintf("Canceled due to error: %v", err)); err != nil {
				submissionLogger.Errorf(ctx, "While canceling runs: %v", err)
			}
		}

		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &empty.Empty{}, nil
}
