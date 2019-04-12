package handlers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	transport "github.com/erikh/go-transport"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	apiSess "github.com/tinyci/ci-agents/api/sessions"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"golang.org/x/net/websocket"
	"golang.org/x/oauth2"
)

// SessionUsername is the name of the session key that contains our username value.
const SessionUsername = "username"

// AllowOrigin changes the scope of CORS requests. By default they are
// insecure; this is intended to be overrided by commands as they set up this
// library along with their other standard init operations.
var AllowOrigin = "*"

var errInvalidCookie = errors.New("cookie was invalid")

var routeTransformer = regexp.MustCompile(`(?:{([^}]+)})+`)

// HandlerConfig provides an interface to managing the HandlerConfig.
type HandlerConfig interface {
	SetRoutes(*H)
	DBConfigure(*H) *errors.Error
	Configure(Routes) *errors.Error
	CustomInit(*H) *errors.Error
	Validate(*H) *errors.Error
}

// H is a series of HTTP handlers for the UI service
type H struct {
	Config HandlerConfig `yaml:"-"`
	Routes Routes        `yaml:"-"`

	config.UserConfig `yaml:",inline"`
	config.Service
}

// GithubClient is a wrapper for config.GithubClient.
func (h *H) GithubClient(token *oauth2.Token) github.Client {
	return h.OAuth.GithubClient(token)
}

// CreateClients creates the clients to be used based on configuration values.
func (h *H) CreateClients() *errors.Error {
	var err *errors.Error
	h.Clients, err = h.ClientConfig.CreateClients(h.Name)

	return err
}

// CreateTransport creates a transport with optional certification information.
func (h *H) CreateTransport() (*transport.HTTP, *errors.Error) {
	var cert *transport.Cert
	if !h.NoTLSServer {
		var err *errors.Error
		cert, err = h.TLS.Load()
		if err != nil {
			return nil, err
		}
	}

	t, err := transport.NewHTTP(cert)
	if err != nil {
		return nil, errors.New(err)
	}

	return t, nil
}

// Boot boots the service. Closing the channel returned will shutdown the service.
func Boot(t *transport.HTTP, handler *H) (chan struct{}, *errors.Error) {
	handler.Formats = strfmt.NewFormats()

	if err := handler.Init(); err != nil {
		return nil, err
	}

	r, err := handler.CreateRouter()
	if err != nil {
		return nil, err
	}

	if handler.Auth.NoAuth {
		if config.DefaultGithubClient == nil {
			config.DefaultGithubClient = handler.Auth.GetNoAuthClient()
		}

		var err *errors.Error
		handler.DefaultUsername, err = config.DefaultGithubClient.MyLogin()
		if err != nil {
			return nil, err
		}
	}

	if t == nil {
		var err *errors.Error
		t, err = handler.CreateTransport()
		if err != nil {
			return nil, err
		}
	}

	var sErr error
	s, l, sErr := t.Server(fmt.Sprintf(":%d", handler.Port), r)
	if err != nil {
		return nil, errors.New(sErr)
	}

	s.IdleTimeout = 30 * time.Second
	doneChan := make(chan struct{})

	go func() {
		<-doneChan
		l.Close()
		s.Close()
	}()

	go s.Serve(l)
	return doneChan, nil
}

// Init initialize the handler and makes it available for requests.
func (h *H) Init() *errors.Error {
	if err := h.Config.Validate(h); err != nil {
		return err
	}

	if err := h.Config.CustomInit(h); err != nil {
		return err
	}

	if err := h.CreateClients(); err != nil {
		return err
	}

	h.Config.SetRoutes(h)

	if err := h.Config.Configure(h.Routes); err != nil {
		return err
	}

	if err := h.Auth.Validate(h.UseSessions); err != nil {
		return err
	}

	return h.dbConnect()
}

func (h *H) dbConnect() *errors.Error {
	if h.UseDB {
		var err *errors.Error
		h.Model, err = model.New(h.DSN)
		if err != nil {
			return err
		}
		return h.Config.DBConfigure(h)
	}

	return nil
}

// CORS primes OPTIONS and normal requests with the appropriate headers and
// also acts like a normal http.Handler so it can be used that way.
func CORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", AllowOrigin)
	ctx.Header("Access-Control-Allow-Methods", "PUT,GET,POST,DELETE")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
}

// GetGithub gets the github user from the session and loads it.
func (h *H) GetGithub(ctx *gin.Context) (*model.User, *errors.Error) {
	sess := sessions.Default(ctx)

	uname, ok := sess.Get(SessionUsername).(string)
	if ok && strings.TrimSpace(uname) != "" {
		// no error, we're already logged in
		return h.Clients.Data.GetUser(uname)
	}

	token := ctx.Request.Header.Get("Authorization")
	if token != "" {
		username, err := h.Clients.Data.ValidateToken(token)
		if err != nil {
			return nil, err
		}

		return h.Clients.Data.GetUser(username)
	}

	return nil, errInvalidCookie
}

