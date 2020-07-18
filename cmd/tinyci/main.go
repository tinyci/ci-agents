package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tinyci/ci-agents/api/hooksvc"
	"github.com/tinyci/ci-agents/api/uisvc/restapi"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/urfave/cli"
)

// TinyCIVersion is the version of tinyci supporting this service.
var TinyCIVersion = "" // to be changed by build processes

func main() {
	app := cli.NewApp()
	app.Version = TinyCIVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Path to configuration file",
			Value:  ".config/services.yaml",
			EnvVar: "TINYCI_CONFIG",
		},
	}

	otherServices := []cli.Command{
		{
			Name:        "hooksvc",
			Usage:       "manage incoming github submissions",
			Description: "manage incoming github submissions",
			Action:      startHooksvc,
		},
		{
			Name:        "uisvc",
			Usage:       "API for the user interface service",
			Description: "API for the user interface service; the service that is directly responsible for presenting data to users.\nThis service typically runs at the border, and leverages session cookies or authentication tokens that we generate for users. It also is responsible for handling the act of oauth and user creation through its login hooks.\nuisvc typically talks to the datasvc and other services to accomplish its goal, it does not save anything locally or carry state.",
			Action:      startUISvc,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "service",
			ShortName:   "s",
			Usage:       "Launch services that power tinyCI",
			Description: "Launch services that power tinyCI",
			Subcommands: mapServers(otherServices),
		},
	}

	if err := app.Run(os.Args); err != nil {
		errors.New(err).Exit()
	}
}

func mapServers(commands []cli.Command) []cli.Command {
	for _, s := range servers {
		commands = s.Make(commands)
	}

	return commands
}

func startHooksvc(ctx *cli.Context) error {
	h := &hooksvc.Handler{}

	if err := config.Parse(ctx.GlobalString("config"), &h.Config); err != nil {
		return errors.New(err)
	}

	if err := h.Init(); err != nil {
		return errors.New(err)
	}

	http.Handle("/hook", h)
	if err := http.ListenAndServe(config.DefaultServices.Hook.String(), http.DefaultServeMux); err != nil {
		return errors.New(err)
	}

	return nil
}

func startUISvc(ctx *cli.Context) error {
	h := &handlers.H{}
	if err := config.Parse(ctx.GlobalString("config"), &h); err != nil {
		return err
	}

	h.Config = restapi.MakeHandlerConfig(h.ServiceConfig)

	finished := make(chan struct{})
	doneChan, err := handlers.Boot(nil, h, finished)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		close(doneChan)
		<-finished
		os.Exit(0)
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {}
}