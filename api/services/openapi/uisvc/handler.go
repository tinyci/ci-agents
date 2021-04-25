package uisvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"

	"github.com/erikh/go-transport"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/legacy"
	gsessions "github.com/gorilla/sessions"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/tinyci/ci-agents/api/sessions"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/clients/log"
)

var errInvalidCookie = errors.New("cookie was invalid")

const (
	// SessionKey is the key used to identify the session in gorilla/sessions.
	SessionKey = "tinyci"
	// SessionUsername is the name of the session key that contains our username value.
	SessionUsername = "username"
)

// H is the stub handler for the service boot.
type H struct {
	Config      config.UserConfig
	clients     *config.Clients
	ServiceName string
	UseTLS      bool
	Swagger     *openapi3.Swagger
	Port        uint16
}

// Boot boots the service
func (h *H) Boot(finished chan struct{}) (chan struct{}, error) {
	t, err := h.createTransport()
	if err != nil {
		return nil, err
	}

	srv, err := h.createServer()
	if err != nil {
		return nil, err
	}

	s, l, err := t.Server(fmt.Sprintf(":%d", h.Port), srv.Handler)
	if err != nil {
		return nil, err
	}

	doneChan := make(chan struct{})

	go func() {
		<-doneChan
		s.Close()
		l.Close()
		h.clients.CloseClients()
		close(finished)
	}()

	go func() {
		if err := s.Serve(l); err != nil {
			h.clients.Log.Error(context.Background(), err)
		}
	}()

	return doneChan, nil
}

func (h *H) customHTTPErrorHandler(err error, ctx echo.Context) {
	go h.clients.Log.WithFields(log.FieldMap{"route": ctx.Request().URL.String()}).Error(context.Background(), err)

	uerr := uisvc.Error{Errors: &[]string{err.Error()}}

	if err := ctx.JSON(500, uerr); err != nil {
		go h.clients.Log.Error(context.Background(), err)
	}
}

// createServer creates the *echo.Server that powers the service.
func (h *H) createServer() (*http.Server, error) {
	var err error

	server := echo.New()

	h.clients, err = h.Config.ClientConfig.CreateClients(h.Config, h.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("while configuring GRPC backend clients: %w", err)
	}

	h.Swagger, err = uisvc.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("while configuring the openapi specification: %w", err)
	}

	h.Swagger.Servers = nil

	server.Use(h.echoLog())
	sessdb := sessions.New(h.clients.Data, nil, h.Config.Auth.ParsedSessionCryptKey())
	server.Use(h.echoSessions(sessdb))
	server.Use(h.echoUser())

	server.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		// FIXME defaults are to configure '*' for the origin bypass; which is pretty insecure.
		AllowHeaders: []string{echo.HeaderContentType},
	}))

	//server.Use(middleware.OapiRequestValidator(h.Swagger))
	server.Use(h.processOpenAPIExtensions())
	server.Use(h.echoCapability())
	server.Use(h.echoOAuthScopes())

	server.HTTPErrorHandler = h.customHTTPErrorHandler
	uisvc.RegisterHandlers(server, h)

	return server.Server, nil
}

// CreateTransport creates a transport with optional certification information.
func (h *H) createTransport() (*transport.HTTP, error) {
	if err := h.Config.TLS.Validate(); err != nil {
		return nil, err
	}

	var cert *transport.Cert
	if h.UseTLS {
		var err error
		cert, err = h.Config.TLS.Load()
		if err != nil {
			return nil, err
		}
	}

	t, err := transport.NewHTTP(cert)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (h *H) echoSessions(sessdb *sessions.SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			sess, err := sessdb.New(ctx.Request(), SessionKey)
			if err == nil && sess != nil {
				ctx.Set("tinyci.Session", sess)
			} else {
				go h.clients.Log.Errorf(context.Background(), "Error retrieving session: %v", err)
			}

			return next(ctx)
		}
	}
}

