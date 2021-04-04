package queuesvc

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	check "github.com/erikh/check"
	"github.com/golang/mock/gomock"
	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

var ctx = context.Background()

func (qs *queuesvcSuite) getMock() *github.MockClientMockRecorder {
	return config.DefaultGithubClient("").(*github.MockClient).EXPECT()
}

func (qs *queuesvcSuite) getUserMock() *github.MockClientMockRecorder {
	return config.DefaultGithubClient("erikh").(*github.MockClient).EXPECT()
}

func (qs *queuesvcSuite) TestBadYAML(c *check.C) {
	// almost the same as testsubmission but with some dependencies logic in it
	_, err := qs.datasvcClient.MakeUser("erikh")
	c.Assert(err, check.IsNil)

	sub := &types.Submission{
		Parent:   "erikh/foobar",
		Fork:     "erikh/foobar2",
		HeadSHA:  "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:  "be3d26c478991039e951097f2c99f56b55396941",
		TicketID: 10,
	}

	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar2", "erikh", false, "erikh/foobar"), check.IsNil)

	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	repoConfigBytes, e := ioutil.ReadFile("../testdata/standard_repoconfig.yml")
	c.Assert(e, check.IsNil)

	taskBytes, e := ioutil.ReadFile("../testdata/bad_task.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetRepository(gomock.Any(), "erikh/foobar2").Return(&gh.Repository{FullName: gh.String("erikh/foobar2")}, nil)
	qs.getMock().GetRepository(gomock.Any(), "erikh/foobar").Return(&gh.Repository{FullName: gh.String("erikh/foobar")}, nil)
	qs.getMock().GetSHA(gomock.Any(), "erikh/foobar2", "heads/master").Return(sub.HeadSHA, nil)
	qs.getMock().GetSHA(gomock.Any(), "erikh/foobar", "heads/master").Return(sub.BaseSHA, nil)
	qs.getMock().GetRefs(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetRefs(gomock.Any(), sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetFile(gomock.Any(), sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)
	qs.getMock().GetDiffFiles(gomock.Any(), sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml"}, nil)
	qs.getMock().GetFileList(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	qs.getMock().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	qs.getMock().CommentError(gomock.Any(), sub.Parent, sub.TicketID, gomock.Any()).Return(nil)

	qs.getMock().GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)

	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.datasvcClient.Client().EnableRepository(ctx, "erikh", sub.Parent), check.IsNil)

	qs.getMock().FinishedStatus(gomock.Any(), "erikh", "foobar", "*global*", "be3d26c478991039e951097f2c99f56b55396940", "url", false, gomock.Any()).Return(nil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.NotNil)
	runs, err := qs.datasvcClient.Client().ListRuns(ctx, "", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 0)
}

func (qs *queuesvcSuite) TestManualSubmissionOfAddedFork(c *check.C) {
	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	c.Assert(qs.queuesvcClient.SetUpSubmissionRepo("erikh/foobar2", "erikh/foobar"), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar2", "erikh", false, "erikh/foobar"), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar3", "erikh", false, "erikh/foobar"), check.IsNil)

	sub := &types.Submission{
		Fork:        "erikh/foobar2",
		HeadSHA:     "heads/master",
		All:         true,
		Manual:      true,
		SubmittedBy: "erikh",
	}

	c.Assert(qs.queuesvcClient.SetMockSubmissionOnFork(qs.getMock(), sub, "erikh/foobar", "be3d26c478991039e951097f2c99f56b55396940", ""), check.IsNil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	c.Assert(qs.queuesvcClient.Client().Submit(ctx, sub), check.IsNil)
	defer cancel()

	runs, err := qs.datasvcClient.Client().ListRuns(ctx, "erikh/foobar2", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 15)

	// not added
	sub = &types.Submission{
		Fork:        "erikh/foobar3",
		HeadSHA:     "heads/master",
		All:         true,
		Manual:      true,
		SubmittedBy: "erikh",
	}

	c.Assert(qs.queuesvcClient.SetMockSubmissionOnFork(qs.getMock(), sub, "erikh/foobar", "be3d26c478991039e951097f2c99f56b55396940", ""), check.IsNil)
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	c.Assert(qs.queuesvcClient.Client().Submit(ctx, sub), check.NotNil)
	defer cancel()
}

func (qs *queuesvcSuite) TestManualSubmission(c *check.C) {
	client := github.NewMockClient(gomock.NewController(c))
	qs.mkGithubClient(client)
	config.SetDefaultGithubClient(github.NewMockClient(gomock.NewController(c)), "erikh")

	sub := &types.Submission{
		Parent:   "erikh/foobar",
		Fork:     "erikh/foobar2",
		HeadSHA:  "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:  "be3d26c478991039e951097f2c99f56b55396941",
		TicketID: 10,
	}

	msub := &types.Submission{
		Fork:        "erikh/foobar2",
		HeadSHA:     "master",
		SubmittedBy: "erikh",
		Manual:      true,
	}

	c.Assert(qs.queuesvcClient.SetUpSubmissionRepo(sub.Parent, ""), check.IsNil)
	qs.getUserMock().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	qs.getUserMock().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{Fork: gh.Bool(true), FullName: gh.String(sub.Fork), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	qs.getMock().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	qs.getMock().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{Fork: gh.Bool(true), FullName: gh.String(sub.Fork), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	qs.getMock().GetSHA(gomock.Any(), sub.Fork, "heads/master").Return(sub.HeadSHA, nil)
	qs.getMock().GetSHA(gomock.Any(), sub.Parent, "heads/master").Return(sub.BaseSHA, nil)
	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub, "heads/master", ""), check.IsNil)
	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), msub), check.IsNil)

	qis, err := qs.datasvcClient.Client().ListRuns(ctx, sub.Fork, sub.HeadSHA, 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(qis), check.Equals, 10)
	for i := len(qis) - 1; i >= 0; i-- {
		// original sha from first run
		qs.getMock().
			ErrorStatus(
				gomock.Any(),
				"erikh",
				"foobar",
				qis[i].Name,
				sub.HeadSHA,
				fmt.Sprintf("url/log/%d", qis[i].ID),
				utils.ErrRunCanceled,
			).Return(nil)
	}

	msub = &types.Submission{
		Fork:        "erikh/foobar2",
		HeadSHA:     "foobar",
		SubmittedBy: "erikh",
		Manual:      true,
	}

	sub = &types.Submission{
		Parent:   "erikh/foobar",
		Fork:     "erikh/foobar2",
		HeadSHA:  "be3d26c478991039e951097f2c99f56b55396942", // note the different sha is a disambiguator here.
		BaseSHA:  "be3d26c478991039e951097f2c99f56b55396941",
		TicketID: 10,
	}

	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub, "heads/foobar", ""), check.IsNil)
	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), msub), check.IsNil)

	qis, err = qs.datasvcClient.Client().ListRuns(ctx, sub.Fork, "be3d26c478991039e951097f2c99f56b55396942", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(qis), check.Equals, 10)

	// cancellation tests

	qs.getMock().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	qs.getMock().GetRepository(gomock.Any(), sub.Fork).Return(&gh.Repository{FullName: gh.String(sub.Fork), Fork: gh.Bool(true), Parent: &gh.Repository{FullName: gh.String(sub.Parent)}}, nil)
	qs.getMock().GetSHA(gomock.Any(), sub.Fork, "heads/foobar").Return("be3d26c478991039e951097f2c99f56b55396942", nil) // also here
	qs.getMock().GetSHA(gomock.Any(), sub.Parent, "heads/master").Return("be3d26c478991039e951097f2c99f56b55396941", nil)

	for i := len(qis) - 1; i >= 0; i-- {
		// original sha from first run
		qs.getMock().
			ErrorStatus(
				gomock.Any(),
				"erikh",
				"foobar",
				qis[i].Name,
				"be3d26c478991039e951097f2c99f56b55396942",
				fmt.Sprintf("url/log/%d", qis[i].ID),
				utils.ErrRunCanceled,
			).Return(nil)
	}

	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub, "heads/foobar", ""), check.IsNil)
	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), msub), check.IsNil)

	qis, err = qs.datasvcClient.Client().ListRuns(ctx, sub.Fork, "be3d26c478991039e951097f2c99f56b55396942", 0, 100)
	c.Assert(err, check.IsNil)
	qis2 := []*model.Run{}

	for _, qi := range qis {
		if !qi.Task.Canceled {
			qis2 = append(qis2, qi)
		}
	}
	c.Assert(len(qis2), check.Equals, 10)
}

