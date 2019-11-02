package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// PostTasksCancelIDValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func PostTasksCancelIDValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	id := ctx.Param("id")

	if len(id) == 0 {
		return errors.New("'/tasks/cancel/{id}': parameter 'id' is empty")
	}

	ctx.Set("id", id)

	return nil
}
