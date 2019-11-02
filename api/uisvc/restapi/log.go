package restapi

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/clients/jsonbuffer"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"golang.org/x/net/websocket"
)

// LogAttach connects to a running logging process and outputs the return data as it arrives.
func LogAttach(pCtx context.Context, h *handlers.H, ctx *gin.Context, conn *websocket.Conn) error {
	id, err := strconv.ParseInt(ctx.GetString("id"), 10, 64)
	if err != nil {
		return errors.New(err)
	}

	if err := h.Clients.Asset.Read(ctx, id, jsonbuffer.NewWrapper(conn)); err != nil {
		return err
	}

	return nil
}
