package uisvc

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/utils"
)

// GetToken obtains a new token from the db. If one is already set, you must
// delete it before this will return a new one.
func (h *H) GetToken(ctx echo.Context) error {
	u, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	token, err := h.clients.Data.GetToken(ctx.Request().Context(), u)
	if err != nil {
		return err
	}

	return ctx.JSON(200, token)
}

// DeleteToken removes the existing token for the user.
func (h *H) DeleteToken(ctx echo.Context) error {
	u, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	if err := h.clients.Data.DeleteToken(context.Background(), u); err != nil {
		return err
	}

	return nil
}
