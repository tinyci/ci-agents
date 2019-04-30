package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/tinyci/ci-agents/clients/tinyci"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"github.com/urfave/cli"
)

var tinyCIConfig = path.Join(os.Getenv("HOME"), ".tinycli")

// Config is the configuration of tinycli for between uses.
type Config struct {
	Endpoint string
	Token    string
}

func getConfigPath(ctx *cli.Context) (string, *errors.Error) {
	if fi, err := os.Stat(tinyCIConfig); err != nil {
		if mkerr := os.MkdirAll(tinyCIConfig, 0700); mkerr != nil {
			return "", errors.New("Could not make config dir").Wrap(mkerr).Wrap(err)
		}
	} else if !fi.IsDir() {
		return "", errors.Errorf("tinycli configuration path %q exists and is not a directory", tinyCIConfig)
	}

	config := ctx.GlobalString("config")
	if config == "" {
		return "", errors.New("invalid config name")
	}

	return path.Join(tinyCIConfig, config), nil
}

func loadConfig(ctx *cli.Context) (*tinyci.Client, *errors.Error) {
	filename, e := getConfigPath(ctx)
	if e != nil {
		return nil, e
	}

	f, err := os.Open(filename) // #nosec
	if err != nil {
		return nil, errors.New(err).Wrapf("Cannot open tinyci configuration file %q", filename)
	}
	defer f.Close()

	c := Config{}

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, errors.New(err).Wrapf("Could not decode tinyCI JSON configuration in %q", filename)
	}

	return c.mkClient()
}

func (c Config) mkClient() (*tinyci.Client, *errors.Error) {
	return tinyci.New(c.Endpoint, c.Token)
}

// Version is the version of this service.
const Version = "1.0.0"

// TinyCIVersion is the version of tinyci supporting this service.
var TinyCIVersion = "" // to be changed by build processes

func main() {
	app := cli.NewApp()
	app.Name = "tinycli"
	app.Description = `
Commandline client to control tinyCI. Useful for a variety of querying and
control operations.

To select a configuation at 'init' time, please specify the configuration you
want to init like you would with other -c / command combinations:

tinycli -c foo init

You can also specify the TINYCLI_CONFIG environment variable.
`

	app.Version = fmt.Sprintf("%s (tinyCI version %s)", Version, TinyCIVersion)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  fmt.Sprintf("Name of configuration to use; comes from %q", tinyCIConfig),
			Value:  "default",
			EnvVar: "TINYCLI_CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "init",
			ShortName:   "i",
			Description: "Initialize the client with a token and endpoint URL",
			Usage:       "Initialize the client with a token and endpoint URL",
			Action:      doInit,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "token, t",
					Usage: "Provide the token on the command-line instead of being prompted for it",
				},
				cli.StringFlag{
					Name:  "url, u",
					Usage: "Provide the URL to access the service",
				},
			},
		},
		{
			Name:        "submit",
			ShortName:   "s",
			Description: "Submit a job to tinyCI",
			Usage:       "Submit a job to tinyCI",
			ArgsUsage:   "[parent or fork repository] [sha]",
			Action:      submit,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "For a test of all task dirs, not just diff-affected ones",
				},
			},
		},
		{
			Name:        "tasks",
			ShortName:   "t",
			Description: "List Tasks",
			Usage:       "List Tasks",
			Action:      tasks,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repository, r",
					Usage: "Repository name for filtering runs",
				},
				cli.StringFlag{
					Name:  "ref, n",
					Usage: "Ref/SHA name for filtering runs. Repository is required if SHA provided, otherwise it is ignored",
				},
				cli.Int64Flag{
					Name:  "page, p",
					Usage: "The page of runs to access",
				},
				cli.Int64Flag{
					Name:  "count, c",
					Usage: "The amount of runs to show",
				},
			},
		},
		{
			Name:        "runs",
			ShortName:   "r",
			Description: "List runs",
			Usage:       "List runs",
			Action:      runs,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repository, r",
					Usage: "Repository name for filtering runs",
				},
				cli.StringFlag{
					Name:  "ref, n",
					Usage: "Ref/SHA name for filtering runs. Repository is required if SHA provided, otherwise it is ignored",
				},
				cli.Int64Flag{
					Name:  "page, p",
					Usage: "The page of runs to access",
				},
				cli.Int64Flag{
					Name:  "count, c",
					Usage: "The amount of runs to show",
				},
			},
		},
		{
			Name:        "log",
			ShortName:   "l",
			Description: "Show a log by Run ID",
			Usage:       "Show a log by Run ID",
			ArgsUsage:   "[run id]",
			Action:      log,
		},
		{
			Name:        "capabilities",
			ShortName:   "c",
			Description: "Manipulate User Capabilities",
			Usage:       "Manipulate User Capabilities",
			Subcommands: []cli.Command{
				{
					Name:        "add",
					ShortName:   "a",
					Description: "Grant a capability to a user",
					Usage:       "Grant a capability to a user",
					ArgsUsage:   "[username] [capability]",
					Action:      addCapability,
				},
				{
					Name:        "remove",
					ShortName:   "r",
					Description: "Remove a capability from a user",
					Usage:       "Remove a capability from a user",
					ArgsUsage:   "[username] [capability]",
					Action:      removeCapability,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		if e, ok := err.(*errors.Error); ok && e == nil {
			return
		}

		fmt.Fprintf(os.Stderr, "%+v\n", err)
	}
}

func stdTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)
}

