package openapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/erikh/go-transport"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/tinyci/ci-agents/api/sessions"
	"github.com/tinyci/ci-agents/config"
)

// RegisterFunc is utilized to register the handler with the echo service.
type RegisterFunc func(h *H, server *echo.Echo)

// H is the stub handler for the service boot.
type H struct {
	Config       config.UserConfig
	Clients      *config.Clients
	ServiceName  string
	UseTLS       bool
	Swagger      *openapi3.Swagger
	RegisterFunc RegisterFunc
	Port         uint16
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

	s.IdleTimeout = 30 * time.Second
	doneChan := make(chan struct{})

	go func() {
		<-doneChan
		s.Close()
		l.Close()
		h.Clients.CloseClients()
		close(finished)
	}()

	go func() {
		if err := s.Serve(l); err != nil {
			h.Clients.Log.Error(context.Background(), err)
		}
	}()

	return doneChan, nil
}

// createServer creates the *echo.Server that powers the service.
func (h *H) createServer() (*http.Server, error) {
	var err error

	server := echo.New()

	h.Clients, err = h.Config.ClientConfig.CreateClients(h.Config, h.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("while configuring GRPC backend clients: %w", err)
	}

	server.Use(echomiddleware.Logger())

	server.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		// FIXME defaults are to configure '*' for the origin bypass; which is pretty insecure.
		AllowHeaders: []string{echo.HeaderContentType},
	}))

	sessdb := sessions.New(h.Clients.Data, nil, h.Config.Auth.ParsedSessionCryptKey())
	server.Use(session.Middleware(sessdb))

	//server.Use(middleware.OapiRequestValidator(h.Swagger))
	server.Use(h.processOpenAPIExtensions())

	h.RegisterFunc(h, server)

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
				ctx.Set(fmt.Sprintf("openapiExt.%s", key), value)
			}

			return next(ctx)
		}
	}
}
