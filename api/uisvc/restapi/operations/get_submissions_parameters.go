package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetSubmissionsValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetSubmissionsValidateURLParams(h *handlers.H, ctx *gin.Context) *errors.Error {
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

	repository := ctx.Query("repository")

	ctx.Set("repository", repository)

	sha := ctx.Query("sha")

	ctx.Set("sha", sha)

	return nil
}
