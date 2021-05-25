package uisvc

import (
	"bytes"
	"context"
	"sort"
	"strings"
	"time"

	"errors"

	check "github.com/erikh/check"
	"github.com/golang/mock/gomock"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"

	gh "github.com/google/go-github/github"
)

func stringp(s string) *string { return &s }
func int64p(i int64) *int64    { return &i }

var ctx = context.Background()

func (us *uisvcSuite) TestCapabilities(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, utc, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(utc.AddCapability(ctx, "erikh2", "modify:user"), check.NotNil)
	c.Assert(tc.AddCapability(ctx, "erikh2", "modify:user"), check.IsNil)
	c.Assert(utc.AddCapability(ctx, "erikh2", "modify:ci"), check.IsNil)

	props, err := utc.GetUserProperties(ctx)
	c.Assert(err, check.IsNil)
	caps := []string{}
	for _, cap := range props["capabilities"].([]interface{}) {
		caps = append(caps, cap.(string))
	}

	sort.Strings(caps)
	c.Assert(caps, check.DeepEquals, []string{"modify:ci", "modify:user"})

	c.Assert(utc.RemoveCapability(ctx, "erikh2", "modify:user"), check.IsNil)
	c.Assert(utc.RemoveCapability(ctx, "erikh2", "modify:ci"), check.NotNil)
}

func (us *uisvcSuite) TestErrors(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, _, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	errs, err := tc.Errors(ctx)
	c.Assert(err, check.IsNil)
	c.Assert(len(errs), check.Equals, 0)
}

type closeBuffer struct {
	*bytes.Buffer
}

func (cb *closeBuffer) Close() error {
	return nil
}

func (us *uisvcSuite) TestLogAttach(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, _, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(us.assetsvcClient.Write(context.Background(), 1, bytes.NewBufferString("this is a log")), check.IsNil)
	time.Sleep(100 * time.Millisecond)
	buf := &closeBuffer{bytes.NewBuffer(nil)}
	c.Assert(tc.LogAttach(ctx, 1, buf), check.IsNil)

	c.Assert(strings.HasPrefix(buf.String(), "this is a log"), check.Equals, true, check.Commentf("buf: %s", buf))
}

func (us *uisvcSuite) TestTokenEndpoints(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, _, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(tc.DeleteToken(ctx), check.IsNil)
	_, err = tc.Errors(ctx)
	c.Assert(err, check.ErrorMatches, ".*invalid authentication")
}

func (us *uisvcSuite) TestDeleteToken(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, _, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(tc.DeleteToken(ctx), check.IsNil)
	_, err = tc.Errors(ctx)
	c.Assert(err, check.ErrorMatches, ".*invalid authentication")
}

