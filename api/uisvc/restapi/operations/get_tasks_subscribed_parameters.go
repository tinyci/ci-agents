package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetTasksSubscribedValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetTasksSubscribedValidateURLParams(h *handlers.H, ctx *gin.Context) *errors.Error {
	page := ctx.Query("page")

	if len(page) == 0 {
		page = "0"
	}

	ctx.Set("page", page)

	perPage := ctx.Query("perPage")

	if len(perPage) == 0 {
		perPage = "100"
	}

	ctx.Set("perPage", perPage)

	return nil
}
