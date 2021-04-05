package restapi

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/handlers"
)

// GetToken obtains a new token from the db. If one is already set, you must
// delete it before this will return a new one.
func GetToken(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	u, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	token, err := h.Clients.Data.GetToken(ctx, u.Username)
	if err != nil {
		return nil, 500, err
	}

	return token, 200, nil
}

// DeleteToken removes the existing token for the user.
func DeleteToken(pCtx context.Context, h *handlers.H, ctx *gin.Context) (interface{}, int, error) {
	u, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	if err := h.Clients.Data.DeleteToken(pCtx, u.Username); err != nil {
		return nil, 500, err
	}

	return nil, 200, nil
}
