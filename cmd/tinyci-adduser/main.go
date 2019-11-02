package main

import (
	"context"
	"fmt"
	"os"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
	"github.com/urfave/cli"
)

// Version is the version of this service.
const Version = "1.0.0"

// TinyCIVersion is the version of tinyci supporting this service.
var TinyCIVersion = "" // to be changed by build processes

func main() {
	app := cli.NewApp()
	app.Name = "tinyci-adduser"
	app.Description = "tinyci-adduser creates a user from a pre-generated token and generates token auth for it"
	app.ArgsUsage = "[github token]"
	app.Action = run
	app.Version = fmt.Sprintf("%s (tinyCI version %s)", Version, TinyCIVersion)

	app.Flags = []cli.Flag{
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
		errors.New(err).(errors.Error).Exit()
	}
}

func run(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return errors.New("See --help for more information on how to use this tool")
	}

	var cert *transport.Cert

	if !(ctx.GlobalString("cacert") == "" && ctx.GlobalString("cert") == "" && ctx.GlobalString("key") == "") {
		var err error
		// last arg is CRL
		cert, err = transport.LoadCert(ctx.GlobalString("cacert"), ctx.GlobalString("cert"), ctx.GlobalString("key"), "")
		if err != nil {
			return errors.New(err).(errors.Error).Wrap("while loading cert")
		}
	}

	client, err := data.New(ctx.GlobalString("datasvc"), cert, false)
	if err != nil {
		return err
	}
	defer client.Close()

	token := ctx.Args()[0]
	tokenStruct, err := inspect(token)
	if err != nil {
		return err
	}

	fmt.Printf("+++ Creating user %s\n", tokenStruct.Username)

	u := &model.User{
		Username: tokenStruct.Username,
		Token:    tokenStruct,
	}

	if _, err := client.PutUser(context.Background(), u); err != nil {
		return err
	}

	tinyCIToken, err := client.GetToken(context.Background(), tokenStruct.Username)
	if err != nil {
		return err
	}

	fmt.Println("+++ Generated tinyCI token is:")
	fmt.Println(tinyCIToken)

	return nil
}

func inspect(token string) (*types.OAuthToken, error) {
	c := github.NewClientFromAccessToken(token)

	login, err := c.MyLogin(context.Background())
	if err != nil {
		return nil, err
	}

	return &types.OAuthToken{
		Token:    token,
		Scopes:   []string{"repo"},
		Username: login,
	}, nil
}
