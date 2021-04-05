package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/handlers"
)

// GetTasksCountValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetTasksCountValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	repository := ctx.Query("repository")

	ctx.Set("repository", repository)

	sha := ctx.Query("sha")

	ctx.Set("sha", sha)

	return nil
}