func (us *uisvcSuite) TestSubmit(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, utc, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	erikhClient := github.NewMockClient(gomock.NewController(c))
	config.SetDefaultGithubClient(erikhClient, "erikh")

	erikhClient.EXPECT().GetRepository(gomock.Any(), "erikh/not-real").Return(nil, errors.New("not found"))
	c.Assert(tc.Submit(ctx, "erikh/not-real", "master", true), check.ErrorMatches, ".* not found")

	erikhClient.EXPECT().MyRepositories(gomock.Any()).Return([]*gh.Repository{{FullName: gh.String("erikh/parent")}}, nil)

	repos, err := tc.LoadRepositories(ctx, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Not(check.Equals), 0)

	c.Assert(us.datasvcClient.Client().EnableRepository(ctx, "erikh", "erikh/parent"), check.IsNil)

	erikhClient.EXPECT().GetRepository(gomock.Any(), "erikh/not-real").Return(nil, errors.New("not found"))
	c.Assert(tc.Submit(ctx, "erikh/not-real", "master", true), check.ErrorMatches, ".* not found")

	sub := &types.Submission{
		Parent:  "erikh/parent",
		Fork:    "erikh/test",
		BaseSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		HeadSHA: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Manual:  true,
		All:     true,
	}

	erikhClient.EXPECT().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	erikhClient.EXPECT().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	erikhClient.EXPECT().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	erikhClient.EXPECT().GetSHA(gomock.Any(), sub.Fork, "heads/master").Return("", errors.New("not found"))
	client.EXPECT().GetSHA(gomock.Any(), sub.Fork, "heads/master").Return("", errors.New("not found"))
	client.EXPECT().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	client.EXPECT().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	client.EXPECT().GetSHA(gomock.Any(), sub.Parent, "heads/master").Return(sub.HeadSHA, nil)
	client.EXPECT().ClearStates(gomock.Any(), "erikh/parent", sub.HeadSHA).Return(nil)

	c.Assert(tc.Submit(ctx, "erikh/test", "master", true), check.ErrorMatches, ".*not found")

	erikhClient.EXPECT().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	erikhClient.EXPECT().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	erikhClient.EXPECT().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	erikhClient.EXPECT().GetSHA(gomock.Any(), "erikh/test", "heads/master").Return("", errors.New("not found"))
	client.EXPECT().GetSHA(gomock.Any(), "erikh/test", "heads/master").Return("", errors.New("not found"))
	client.EXPECT().GetSHA(gomock.Any(), "erikh/parent", "heads/master").Return("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil)
	client.EXPECT().ClearStates(gomock.Any(), "erikh/parent", sub.HeadSHA).Return(nil)

	c.Assert(tc.Submit(ctx, "erikh/test", "master", true), check.ErrorMatches, ".*not found")

	erikhClient.EXPECT().ClearStates(gomock.Any(), "erikh/parent", sub.HeadSHA).Return(nil)
	c.Assert(us.queuesvcClient.SetMockSubmissionSuccess(erikhClient.EXPECT(), sub, "heads/master", ""), check.IsNil)

	c.Assert(utc.Submit(ctx, "erikh/test", "master", true), check.NotNil)
	c.Assert(tc.Submit(ctx, "erikh/test", "master", true), check.IsNil)

	tasks, err := tc.Tasks(ctx, stringp("erikh/test"), &sub.HeadSHA, nil, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(tasks), check.Not(check.Equals), 0)
	count, err := tc.TaskCount(ctx, stringp("erikh/test"), &sub.HeadSHA)
	c.Assert(err, check.IsNil)
	c.Assert(int(count), check.Equals, len(tasks))

	for _, task := range tasks {
		runs, err := tc.RunsForTask(ctx, *task.Id, nil, nil)
		c.Assert(err, check.IsNil)
		count, err := tc.RunsForTaskCount(ctx, *task.Id)
		c.Assert(err, check.IsNil)
		c.Assert(len(runs), check.Equals, int(count))
	}

	runs, err := tc.Runs(ctx, stringp("erikh/test"), &sub.HeadSHA, nil, int64p(200)) // there should be 15 runs here.
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Not(check.Equals), 0)

	count, err = tc.RunsCount(ctx, stringp("erikh/test"), &sub.HeadSHA)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, int(count))

	runs, err = tc.RunsForTask(ctx, *runs[0].Task.Id, nil, nil)
	c.Assert(err, check.IsNil)

	for i := 0; i < len(runs); i++ {
		client.EXPECT().ErrorStatus(
			gomock.Any(),
			"erikh",
			"parent",
			gomock.Any(),
			sub.HeadSHA,
			gomock.Any(),
			gomock.Any()).Return(nil)
	}

	c.Assert(*runs[0].Task.Canceled, check.Equals, false)
	c.Assert(tc.CancelRun(ctx, *runs[0].Id), check.IsNil)
	run, err := tc.GetRun(ctx, *runs[0].Id)
	c.Assert(err, check.IsNil)
	c.Assert(*run.Task.Canceled, check.Equals, true)
}

