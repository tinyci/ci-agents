package uisvc

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h *H) getSession(ctx echo.Context) (*sessions.Session, error) {
	return session.Get(SessionKey, ctx)
}
