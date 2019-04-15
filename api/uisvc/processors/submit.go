package processors

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/types"
)

// Submit powers a manual submission to the queuesvc.
func Submit(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	repo := ctx.GetString("repository")
	sha := ctx.GetString("sha")

	all, perr := strconv.ParseBool(ctx.GetString("all"))
	if perr != nil {
		all = false
	}

	user, err := getUser(h, ctx)
	if err != nil {
		return nil, 500, err
	}

	err = h.Clients.Queue.Submit(&types.Submission{
		Fork:        repo,
		HeadSHA:     sha,
		SubmittedBy: user.Username,
		All:         all,
		Manual:      true,
	})
	if err != nil {
		return nil, 500, err
	}

	return nil, 200, nil
}
