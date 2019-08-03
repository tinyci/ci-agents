package main

import (
	"fmt"
	"os"
	"time"

	transport "github.com/erikh/go-transport"
	"github.com/sirupsen/logrus"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/urfave/cli"
)

// Version is the version of this service.
const Version = "1.0.0"

const walkIncrement = 100

// TinyCIVersion is the version of tinyci supporting this service.
var TinyCIVersion = "" // to be changed by build processes

func main() {
	app := cli.NewApp()
	app.Name = "cancelbot"
	app.Description = "cancelbot is a cron-based observer for canceling dangling jobs"
	app.Action = run
	app.Version = fmt.Sprintf("%s (tinyCI version %s)", Version, TinyCIVersion)

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run, n",
			Usage: "Just print what runs would be canceled, but don't do anything",
		},
		cli.DurationFlag{
			Name:  "timeout, t",
			Usage: "After this time, cancel the run",
			Value: 3 * time.Hour,
		},
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "Limit to last N runs: set to 0 to not limit",
			Value: 1000,
		},
		cli.StringFlag{
			Name:  "datasvc, d",
			Usage: "Location of datasvc",
			Value: config.DefaultServices.Data.String(),
		},
		cli.StringFlag{
			Name:  "cacert, ca",
			Usage: "Location of CA certificate for encrypted connections",
		},
		cli.StringFlag{
			Name:  "cert, c",
			Usage: "Client cert used to connect to datasvc",
		},
		cli.StringFlag{
			Name:  "key, k",
			Usage: "Client key used to connect to datasvc",
		},
	}

	if err := app.Run(os.Args); err != nil {
		errors.New(err).Exit()
	}
}

func run(ctx *cli.Context) error {
	var cert *transport.Cert

	if !(ctx.GlobalString("cacert") == "" && ctx.GlobalString("cert") == "" && ctx.GlobalString("key") == "") {
		var err error
		// last arg is CRL
		cert, err = transport.LoadCert(ctx.GlobalString("cacert"), ctx.GlobalString("cert"), ctx.GlobalString("key"), "")
		if err != nil {
			return errors.New(err).Wrap("while loading cert")
		}
	}

	client, err := data.New(ctx.GlobalString("datasvc"), cert, false)
	if err != nil {
		return err
	}
	defer client.Close()

	for count := ctx.GlobalInt("limit"); count >= 0; count -= walkIncrement {
		runs, err := client.ListRuns("", "", int64(count/walkIncrement), walkIncrement)
		if err != nil {
			return err
		}

		for _, run := range runs {
			if run.Status == nil && time.Since(run.CreatedAt) > ctx.GlobalDuration("timeout") {
				if ctx.GlobalBool("dry-run") {
					logrus.Infof("Would cancel run %d, repository %v, ref %v, name %v -- %v old", run.ID, run.Task.Parent.Name, run.Task.Ref.RefName, run.Name, time.Since(run.CreatedAt))
				} else {
					if err := client.SetCancel(run.ID); err != nil {
						return err
					}
					logrus.Infof("Canceled run %d, repository %v, ref %v, name %v -- %v old", run.ID, run.Task.Parent.Name, run.Task.Ref.RefName, run.Name, time.Since(run.CreatedAt))
				}
			}
		}
	}

	return nil
}
