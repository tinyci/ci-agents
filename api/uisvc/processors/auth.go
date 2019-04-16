package processors

import (
	"context"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/model"
	topUtils "github.com/tinyci/ci-agents/utils"
	"golang.org/x/oauth2"
)

const sessionUsername = "username"

var (
	errRedirect = errors.New("redirection")
)

func getUser(h *handlers.H, ctx *gin.Context) (*model.User, *errors.Error) {
	client, err := getClient(h, ctx)
	if err != nil {
		return nil, err
	}

	var name string
	sess := h.Session(ctx)

	username := sess.Get(sessionUsername)

	// FIXME clean up this spaghetti -erikh
	var u *model.User

	if username == nil {
		if token := ctx.Request.Header.Get("Authorization"); token != "" {
			token := ctx.Request.Header.Get("Authorization")
			if token != "" {
				u, err = h.Clients.Data.ValidateToken(token)
				if err != nil {
					return nil, err
				}
			}
		} else {
			var err *errors.Error
			name, err = client.MyLogin()
			if err != nil {
				return nil, errors.New(err)
			}

			sess.Set(sessionUsername, name)
			if err := sess.Save(); err != nil {
				return nil, errors.New(err)
			}
		}
	} else {
		var ok bool
		name, ok = username.(string)
		if !ok {
			return nil, errors.ErrInvalidAuth
		}
	}

	if u == nil && name != "" {
		u, err = h.Clients.Data.GetUser(name)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func getClient(h *handlers.H, ctx *gin.Context) (github.Client, *errors.Error) {
	user, err := h.GetGithub(ctx)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}

	if err := topUtils.JSONIO(user.Token, token); err != nil {
		return nil, err
	}

	return h.GithubClient(token), nil
}

func handleOAuth(h *handlers.H, code string) (*oauth2.Token, string, *errors.Error) {
	conf := h.OAuth.Config()

	tok, err := conf.Exchange(context.Background(), code)
	if err != nil {
		switch err.(type) {
		case *oauth2.RetrieveError:
			return nil, "", errors.New(err)
		default:
			h.Clients.Log.Error(err)
			return nil, "", errRedirect
		}
	}

	client := conf.Client(context.Background(), tok)
	c := gh.NewClient(client)
	u, _, err := c.Users.Get(context.Background(), "")
	if err != nil {
		return nil, "", errors.New(err)
	}

	return tok, u.GetLogin(), nil
}

func getOAuthURL(h *handlers.H, ctx *gin.Context) (string, *errors.Error) {
	conf := h.OAuth.Config()

	state := strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)), "=")
	// FIXME finish
	if err := h.Clients.Data.OAuthRegisterState(state); err != nil {
		return "", err
	}

	return conf.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	), nil
}

func oauthRedirect(h *handlers.H, ctx *gin.Context) *errors.Error {
	url, err := getOAuthURL(h, ctx)
	if err != nil {
		return err
	}

	ctx.Redirect(302, url)
	return nil
}

// LoggedIn handles the process of signaling javascript whether or not to login.
func LoggedIn(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	handlers.CORS(ctx)

	res := "true"

	_, err := h.GetGithub(ctx)
	if err != nil {
		var err *errors.Error
		res, err = getOAuthURL(h, ctx)
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
	tok, username, err := handleOAuth(h, ctx.Query("code"))
	if err != nil {
		switch err {
		case errRedirect:
			return nil, 302, oauthRedirect(h, ctx)
		default:
			return nil, 500, err
		}
	}

	if err := h.Clients.Data.OAuthValidateState(ctx.Query("state")); err != nil {
		return nil, 500, err
	}

	user, err := h.Clients.Data.GetUser(username)
	if err != nil {
		var createErr *errors.Error
		_, createErr = h.Clients.Data.PutUser(&model.User{Username: username, Token: tok})
		if createErr != nil {
			return nil, 500, errors.New(fmt.Sprintf("could not create (%v) or read (%v) user %s after oauth challenge", createErr, err, username))
		}
	} else {
		user.Token = tok
		if err := h.Clients.Data.PatchUser(user); err != nil {
			return nil, 500, errors.New(err).Wrapf("could not update oauth token for %s", username)
		}
	}

	sess := sessions.Default(ctx)
	sess.Set(handlers.SessionUsername, username)
	err2 := sess.Save()
	if err2 != nil {
		return nil, 500, errors.New(err2).Wrapf("could not persist session for %s", username)
	}

	ctx.Redirect(302, "/")

	return nil, 302, nil
}

// GetUserProperties gives an object containing information about the user.
func GetUserProperties(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := getUser(h, ctx)
	if err != nil {
		return nil, 500, err
	}

	ret := map[string]string{}
	ret["username"] = user.Username
	return ret, 200, nil
}
