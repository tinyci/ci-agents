// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetTasksRunsID swagger:route GET /tasks/runs/{id} getTasksRunsId
// Obtain the run list based on the task ID.
// The queue list only contains: * stuff * other junk
//
func GetTasksRunsID(h *handlers.H, ctx *gin.Context, processingHandler handlers.HandlerFunc) *errors.Error {
	if h.RequestLogging {
		start := time.Now()
		u := uuid.New()

		content, jsonErr := json.Marshal(ctx.Params)
		if jsonErr != nil {
			h.Clients.Log.Error(ctx.Request.Context(), errors.New(jsonErr).Wrap("encoding params for log message"))
		}

		logger := h.Clients.Log.WithRequest(ctx.Request).WithFields(log.FieldMap{
			"params":       string(content),
			"request_uuid": u.String(),
		})

		user, err := h.GetGithub(ctx)
		if err == nil {
			logger = logger.WithUser(user)
		}

		logger.Debug(ctx.Request.Context(), "incoming request")

		defer func() {
			logger.WithFields(log.FieldMap{
				"duration": time.Since(start).String(),
			}).Debug(ctx.Request.Context(), "request completed")
		}()
	}

	if err := GetTasksRunsIDValidateURLParams(h, ctx); err != nil {
		return errors.New(err)
	}

	if processingHandler == nil {
		return errors.New("'/tasks/runs/{id}': no processor defined")
	}

	processingContext, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	resp, code, err := processingHandler(processingContext, h, ctx)
	return GetTasksRunsIDResponse(h, ctx, resp, code, err)
}
