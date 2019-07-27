package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetToken obtains a new token from the db. If one is already set, you must
// delete it before this will return a new one.
func GetToken(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	u, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	token, err := h.Clients.Data.GetToken(u.Username)
	if err != nil {
		return nil, 500, err
	}

	return token, 200, nil
}

// DeleteToken removes the existing token for the user.
func DeleteToken(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	u, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	if err := h.Clients.Data.DeleteToken(u.Username); err != nil {
		return nil, 500, err
	}

	return nil, 200, nil
}
