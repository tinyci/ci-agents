package hooksvc

// based on konboi/ghooks.

import (
	"context"
	"crypto/hmac"
	"crypto/sha1" // #nosec
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tinyci/ci-agents/clients/log"
)

// this ugly as sin setup allows me to pass around the http writer and not have
// to care about where the errors get set. I am 100% certain this will break me
// in some way later, yet I am a frustrated, overworked engineer that just
// wants to get shit done. I love all of you. -erikh
func (h *Handler) handleBasic(w http.ResponseWriter, logger *log.SubLogger, event string) bool {
	if event == "" {
		logger.Error(context.Background(), "Rejected hook event because the event was missing")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return true
	}

	if event == eventPing {
		logger.Infof(context.Background(), "responding to hook ping from remote host")
		return true
	}

	return false
}

func (h *Handler) parseBody(w http.ResponseWriter, logger *log.SubLogger, body []byte, event string) interface{} {
	converter, ok := h.converter[event]
	if !ok {
		logger.Errorf(context.Background(), "Could not find converter for event %q", event)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return nil
	}

	obj, err := converter(body)
	if err != nil {
		logger.Errorf(context.Background(), "Rejected hook event because the JSON converter failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return obj
}

func (h *Handler) checkRepo(w http.ResponseWriter, req *http.Request, logger *log.SubLogger, body []byte, obj interface{}, event string) bool {
	getRepo, ok := h.getRepo[event]
	if !ok {
		logger.Errorf(context.Background(), "Could not find repository reader for event %q", event)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return true
	}

	repo, err := getRepo(obj)
	if err != nil {
		logger.Errorf(context.Background(), "Rejected hook event because we could not read from the repository: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}

	if repo.HookSecret == "" {
		logger.Error(context.Background(), "Rejected hook event because the hook secret is missing")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return true
	}

	signature := req.Header.Get("X-Hub-Signature")
	if !h.isValidSignature(body, repo.HookSecret, signature) {
		logger.Error(context.Background(), "Rejected hook event because the request signature is invalid")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return true
	}

	return false
}

func (h *Handler) doSubmit(w http.ResponseWriter, logger *log.SubLogger, obj interface{}, event string) bool {
	dispatch, ok := h.dispatch[event]
	if !ok {
		logger.Errorf(context.Background(), "Could not find dispatcher for event %q", event)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return true
	}
	sub, err := dispatch(obj)
	if err != nil {
		switch err := err.(type) {
		case *ErrCancelPR:
			logger.WithFields(log.FieldMap{"pull_request": fmt.Sprintf("%v", err.PRID), "repository": err.Repository}).Infof(context.Background(), "Canceling PR")
			if err := h.dataClient.CancelTasksByPR(err.Repository, err.PRID); err != nil {
				return true
			}

			return false
		default:
			logger.Errorf(context.Background(), "Rejected hook event because we could not dispatch the submission: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}
	}

	go func() {
		now := time.Now()
		if err := h.queueClient.Submit(context.Background(), sub); err != nil {
			logger.Error(context.Background(), "Rejected hook event because the submission failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.WithFields(log.FieldMap{"submission": fmt.Sprintf("%#v", *sub)}).Infof(context.Background(), "Submission took %v", time.Since(now))
	}()

	return false
}

func (h *Handler) isValidSignature(body []byte, secret, signature string) bool {
	if !strings.HasPrefix(signature, "sha1=") {
		return false
	}

	mac := hmac.New(sha1.New, []byte(secret))
	if _, err := mac.Write(body); err != nil {
		return false
	}
	actual := mac.Sum(nil)

	expected, err := hex.DecodeString(signature[5:])
	if err != nil {
		return false
	}

	return hmac.Equal(actual, expected)
}

// ServeHTTP is the primary handler returned by the server.
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	u := uuid.New()
	logger := h.getLog(req, u)

	since := time.Now()
	logger.Info(context.Background(), "Received hook request")

	if req.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	event := req.Header.Get("X-GitHub-Event")
	logger = logger.WithFields(log.FieldMap{"event": event})
	go func() { logger.Infof(context.Background(), "Full RTT for hook was %v", time.Since(since)) }()

	if h.handleBasic(w, logger, event) {
		return
	}

	if req.Body == nil {
		logger.Error(context.Background(), "Rejected hook event because the body was missing")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Errorf(context.Background(), "Rejected hook event because we could not read the body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	obj := h.parseBody(w, logger, body, event)
	if obj == nil {
		return // hit an error
	}

	if h.checkRepo(w, req, logger, body, obj, event) {
		return // hit an error
	}

	if h.doSubmit(w, logger, obj, event) {
		return // error
	}

	w.WriteHeader(http.StatusOK)
}
