package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/erikh/colorwriter"
	transport "github.com/erikh/go-transport"
	"github.com/fatih/color"
	"github.com/tinyci/ci-agents/clients/tinyci"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var tinyCIConfig = path.Join(os.Getenv("HOME"), ".tinycli")

// Config is the configuration of tinycli for between uses.
type Config struct {
	Endpoint string
	Token    string
}

func getConfigPath(ctx *cli.Context) (string, error) {
	if fi, err := os.Stat(tinyCIConfig); err != nil {
		if mkerr := os.MkdirAll(tinyCIConfig, 0700); mkerr != nil {
			return "", utils.WrapError(mkerr, "Could not make config dir: %v", err)
		}
	} else if !fi.IsDir() {
		return "", fmt.Errorf("tinycli configuration path %q exists and is not a directory", tinyCIConfig)
	}

	config := ctx.String("config")
	if config == "" {
		return "", errors.New("invalid config name")
	}

	return path.Join(tinyCIConfig, config), nil
}

func getCert(ctx *cli.Context) (*transport.Cert, error) {
	ca, certStr, keyStr := ctx.String("ca"),
		ctx.String("cert"),
		ctx.String("key")

	if ca == "" && certStr == "" && keyStr == "" {
		return nil, nil
	}

	cert, err := transport.LoadCert(ca, certStr, keyStr, "")
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func loadConfig(ctx *cli.Context) (*tinyci.Client, error) {
	filename, e := getConfigPath(ctx)
	if e != nil {
		return nil, e
	}

	f, err := os.Open(filename) // #nosec
	if err != nil {
		return nil, utils.WrapError(err, "Cannot open tinyci configuration file %q", filename)
	}
	defer f.Close()

	c := Config{}

	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, utils.WrapError(err, "Could not decode tinyCI JSON configuration in %q", filename)
	}

	return c.mkClient(ctx)
}

