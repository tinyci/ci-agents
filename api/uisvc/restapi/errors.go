package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// Errors processes the /errors GET endpoint
func Errors(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	errs, err := h.Clients.Data.GetErrors(user.Username)
	if err != nil && !err.Contains(errors.New("record not found")) {
		return nil, 500, err
	}

	for _, err := range errs {
		if err := h.Clients.Data.DeleteError(err.ID, user.ID); err != nil && !err.Contains(errors.New("record not found")) {
			return nil, 500, err
		}
	}

	return errs, 200, nil
}
