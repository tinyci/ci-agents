package uisvc

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/utils"
)

// GetSubmissionId retrieves a submission by id
func (h *H) GetSubmissionId(ctx echo.Context, id int64) error {
	sub, err := h.clients.Data.GetSubmissionByID(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(200, sub)
}

// GetSubmissionIdRuns retrieves a submission's runs from the submission id
func (h *H) GetSubmissionIdRuns(ctx echo.Context, id int64, params uisvc.GetSubmissionIdRunsParams) error {
	sub, err := h.clients.Data.GetSubmissionByID(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	page, perPage, err := utils.ScopePaginationInt(params.Page, params.PerPage)
	if err != nil {
		return err
	}

	runs, err := h.clients.Data.GetRunsForSubmission(ctx.Request().Context(), sub, int64(page), int64(perPage))
	if err != nil {
		return err
	}

	return ctx.JSON(200, sanitizeRuns(runs.List))
}

// GetSubmissionIdTasks retrieves a submission's task from the submission id
func (h *H) GetSubmissionIdTasks(ctx echo.Context, id int64, params uisvc.GetSubmissionIdTasksParams) error {
	sub, err := h.clients.Data.GetSubmissionByID(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	page, perPage, err := utils.ScopePaginationInt(params.Page, params.PerPage)
	if err != nil {
		return err
	}

	tasks, err := h.clients.Data.GetTasksForSubmission(ctx.Request().Context(), sub, int64(page), int64(perPage))
	if err != nil {
		return err
	}

	return ctx.JSON(200, sanitizeTasks(tasks.Tasks))
}

// GetSubmissions lists the submissions with optional repository/sha filtering and pagination.
func (h *H) GetSubmissions(ctx echo.Context, params uisvc.GetSubmissionsParams) error {
	page, perPage, err := utils.ScopePaginationInt(params.Page, params.PerPage)
	if err != nil {
		return err
	}

	list, err := h.clients.Data.ListSubmissions(ctx.Request().Context(), int64(page), int64(perPage), stringDeref(params.Repository), stringDeref(params.Sha))
	if err != nil {
		return err
	}

	return ctx.JSON(200, sanitizeSubmissions(list.Submissions))
}

// GetSubmissionsCount counts the submissions with optional repository/sha filtering.
func (h *H) GetSubmissionsCount(ctx echo.Context, params uisvc.GetSubmissionsCountParams) error {
	count, err := h.clients.Data.CountSubmissions(ctx.Request().Context(), stringDeref(params.Repository), stringDeref(params.Sha))
	if err != nil {
		return err
	}

	return ctx.JSON(200, count)
}

// PostSubmissionIdCancel cancels a submission by ID.
func (h *H) PostSubmissionIdCancel(ctx echo.Context, id int64) error {
	if err := h.clients.Data.CancelSubmission(context.Background(), id); err != nil {
		return err
	}

	return ctx.NoContent(200)
}
