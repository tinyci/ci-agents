package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/urfave/cli"

	"github.com/tinyci/ci-agents/ci-gen/gen/svc/uisvc/restapi"
)

// Version of API service
const Version = "1.0.0"

// TinyCIVersion is the global version of tinyCI
var TinyCIVersion = "" // to be changed by build processes

func main() {
	app := cli.NewApp()
	app.Name = "uisvc"
	app.Description = "API for the user interface service; the service that is directly responsible for presenting data to users.\nThis service typically runs at the border, and leverages session cookies or authentication tokens that we generate for users. It also is responsible for handling the act of oauth and user creation through its login hooks.\nuisvc typically talks to the datasvc and other services to accomplish its goal, it does not save anything locally or carry state.\n"
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
		if e, ok := err.(*errors.Error); ok && e == nil {
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve(ctx *cli.Context) error {
	h := &handlers.H{}
	if err := config.Parse(ctx.String("config"), &h); err != nil {
		return err
	}

	h.Config = restapi.MakeHandlerConfig(h.ServiceConfig)

	doneChan, err := handlers.Boot(nil, h)
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