func (us *uisvcSuite) TestAddDeleteCI(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	h, doneChan, tc, utc, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(tc.AddToCI(ctx, "f7u12"), check.NotNil)

	config.SetDefaultGithubClient(client, "erikh")

	client.EXPECT().MyRepositories(gomock.Any()).Return([]*gh.Repository{{FullName: gh.String("erikh/test")}}, nil)

	repos, err := tc.LoadRepositories(ctx, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Not(check.Equals), 0)

	c.Assert(tc.AddToCI(ctx, "erikh/not-real"), check.NotNil) // not saved

	client.EXPECT().TeardownHook(gomock.Any(), "erikh", "test", h.Config.HookURL).Return(errors.New("wat's up"))

	c.Assert(tc.AddToCI(ctx, "erikh/test"), check.ErrorMatches, "wat's up")

	client.EXPECT().TeardownHook(gomock.Any(), "erikh", "test", h.Config.HookURL).Return(nil)
	client.EXPECT().SetupHook(gomock.Any(), "erikh", "test", h.Config.HookURL, gomock.Any()).Return(errors.New("yep"))

	c.Assert(tc.AddToCI(ctx, "erikh/test"), check.ErrorMatches, "yep")

	client.EXPECT().TeardownHook(gomock.Any(), "erikh", "test", h.Config.HookURL).Return(nil)
	client.EXPECT().SetupHook(gomock.Any(), "erikh", "test", h.Config.HookURL, gomock.Any()).Return(nil)

	c.Assert(utc.AddToCI(ctx, "erikh/test"), check.ErrorMatches, utils.ErrInvalidAuth.Error())
	c.Assert(tc.AddToCI(ctx, "erikh/test"), check.IsNil)

	client.EXPECT().MyRepositories(gomock.Any()).Return([]*gh.Repository{{FullName: gh.String("erikh/test")}}, nil)

	visible, err := tc.Visible(ctx, stringp("erikh/test"))
	c.Assert(err, check.IsNil)
	c.Assert(len(visible), check.Equals, 1)
	c.Assert(tc.DeleteFromCI(ctx, "erikh/not-real"), check.NotNil)

	client.EXPECT().TeardownHook(gomock.Any(), "erikh", "test", h.Config.HookURL).Return(errors.New("wat's up"))
	c.Assert(tc.DeleteFromCI(ctx, "erikh/test"), check.NotNil)

	client.EXPECT().TeardownHook(gomock.Any(), "erikh", "test", h.Config.HookURL).Return(nil)
	c.Assert(utc.DeleteFromCI(ctx, "erikh/test"), check.NotNil)
	c.Assert(tc.DeleteFromCI(ctx, "erikh/test"), check.IsNil)
	c.Assert(tc.DeleteFromCI(ctx, "erikh/test"), check.ErrorMatches, "repo is not enabled")
}

func (us *uisvcSuite) TestSubscriptions(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, _, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	erikhClient := github.NewMockClient(gomock.NewController(c))
	config.SetDefaultGithubClient(erikhClient, "erikh")

	c.Assert(tc.Subscribe(ctx, "erikh/test"), check.NotNil)
	c.Assert(tc.Unsubscribe(ctx, "erikh/test"), check.NotNil)

	erikhClient.EXPECT().MyRepositories(gomock.Any()).Return([]*gh.Repository{{FullName: gh.String("erikh/test")}}, nil)

	repos, err := tc.LoadRepositories(ctx, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Not(check.Equals), 0)

	repos, err = tc.Subscribed(ctx, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 0)

	c.Assert(tc.Subscribe(ctx, "erikh/test"), check.IsNil)

	repos, err = tc.Subscribed(ctx, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 1)
	c.Assert(*repos[0].Name, check.Equals, "erikh/test")

	c.Assert(tc.Unsubscribe(ctx, "erikh/test"), check.IsNil)

	repos, err = tc.Subscribed(ctx, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 0)
}

func (us *uisvcSuite) TestVisibility(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, _, err := MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	_, err = us.datasvcClient.MakeUser("not-erikh")
	c.Assert(err, check.IsNil)
	_, err = us.datasvcClient.MakeUser("erikh-the-third")
	c.Assert(err, check.IsNil)

	c.Assert(us.datasvcClient.MakeRepo("not-erikh/private-test", "not-erikh", true, ""), check.IsNil)
	c.Assert(us.datasvcClient.MakeRepo("erikh/private-test", "erikh", true, ""), check.IsNil)
	c.Assert(us.datasvcClient.MakeRepo("erikh/public", "erikh", false, ""), check.IsNil)
	c.Assert(us.datasvcClient.MakeRepo("erikh-the-third/public", "erikh-the-third", false, ""), check.IsNil)

	repos, err := tc.Visible(ctx, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 3)

	for _, repo := range repos {
		c.Assert(repo.Name, check.Not(check.Equals), "not-erikh/private-test")
	}
}