func (h *H) echoUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			u, err := h.findUser(ctx)
			if err != nil {
				if !errors.As(err, &utils.ErrInvalidAuth) {
					go h.clients.Log.Errorf(context.Background(), "Error retrieving user: %v", err)
				}

				return err
			} else if u != nil {
				ctx.Set("tinyci.User", u)
				ctx.Set("tinyci.Username", u.Username)
			}

			return next(ctx)
		}
	}
}

func (h *H) echoLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			fields := log.FieldMap{}
			u, ok := h.getUsername(ctx)
			if ok {
				fields["user"] = u
			}

			go h.clients.Log.WithFields(fields).Info(context.Background(), ctx.Request().URL)
			return next(ctx)
		}
	}
}

func (h *H) echoCapability() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cap, ok := ctx.Get("openapiExt.x-capability").(string)

			if ok {
				if u, ok := h.getUser(ctx); ok {
					caps, err := h.clients.Data.GetCapabilities(ctx.Request().Context(), u)
					if err != nil {
						return err // error fetching caps
					}

					for _, c := range caps {
						if string(c) == cap {
							return next(ctx) // cap matched
						}
					}
				}

				return utils.ErrInvalidAuth // cap not matched
			}

			return next(ctx) // no caps; this request can proceed
		}
	}
}

func (h *H) echoOAuthScopes() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			scope, ok := ctx.Get("openapiExt.x-token-scope").(string)
			if ok {
				if u, ok := h.getUser(ctx); ok {
					for _, s := range u.Token.Scopes {
						if s == scope {
							return next(ctx)
						}
					}
				}

				return utils.ErrInvalidAuth
			}

			return next(ctx)
		}
	}
}

func (h *H) processOpenAPIExtensions() echo.MiddlewareFunc {
	router, err := legacy.NewRouter(h.Swagger)
	if err != nil {
		panic(err)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			route, _, err := router.FindRoute(req)
			if err != nil {
				return err
			}

			op := route.PathItem.GetOperation(req.Method)
			if op == nil {
				return ctx.NoContent(500)
			}

			for key, value := range op.Extensions {
				var v interface{}
				switch value := value.(type) {
				case json.RawMessage:
					json.Unmarshal(value, &v)
				default:
					v = value
				}

				ctx.Set(fmt.Sprintf("openapiExt.%s", key), v)
			}

			return next(ctx)
		}
	}
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

// getUser retreives the user from the context, and returns true if it exists.
func (h *H) getUser(ctx echo.Context) (*model.User, bool) {
	u, ok := ctx.Get("tinyci.User").(*model.User)
	return u, ok
}

// getUsername retreives the user from the context, and returns true if it exists.
func (h *H) getUsername(ctx echo.Context) (string, bool) {
	u, ok := ctx.Get("tinyci.Username").(string)
	return u, ok
}

func (h *H) getSession(ctx echo.Context) (*gsessions.Session, bool) {
	sess, ok := ctx.Get("tinyci.Session").(*gsessions.Session)
	return sess, ok
}

// findUser retrieves the user based on information in the gin context.
func (h *H) findUser(ctx echo.Context) (*model.User, error) {
	var u *model.User

	req := ctx.Request()
	reqCtx := req.Context()

	if token := req.Header.Get("Authorization"); token != "" {
		if token != "" {
			var err error
			u, err = h.clients.Data.ValidateToken(reqCtx, token)
			if err != nil {
				return nil, err
			}
		}
	} else {
		var err error

		sess, ok := h.getSession(ctx)
		if !ok || sess == nil {
			return nil, nil
		}

		username, ok := sess.Values[SessionUsername].(string)
		if !ok {
			return nil, nil
		}

		u, err = h.clients.Data.GetUser(reqCtx, username)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

// GetClient returns a github client that works with the credentials in the given context.
func (h *H) getClient(ctx echo.Context) (github.Client, error) {
	user, err := h.getGithub(ctx)
	if err != nil {
		return nil, err
	}

	return h.Config.OAuth.GithubClient(user.Token.Username, user.Token.Token), nil
}

// GetGithub gets the github user from the session and loads it.
func (h *H) getGithub(ctx echo.Context) (u *model.User, outErr error) {
	sess, ok := h.getSession(ctx)
	if !ok {
		return nil, utils.ErrInvalidAuth
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
