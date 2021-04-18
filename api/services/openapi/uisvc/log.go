package uisvc

import (
	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/clients/jsonbuffer"
	"golang.org/x/net/websocket"
)

// GetLogAttachId connects to a running logging process and outputs the return data as it arrives.
func (h *H) GetLogAttachId(ctx echo.Context, id int64) error {
	var retErr error
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		retErr = h.clients.Asset.Read(ctx.Request().Context(), id, jsonbuffer.NewWriteWrapper(ws))
	}).ServeHTTP(ctx.Response(), ctx.Request())

	return retErr
}
