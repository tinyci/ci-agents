package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetSubmitValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetSubmitValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	all := ctx.Query("all")

	ctx.Set("all", all)

	repository := ctx.Query("repository")

	if len(repository) == 0 {
		return errors.New("'/submit': parameter 'repository' is empty")
	}

	ctx.Set("repository", repository)

	sha := ctx.Query("sha")

	if len(sha) == 0 {
		return errors.New("'/submit': parameter 'sha' is empty")
	}

	ctx.Set("sha", sha)

	return nil
}