func (h *H) authed(gatewayFunc func(*H, *gin.Context, HandlerFunc) *errors.Error) func(h *H, ctx *gin.Context, processor HandlerFunc) *errors.Error {
	return func(h *H, ctx *gin.Context, processor HandlerFunc) *errors.Error {
		if !h.Auth.NoAuth {
			token := ctx.Request.Header.Get("Authorization")
			if token != "" {
				if _, err := h.Clients.Data.ValidateToken(token); err != nil {
					return err
				}
			} else {
				_, err := h.GetGithub(ctx)
				if err != nil {
					return errors.New(err)
				}
			}
		}

		return gatewayFunc(h, ctx, processor)
	}
}

func (h *H) inWebsocket(paramHandler func(*H, *gin.Context) *errors.Error, handler func(h *H, ctx *gin.Context, conn *websocket.Conn) *errors.Error) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		outerHandler := func(conn *websocket.Conn) {
			if err := paramHandler(h, ctx); err != nil {
				conn.Close()
			}

			if err := handler(h, ctx, conn); err != nil {
				conn.Close()
			}
		}

		websocket.Handler(outerHandler).ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// TransformSwaggerRoute merely translates url params from {thisformat} to :thisformat
func TransformSwaggerRoute(route string) string {
	return routeTransformer.ReplaceAllStringFunc(route, func(input string) string {
		return string(routeTransformer.ExpandString([]byte{}, ":$1", input, routeTransformer.FindStringSubmatchIndex(input)))
	})
}

func (h *H) configureSessions(r *gin.Engine) *errors.Error {
	sessdb := apiSess.New(h.Clients.Data, nil, h.Auth.ParsedSessionCryptKey())
	r.Use(sessions.Sessions(config.SessionKey, sessdb))

	return nil
}

func (h *H) configureRestHandler(r *gin.Engine, key string, route *Route, optionsRoutes map[string]struct{}) {

	var dispatchFunc func(string, ...gin.HandlerFunc) gin.IRoutes

	switch route.Method {
	case "GET":
		dispatchFunc = r.GET
	case "POST":
		dispatchFunc = r.POST
	case "DELETE":
		dispatchFunc = r.DELETE
	case "PATCH":
		dispatchFunc = r.PATCH
	case "PUT":
		dispatchFunc = r.PUT
	case "OPTIONS":
		dispatchFunc = r.OPTIONS
	case "HEAD":
		dispatchFunc = r.HEAD
	}

	var handler func(*H, *gin.Context, HandlerFunc) *errors.Error = route.Handler

	if route.UseAuth {
		handler = h.authed(handler)
	}

	if route.UseCORS {
		if _, ok := optionsRoutes[key]; !ok {
			r.OPTIONS(key, CORS)
			optionsRoutes[key] = struct{}{}
		}
	}
	dispatchFunc(key, h.wrapHandler(handler, route.Processor))
}

// CreateRouter creates a *mux.Router capable of serving the UI server.
func (h *H) CreateRouter() (*gin.Engine, *errors.Error) {
	r := gin.New()

	if h.UseSessions {
		if err := h.configureSessions(r); err != nil {
			return nil, err
		}
	}

	optionsRoutes := map[string]struct{}{}

	for key, methodRoutes := range h.Routes {
		for _, route := range methodRoutes {
			if route.WebsocketProcessor != nil {
				r.GET(key, h.inWebsocket(route.ParamValidator, route.WebsocketProcessor))
			} else {
				h.configureRestHandler(r, key, route, optionsRoutes)
			}
		}
	}

	return r, nil
}

func (h *H) wrapHandler(handler func(*H, *gin.Context, HandlerFunc) *errors.Error, processor HandlerFunc) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		err := handler(h, ctx, processor)
		if err != nil {
			h.WriteError(ctx, err)
		}
	}
}

// WriteError standardizes the writing of error states for easier typing. It is
// not intended to be used to write specific statuses, only 500 errors with JSON output.
// If UseSessions is on, it will populate the errors session store.
func (h *H) WriteError(ctx *gin.Context, err *errors.Error) {
	ctx.AbortWithStatusJSON(500, err)
}

// Session returns the current user session.
func (h *H) Session(ctx *gin.Context) sessions.Session {
	return sessions.Default(ctx)
}

// LogError logs an HTTP error to the client.
func (h *H) LogError(err error, ctx *gin.Context, code int) {
	logger := h.Clients.Log.WithRequest(ctx.Request).WithFields(log.FieldMap{"code": fmt.Sprintf("%v", code)})
	user, gitErr := h.GetGithub(ctx)
	if gitErr == nil {
		logger = logger.WithUser(user)
	}

	content, jsonErr := json.Marshal(ctx.Params)
	if jsonErr != nil {
		logger.Error(errors.New(jsonErr).Wrap("encoding params for log message"))
	}

	var doLog bool

	switch err := err.(type) {
	case *errors.Error:
		if err.GetLog() {
			doLog = true
		}
	default:
		doLog = true
	}

	if doLog {
		logger.WithFields(log.FieldMap{"params": string(content)}).Error(err)
	}
}
