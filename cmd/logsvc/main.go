package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/api/logsvc/processors"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/log"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

// Version is the version of this service.
const Version = "1.0.0"

// TinyCIVersion is the version of tinyci supporting this service.
var TinyCIVersion = "" // to be changed by build processes

func main() {
	app := cli.NewApp()
	app.Name = "logsvc"
	app.Description = "Centralized logging for tinyCI\n"
	app.Action = serve
	app.Version = fmt.Sprintf("%s (tinyCI version %s)", Version, TinyCIVersion)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path to configuration file",
			Value: ".config/services.yaml",
		},
	}

	if err := app.Run(os.Args); err != nil {
		errors.New(err).Exit()
	}
}

func serve(ctx *cli.Context) error {
	h := &handler.H{}
	if err := config.Parse(ctx.String("config"), &h); err != nil {
		return err
	}

	h.Name = "logsvc"

	cert, certErr := h.TLS.Load()
	if certErr != nil {
		return certErr
	}

	t, transportErr := transport.Listen(cert, "tcp", fmt.Sprintf(":%d", 6005)) // FIXME parameterize
	if transportErr != nil {
		return transportErr
	}

	s := grpc.NewServer()
	log.RegisterLogServer(s, processors.New(nil))

	doneChan, err := h.Boot(t, s)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		close(doneChan)
		os.Exit(0)
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {}
}
