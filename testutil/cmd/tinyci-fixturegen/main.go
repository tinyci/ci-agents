package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/testutil/testclients"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
	"github.com/urfave/cli"
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
		cli.UintFlag{
			Name:  "maxlen, m",
			Usage: "Max length of strings (repo, user names, etc)",
			Value: 10,
		},
		cli.UintFlag{
			Name:  "minlen, n",
			Usage: "Minimum length of strings (repo, user names, etc)",
			Value: 5,
		},
		cli.UintFlag{
			Name:  "tasks, t",
			Usage: "Upper bound of tasks to generate per SHA",
			Value: 10,
		},
		cli.UintFlag{
			Name:  "runs, r",
			Usage: "Upper bound of runs to generate per task",
			Value: 10,
		},
		cli.UintFlag{
			Name:  "shas, s",
			Usage: "Upper bound of shas to generate per ref",
			Value: 10,
		},
		cli.UintFlag{
			Name:  "refs, f",
			Usage: "Upper bound of refs to generate per repository",
			Value: 10,
		},
		cli.UintFlag{
			Name:  "forks, k",
			Usage: "Upper bound of fork repositories to generate for each parent",
			Value: 10,
		},
		cli.UintFlag{
			Name:  "repositories, p",
			Usage: "Upper bound of repositories to generate",
			Value: 10,
		},
		cli.UintFlag{
			Name:  "owners, o",
			Usage: "Upper bound of owners (users) to generate to manage repositories",
			Value: 10,
		},
		cli.BoolFlag{
			Name:  "private",
			Usage: "Make repostories private to the owner",
		},
		cli.BoolFlag{
			Name:  "disable",
			Usage: "Leave newly created repos disabled",
		},
	}

	app.Action = generate

	if err := app.Run(os.Args); err != nil {
		errors.New(err).Exit()
	}
}

func (c *cmd) getString() string {
	return testutil.RandString(rand.Intn(c.max-c.min) + c.min)
}

