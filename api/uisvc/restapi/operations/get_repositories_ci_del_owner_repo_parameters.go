package operations

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// GetRepositoriesCiDelOwnerRepoValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetRepositoriesCiDelOwnerRepoValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	owner := ctx.Param("owner")

	if len(owner) == 0 {
		return errors.New("'/repositories/ci/del/{owner}/{repo}': parameter 'owner' is empty")
	}

	ctx.Set("owner", owner)
	repo := ctx.Param("repo")

	if len(repo) == 0 {
		return errors.New("'/repositories/ci/del/{owner}/{repo}': parameter 'repo' is empty")
	}

	ctx.Set("repo", repo)

	return nil
}
