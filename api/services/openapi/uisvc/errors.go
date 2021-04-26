package uisvc

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/utils"
)

// GetErrors processes the /errors GET endpoint
func (h *H) GetErrors(ctx echo.Context) error {
	user, ok := h.getUser(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	errs, err := h.clients.Data.GetErrors(ctx.Request().Context(), user.Username)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}

	for _, err := range errs.Errors {
		if err := h.clients.Data.DeleteError(context.Background(), err.Id, user.Id); err != nil && err != utils.ErrNotFound {
			return err
		}
	}

	return ctx.JSON(200, errs.Errors)
}
