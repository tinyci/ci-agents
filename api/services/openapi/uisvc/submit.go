package uisvc

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// GetSubmit powers a manual submission to the queuesvc.
func (h *H) GetSubmit(ctx echo.Context, params uisvc.GetSubmitParams) error {
	repo := params.Repository
	sha := params.Sha
	all := false

	if params.All != nil {
		all = *params.All
	}

	user, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	err := h.clients.Queue.Submit(context.Background(), &types.Submission{
		Fork:        repo,
		HeadSHA:     sha,
		SubmittedBy: user,
		All:         all,
		Manual:      true,
	})
	if err != nil {
		return err
	}

	return ctx.NoContent(200)
}