func (c Config) mkClient(ctx *cli.Context) (*tinyci.Client, error) {
	cert, err := getCert(ctx)
	if err != nil {
		return nil, err
	}
	return tinyci.New(c.Endpoint, c.Token, cert)
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
		&cli.StringFlag{
			Name:    "config, c",
			Usage:   fmt.Sprintf("Name of configuration to use; comes from %q", tinyCIConfig),
			Value:   "default",
			EnvVars: []string{"TINYCLI_CONFIG"},
		},
		&cli.StringFlag{
			Name:    "ca, a",
			Usage:   "CA certificate used to contact remote service",
			EnvVars: []string{"TINYCLI_CA_CERT"},
		},
		&cli.StringFlag{
			Name:    "cert, t",
			Usage:   "TLS certificate to use to contact remote service (ecdsa only)",
			EnvVars: []string{"TINYCLI_CERT"},
		},
		&cli.StringFlag{
			Name:    "key, k",
			Usage:   "TLS private key to use to contact remote service (ecdsa only)",
			EnvVars: []string{"TINYCLI_KEY"},
		},
		&cli.BoolFlag{
			Name:    "no-color, nc",
			Usage:   "Turn off coloring for output",
			EnvVars: []string{"TINYCLI_NOCOLOR"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:        "init",
			Aliases:     []string{"i"},
			Description: "Initialize the client with a token and endpoint URL",
			Usage:       "Initialize the client with a token and endpoint URL",
			Action:      doInit,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "token, t",
					Usage: "Provide the token on the command-line instead of being prompted for it",
				},
				&cli.StringFlag{
					Name:  "url, u",
					Usage: "Provide the URL to access the service",
				},
			},
		},
		{
			Name:        "submit",
			Aliases:     []string{"sub"},
			Description: "Submit a job to tinyCI",
			Usage:       "Submit a job to tinyCI",
			ArgsUsage:   "[parent or fork repository] [sha]",
			Action:      submit,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "all, a",
					Usage: "For a test of all task dirs, not just diff-affected ones",
				},
			},
		},
		{
			Name:        "submissions",
			Aliases:     []string{"s"},
			Description: "List Submissions",
			Usage:       "List Submissions",
			Action:      submissions,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "repository, r",
					Usage: "Repository name for filtering runs",
				},
				&cli.StringFlag{
					Name:  "ref, n",
					Usage: "Ref/SHA name for filtering runs. Repository is required if SHA provided, otherwise it is ignored",
				},
				&cli.Int64Flag{
					Name:  "page, p",
					Usage: "The page of runs to access",
				},
				&cli.Int64Flag{
					Name:  "count, c",
					Usage: "The amount of runs to show",
				},
			},
		},
		{
			Name:        "tasks",
			Aliases:     []string{"t"},
			Description: "List Tasks",
			Usage:       "List Tasks",
			Action:      tasks,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "repository, r",
					Usage: "Repository name for filtering runs",
				},
				&cli.StringFlag{
					Name:  "ref, n",
					Usage: "Ref/SHA name for filtering runs. Repository is required if SHA provided, otherwise it is ignored",
				},
				&cli.Int64Flag{
					Name:  "page, p",
					Usage: "The page of runs to access",
				},
				&cli.Int64Flag{
					Name:  "count, c",
					Usage: "The amount of runs to show",
				},
			},
		},
		{
			Name:        "runs",
			Aliases:     []string{"r"},
			Description: "List runs",
			Usage:       "List runs",
			Action:      runs,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "repository, r",
					Usage: "Repository name for filtering runs",
				},
				&cli.StringFlag{
					Name:  "ref, n",
					Usage: "Ref/SHA name for filtering runs. Repository is required if SHA provided, otherwise it is ignored",
				},
				&cli.Int64Flag{
					Name:  "page, p",
					Usage: "The page of runs to access",
				},
				&cli.Int64Flag{
					Name:  "count, c",
					Usage: "The amount of runs to show",
				},
			},
		},
		{
			Name:        "log",
			Aliases:     []string{"l"},
			Description: "Show a log by Run ID",
			Usage:       "Show a log by Run ID",
			ArgsUsage:   "[run id]",
			Action:      log,
		},
		{
			Name:        "capabilities",
			Aliases:     []string{"c"},
			Description: "Manipulate User Capabilities",
			Usage:       "Manipulate User Capabilities",
			Subcommands: []*cli.Command{
				{
					Name:        "add",
					Aliases:     []string{"a"},
					Description: "Grant a capability to a user",
					Usage:       "Grant a capability to a user",
					ArgsUsage:   "[username] [capability]",
					Action:      addCapability,
				},
				{
					Name:        "remove",
					Aliases:     []string{"r"},
					Description: "Remove a capability from a user",
					Usage:       "Remove a capability from a user",
					ArgsUsage:   "[username] [capability]",
					Action:      removeCapability,
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func stdTabWriter(ctx *cli.Context) *colorwriter.Writer {
	if ctx.IsSet("no-color") {
		color.NoColor = ctx.Bool("no-color")
	} else if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		color.NoColor = true
	}

	return colorwriter.NewWriter(os.Stdout, 2, 2, 4, ' ', 0)
}

func getHeaderColorFunc() func(string, ...interface{}) string {
	return color.CyanString
}

func getRowColorFunc(i int) func(string, ...interface{}) string {
	colorFunc := color.HiBlackString
	if i%2 == 0 {
		colorFunc = color.WhiteString
	}

	return colorFunc
}

func doInit(ctx *cli.Context) error {
	token := ctx.String("token")
	u := ctx.String("url")

	cert, err := getCert(ctx)
	if err != nil {
		return err
	}

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

	client, err := tinyci.New(u, token, cert)
	if err != nil {
		return err
	}

	if _, err := client.Errors(context.Background()); err != nil {
		return utils.WrapError(err, "Could not retrieve with the client, token or URL issue")
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
		return utils.WrapError(ferr, "Could not create configuration file %v", filename)
	}
	defer f.Close()
	defer fmt.Printf("Created configuration file %q\n", filename)

	return json.NewEncoder(f).Encode(c)
}

func submit(ctx *cli.Context) error {
	if ctx.Args().Len() != 2 {
		return errors.New("Invalid arguments: [repository] [sha] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	owner, repo, all :=
		ctx.Args().Get(0),
		ctx.Args().Get(1),
		ctx.Bool("all")

	fmt.Printf("Submitting %s / %s (all tasks: %v) -- this may take a few seconds to complete.\n", owner, repo, all)

	if err := client.Submit(context.Background(), owner, repo, all); err != nil {
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

func mkSubRunCounts(ctx context.Context, client *tinyci.Client, sub *model.Submission) (int64, int64, int64, error) {
	var (
		page = int64(0)
		runs = []*model.Run{}
	)

	for {
		tmp, err := client.RunsForSubmission(ctx, sub.ID, page, utils.MaxPerPage)
		if err != nil || len(tmp) == 0 {
			break
		}

		runs = append(runs, tmp...)
		page++
	}

	var runningCount, finishedCount, totalCount int64

	for _, run := range runs {
		if run.FinishedAt != nil {
			finishedCount++
		} else if run.StartedAt != nil {
			runningCount++
		}

		totalCount++
	}

	return runningCount, finishedCount, totalCount, nil
}

func mkTaskRunCounts(ctx context.Context, client *tinyci.Client, task *model.Task) (int64, int64, int64, error) {
	runs := []*model.Run{}

	var i int64
	for {
		tmp, err := client.RunsForTask(ctx, task.ID, &i, int64p(utils.MaxPerPage))
		if err != nil || len(tmp) == 0 {
			break
		}

		runs = append(runs, tmp...)
		i++
	}

	var totalCount, runningCount, finishedCount int64
	for _, run := range runs {
		if run.FinishedAt != nil {
			finishedCount++
		} else if run.StartedAt != nil {
			runningCount++
		}

		totalCount++
	}

	return runningCount, finishedCount, totalCount, nil
}

func stringp(s string) *string {
	return &s
}

func int64p(i int64) *int64 {
	return &i
}

func submissions(ctx *cli.Context) error {
	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	ct := context.Background()
	subs, err := client.Submissions(ct, stringp(ctx.String("repository")), stringp(ctx.String("ref")), int64p(ctx.Int64("page")), int64p(ctx.Int64("count")))
	if err != nil {
		return err
	}

	w := stdTabWriter(ctx)
	if _, err := w.Write([]byte(getHeaderColorFunc()("SUB ID\tREPOSITORY\tREF\tSHA\tRUN/FIN/TOT\tSTATE\tDURATION\n"))); err != nil {
		return err
	}

	for i, sub := range subs {
		running, finished, total, err := mkSubRunCounts(ct, client, sub)
		if err != nil {
			return err
		}

		status := "created"
		duration := time.Since(sub.CreatedAt)

		if sub.Status != nil {
			if *sub.Status {
				status = "success"
			} else {
				status = "failed"
			}

			if sub.FinishedAt == nil || sub.StartedAt == nil {
				duration = 0
			} else {
				duration = (*sub.FinishedAt).Sub(*sub.StartedAt)
			}
		} else if sub.StartedAt != nil {
			status = "started"
			duration = time.Since(*sub.StartedAt)
		}

		_, eErr := fmt.Fprintf(w,
			getRowColorFunc(i)("%d\t%s\t%s\t%s\t%d/%d/%d\t%v\t%v\n"),
			sub.ID,
			sub.HeadRef.Repository.Name,
			strings.TrimPrefix(sub.HeadRef.RefName, "heads/"),
			sub.HeadRef.SHA[:12],
			running, finished, total,
			status,
			duration,
		)
		if eErr != nil {
			return eErr
		}
	}

	return w.Flush()
}

func tasks(ctx *cli.Context) error {
	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	ct := context.Background()

	tasks, err := client.Tasks(ct, stringp(ctx.String("repository")), stringp(ctx.String("ref")), int64p(ctx.Int64("page")), int64p(ctx.Int64("count")))
	if err != nil {
		return err
	}

	w := stdTabWriter(ctx)
	if _, err := w.Write([]byte(getHeaderColorFunc()("TASK ID\tREPOSITORY\tREF\tSHA\tPATH\tRUN/FIN/TOT\tSTATE\tDURATION\n"))); err != nil {
		return err
	}

	for i, task := range tasks {
		statusStr := mkTaskStatus(task)

		duration := ""

		if task.StartedAt != nil && task.FinishedAt != nil {
			d := task.FinishedAt.Sub(*task.StartedAt)
			duration = d.Round(time.Millisecond).String()
		} else if task.StartedAt != nil {
			duration = time.Since(*task.StartedAt).Round(time.Millisecond).String()
		}

		refName := task.Submission.HeadRef.RefName
		sha := task.Submission.HeadRef.SHA[:12]

		runningCount, finishedCount, totalCount, err := mkTaskRunCounts(ct, client, task)
		if err != nil {
			return err
		}

		path := task.Path
		if path == "." {
			path = "*root*"
		}

		if _, err := w.Write([]byte(getRowColorFunc(i)(fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d/%d/%d\t%s\t%s\n", task.ID, task.Submission.HeadRef.Repository.Name, refName, sha, path, runningCount, finishedCount, totalCount, statusStr, duration)))); err != nil {
			return err
		}
	}
	w.Flush()

	return nil
}

func runs(ctx *cli.Context) error {
	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	runs, err := client.Runs(context.Background(), stringp(ctx.String("repository")), stringp(ctx.String("ref")), int64p(ctx.Int64("page")), int64p(ctx.Int64("count")))
	if err != nil {
		return err
	}

	w := stdTabWriter(ctx)
	if _, err := w.Write([]byte(getHeaderColorFunc()("RUN ID\tREPOSITORY\tREF\tSHA\tRUN\tTASK ID\tSTATE\tDURATION\n"))); err != nil {
		return err
	}
	for i, run := range runs {
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

		refName := run.Task.Submission.HeadRef.RefName
		sha := run.Task.Submission.HeadRef.SHA[:12]

		if _, err := w.Write([]byte(getRowColorFunc(i)(fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d\t%s\t%s\n", run.ID, run.Task.Submission.HeadRef.Repository.Name, refName, sha, run.Name, run.Task.ID, statusStr, duration)))); err != nil {
			return err
		}
	}
	w.Flush()

	return nil
}

func log(ctx *cli.Context) error {
	if ctx.Args().Len() != 1 {
		return errors.New("Invalid arguments: [run id] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	id, convErr := strconv.ParseInt(ctx.Args().First(), 10, 64)
	if convErr != nil {
		return utils.WrapError(convErr, "Invalid ID")
	}

	return client.LogAttach(context.Background(), id, os.Stdout)
}

func addCapability(ctx *cli.Context) error {
	if ctx.Args().Len() != 2 {
		return errors.New("Invalid arguments: [username] [capability] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	err = client.AddCapability(context.Background(), ctx.Args().Get(0), model.Capability(ctx.Args().Get(1)))
	if err != nil {
		return err
	}

	return nil
}

func removeCapability(ctx *cli.Context) error {
	if ctx.Args().Len() != 2 {
		return errors.New("Invalid arguments: [username] [capability] required")
	}

	client, err := loadConfig(ctx)
	if err != nil {
		return err
	}

	err = client.RemoveCapability(context.Background(), ctx.Args().Get(0), model.Capability(ctx.Args().Get(1)))
	if err != nil {
		return err
	}

	return nil
}
