package restapi

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// Errors processes the /errors GET endpoint
func Errors(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	errs, err := h.Clients.Data.GetErrors(ctx, user.Username)
	if err != nil && !err.(errors.Error).Contains(errors.New("record not found")) {
		return nil, 500, err
	}

	for _, err := range errs {
		if err := h.Clients.Data.DeleteError(pCtx, err.ID, user.ID); err != nil && !err.(errors.Error).Contains(errors.New("record not found")) {
			return nil, 500, err
		}
	}

	return errs, 200, nil
}
