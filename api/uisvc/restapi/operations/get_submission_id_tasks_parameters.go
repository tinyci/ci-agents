package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetSubmissionIDTasksValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetSubmissionIDTasksValidateURLParams(h *handlers.H, ctx *gin.Context) *errors.Error {
	id := ctx.Param("id")

	if len(id) == 0 {
		return errors.New("'/submission/{id}/tasks': parameter 'id' is empty")
	}

	ctx.Set("id", id)

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
