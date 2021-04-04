package operations

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/handlers"
)

// DeleteCapabilitiesUsernameCapabilityValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func DeleteCapabilitiesUsernameCapabilityValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	capability := ctx.Param("capability")

	if len(capability) == 0 {
		return errors.New("'/capabilities/{username}/{capability}': parameter 'capability' is empty")
	}

	ctx.Set("capability", capability)
	username := ctx.Param("username")

	if len(username) == 0 {
		return errors.New("'/capabilities/{username}/{capability}': parameter 'username' is empty")
	}

	ctx.Set("username", username)

	return nil
}
