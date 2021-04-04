package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/handlers"
)

// GetRepositoriesMyValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetRepositoriesMyValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	search := ctx.Query("search")

	ctx.Set("search", search)

	return nil
}
