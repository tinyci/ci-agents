package test

import (
	"bytes"
	"io"
	"strings"
	"time"

	check "github.com/erikh/check"
	"github.com/golang/mock/gomock"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/testutil/testservers"
	"github.com/tinyci/ci-agents/types"

	gh "github.com/google/go-github/github"
)

func (us *uisvcSuite) TestErrors(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	errs, err := tc.Errors()
	c.Assert(err, check.IsNil)
	c.Assert(len(errs), check.Equals, 0)
}

func (us *uisvcSuite) TestLogAttach(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(us.assetsvcClient.Write(1, bytes.NewBufferString("this is a log")), check.IsNil)
	time.Sleep(100 * time.Millisecond)

	pr, pw := io.Pipe()
	buf := bytes.NewBuffer(nil)

	finished := make(chan struct{})
	go func() {
		defer close(finished)
		io.Copy(buf, pr)
	}()

	c.Assert(tc.LogAttach(1, pw), check.IsNil)
	pw.Close()
	<-finished
	c.Assert(strings.HasPrefix(buf.String(), "this is a log"), check.Equals, true)

	// XXX LogAttach does not error on missing ids -- https://github.com/tinyci/ci-agents/issues/270
}

func (us *uisvcSuite) TestDeleteToken(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(tc.DeleteToken(), check.IsNil)
	_, err = tc.Errors()
	c.Assert(err, check.ErrorMatches, ".*invalid authentication")
}

func (us *uisvcSuite) TestSubmit(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	client.EXPECT().GetRepository("erikh/not-real").Return(nil, errors.New("not found"))

	c.Assert(tc.Submit("erikh/not-real", "master", true), check.ErrorMatches, ".* not found")

	client.EXPECT().MyRepositories().Return([]*gh.Repository{{FullName: gh.String("erikh/parent")}}, nil)

	repos, err := tc.LoadRepositories()
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Not(check.Equals), 0)

	c.Assert(us.datasvcClient.Client().EnableRepository("erikh", "erikh/parent"), check.IsNil)

	client.EXPECT().GetRepository("erikh/not-real").Return(nil, errors.New("not found"))
	c.Assert(tc.Submit("erikh/not-real", "master", true), check.ErrorMatches, ".* not found")

	sub := &types.Submission{
		Parent:  "erikh/parent",
		Fork:    "erikh/test",
		BaseSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		HeadSHA: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Manual:  true,
		All:     true,
	}

	c.Assert(us.queuesvcClient.SetMockSubmissionSuccess(client.EXPECT(), sub), check.IsNil)

	client.EXPECT().GetSHA("erikh/parent", "heads/master").Return("", errors.New("not found"))
	client.EXPECT().GetSHA("erikh/test", "heads/master").Return("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", nil)

	c.Assert(tc.Submit("erikh/test", "master", true), check.ErrorMatches, ".*not found")

	client.EXPECT().GetSHA("erikh/parent", "heads/master").Return("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil)
	client.EXPECT().GetSHA("erikh/test", "heads/master").Return("", errors.New("not found"))

	c.Assert(tc.Submit("erikh/test", "master", true), check.ErrorMatches, ".*not found")

	c.Assert(us.queuesvcClient.SetMockSubmissionSuccess(client.EXPECT(), sub), check.IsNil)
	client.EXPECT().GetSHA("erikh/parent", "heads/master").Return("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", nil)
	client.EXPECT().GetSHA("erikh/test", "heads/master").Return("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", nil)
	client.EXPECT().ClearStates("erikh/parent", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb").Return(nil)

	c.Assert(tc.Submit("erikh/test", "master", true), check.IsNil)

	tasks, err := tc.Tasks("erikh/test", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 0, 200)
	c.Assert(err, check.IsNil)
	c.Assert(len(tasks), check.Not(check.Equals), 0)
	count, err := tc.TaskCount("erikh/test", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	c.Assert(err, check.IsNil)
	c.Assert(int(count), check.Equals, len(tasks))

	for _, task := range tasks {
		c.Assert(task.Parent.HookSecret, check.Equals, "")
		c.Assert(task.Ref.Repository.HookSecret, check.Equals, "")
		runs, err := tc.RunsForTask(task.ID, 0, 200)
		c.Assert(err, check.IsNil)
		count, err := tc.RunsForTaskCount(task.ID)
		c.Assert(err, check.IsNil)
		c.Assert(len(runs), check.Equals, int(count))
	}

	runs, err := tc.Runs("erikh/test", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 0, 200)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Not(check.Equals), 0)

	count, err = tc.RunsCount("erikh/test", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, int(count))

	for _, run := range runs {
		c.Assert(run.Task.Parent.HookSecret, check.Equals, "")
		c.Assert(run.Task.Ref.Repository.HookSecret, check.Equals, "")
	}

	runs, err = tc.RunsForTask(runs[0].Task.ID, 0, 200)
	c.Assert(err, check.IsNil)

	for i := 0; i < len(runs); i++ {
		client.EXPECT().ErrorStatus(
			"erikh",
			"parent",
			gomock.Any(),
			"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			gomock.Any(),
			gomock.Any()).Return(nil)
	}

	c.Assert(runs[0].Task.Canceled, check.Equals, false)
	c.Assert(tc.CancelRun(runs[0].ID), check.IsNil)
	run, err := tc.GetRun(runs[0].ID)
	c.Assert(err, check.IsNil)
	c.Assert(run.Task.Canceled, check.Equals, true)
}

func (us *uisvcSuite) TestAddDeleteCI(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	h, doneChan, tc, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(tc.AddToCI("f7u12"), check.NotNil)

	client.EXPECT().MyRepositories().Return([]*gh.Repository{{FullName: gh.String("erikh/test")}}, nil)

	repos, err := tc.LoadRepositories()
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Not(check.Equals), 0)

	c.Assert(tc.AddToCI("erikh/not-real"), check.NotNil) // not saved

	client.EXPECT().TeardownHook("erikh", "test", h.HookURL).Return(errors.New("wat's up"))

	c.Assert(tc.AddToCI("erikh/test"), check.ErrorMatches, "wat's up")

	client.EXPECT().TeardownHook("erikh", "test", h.HookURL).Return(nil)
	client.EXPECT().SetupHook("erikh", "test", h.HookURL, gomock.Any()).Return(errors.New("yep"))

	c.Assert(tc.AddToCI("erikh/test"), check.ErrorMatches, "yep")

	client.EXPECT().TeardownHook("erikh", "test", h.HookURL).Return(nil)
	client.EXPECT().SetupHook("erikh", "test", h.HookURL, gomock.Any()).Return(nil)

	c.Assert(tc.AddToCI("erikh/test"), check.IsNil)

	c.Assert(tc.DeleteFromCI("erikh/not-real"), check.NotNil)

	client.EXPECT().TeardownHook("erikh", "test", h.HookURL).Return(errors.New("wat's up"))
	c.Assert(tc.DeleteFromCI("erikh/test"), check.NotNil)

	client.EXPECT().TeardownHook("erikh", "test", h.HookURL).Return(nil)
	c.Assert(tc.DeleteFromCI("erikh/test"), check.IsNil)
	c.Assert(tc.DeleteFromCI("erikh/test"), check.ErrorMatches, "repo is not enabled")
}

