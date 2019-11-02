package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetTasksRunsIDCountValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetTasksRunsIDCountValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	id := ctx.Param("id")

	if len(id) == 0 {
		return errors.New("'/tasks/runs/{id}/count': parameter 'id' is empty")
	}

	ctx.Set("id", id)

	return nil
}
