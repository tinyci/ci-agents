package cmdlib

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/config"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

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
func (s *GRPCServer) Make(commands []cli.Command) []cli.Command {
	return append(commands, cli.Command{
		Name:        s.Name,
		Usage:       s.Description,
		Description: s.Description,
		UsageText:   s.Name + " [flags]",
		Action:      s.serve,
	})
}

func (s *GRPCServer) serve(ctx *cli.Context) error {
	h := &handler.H{}
	if err := config.Parse(ctx.GlobalString("config"), &h); err != nil {
		return err
	}

	h.Name = s.Name
	h.UseDB = s.UseDB
	h.UseSessions = s.UseSessions

	cert, certErr := h.TLS.Load()
	if certErr != nil {
		return certErr
	}

	t, transportErr := transport.Listen(cert, "tcp", fmt.Sprintf(":%v", s.DefaultService.Port)) // FIXME parameterize
	if transportErr != nil {
		return transportErr
	}

	grpc, closer, err := h.CreateServer()
	if err != nil {
		return err
	}

	if err := s.RegisterService(grpc, h); err != nil {
		return err
	}

	finished := make(chan struct{})
	doneChan, err := h.Boot(t, grpc, finished)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		close(doneChan)
		<-finished
		if closer != nil {
			closer.Close()
		}
		os.Exit(0)
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {}
}
