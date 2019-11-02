package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// PostSubmissionIDCancelValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func PostSubmissionIDCancelValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	id := ctx.Param("id")

	if len(id) == 0 {
		return errors.New("'/submission/{id}/cancel': parameter 'id' is empty")
	}

	ctx.Set("id", id)

	return nil
}
