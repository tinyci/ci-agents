package restapi

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/model"
)

// AddCapability adds a capability for a user.
func AddCapability(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	u, err := h.Clients.Data.GetUser(pCtx, ctx.GetString("username"))
	if err != nil {
		return nil, 500, err
	}

	return nil, 200, h.Clients.Data.AddCapability(pCtx, u, model.Capability(ctx.GetString("capability")))
}

// RemoveCapability removes a capability from a user.
func RemoveCapability(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	u, err := h.Clients.Data.GetUser(pCtx, ctx.GetString("username"))
	if err != nil {
		return nil, 500, err
	}

	return nil, 200, h.Clients.Data.RemoveCapability(pCtx, u, model.Capability(ctx.GetString("capability")))
}
