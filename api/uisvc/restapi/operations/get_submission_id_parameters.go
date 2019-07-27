package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetSubmissionIDValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetSubmissionIDValidateURLParams(h *handlers.H, ctx *gin.Context) *errors.Error {
	id := ctx.Param("id")

	if len(id) == 0 {
		return errors.New("'/submission/{id}': parameter 'id' is empty")
	}

	ctx.Set("id", id)

	return nil
}