func (c *cmd) mkUsers() ([]*model.User, *errors.Error) {
	users := []*model.User{}

	for i := rand.Intn(int(c.ctx.GlobalUint("owners"))) + 1; i >= 0; i-- {
		u, err := c.dc.MakeUser(c.getString())
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (c *cmd) mkParents(users []*model.User) (model.RepositoryList, *errors.Error) {
	parents := model.RepositoryList{}

	for i := rand.Intn(int(c.ctx.GlobalUint("repositories"))) + 1; i >= 0; i-- {
		ou := users[rand.Intn(len(users))]
		name := strings.Join([]string{c.getString(), c.getString()}, "/")
		if err := c.dc.MakeRepo(name, ou.Username, c.ctx.GlobalBool("private")); err != nil {
			return nil, err
		}

		if !c.ctx.GlobalBool("disable") {
			if err := c.dc.Client().EnableRepository(ou.Username, name); err != nil {
				return nil, err
			}
		}

		r, err := c.dc.Client().GetRepository(name)
		if err != nil {
			return nil, err
		}

		parents = append(parents, r)
	}

	return parents, nil
}

func (c *cmd) mkForks(users []*model.User, parents model.RepositoryList) (map[string]*model.Repository, model.RepositoryList, *errors.Error) {
	forkParents := map[string]*model.Repository{}
	forks := model.RepositoryList{}

	for i := rand.Intn(int(c.ctx.GlobalUint("forks"))) + 1; i >= 0; i-- {
		ou := users[rand.Intn(len(users))]
		pr := parents[rand.Intn(len(parents))]
		name := strings.Join([]string{c.getString(), c.getString()}, "/")

		repos := []interface{}{map[string]interface{}{
			"full_name": name,
			"private":   c.ctx.GlobalBool("private"),
			"fork":      true,
			"parent": map[string]interface{}{
				"full_name": pr.Name,
				"private":   c.ctx.GlobalBool("private"),
			},
		}}

		ghRepos := []*github.Repository{}

		if err := utils.JSONIO(repos, &ghRepos); err != nil {
			return nil, nil, err
		}

		if err := c.dc.Client().PutRepositories(ou.Username, ghRepos, true); err != nil {
			if err != nil {
				return nil, nil, err
			}
		}

		repo, err := c.dc.Client().GetRepository(name)
		if err != nil {
			return nil, nil, err
		}

		forks = append(forks, repo)
		forkParents[repo.Name] = pr
	}

	return forkParents, forks, nil
}

func (c *cmd) mkRefs(forks model.RepositoryList) ([]*model.Ref, *errors.Error) {
	refs := []*model.Ref{}

	for _, fork := range forks {
		for refC := rand.Intn(int(c.ctx.GlobalUint("refs"))) + 1; refC >= 0; refC-- {
			refName := "refs/heads/" + c.getString()

			for shaC := rand.Intn(int(c.ctx.GlobalUint("shas"))) + 1; shaC >= 0; shaC-- {
				sha := ""
				for i := 0; i < 40; i++ {
					sha += fmt.Sprintf("%x", rune(rand.Intn(16)))
				}

				ref := &model.Ref{
					Repository: fork,
					RefName:    refName,
					SHA:        sha,
				}

				var err *errors.Error
				ref.ID, err = c.dc.Client().PutRef(ref)
				if err != nil {
					return nil, err
				}

				refs = append(refs, ref)
			}
		}
	}

	return refs, nil
}

func (c *cmd) mkTask(ref *model.Ref, forkParents map[string]*model.Repository) (*model.Task, *errors.Error) {
	basesha := ""
	for i := 0; i < 40; i++ {
		basesha += fmt.Sprintf("%x", rune(rand.Intn(16)))
	}

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
	if err := ts.Validate(true); err != nil {
		return nil, err
	}

	task := &model.Task{
		Path:          c.getString(),
		Parent:        forkParents[ref.Repository.Name],
		Ref:           ref,
		BaseSHA:       basesha,
		PullRequestID: rand.Int63n(9000),
		CreatedAt:     createdAt,
		FinishedAt:    finishedAt,
		StartedAt:     startedAt,
		Status:        status,
		TaskSettings:  ts,
	}

	return c.dc.Client().PutTask(task)
}

func (c *cmd) mkTasks(refs []*model.Ref, forkParents map[string]*model.Repository) *errors.Error {
	for _, ref := range refs {
		for taskC := rand.Intn(int(c.ctx.GlobalUint("tasks"))) + 1; taskC >= 0; taskC-- {
			task, err := c.mkTask(ref, forkParents)
			if err != nil {
				return err
			}

			qis := []*model.QueueItem{}
			for runC := rand.Intn(int(c.ctx.GlobalUint("runs"))) + 1; runC >= 0; runC-- {
				run := &model.Run{
					RunSettings: dummyRun,
					CreatedAt:   task.CreatedAt,
					Task:        task,
					Name:        c.getString(),
				}
				qis = append(qis, &model.QueueItem{
					Run:       run,
					QueueName: "default",
				})
			}

			if _, err := c.dc.Client().PutQueue(qis); err != nil {
				return err
			}
		}
	}

	return nil
}

func generate(ctx *cli.Context) error {
	dc, err := testclients.NewDataClient()
	if err != nil {
		return err
	}

	max := int(ctx.GlobalUint("maxlen"))
	min := int(ctx.GlobalUint("minlen"))

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

	parents, err := c.mkParents(users)
	if err != nil {
		return err
	}

	forkParents, forks, err := c.mkForks(users, parents)
	if err != nil {
		return err
	}

	refs, err := c.mkRefs(forks)
	if err != nil {
		return err
	}

	if err := c.mkTasks(refs, forkParents); err != nil {
		return err
	}

	return nil
}
