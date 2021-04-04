package cmdlib

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/utils"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// ServerStatus represents multiple states within the server
type ServerStatus struct {
	// Alive is for closing when you want the server to stop.
	Alive chan struct{}
	// Finished is closed when the server has gracefully stopped.
	Finished chan struct{}
	// TracingCloser is the io.Closer that controls tracing.
	TracingCloser io.Closer
}

// HandlerFunc is a function that launches a service
type HandlerFunc func() (*ServerStatus, error)

// GRPCServer is the server to run
type GRPCServer struct {
	Name            string
	Description     string
	DefaultService  config.ServiceAddress
	RegisterService func(*grpc.Server, *handler.H) error
	UseDB           bool
	UseSessions     bool
}

// Make makes a command-line server out of the provided parameters
func (s *GRPCServer) Make(commands []*cli.Command) []*cli.Command {
	return append(commands, &cli.Command{
		Name:        s.Name,
		Usage:       s.Description,
		Description: s.Description,
		UsageText:   s.Name + " [flags]",
		Action:      s.serve,
	})
}

// MakeHandlerFunc returns a function with a channel to close to stop it, or any
// error it received while trying to create the server. It accepts a
// string to get the configuration filename.
func (s *GRPCServer) MakeHandlerFunc(configFile string) (HandlerFunc, error) {
	h := &handler.H{}
	if err := config.Parse(configFile, &h); err != nil {
		return nil, err
	}

	h.Name = s.Name
	h.UseDB = s.UseDB
	h.UseSessions = s.UseSessions

	cert, certErr := h.TLS.Load()
	if certErr != nil {
		return nil, certErr
	}

	t, transportErr := transport.Listen(cert, "tcp", fmt.Sprintf(":%v", s.DefaultService.Port)) // FIXME parameterize
	if transportErr != nil {
		return nil, transportErr
	}

	grpc, closer, err := h.CreateServer()
	if err != nil {
		return nil, err
	}

	if err := s.RegisterService(grpc, h); err != nil {
		return nil, err
	}

	return func() (*ServerStatus, error) {
		finished := make(chan struct{})
		doneChan, err := h.Boot(t, grpc, finished)
		if err != nil {
			return nil, err
		}

		return &ServerStatus{
			Finished:      finished,
			Alive:         doneChan,
			TracingCloser: closer,
		}, nil
	}, nil
}

func (s *GRPCServer) serve(ctx *cli.Context) error {
	fun, err := s.MakeHandlerFunc(ctx.String("config"))
	if err != nil {
		return utils.WrapError(err, "while constructing GRPC handler")
	}

	status, err := fun()
	if err != nil {
		return utils.WrapError(err, "while booting service")
	}

	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		close(status.Alive)
		<-status.Finished
		if status.TracingCloser != nil {
			status.TracingCloser.Close()
		}
		os.Exit(0)
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {}
}