func doInit(ctx *cli.Context) error {
	token := ctx.String("token")
	u := ctx.String("url")

	if u == "" {
		fmt.Print("Paste in your tinyCI ui service URL endpoint: ")
		s := bufio.NewScanner(os.Stdin)
		if s.Scan() {
			u = strings.TrimSpace(s.Text())
		} else {
			return errors.New("Could not scan url; will not continue")
		}
	}

	if token == "" {
		fmt.Println("THE FOLLOWING URL WILL ONLY SHOW THE TOKEN ONCE. DO NOT RELOAD IT!")
		fmt.Printf("Go here: %s/token\n", u)
		fmt.Print("Paste in the key as provided: ")

		s := bufio.NewScanner(os.Stdin)
		if s.Scan() {
			token = strings.TrimSpace(s.Text())
		} else {
			return errors.New("Could not scan url; will not continue")
		}

		if token[0] == '"' && token[len(token)-1] == '"' {
			token = token[1 : len(token)-1]
		}
	}

	client, err := tinyci.New(u, token)
	if err != nil {
		return errors.New(err)
	}

	if _, err := client.Errors(); err != nil {
		return err.Wrap("Could not retrieve with the client, token or URL issue")
	}

	c := Config{
		Endpoint: u,
		Token:    token,
	}

	filename, err := getConfigPath(ctx)
	if err != nil {
		return err
	}
	f, ferr := os.Create(filename)
	if ferr != nil {
		return errors.New(ferr).Wrapf("Could not create configuration file %v", filename)
	}
	defer f.Close()
	defer fmt.Printf("Created configuration file %q\n", filename)

	return json.NewEncoder(f).Encode(c)
}

func submit(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return errors.New("Invalid arguments: [repository] [sha] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Submitting %s / %s (all tasks: %v) -- this may take a few seconds to complete.\n", ctx.Args()[0], ctx.Args()[1], ctx.Bool("all"))

	if err := client.Submit(ctx.Args()[0], ctx.Args()[1], ctx.Bool("all")); err != nil {
		return err
	}

	fmt.Println("Successfully submitted!")
	return nil
}

func mkTaskStatus(task *model.Task) string {
	statusStr := "queued"
	if task.Canceled {
		statusStr = "canceled"
	} else if task.Status != nil {
		if *task.Status {
			statusStr = "success"
		} else {
			statusStr = "failure"
		}
	} else if task.StartedAt != nil && task.FinishedAt == nil {
		statusStr = "running"
	}

	return statusStr
}

