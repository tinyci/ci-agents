package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// PostCancelRunIDValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func PostCancelRunIDValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	runID := ctx.Param("run_id")

	if len(runID) == 0 {
		return errors.New("'/cancel/{run_id}': parameter 'run_id' is empty")
	}

	ctx.Set("run_id", runID)

	return nil
}
