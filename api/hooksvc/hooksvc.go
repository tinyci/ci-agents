package hooksvc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/clients/queue"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
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
	dispatchFunc  map[string]func(interface{}) (*types.Submission, error)
	converterFunc map[string]func([]byte) (interface{}, error)
	getRepoFunc   map[string]func(interface{}) (*model.Repository, error)
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
func (h *Handler) Init() *errors.Error {
	cert, err := h.Config.TLS.Load()
	if err != nil {
		return err
	}

	h.queueClient, err = queue.New(h.Config.QueueEndpoint, cert)
	if err != nil {
		return err
	}

	h.dataClient, err = data.New(h.Config.DataEndpoint, cert)
	if err != nil {
		return err
	}

	if h.Config.LogEndpoint != "" {
		log.ConfigureRemote(h.Config.LogEndpoint, cert)
	}

	h.logClient = log.NewWithData("hooksvc", nil)
	h.logClient.Info("Initializing logger")

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

func (h *Handler) pushDispatch(obj interface{}) (*types.Submission, error) {
	push, ok := obj.(*github.PushEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	return &types.Submission{
		Parent:  push.GetRepo().GetFullName(),
		Fork:    push.GetRepo().GetFullName(),
		HeadSHA: push.GetAfter(),
		BaseSHA: push.GetBefore(),
	}, nil
}

func (h *Handler) prDispatch(obj interface{}) (*types.Submission, error) {
	pr, ok := obj.(*github.PullRequestEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	action := pr.GetAction()

	switch action {
	case actionOpened, actionSynchronize:
		return &types.Submission{
			Parent:      pr.PullRequest.Base.Repo.GetFullName(),
			Fork:        pr.PullRequest.Head.Repo.GetFullName(),
			HeadSHA:     pr.PullRequest.Head.GetSHA(),
			BaseSHA:     pr.PullRequest.Base.GetSHA(),
			PullRequest: int64(pr.PullRequest.GetNumber()),
		}, nil
	case actionClosed:
		return nil, &ErrCancelPR{Repository: pr.PullRequest.Base.Repo.GetFullName(), PRID: int64(pr.PullRequest.GetNumber())}
	default:
		return nil, errors.Errorf("cannot submit; entered %s state\n", action)
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

func (h *Handler) pushGetRepo(obj interface{}) (*model.Repository, error) {
	push, ok := obj.(*github.PushEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	_, _, err := utils.OwnerRepo(push.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	repo, err := h.dataClient.GetRepository(push.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (h *Handler) prGetRepo(obj interface{}) (*model.Repository, error) {
	pr, ok := obj.(*github.PullRequestEvent)
	if !ok {
		return nil, errors.New("cast failed")
	}

	_, _, err := utils.OwnerRepo(pr.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	repo, err := h.dataClient.GetRepository(pr.GetRepo().GetFullName())
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (h *Handler) getLog(req *http.Request, reqUUID uuid.UUID) *log.SubLogger {
	return h.logClient.WithRequest(req).WithFields(log.FieldMap{"request_uuid": reqUUID.String()})
}
