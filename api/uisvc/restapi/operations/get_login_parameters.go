package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetLoginValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetLoginValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	code := ctx.Query("code")

	if len(code) == 0 {
		return errors.New("'/login': parameter 'code' is empty")
	}

	ctx.Set("code", code)

	state := ctx.Query("state")

	if len(state) == 0 {
		return errors.New("'/login': parameter 'state' is empty")
	}

	ctx.Set("state", state)

	return nil
}
