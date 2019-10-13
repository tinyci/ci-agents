package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/model"
)

// AddCapability adds a capability for a user.
func AddCapability(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	u, err := h.Clients.Data.GetUser(ctx, ctx.GetString("username"))
	if err != nil {
		return nil, 500, err
	}

	return nil, 200, h.Clients.Data.AddCapability(ctx, u, model.Capability(ctx.GetString("capability")))
}

// RemoveCapability removes a capability from a user.
func RemoveCapability(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	u, err := h.Clients.Data.GetUser(ctx, ctx.GetString("username"))
	if err != nil {
		return nil, 500, err
	}

	return nil, 200, h.Clients.Data.RemoveCapability(ctx, u, model.Capability(ctx.GetString("capability")))
}
