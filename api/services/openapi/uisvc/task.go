package uisvc

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/utils"
)

// GetTasks retrieves the task list.
func (h *H) GetTasks(ctx echo.Context, params uisvc.GetTasksParams) error {
	page, perPage, err := utils.ScopePaginationInt(params.Page, params.PerPage)
	if err != nil {
		return err
	}

	tasks, err := h.clients.Data.ListTasks(ctx.Request().Context(), stringDeref(params.Repository), stringDeref(params.Sha), int64(page), int64(perPage))
	if err != nil {
		return err
	}

	return ctx.JSON(200, sanitizeTasks(tasks.Tasks))
}

// GetTasksCount counts the task list with the supplied repo/sha filtering.
func (h *H) GetTasksCount(ctx echo.Context, params uisvc.GetTasksCountParams) error {
	count, err := h.clients.Data.CountTasks(ctx.Request().Context(), stringDeref(params.Repository), stringDeref(params.Sha))
	if err != nil {
		return err
	}

	return ctx.JSON(200, count)
}

// GetTasksRunsId retrieves all the runs by task id. Pagination rules are in effect.
func (h *H) GetTasksRunsId(ctx echo.Context, id int64, params uisvc.GetTasksRunsIdParams) error {
	page, perPage, err := utils.ScopePaginationInt(params.Page, params.PerPage)
	if err != nil {
		return err
	}

	runs, err := h.clients.Data.GetRunsForTask(ctx.Request().Context(), id, int64(page), int64(perPage))
	if err != nil {
		return err
	}

	return ctx.JSON(200, sanitizeRuns(runs.List))
}

// GetTasksRunsIdCount counts all the runs by task ID.
func (h *H) GetTasksRunsIdCount(ctx echo.Context, id int64) error {
	count, err := h.clients.Data.CountRunsForTask(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(200, count)
}

// GetTasksSubscribed lists only the tasks for the repositories the user is subscribed to.
func (h *H) GetTasksSubscribed(ctx echo.Context, params uisvc.GetTasksSubscribedParams) error {
	return ctx.JSON(200, []interface{}{}) // FIXME remove this method
}

// PostTasksCancelId cancels a task by ID.
func (h *H) PostTasksCancelId(ctx echo.Context, id int64) error {
	if err := h.clients.Data.CancelTask(context.Background(), id); err != nil {
		return err
	}

	return ctx.NoContent(200)
}
