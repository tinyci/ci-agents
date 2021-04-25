package uisvc

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/utils"
)

// GetLoginUpgrade upgrades the user's api keys.
func (h *H) GetLoginUpgrade(ctx echo.Context) error {
	return h.oauthRedirect(ctx, config.OAuthRepositoryScope)
}

// GetLoggedin handles the process of signaling javascript whether or not to login.
func (h *H) GetLoggedin(ctx echo.Context) error {
	res := "true"

	sess, ok := h.getSession(ctx)
	if !ok {
		// if there are any errors retrieving the session, just ask the user to relog
		url, err := h.clients.Auth.GetOAuthURL(ctx.Request().Context(), nil)
		if err != nil {
			return fmt.Errorf("while retrieving the oauth url: %w", err)
		}

		return ctx.JSON(200, url)
	}

	// guard against empty usernames as well as missing keys
	u, ok := sess.Values[SessionUsername]
	if su, sok := u.(string); !ok || (sok && strings.TrimSpace(su) == "") {
		url, err := h.clients.Auth.GetOAuthURL(ctx.Request().Context(), nil)
		if err != nil {
			return fmt.Errorf("while retrieving the oauth url: %w", err)
		}
		return ctx.JSON(200, url)
	}

	return ctx.JSON(200, res)
}

// GetLogout logs the user out of the tinyCI system.
func (h *H) GetLogout(ctx echo.Context) error {
	sess, ok := h.getSession(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	delete(sess.Values, SessionUsername)

	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return utils.WrapError(err, "could not persist session while logging out")
	}

	return ctx.Redirect(302, "/")
}

// GetLogin processes the oauth response and optionally redirects the user if not
// logged in already.
func (h *H) GetLogin(ctx echo.Context, params uisvc.GetLoginParams) error {
	oauthinfo, err := h.clients.Auth.OAuthChallenge(ctx.Request().Context(), params.State, params.Code)
	if err != nil {
		return fmt.Errorf("while generating oauth challenge url: %w", err)
	}

	if oauthinfo.Redirect {
		return ctx.Redirect(302, oauthinfo.Url)
	}

	sess, ok := h.getSession(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	sess.Values[SessionUsername] = oauthinfo.Username

	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return utils.WrapError(err, "could not persist session for %s", oauthinfo.Username)
	}

	return ctx.Redirect(302, "/")
}

// GetUserProperties gives an object containing information about the user.
func (h *H) GetUserProperties(ctx echo.Context) error {
	user, ok := h.getUser(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	caps, err := h.clients.Data.GetCapabilities(ctx.Request().Context(), user)
	if err != nil {
		return err
	}

	ret := map[string]interface{}{}
	ret["username"] = user.Username
	ret["capabilities"] = caps
	return ctx.JSON(200, ret)
}
