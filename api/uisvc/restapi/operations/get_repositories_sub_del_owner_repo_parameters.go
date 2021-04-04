package operations

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/handlers"
)

// GetRepositoriesSubDelOwnerRepoValidateURLParams validates the parameters in the
// URL according to the swagger specification.
func GetRepositoriesSubDelOwnerRepoValidateURLParams(h *handlers.H, ctx *gin.Context) error {
	owner := ctx.Param("owner")

	if len(owner) == 0 {
		return errors.New("'/repositories/sub/del/{owner}/{repo}': parameter 'owner' is empty")
	}

	ctx.Set("owner", owner)
	repo := ctx.Param("repo")

	if len(repo) == 0 {
		return errors.New("'/repositories/sub/del/{owner}/{repo}': parameter 'repo' is empty")
	}

	ctx.Set("repo", repo)

	return nil
}
