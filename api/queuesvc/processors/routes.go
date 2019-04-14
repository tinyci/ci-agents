package processors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/handler"
	"github.com/tinyci/ci-agents/grpc/services/queue"
	gtypes "github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
	"golang.org/x/oauth2"
)

// QueueServer encapsulates a GRPC server for the queuesvc.
type QueueServer struct {
	H *handler.H
}

// SetCancel mirrors the cancel in datasvc -- just easier to access by runners.
func (qs *QueueServer) SetCancel(ctx context.Context, id *gtypes.IntID) (*empty.Empty, error) {
	if err := qs.H.Clients.Data.SetCancel(id.ID); err != nil {
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}

// GetCancel mirrors the GetCancel in datasvc -- just easier to access by runners.
func (qs *QueueServer) GetCancel(ctx context.Context, id *gtypes.IntID) (*gtypes.Status, error) {
	state, err := qs.H.Clients.Data.GetCancel(id.ID)
	if err != nil {
		return &gtypes.Status{}, err
	}

	return &gtypes.Status{Status: state}, nil
}

// PutStatus pushes the finished run's status out to github and back into the
// datasvc.
func (qs *QueueServer) PutStatus(ctx context.Context, status *gtypes.Status) (*empty.Empty, error) {
	if err := qs.H.Clients.Data.PutStatus(status.Id, status.Status, status.AdditionalMessage); err != nil {
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}

// NextQueueItem gathers the next available item from the queue, if any, and
// returns it. If there is any failure, the queue could not be read and there
// is a need to retry after a wait.
func (qs *QueueServer) NextQueueItem(ctx context.Context, qr *gtypes.QueueRequest) (*gtypes.QueueItem, error) {
	qi, err := qs.H.Clients.Data.NextQueueItem(qr.QueueName, qr.RunningOn)
	if err != nil {
		if err.Contains(errors.ErrNotFound) {
			err.SetLog(false)
		}
		return &gtypes.QueueItem{}, err
	}

	if qi.Run.Task.Parent.Owner == nil {
		err := errors.New("No owner for repository for queued run; skipping")
		qs.H.Clients.Log.WithFields(log.FieldMap{
			"repository": qi.Run.Task.Parent.Name,
			"run_id":     fmt.Sprintf("%d", qi.Run.ID),
		}).Error(err)

		return nil, err
	}

	token := &oauth2.Token{}
	if err := utils.JSONIO(qi.Run.Task.Parent.Owner.Token, token); err != nil {
		return &gtypes.QueueItem{}, err
	}

	github := qs.H.OAuth.GithubClient(token)
	parts := strings.SplitN(qi.Run.Task.Parent.Name, "/", 2)
	if len(parts) != 2 {
		return &gtypes.QueueItem{}, errors.New("invalid repository")
	}

	go func() {
		if err := github.StartedStatus(parts[0], parts[1], qi.Run.Name, qi.Run.Task.Ref.SHA, fmt.Sprintf("%s/log/%d", qs.H.URL, qi.Run.ID)); err != nil {
			fmt.Println(err)
		}
	}()

	return qi.ToProto(), nil
}

func doSubmit(h *handler.H, qis []*model.QueueItem) (retErr *errors.Error) {
	since := time.Now()
	defer func() {
		if retErr == nil {
			h.Clients.Log.Infof("Successful submission took %v", time.Since(since))
		} else {
			h.Clients.Log.Errorf("Submission failed with errors: %v", retErr)
		}
	}()

	if _, err := h.Clients.Data.PutQueue(qis); err != nil {
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
		PullRequest: sub.PullRequest,
		SubmittedBy: sub.SubmittedBy,
		All:         sub.All,
		Manual:      sub.Manual,
	}

	qis, err := Process(qs.H, submission)
	if err != nil {
		return &empty.Empty{}, err
	}

	if err := doSubmit(qs.H, qis); err != nil {
		for _, qi := range qis {
			qs.H.Clients.Data.PutStatus(qi.Run.ID, false, fmt.Sprintf("Canceled due to error: %v", err))
		}
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}