func (us *uisvcSuite) TestSubscriptions(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	c.Assert(tc.Subscribe("erikh/test"), check.NotNil)
	c.Assert(tc.Unsubscribe("erikh/test"), check.NotNil)

	client.EXPECT().MyRepositories().Return([]*gh.Repository{{FullName: gh.String("erikh/test")}}, nil)

	repos, err := tc.LoadRepositories()
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Not(check.Equals), 0)

	repos, err = tc.Subscribed()
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 0)

	c.Assert(tc.Subscribe("erikh/test"), check.IsNil)

	repos, err = tc.Subscribed()
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 1)
	c.Assert(repos[0].Name, check.Equals, "erikh/test")

	c.Assert(tc.Unsubscribe("erikh/test"), check.IsNil)

	repos, err = tc.Subscribed()
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 0)
}

func (us *uisvcSuite) TestVisibility(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, tc, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	_, err = us.datasvcClient.MakeUser("not-erikh")
	c.Assert(err, check.IsNil)
	_, err = us.datasvcClient.MakeUser("erikh-the-third")
	c.Assert(err, check.IsNil)

	c.Assert(us.datasvcClient.MakeRepo("not-erikh/private-test", "not-erikh", true), check.IsNil)
	c.Assert(us.datasvcClient.MakeRepo("erikh/private-test", "erikh", true), check.IsNil)
	c.Assert(us.datasvcClient.MakeRepo("erikh/public", "erikh", false), check.IsNil)
	c.Assert(us.datasvcClient.MakeRepo("erikh-the-third/public", "erikh-the-third", false), check.IsNil)

	repos, err := tc.Visible()
	c.Assert(err, check.IsNil)
	c.Assert(len(repos), check.Equals, 3)

	for _, repo := range repos {
		c.Assert(repo.Name, check.Not(check.Equals), "not-erikh/private-test")
	}
}

func (us *uisvcSuite) TestNoAuthUserCreation(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	_, doneChan, _, err := testservers.MakeUIServer(client)
	c.Assert(err, check.IsNil)
	defer close(doneChan)

	u, err := us.datasvcClient.Client().GetUser("erikh")
	c.Assert(err, check.IsNil)
	c.Assert(u.Username, check.Equals, "erikh")
}
