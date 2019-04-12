package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers/hooksvc"
	"github.com/urfave/cli"
)

const listenPort = 2020

func main() {
	app := cli.NewApp()
	app.Description = "manage incoming github submissions"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: ".config/hooksvc.yaml",
		},
	}

	app.Action = serve

	err := app.Run(os.Args)
	if e, _ := err.(*errors.Error); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}

func serve(ctx *cli.Context) error {
	h := &hooksvc.Handler{}

	if err := config.Parse(ctx.String("config"), &h.Config); err != nil {
		return errors.New(err)
	}

	if err := h.Init(); err != nil {
		return errors.New(err)
	}

	http.Handle("/hook", h)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", listenPort), http.DefaultServeMux); err != nil {
		return errors.New(err)
	}

	return nil
}
