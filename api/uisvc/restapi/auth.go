package restapi

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// Upgrade upgrades the user's api keys.
func Upgrade(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	handlers.CORS(ctx)
	if err := h.OAuthRedirect(ctx, config.OAuthRepositoryScope); err != nil {
		return nil, 500, err
	}

	return nil, 302, nil
}

// LoggedIn handles the process of signaling javascript whether or not to login.
func LoggedIn(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	handlers.CORS(ctx)

	res := "true"

	_, err := h.GetGithub(ctx)
	if err != nil {
		var err *errors.Error
		res, err = h.Clients.Auth.GetOAuthURL(ctx, nil)
		if err != nil {
			return nil, 500, err
		}
	}

	return res, 200, nil
}

// Logout logs the user out of the tinyCI system.
func Logout(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	sess := sessions.Default(ctx)
	sess.Delete(handlers.SessionUsername)

	if err := sess.Save(); err != nil {
		return nil, 500, errors.New(err).Wrap("could not persist session while logging out")
	}

	ctx.Redirect(302, "/")
	return nil, 302, nil
}

// Login processes the oauth response and optionally redirects the user if not
// logged in already.
func Login(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	oauthinfo, err := h.Clients.Auth.OAuthChallenge(ctx, ctx.Query("state"), ctx.Query("code"))
	if err != nil {
		return nil, 500, err
	}

	if oauthinfo.Redirect {
		ctx.Redirect(302, oauthinfo.Url)
		return nil, 302, nil
	}

	sess := sessions.Default(ctx)
	sess.Set(handlers.SessionUsername, oauthinfo.Username)
	err2 := sess.Save()
	if err2 != nil {
		return nil, 500, errors.New(err2).Wrapf("could not persist session for %s", oauthinfo.Username)
	}

	ctx.Redirect(302, "/")

	return nil, 302, nil
}

// GetUserProperties gives an object containing information about the user.
func GetUserProperties(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	caps, err := h.Clients.Data.GetCapabilities(ctx, user)
	if err != nil {
		return nil, 500, err
	}

	ret := map[string]interface{}{}
	ret["username"] = user.Username
	ret["capabilities"] = caps
	return ret, 200, nil
}