func mkTaskRunCounts(client *tinyci.Client, task *model.Task) (int64, int64, int64, error) {
	totalCount, err := client.RunsForTaskCount(task.ID)
	if err != nil {
		return 0, 0, 0, err
	}

	runs := []*model.Run{}

	for i := int64(0); i <= totalCount/utils.MaxPerPage; i++ {
		tmp, err := client.RunsForTask(task.ID, i, utils.MaxPerPage)
		if err != nil {
			return 0, 0, 0, err
		}

		runs = append(runs, tmp...)
	}

	var runningCount, finishedCount int64
	for _, run := range runs {
		if run.FinishedAt != nil {
			finishedCount++
		} else if run.StartedAt != nil {
			runningCount++
		}
	}

	return runningCount, finishedCount, totalCount, nil
}

func tasks(ctx *cli.Context) error {
	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	tasks, err := client.Tasks(ctx.String("repository"), ctx.String("ref"), ctx.Int64("page"), ctx.Int64("count"))
	if err != nil {
		return err
	}

	w := stdTabWriter()
	w.Write([]byte("TASK ID\tREPOSITORY\tREF\tSHA\tPATH\tRUN/FIN/TOT\tSTATE\tDURATION\n"))

	for _, task := range tasks {
		statusStr := mkTaskStatus(task)

		duration := ""

		if task.StartedAt != nil && task.FinishedAt != nil {
			d := task.FinishedAt.Sub(*task.StartedAt)
			duration = d.Round(time.Millisecond).String()
		} else if task.StartedAt != nil {
			duration = time.Since(*task.StartedAt).Round(time.Millisecond).String()
		}

		refName := task.Ref.RefName
		sha := task.Ref.SHA[:12]

		runningCount, finishedCount, totalCount, err := mkTaskRunCounts(client, task)
		if err != nil {
			return err
		}

		path := task.Path
		if path == "." {
			path = "*root*"
		}

		w.Write([]byte(fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d/%d/%d\t%s\t%s\n", task.ID, task.Ref.Repository.Name, refName, sha, path, runningCount, finishedCount, totalCount, statusStr, duration)))
	}
	w.Flush()

	return nil
}

func runs(ctx *cli.Context) error {
	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	runs, err := client.Runs(ctx.String("repository"), ctx.String("ref"), ctx.Int64("page"), ctx.Int64("count"))
	if err != nil {
		return err
	}

	w := stdTabWriter()
	w.Write([]byte("RUN ID\tREPOSITORY\tREF\tSHA\tRUN\tTASK ID\tSTATE\tDURATION\n"))
	for _, run := range runs {
		statusStr := "queued"
		if run.Task.Canceled {
			statusStr = "canceled"
		} else if run.Status != nil {
			if *run.Status {
				statusStr = "success"
			} else {
				statusStr = "failure"
			}
		} else if run.StartedAt != nil && run.FinishedAt == nil {
			statusStr = "running"
		}

		duration := ""

		if run.StartedAt != nil && run.FinishedAt != nil {
			d := run.FinishedAt.Sub(*run.StartedAt).Round(time.Millisecond)
			duration = d.String()
		}

		refName := run.Task.Ref.RefName
		sha := run.Task.Ref.SHA[:12]

		w.Write([]byte(fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d\t%s\t%s\n", run.ID, run.Task.Ref.Repository.Name, refName, sha, run.Name, run.Task.ID, statusStr, duration)))
	}
	w.Flush()

	return nil
}

func log(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return errors.New("Invalid arguments: [run id] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	id, convErr := strconv.ParseInt(ctx.Args()[0], 10, 64)
	if convErr != nil {
		return errors.New(convErr).Wrap("Invalid ID")
	}

	return client.LogAttach(id, os.Stdout)
}

func addCapability(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return errors.New("Invalid arguments: [username] [capability] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	return client.AddCapability(ctx.Args()[0], model.Capability(ctx.Args()[1]))
}

func removeCapability(ctx *cli.Context) error {
	if len(ctx.Args()) != 2 {
		return errors.New("Invalid arguments: [username] [capability] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	return client.RemoveCapability(ctx.Args()[0], model.Capability(ctx.Args()[1]))
}
