package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tinyci/ci-agents/api/services/hybrid/hooksvc"
	httpHandler "github.com/tinyci/ci-agents/api/services/openapi/uisvc"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/cmdlib"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/db"
	"github.com/tinyci/ci-agents/db/migrations"
	"github.com/tinyci/ci-agents/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/unix"
)

// TinyCIVersion is the version of tinyci supporting this service.
var TinyCIVersion = "" // to be changed by build processes

func main() {
	app := cli.NewApp()
	app.Version = TinyCIVersion

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config, c",
			Usage:   "Path to configuration file",
			Value:   ".config/services.yaml",
			EnvVars: []string{"TINYCI_CONFIG"},
		},
	}

	otherServices := []*cli.Command{
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

	app.Commands = []*cli.Command{
		{
			Name:        "migrate",
			Usage:       "Migrate the database to the latest version of tinyCI; should be done once when upgrading",
			Description: "Migrate the database to the latest version of tinyCI; should be done once when upgrading",
			Action:      migrate,
		},
		{
			Name:        "service",
			Aliases:     []string{"s"},
			Usage:       "Launch services that power tinyCI",
			Description: "Launch services that power tinyCI",
			Subcommands: mapServers(otherServices),
		},
		{
			Name:        "launch",
			Aliases:     []string{"l"},
			Usage:       "Launch all services to power tinyCI",
			Description: "Launch all services to power tinyCI",
			Action:      launch,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mapServers(commands []*cli.Command) []*cli.Command {
	for _, s := range servers {
		commands = s.Make(commands)
	}

	return commands
}

func migrate(ctx *cli.Context) error {
	configFile := ctx.String("config")
	conf := config.UserConfig{}

	if err := config.Parse(configFile, &conf); err != nil {
		return err
	}

	sqlDB, err := db.NewConn(conf.DSN)
	if err != nil {
		return err
	}

	return migrations.Migrate(sqlDB)
}

func launch(ctx *cli.Context) error {
	configFile := ctx.String("config")
	handlers := []cmdlib.HandlerFunc{}

	for _, s := range servers {
		handler, err := s.MakeHandlerFunc(configFile)
		if err != nil {
			return utils.WrapError(err, "while constructing handler for %s", s.Name)
		}

		handlers = append(handlers, handler)
	}

	uisvc, err := makeUISvcHandler(configFile)
	if err != nil {
		return err
	}
	handlers = append(handlers, uisvc)

	statuses := []*cmdlib.ServerStatus{}

	for _, h := range handlers {
		status, err := h()
		if err != nil {
			return err
		}

		statuses = append(statuses, status)
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, unix.SIGINT, unix.SIGTERM)

	<-sigChan
	for _, status := range statuses {
		close(status.Alive)
		<-status.Finished
		if status.TracingCloser != nil {
			status.TracingCloser.Close()
		}
	}

	return nil
}

func startHooksvc(ctx *cli.Context) error {
	h := &hooksvc.Handler{}

	if err := config.Parse(ctx.String("config"), &h.Config); err != nil {
		return err
	}

	if err := h.Init(); err != nil {
		return err
	}

	http.Handle("/hook", h)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.DefaultServices.Hook.Port), http.DefaultServeMux); err != nil {
		return err
	}

	return nil
}

func makeUISvcHandler(configFile string) (cmdlib.HandlerFunc, error) {
	h := &httpHandler.H{
		Port:        6010,
		ServiceName: "uisvc",
	}
	if err := config.Parse(configFile, &h.Config); err != nil {
		return nil, err
	}

	if err := h.Config.Auth.Validate(true); err != nil {
		return nil, err
	}

	var err error
	h.Swagger, err = uisvc.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("while configuring the openapi specification: %w", err)
	}

	h.Swagger.Servers = nil
	h.Swagger.Security = nil

	return func() (*cmdlib.ServerStatus, error) {
		finished := make(chan struct{})

		doneChan, err := h.Boot(finished)
		if err != nil {
			return nil, err
		}

		return &cmdlib.ServerStatus{
			Alive:    doneChan,
			Finished: finished,
		}, nil
	}, nil

}

func startUISvc(ctx *cli.Context) error {
	fun, err := makeUISvcHandler(ctx.String("config"))
	if err != nil {
		return err
	}

	status, err := fun()
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 2)
	go func() {
		<-sigChan
		close(status.Alive)
		<-status.Finished
		os.Exit(0)
	}()
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {}
}
