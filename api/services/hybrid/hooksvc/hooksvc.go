package hooksvc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"errors"

	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/clients/queue"
	"github.com/tinyci/ci-agents/config"
	topTypes "github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

const (
	eventPush        = "push"
	eventPullRequest = "pull_request"
	eventPing        = "ping"

	actionOpened      = "opened"
	actionSynchronize = "synchronize"
	actionClosed      = "closed"
)

// ErrCancelPR is the error returned when a pr should be canceled.
type ErrCancelPR struct {
	PRID       int64
	Repository string
}

func (ec *ErrCancelPR) Error() string {
	return fmt.Sprintf("PR ID %d is requested to be canceled", ec.PRID)
}

// HandlerConfig configures the hooksvc handler.
type HandlerConfig struct {
	TLS           config.CertConfig `yaml:"tls"`
	QueueEndpoint string            `yaml:"queue_service"` // endpoint of queuesvc to submit to
	DataEndpoint  string            `yaml:"data_service"`
	LogEndpoint   string            `yaml:"log_service"`
}

type (
	dispatchFunc  map[string]func(interface{}) (*topTypes.Submission, error)
	converterFunc map[string]func([]byte) (interface{}, error)
	getRepoFunc   map[string]func(interface{}) (*types.Repository, error)
)

// Handler is the hooksvc handler.
type Handler struct {
	Config HandlerConfig `yaml:",inline"`

	logClient   *log.SubLogger
	queueClient *queue.Client
	dataClient  *data.Client
	dispatch    dispatchFunc
	converter   converterFunc
	getRepo     getRepoFunc
}

// Init initializes the handler.
func (h *Handler) Init() error {
	cert, err := h.Config.TLS.Load()
	if err != nil {
		return err
	}

	h.queueClient, err = queue.New(h.Config.QueueEndpoint, cert, false)
	if err != nil {
		return err
	}

	h.dataClient, err = data.New(h.Config.DataEndpoint, cert, false)
	if err != nil {
		return err
	}

	if h.Config.LogEndpoint != "" {
		if err := log.ConfigureRemote(h.Config.LogEndpoint, cert, false); err != nil {
			fmt.Fprintf(os.Stderr, "Could not configure remote logger: %v\n", err)
		}
	}

	h.logClient = log.NewWithData("hooksvc", nil)
	h.logClient.Info(context.Background(), "Initializing logger")

	h.dispatch = dispatchFunc{
		eventPush:        h.pushDispatch,
		eventPullRequest: h.prDispatch,
	}

	h.converter = converterFunc{
		eventPush:        h.pushConvert,
		eventPullRequest: h.prConvert,
	}

	h.getRepo = getRepoFunc{
		eventPush:        h.pushGetRepo,
		eventPullRequest: h.prGetRepo,
	}

	return nil
}

func (h *Handler) pushDispatch(obj interface{}) (*topTypes.Submission, error) {
	push, ok := obj.(*github.PushEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	return &topTypes.Submission{
		Parent:  push.GetRepo().GetFullName(),
		Fork:    push.GetRepo().GetFullName(),
		HeadSHA: push.GetAfter(),
		BaseSHA: push.GetBefore(),
	}, nil
}

func (h *Handler) prDispatch(obj interface{}) (*topTypes.Submission, error) {
	pr, ok := obj.(*github.PullRequestEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	action := pr.GetAction()

	switch action {
	case actionOpened, actionSynchronize:
		return &topTypes.Submission{
			Parent:   pr.PullRequest.Base.Repo.GetFullName(),
			Fork:     pr.PullRequest.Head.Repo.GetFullName(),
			HeadSHA:  pr.PullRequest.Head.GetSHA(),
			BaseSHA:  pr.PullRequest.Base.GetSHA(),
			TicketID: int64(pr.PullRequest.GetNumber()),
		}, nil
	case actionClosed:
		return nil, &ErrCancelPR{Repository: pr.PullRequest.Base.Repo.GetFullName(), PRID: int64(pr.PullRequest.GetNumber())}
	default:
		return nil, fmt.Errorf("cannot submit; entered %s state", action)
	}
}

func (h *Handler) pushConvert(data []byte) (interface{}, error) {
	obj := &github.PushEvent{}
	return obj, json.Unmarshal(data, obj)
}

func (h *Handler) prConvert(data []byte) (interface{}, error) {
	obj := &github.PullRequestEvent{}
	return obj, json.Unmarshal(data, obj)
}

func (h *Handler) pushGetRepo(obj interface{}) (*types.Repository, error) {
	push, ok := obj.(*github.PushEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	_, _, err := utils.OwnerRepo(push.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	repo, err := h.dataClient.GetRepository(context.Background(), push.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (h *Handler) prGetRepo(obj interface{}) (*types.Repository, error) {
	pr, ok := obj.(*github.PullRequestEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	_, _, err := utils.OwnerRepo(pr.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	repo, err := h.dataClient.GetRepository(context.Background(), pr.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (h *Handler) getLog(req *http.Request, reqUUID uuid.UUID) *log.SubLogger {
	return h.logClient.WithRequest(req).WithFields(log.FieldMap{"request_uuid": reqUUID.String()})
}
