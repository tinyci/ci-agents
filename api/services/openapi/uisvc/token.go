package uisvc

import (
	"context"

	"github.com/labstack/echo/v4"
)

// GetToken obtains a new token from the db. If one is already set, you must
// delete it before this will return a new one.
func (h *H) GetToken(ctx echo.Context) error {
	u, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	token, err := h.clients.Data.GetToken(ctx.Request().Context(), u.Username)
	if err != nil {
		return err
	}

	return ctx.JSON(200, token)
}

// DeleteToken removes the existing token for the user.
func (h *H) DeleteToken(ctx echo.Context) error {
	u, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	if err := h.clients.Data.DeleteToken(context.Background(), u.Username); err != nil {
		return err
	}

	return nil
}
