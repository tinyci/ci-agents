package operations

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/handlers"
)

// GetTasksRunsIDValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetTasksRunsIDValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	id := ctx.Param("id")

	if len(id) == 0 {
		return errors.New("'/tasks/runs/{id}': parameter 'id' is empty")
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
