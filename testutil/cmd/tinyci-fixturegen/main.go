package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"errors"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/testutil/testclients"
	"github.com/tinyci/ci-agents/utils"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var dummyRun = &types.RunSettings{
	Command: []string{"dummy"},
	Image:   "dummy",
	Queue:   "default",
}

type cmd struct {
	dc       *testclients.DataClient
	ctx      *cli.Context
	min, max int
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	app := cli.NewApp()
	app.Usage = "Generate fake data for tinyci UI and data testing"
	app.Description = "You can just re-run this to generate more data."

	app.Flags = []cli.Flag{
		&cli.UintFlag{
			Name:  "maxlen, m",
			Usage: "Max length of strings (repo, user names, etc)",
			Value: 10,
		},
		&cli.UintFlag{
			Name:  "minlen, n",
			Usage: "Minimum length of strings (repo, user names, etc)",
			Value: 5,
		},
		&cli.UintFlag{
			Name:  "tasks, t",
			Usage: "Upper bound of tasks to generate per SHA",
			Value: 10,
		},
		&cli.UintFlag{
			Name:  "runs, r",
			Usage: "Upper bound of runs to generate per task",
			Value: 10,
		},
		&cli.UintFlag{
			Name:  "shas, s",
			Usage: "Upper bound of shas to generate per ref",
			Value: 10,
		},
		&cli.UintFlag{
			Name:  "refs, f",
			Usage: "Upper bound of refs to generate per repository",
			Value: 10,
		},
		&cli.UintFlag{
			Name:  "forks, k",
			Usage: "Upper bound of fork repositories to generate for each parent",
			Value: 10,
		},
		&cli.UintFlag{
			Name:  "repositories, p",
			Usage: "Upper bound of repositories to generate",
			Value: 10,
		},
		&cli.UintFlag{
			Name:  "owners, o",
			Usage: "Upper bound of owners (users) to generate to manage repositories",
			Value: 10,
		},
		&cli.BoolFlag{
			Name:  "private",
			Usage: "Make repostories private to the owner",
		},
		&cli.BoolFlag{
			Name:  "disable",
			Usage: "Leave newly created repos disabled",
		},
	}

	app.Action = generate

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (c *cmd) getString() string {
	return testutil.RandString(rand.Intn(c.max-c.min) + c.min)
}

func (c *cmd) mkUsers() ([]*types.User, error) {
	users := []*types.User{}

	for i := rand.Intn(int(c.ctx.Uint("owners"))) + 1; i >= 0; i-- {
		u, err := c.dc.MakeUser(c.getString())
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (c *cmd) mkParents(ctx context.Context, users []*types.User) (*types.RepositoryList, error) {
	parents := &types.RepositoryList{}

	for i := rand.Intn(int(c.ctx.Uint("repositories"))) + 1; i >= 0; i-- {
		ou := users[rand.Intn(len(users))]
		name := strings.Join([]string{c.getString(), c.getString()}, "/")
		if err := c.dc.MakeRepo(name, ou.Username, c.ctx.Bool("private"), ""); err != nil {
			return parents, err
		}

		if !c.ctx.Bool("disable") {
			if err := c.dc.Client().EnableRepository(ctx, ou.Username, name); err != nil {
				return parents, err
			}
		}

		r, err := c.dc.Client().GetRepository(ctx, name)
		if err != nil {
			return parents, err
		}

		parents.List = append(parents.List, r)
	}

	return parents, nil
}

func (c *cmd) mkForks(ctx context.Context, users []*types.User, parents *types.RepositoryList) (map[string]*types.Repository, error) {
	forkParents := map[string]*types.Repository{}

	for i := rand.Intn(int(c.ctx.Uint("forks"))) + 1; i >= 0; i-- {
		ou := users[rand.Intn(len(users))]
		pr := parents.List[rand.Intn(len(parents.List))]
		name := strings.Join([]string{c.getString(), c.getString()}, "/")

		repos := []interface{}{map[string]interface{}{
			"full_name": name,
			"private":   c.ctx.Bool("private"),
			"fork":      true,
			"parent": map[string]interface{}{
				"full_name": pr.Name,
				"private":   c.ctx.Bool("private"),
			},
		}}

		ghRepos := []*github.Repository{}

		if err := utils.JSONIO(repos, &ghRepos); err != nil {
			return nil, err
		}

		if err := c.dc.Client().PutRepositories(ctx, ou.Username, ghRepos, true); err != nil {
			return nil, err
		}

		repo, err := c.dc.Client().GetRepository(ctx, name)
		if err != nil {
			return nil, err
		}

		forkParents[repo.Name] = pr
	}

	return forkParents, nil
}

func (c *cmd) mkRefs(ctx context.Context, forkParents map[string]*types.Repository) ([]*types.Ref, []*types.Ref, error) {
	headrefs := []*types.Ref{}
	baserefs := []*types.Ref{}

	for fork, parent := range forkParents {
		for refC := rand.Intn(int(c.ctx.Uint("refs"))) + 1; refC >= 0; refC-- {
			refName := "heads/" + c.getString()

			for shaC := rand.Intn(int(c.ctx.Uint("shas"))) + 1; shaC >= 0; shaC-- {
				sha := ""
				for i := 0; i < 40; i++ {
					sha += fmt.Sprintf("%x", rune(rand.Intn(16)))
				}

				f, err := c.dc.Client().GetRepository(ctx, fork)
				if err != nil {
					return nil, nil, err
				}

				ref := &types.Ref{
					Repository: f,
					RefName:    refName,
					Sha:        sha,
				}

				ref.Id, err = c.dc.Client().PutRef(ctx, ref)
				if err != nil {
					return nil, nil, err
				}

				headrefs = append(headrefs, ref)
				sha = ""
				for i := 0; i < 40; i++ {
					sha += fmt.Sprintf("%x", rune(rand.Intn(16)))
				}

				ref = &types.Ref{
					Repository: parent,
					RefName:    "heads/master",
					Sha:        sha,
				}

				ref.Id, err = c.dc.Client().PutRef(ctx, ref)
				if err != nil {
					return nil, nil, err
				}

				baserefs = append(baserefs, ref)
			}
		}
	}

	return headrefs, baserefs, nil
}

func (c *cmd) mkTask(ctx context.Context, sub *types.Submission) (*types.Task, error) {
	started := rand.Intn(2) == 0
	finished := started && rand.Intn(2) == 0

	createdAt := time.Now().Add(time.Duration(-rand.Int63n(30 * 24 * int64(time.Hour))))
	var startedAt, finishedAt *time.Time
	var status *bool

	if started {
		tmp := createdAt.Add(time.Duration(rand.Int63n(int64(time.Hour))))
		startedAt = &tmp
	}

	if started && finished {
		tmp := startedAt.Add(time.Duration(rand.Int63n(int64(72 * time.Hour))))
		finishedAt = &tmp
		tmp2 := rand.Intn(2) == 0
		status = &tmp2
	}

	ts := &types.TaskSettings{
		Mountpoint: "/",
		Runs:       map[string]*types.RunSettings{"dummy": dummyRun},
	}

	var finishedTs *timestamppb.Timestamp

	if finishedAt != nil {
		finishedTs = timestamppb.New(*finishedAt)
	}

	var startedTs *timestamppb.Timestamp

	if startedAt != nil {
		startedTs = timestamppb.New(*startedAt)
	}

	var sbool bool
	if status != nil {
		sbool = *status
	}

	task := &types.Task{
		Path:       c.getString(),
		CreatedAt:  timestamppb.New(createdAt),
		FinishedAt: finishedTs,
		StartedAt:  startedTs,
		StatusSet:  status != nil,
		Status:     sbool,
		Settings:   ts,
		Submission: sub,
	}

	return c.dc.Client().PutTask(ctx, task)
}

func (c *cmd) mkTasks(ctx context.Context, subs []*types.Submission) error {
	for _, sub := range subs {
		for taskC := rand.Intn(int(c.ctx.Uint("tasks"))) + 1; taskC >= 0; taskC-- {
			task, err := c.mkTask(ctx, sub)
			if err != nil {
				return err
			}

			qis := []*types.QueueItem{}
			for runC := rand.Intn(int(c.ctx.Uint("runs"))) + 1; runC >= 0; runC-- {
				run := &types.Run{
					Settings:  dummyRun,
					CreatedAt: task.CreatedAt,
					Task:      task,
					Name:      c.getString(),
				}
				qis = append(qis, &types.QueueItem{
					Run:       run,
					QueueName: "default",
				})
			}

			if _, err := c.dc.Client().PutQueue(ctx, qis); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *cmd) mkSubmissions(ctx context.Context, u *types.User, baserefs []*types.Ref, headrefs []*types.Ref) ([]*types.Submission, error) {
	if len(headrefs) != len(baserefs) {
		return nil, errors.New("refs count is not equal")
	}

	subs := []*types.Submission{}

	for i := 0; i < len(baserefs); i++ {
		sub := &types.Submission{
			HeadRef: headrefs[i],
			BaseRef: baserefs[i],
			User:    u,
		}

		var err error
		sub, err = c.dc.Client().PutSubmission(ctx, sub)
		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func generate(ctx *cli.Context) error {
	dc, err := testclients.NewDataClient()
	if err != nil {
		return err
	}

	max := int(ctx.Uint("maxlen"))
	min := int(ctx.Uint("minlen"))

	if max < min {
		return errors.New("maxlen is smaller than minlen")
	}

	c := &cmd{
		dc:  dc,
		ctx: ctx,
		min: min,
		max: max,
	}

	users, err := c.mkUsers()
	if err != nil {
		return err
	}

	ct := context.Background()

	parents, err := c.mkParents(ct, users)
	if err != nil {
		return err
	}

	forkParents, err := c.mkForks(ct, users, parents)
	if err != nil {
		return err
	}

	headrefs, baserefs, err := c.mkRefs(ct, forkParents)
	if err != nil {
		return err
	}

	subs, err := c.mkSubmissions(ct, users[0], headrefs, baserefs)
	if err != nil {
		return err
	}

	if err := c.mkTasks(ct, subs); err != nil {
		return err
	}

	return nil
}