func (qs *queuesvcSuite) TestSubmission2(c *check.C) {
	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	sub := &types.Submission{
		Parent:   "erikh/foobar",
		Fork:     "erikh/foobar2",
		HeadSHA:  "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:  "be3d26c478991039e951097f2c99f56b55396941",
		TicketID: 10,
	}

	c.Assert(qs.queuesvcClient.SetUpSubmissionRepo(sub.Parent, ""), check.IsNil)
	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub, "heads/master", ""), check.IsNil)
	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err := qs.datasvcClient.Client().ListRuns(ctx, "", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 10)

	sub = &types.Submission{
		Parent:   path.Join(testutil.RandString(8), testutil.RandString(8)),
		Fork:     path.Join(testutil.RandString(8), testutil.RandString(8)),
		HeadSHA:  "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:  "be3d26c478991039e951097f2c99f56b55396941",
		TicketID: 10,
	}

	c.Assert(qs.queuesvcClient.SetUpSubmissionRepo(sub.Parent, ""), check.IsNil)
	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub, "heads/master", ""), check.IsNil)
	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err = qs.datasvcClient.Client().ListRuns(ctx, "", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 20)
}

func (qs *queuesvcSuite) TestSubmission(c *check.C) {
	_, err := qs.datasvcClient.MakeUser("erikh")
	c.Assert(err, check.IsNil)

	sub := &types.Submission{
		Parent:   "erikh/foobar",
		Fork:     "erikh/foobar2",
		HeadSHA:  "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:  "be3d26c478991039e951097f2c99f56b55396941",
		TicketID: 10,
	}

	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar2", "erikh", false, "erikh/foobar"), check.IsNil)

	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	repoConfigBytes, e := ioutil.ReadFile("../testdata/standard_repoconfig.yml")
	c.Assert(e, check.IsNil)

	taskBytes, e := ioutil.ReadFile("../testdata/standard_task.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetRepository(gomock.Any(), "erikh/foobar2").Return(&gh.Repository{FullName: gh.String("erikh/foobar2")}, nil)
	qs.getMock().GetRepository(gomock.Any(), "erikh/foobar").Return(&gh.Repository{FullName: gh.String("erikh/foobar")}, nil)
	qs.getMock().GetSHA(gomock.Any(), sub.Fork, "heads/master").Return(sub.HeadSHA, nil) // also here
	qs.getMock().GetSHA(gomock.Any(), sub.Parent, "heads/master").Return(sub.BaseSHA, nil)
	qs.getMock().GetRefs(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetRefs(gomock.Any(), sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetFile(gomock.Any(), sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)
	qs.getMock().GetDiffFiles(gomock.Any(), sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar"}, nil)
	qs.getMock().GetFileList(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	qs.getMock().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)

	qs.getMock().GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "foo/task.yml").Return(taskBytes, nil)
	qs.getMock().GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)

	for _, name := range []string{"*root*", "foo"} {
		for x := 1; x <= 5; x++ {
			qs.getMock().PendingStatus(gomock.Any(), "erikh", "foobar", fmt.Sprintf("%s:%d", name, x), sub.HeadSHA, "url")
		}
	}

	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.datasvcClient.Client().EnableRepository(ctx, "erikh", sub.Parent), check.IsNil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err := qs.datasvcClient.Client().ListRuns(ctx, "", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 10)

	tasks, err := qs.datasvcClient.Client().ListTasks(ctx, "", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(tasks[0].Runs, check.Not(check.Equals), int64(0))
	c.Assert(tasks[1].Runs, check.Not(check.Equals), int64(0))

	dirs := map[string]struct{}{}

	for _, task := range tasks {
		dirs[task.Path] = struct{}{}
	}

	for _, dir := range []string{".", "foo"} {
		_, ok := dirs[dir]
		c.Assert(ok, check.Equals, true, check.Commentf("%v", dir))
	}
}

func (qs *queuesvcSuite) TestDependencies(c *check.C) {
	// almost the same as testsubmission but with some dependencies logic in it
	_, err := qs.datasvcClient.MakeUser("erikh")
	c.Assert(err, check.IsNil)

	sub := &types.Submission{
		Parent:   "erikh/foobar",
		Fork:     "erikh/foobar2",
		HeadSHA:  "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:  "be3d26c478991039e951097f2c99f56b55396941",
		TicketID: 10,
	}

	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar2", "erikh", false, "erikh/foobar"), check.IsNil)

	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	repoConfigBytes, e := ioutil.ReadFile("../testdata/standard_repoconfig.yml")
	c.Assert(e, check.IsNil)

	taskBytes, e := ioutil.ReadFile("../testdata/task_with_dependencies.yml")
	c.Assert(e, check.IsNil)

	depTaskBytes, e := ioutil.ReadFile("../testdata/deps_only.yml")
	c.Assert(e, check.IsNil)

	standardTaskBytes, e := ioutil.ReadFile("../testdata/standard_task.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetRepository(gomock.Any(), "erikh/foobar2").Return(&gh.Repository{FullName: gh.String("erikh/foobar2")}, nil)
	qs.getMock().GetRepository(gomock.Any(), "erikh/foobar").Return(&gh.Repository{FullName: gh.String("erikh/foobar")}, nil)
	qs.getMock().GetSHA(gomock.Any(), "erikh/foobar2", "heads/master").Return(sub.HeadSHA, nil)
	qs.getMock().GetSHA(gomock.Any(), "erikh/foobar", "heads/master").Return(sub.BaseSHA, nil)
	qs.getMock().GetRefs(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetRefs(gomock.Any(), sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetFile(gomock.Any(), sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)
	qs.getMock().GetDiffFiles(gomock.Any(), sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml"}, nil)
	qs.getMock().GetFileList(gomock.Any(), sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	qs.getMock().GetRepository(gomock.Any(), sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)

	qs.getMock().GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)
	qs.getMock().GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "bar/task.yml").Return(depTaskBytes, nil)
	qs.getMock().GetFile(gomock.Any(), sub.Fork, sub.HeadSHA, "foo/task.yml").Return(standardTaskBytes, nil)

	qs.getMock().PendingStatus(gomock.Any(), "erikh", "foobar", "*root*:1", sub.HeadSHA, "url")
	for x := 1; x <= 5; x++ {
		qs.getMock().PendingStatus(gomock.Any(), "erikh", "foobar", fmt.Sprintf("%s:%d", "foo", x), sub.HeadSHA, "url")
	}

	qs.getMock().ClearStates(gomock.Any(), sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.datasvcClient.Client().EnableRepository(ctx, "erikh", sub.Parent), check.IsNil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err := qs.datasvcClient.Client().ListRuns(ctx, "", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 6)

	tasks, err := qs.datasvcClient.Client().ListTasks(ctx, "", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(tasks), check.Equals, 2)
	c.Assert(tasks[0].Runs, check.Not(check.Equals), int64(0))
	c.Assert(tasks[1].Runs, check.Not(check.Equals), int64(0))
}
