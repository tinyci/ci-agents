package restapi

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/utils"
)

// ListTasks retrieves the task list.
func ListTasks(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	page, perPage, err := utils.ScopePagination(ctx.GetString("page"), ctx.GetString("perPage"))
	if err != nil {
		return nil, 500, err
	}

	tasks, err := h.Clients.Data.ListTasks(ctx, ctx.GetString("repository"), ctx.GetString("sha"), page, perPage)
	if err != nil {
		return nil, 500, err
	}

	return tasks, 200, nil
}

// CountTasks counts the task list with the supplied repo/sha filtering.
func CountTasks(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	count, err := h.Clients.Data.CountTasks(ctx, ctx.GetString("repository"), ctx.GetString("sha"))
	if err != nil {
		return nil, 500, err
	}

	return count, 200, nil
}

// GetRunsForTask retrieves all the runs by task id. Pagination rules are in effect.
func GetRunsForTask(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	id, eErr := strconv.ParseInt(ctx.GetString("id"), 10, 64)
	if eErr != nil {
		return nil, 500, errors.New(eErr)
	}

	page, perPage, err := utils.ScopePagination(ctx.GetString("page"), ctx.GetString("perPage"))
	if err != nil {
		return nil, 500, err
	}

	runs, err := h.Clients.Data.GetRunsForTask(ctx, id, page, perPage)
	if err != nil {
		return nil, 500, err
	}

	return runs, 200, nil
}

// CountRunsForTask counts all the runs by task ID.
func CountRunsForTask(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	id, eErr := strconv.ParseInt(ctx.GetString("id"), 10, 64)
	if eErr != nil {
		return nil, 500, errors.New(eErr)
	}

	count, err := h.Clients.Data.CountRunsForTask(ctx, id)
	if err != nil {
		return nil, 500, err
	}

	return count, 200, nil
}

// ListSubscribedTasksForUser lists only the tasks for the repositories the user is subscribed to.
func ListSubscribedTasksForUser(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	page, perPage, err := utils.ScopePagination(ctx.GetString("page"), ctx.GetString("perPage"))
	if err != nil {
		return nil, 500, err
	}

	uname := h.Session(ctx).Get(handlers.SessionUsername)
	u, ok := uname.(string)
	if !ok {
		return nil, 500, errors.New("invalid cookie")
	}

	user, err := h.Clients.Data.GetUser(ctx, u)
	if err != nil {
		return nil, 500, err
	}

	tasks, err := h.Clients.Data.ListSubscribedTasksForUser(ctx, user.ID, page, perPage)
	if err != nil {
		return nil, 500, err
	}

	return tasks, 200, nil
}

// CancelTask cancels a task by ID.
func CancelTask(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	id, eErr := strconv.ParseInt(ctx.GetString("id"), 10, 64)
	if eErr != nil {
		return nil, 500, errors.New(eErr)
	}

	if err := h.Clients.Data.CancelTask(pCtx, id); err != nil {
		return nil, 500, err
	}

	return nil, 200, nil
}
