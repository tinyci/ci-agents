package operations

import (
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"

	"github.com/gin-gonic/gin"
)

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// GetSubmitResponse responds for GetSubmit.
func GetSubmitResponse(h *handlers.H, ctx *gin.Context, resp interface{}, code int, err error) *errors.Error {
	var isError bool
	switch err := err.(type) {
	case *errors.Error:
		if err != nil {
			isError = true
		}
	case error:
		if err != nil {
			isError = true
		}
	default:
	}

	if isError {
		h.LogError(err, ctx, code)
		return errors.New(err)
	}

	ctx.JSON(code, resp)

	return nil
}
