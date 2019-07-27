package restapi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/utils"
)

// GetSubmission retrieves a submission by id
func GetSubmission(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	id, eErr := strconv.ParseInt(ctx.GetString("id"), 10, 64)
	if eErr != nil {
		return nil, 500, errors.New(eErr)
	}

	sub, err := h.Clients.Data.GetSubmissionByID(id)
	if err != nil {
		return nil, 500, err
	}

	return sub, 200, nil
}

// GetSubmissionTasks retrieves a submission's task from the submission id
func GetSubmissionTasks(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	id, eErr := strconv.ParseInt(ctx.GetString("id"), 10, 64)
	if eErr != nil {
		return nil, 500, errors.New(eErr)
	}

	sub, err := h.Clients.Data.GetSubmissionByID(id)
	if err != nil {
		return nil, 500, err
	}

	page, perPage, err := utils.ScopePagination(ctx.GetString("page"), ctx.GetString("perPage"))
	if err != nil {
		return nil, 500, err
	}

	tasks, err := h.Clients.Data.GetTasksForSubmission(sub, page, perPage)
	if err != nil {
		return nil, 500, err
	}

	return tasks, 200, nil
}

// ListSubmissions lists the submissions with optional repository/sha filtering and pagination.
func ListSubmissions(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	page, perPage, err := utils.ScopePagination(ctx.GetString("page"), ctx.GetString("perPage"))
	if err != nil {
		return nil, 500, err
	}

	list, err := h.Clients.Data.ListSubmissions(page, perPage, ctx.GetString("repository"), ctx.GetString("sha"))
	if err != nil {
		return nil, 500, err
	}

	return list, 200, nil
}
