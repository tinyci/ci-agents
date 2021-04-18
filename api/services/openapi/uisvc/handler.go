package uisvc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
)

var errInvalidCookie = errors.New("cookie was invalid")

const (
	// SessionKey is the key used to identify the session in gorilla/sessions.
	SessionKey = "tinyci"
	// SessionUsername is the name of the session key that contains our username value.
	SessionUsername = "username"
)

// H is the top-level handler for all uisvc methods.
type H struct {
	config  config.UserConfig
	clients *config.Clients
}

// NewHandler creates a new uisvc top-level handler
func NewHandler(config config.UserConfig, clients *config.Clients) *H {
	return &H{config: config, clients: clients}
}

// OAuthRedirect redirects the user to the oauth confirmation screen, requesting the additional scopes.
func (h *H) oauthRedirect(ctx echo.Context, scopes []string) error {
	url, err := h.clients.Auth.GetOAuthURL(ctx.Request().Context(), scopes)
	if err != nil {
		return err
	}

	ctx.Redirect(302, url)
	return nil
}

func (h *H) getSession(ctx echo.Context) (*sessions.Session, error) {
	return session.Get(SessionKey, ctx)
}

// GetUser retrieves the user based on information in the gin context.
func (h *H) getUser(ctx echo.Context) (*model.User, error) {
	var u *model.User

	req := ctx.Request()
	reqCtx := req.Context()

	if token := ctx.Request().Header.Get("Authorization"); token != "" {
		if token != "" {
			var err error
			u, err = h.clients.Data.ValidateToken(reqCtx, token)
			if err != nil {
				return nil, err
			}
		}
	} else {
		sess, err := h.getSession(ctx)
		if err != nil {
			return nil, fmt.Errorf("While retrieving session: %w", err)
		}

		username, ok := sess.Values[SessionUsername].(string)
		if !ok {
			return nil, utils.ErrInvalidAuth
		}

		u, err = h.clients.Data.GetUser(reqCtx, username)
		if err != nil {
			return nil, err
		}
	}

	if u == nil {
		return nil, utils.ErrInvalidAuth
	}

	return u, nil
}

// GetClient returns a github client that works with the credentials in the given context.
func (h *H) getClient(ctx echo.Context) (github.Client, error) {
	user, err := h.getGithub(ctx)
	if err != nil {
		return nil, err
	}

	return h.config.OAuth.GithubClient(user.Token.Username, user.Token.Token), nil
}

// GetGithub gets the github user from the session and loads it.
func (h *H) getGithub(ctx echo.Context) (u *model.User, outErr error) {
	sess, err := h.getSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("while retrieving github information about the user: %w", err)
	}

	defer func() {
		if outErr != nil {
			if sess != nil {
				sess.Values = map[interface{}]interface{}{}
				sess.Save(ctx.Request(), ctx.Response())
			}
		}
	}()

	reqCtx := ctx.Request().Context()

	uname, ok := sess.Values[SessionUsername].(string)
	if ok && strings.TrimSpace(uname) != "" {
		// no error, we're already logged in
		return h.clients.Data.GetUser(reqCtx, uname)
	}

	token := ctx.Request().Header.Get("Authorization")
	if token != "" {
		return h.clients.Data.ValidateToken(reqCtx, token)
	}

	return nil, errInvalidCookie
}
