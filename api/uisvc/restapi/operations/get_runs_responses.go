package operations

import (
	"github.com/tinyci/ci-agents/handlers"

	"github.com/gin-gonic/gin"
)

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// GetRunsResponse responds for GetRuns.
func GetRunsResponse(h *handlers.H, ctx *gin.Context, resp interface{}, code int, err error) error {
	if err != nil {
		h.LogError(err, ctx, code)
		return err
	}

	ctx.JSON(code, resp)

	return nil
}
