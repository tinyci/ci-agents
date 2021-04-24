package uisvc

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/model"
)

// PostCapabilitiesUsernameCapability adds a capability for a user.
func (h *H) PostCapabilitiesUsernameCapability(ctx echo.Context, capability, username string) error {
	fmt.Println(username, capability)

	u, err := h.clients.Data.GetUser(ctx.Request().Context(), username)
	if err != nil {
		return err
	}

	return h.clients.Data.AddCapability(context.Background(), u, model.Capability(capability))
}

// DeleteCapabilitiesUsernameCapability removes a capability from a user.
func (h *H) DeleteCapabilitiesUsernameCapability(ctx echo.Context, capability, username string) error {
	u, err := h.clients.Data.GetUser(ctx.Request().Context(), username)
	if err != nil {
		return err
	}

	return h.clients.Data.RemoveCapability(context.Background(), u, model.Capability(capability))
}
