package restapi

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/utils"
)

// CountRuns returns a count of the queue items by asking the datasvc for it.
func CountRuns(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	count, err := h.Clients.Data.RunCount(ctx, ctx.GetString("repository"), ctx.GetString("sha"))
	if err != nil {
		return nil, 500, err
	}

	return count, 200, nil
}

// ListRuns lists all the runs that were requested by the page/perPage parameters.
func ListRuns(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	page, perPage, err := utils.ScopePagination(ctx.GetString("page"), ctx.GetString("perPage"))
	if err != nil {
		return nil, 500, err
	}

	list, err := h.Clients.Data.ListRuns(ctx, ctx.GetString("repository"), ctx.GetString("sha"), page, perPage)
	if err != nil {
		return nil, 500, err
	}

	return list, 200, nil
}

// GetRun retrieves a run by id.
func GetRun(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	runID, err := strconv.ParseInt(ctx.GetString("run_id"), 10, 64)
	if err != nil {
		return nil, 500, errors.New(err)
	}

	run, eErr := h.Clients.Data.GetRunUI(ctx, runID)
	if eErr != nil {
		return nil, 500, eErr
	}

	return run, 200, nil
}

// CancelRun cancels a run by id.
func CancelRun(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	runID, err := strconv.ParseInt(ctx.GetString("run_id"), 10, 64)
	if err != nil {
		return nil, 500, errors.New(err)
	}

	if err := h.Clients.Data.SetCancel(pCtx, runID); err != nil {
		return nil, 500, errors.New(err)
	}

	return nil, 200, nil
}
