package processors

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	check "github.com/erikh/check"
	"github.com/golang/mock/gomock"
	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/api/queuesvc/processors"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

func (qs *queuesvcSuite) getMock() *github.MockClientMockRecorder {
	return config.DefaultGithubClient.(*github.MockClient).EXPECT()
}

func (qs *queuesvcSuite) TestBadYAML(c *check.C) {
	// almost the same as testsubmission but with some dependencies logic in it
	_, err := qs.datasvcClient.MakeUser("erikh")
	c.Assert(err, check.IsNil)

	sub := &types.Submission{
		Parent:      "erikh/foobar",
		Fork:        "erikh/foobar2",
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar2", "erikh", false, "erikh/foobar"), check.IsNil)

	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	repoConfigBytes, e := ioutil.ReadFile("../../testdata/standard_repoconfig.yml")
	c.Assert(e, check.IsNil)

	taskBytes, e := ioutil.ReadFile("../../testdata/bad_task.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetRepository("erikh/foobar2").Return(&gh.Repository{FullName: gh.String("erikh/foobar2")}, nil)
	qs.getMock().GetRepository("erikh/foobar").Return(&gh.Repository{FullName: gh.String("erikh/foobar")}, nil)
	qs.getMock().GetRefs(sub.Fork, sub.HeadSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetRefs(sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetFile(sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)
	qs.getMock().GetDiffFiles(sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml"}, nil)
	qs.getMock().GetFileList(sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	qs.getMock().GetRepository(sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	qs.getMock().CommentError(sub.Parent, sub.PullRequest, gomock.Any()).Return(nil)

	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)

	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.datasvcClient.Client().EnableRepository("erikh", sub.Parent), check.IsNil)

	qs.getMock().FinishedStatus("erikh", "foobar", "*global*", "be3d26c478991039e951097f2c99f56b55396940", "url", false, gomock.Any()).Return(nil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.NotNil)
	runs, err := qs.datasvcClient.Client().ListRuns("", "", 0, 100)
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

	c.Assert(qs.queuesvcClient.SetMockSubmissionOnFork(qs.getMock(), sub, "erikh/foobar", "be3d26c478991039e951097f2c99f56b55396940"), check.IsNil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	c.Assert(qs.queuesvcClient.Client().Submit(ctx, sub), check.IsNil)
	defer cancel()

	runs, err := qs.datasvcClient.Client().ListRuns("erikh/foobar2", "", 0, 100)
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

	c.Assert(qs.queuesvcClient.SetMockSubmissionOnFork(qs.getMock(), sub, "erikh/foobar", "be3d26c478991039e951097f2c99f56b55396940"), check.IsNil)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	c.Assert(qs.queuesvcClient.Client().Submit(ctx, sub), check.NotNil)
	defer cancel()
}

func (qs *queuesvcSuite) TestManualSubmission(c *check.C) {
	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	sub := &types.Submission{
		Parent:      "erikh/foobar",
		Fork:        "erikh/foobar2",
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	msub := &types.Submission{
		Fork:        "erikh/foobar2",
		HeadSHA:     "master",
		SubmittedBy: "erikh",
		Manual:      true,
	}

	c.Assert(qs.queuesvcClient.SetUpSubmissionRepo(sub.Parent, ""), check.IsNil)
	qs.getMock().GetSHA(sub.Fork, "heads/master").Return("be3d26c478991039e951097f2c99f56b55396940", nil)
	qs.getMock().GetSHA(sub.Parent, "heads/master").Return("be3d26c478991039e951097f2c99f56b55396941", nil)
	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub), check.IsNil)
	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), msub), check.IsNil)

	qis, err := qs.datasvcClient.Client().ListRuns(sub.Fork, "be3d26c478991039e951097f2c99f56b55396940", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(qis), check.Equals, 10)
	for i := len(qis) - 1; i >= 0; i-- {
		// original sha from first run
		qs.getMock().
			ErrorStatus(
				"erikh",
				"foobar",
				qis[i].Name,
				"be3d26c478991039e951097f2c99f56b55396940",
				fmt.Sprintf("url/log/%d", qis[i].ID),
				errors.ErrRunCanceled,
			).Return(nil)
	}

	msub = &types.Submission{
		Fork:        "erikh/foobar2",
		HeadSHA:     "foobar",
		SubmittedBy: "erikh",
		Manual:      true,
	}

	sub = &types.Submission{
		Parent:      "erikh/foobar",
		Fork:        "erikh/foobar2",
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396942", // note the different sha is a disambiguator here.
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	qs.getMock().GetSHA(sub.Fork, "heads/foobar").Return("be3d26c478991039e951097f2c99f56b55396942", nil) // also here
	qs.getMock().GetSHA(sub.Parent, "heads/master").Return("be3d26c478991039e951097f2c99f56b55396941", nil)
	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub), check.IsNil)
	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), msub), check.IsNil)

	qis, err = qs.datasvcClient.Client().ListRuns(sub.Fork, "be3d26c478991039e951097f2c99f56b55396942", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(qis), check.Equals, 10)

	// cancellation tests

	qs.getMock().GetSHA(sub.Fork, "heads/foobar").Return("be3d26c478991039e951097f2c99f56b55396942", nil) // also here
	qs.getMock().GetSHA(sub.Parent, "heads/master").Return("be3d26c478991039e951097f2c99f56b55396941", nil)

	for i := len(qis) - 1; i >= 0; i-- {
		// original sha from first run
		qs.getMock().
			ErrorStatus(
				"erikh",
				"foobar",
				qis[i].Name,
				"be3d26c478991039e951097f2c99f56b55396942",
				fmt.Sprintf("url/log/%d", qis[i].ID),
				errors.ErrRunCanceled,
			).Return(nil)
	}

	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub), check.IsNil)
	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), msub), check.IsNil)

	qis, err = qs.datasvcClient.Client().ListRuns(sub.Fork, "be3d26c478991039e951097f2c99f56b55396942", 0, 100)
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
		Parent:      "erikh/foobar",
		Fork:        "erikh/foobar2",
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	c.Assert(qs.queuesvcClient.SetUpSubmissionRepo(sub.Parent, ""), check.IsNil)
	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub), check.IsNil)
	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err := qs.datasvcClient.Client().ListRuns("", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 10)

	sub = &types.Submission{
		Parent:      path.Join(testutil.RandString(8), testutil.RandString(8)),
		Fork:        path.Join(testutil.RandString(8), testutil.RandString(8)),
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	c.Assert(qs.queuesvcClient.SetUpSubmissionRepo(sub.Parent, ""), check.IsNil)
	c.Assert(qs.queuesvcClient.SetMockSubmissionSuccess(qs.getMock(), sub), check.IsNil)
	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err = qs.datasvcClient.Client().ListRuns("", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 20)
}

func (qs *queuesvcSuite) TestSubmission(c *check.C) {
	_, err := qs.datasvcClient.MakeUser("erikh")
	c.Assert(err, check.IsNil)

	sub := &types.Submission{
		Parent:      "erikh/foobar",
		Fork:        "erikh/foobar2",
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar2", "erikh", false, "erikh/foobar"), check.IsNil)

	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	repoConfigBytes, e := ioutil.ReadFile("../../testdata/standard_repoconfig.yml")
	c.Assert(e, check.IsNil)

	taskBytes, e := ioutil.ReadFile("../../testdata/standard_task.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetRepository("erikh/foobar2").Return(&gh.Repository{FullName: gh.String("erikh/foobar2")}, nil)
	qs.getMock().GetRepository("erikh/foobar").Return(&gh.Repository{FullName: gh.String("erikh/foobar")}, nil)
	qs.getMock().GetRefs(sub.Fork, sub.HeadSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetRefs(sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetFile(sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)
	qs.getMock().GetDiffFiles(sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar"}, nil)
	qs.getMock().GetFileList(sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	qs.getMock().GetRepository(sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)

	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "foo/task.yml").Return(taskBytes, nil)
	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)

	for _, name := range []string{"*root*", "foo"} {
		for x := 1; x <= 5; x++ {
			qs.getMock().PendingStatus("erikh", "foobar", fmt.Sprintf("%s:%d", name, x), sub.HeadSHA, "url")
		}
	}

	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.datasvcClient.Client().EnableRepository("erikh", sub.Parent), check.IsNil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err := qs.datasvcClient.Client().ListRuns("", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 10)

	tasks, err := qs.datasvcClient.Client().ListTasks("", "", 0, 100)
	c.Assert(err, check.IsNil)

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
		Parent:      "erikh/foobar",
		Fork:        "erikh/foobar2",
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar2", "erikh", false, "erikh/foobar"), check.IsNil)

	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	repoConfigBytes, e := ioutil.ReadFile("../../testdata/standard_repoconfig.yml")
	c.Assert(e, check.IsNil)

	taskBytes, e := ioutil.ReadFile("../../testdata/task_with_dependencies.yml")
	c.Assert(e, check.IsNil)

	depTaskBytes, e := ioutil.ReadFile("../../testdata/deps_only.yml")
	c.Assert(e, check.IsNil)

	standardTaskBytes, e := ioutil.ReadFile("../../testdata/standard_task.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetRepository("erikh/foobar2").Return(&gh.Repository{FullName: gh.String("erikh/foobar2")}, nil)
	qs.getMock().GetRepository("erikh/foobar").Return(&gh.Repository{FullName: gh.String("erikh/foobar")}, nil)
	qs.getMock().GetRefs(sub.Fork, sub.HeadSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetRefs(sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil)
	qs.getMock().GetFile(sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)
	qs.getMock().GetDiffFiles(sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml"}, nil)
	qs.getMock().GetFileList(sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil)
	qs.getMock().GetRepository(sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)

	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)
	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "bar/task.yml").Return(depTaskBytes, nil)
	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "foo/task.yml").Return(standardTaskBytes, nil)

	qs.getMock().PendingStatus("erikh", "foobar", "*root*:1", sub.HeadSHA, "url")
	for x := 1; x <= 5; x++ {
		qs.getMock().PendingStatus("erikh", "foobar", fmt.Sprintf("%s:%d", "foo", x), sub.HeadSHA, "url")
	}

	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)
	c.Assert(qs.datasvcClient.Client().EnableRepository("erikh", sub.Parent), check.IsNil)

	c.Assert(qs.queuesvcClient.Client().Submit(context.Background(), sub), check.IsNil)
	runs, err := qs.datasvcClient.Client().ListRuns("", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs), check.Equals, 6)

	tasks, err := qs.datasvcClient.Client().ListTasks("", "", 0, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(tasks), check.Equals, 2)
}

func (qs *queuesvcSuite) TestBasic(c *check.C) {
	_, err := qs.datasvcClient.MakeUser("erikh")
	c.Assert(err, check.IsNil)

	sub := &types.Submission{
		Parent:      "erikh/foobar",
		Fork:        "erikh/foobar2",
		HeadSHA:     "be3d26c478991039e951097f2c99f56b55396940",
		BaseSHA:     "be3d26c478991039e951097f2c99f56b55396941",
		PullRequest: 10,
	}

	_, _, _, err = processors.ManageRepositories(qs.queueHandler, sub)
	c.Assert(err, check.NotNil)

	qs.mkGithubClient(github.NewMockClient(gomock.NewController(c)))

	c.Assert(qs.datasvcClient.MakeRepo("erikh/foobar", "erikh", false, ""), check.IsNil)
	c.Assert(qs.datasvcClient.Client().EnableRepository("erikh", "erikh/foobar"), check.IsNil)

	gomock.InOrder(
		qs.getMock().GetRepository("erikh/foobar2").Return(&gh.Repository{FullName: gh.String("erikh/foobar2")}, nil),
		qs.getMock().GetRepository("erikh/foobar").Return(&gh.Repository{FullName: gh.String("erikh/foobar")}, nil),
	)

	parent, fork, _, err := processors.ManageRepositories(qs.queueHandler, sub)
	c.Assert(err, check.IsNil)

	c.Assert(parent.ID, check.Not(check.Equals), int64(0))
	c.Assert(parent.Name, check.Equals, "erikh/foobar")
	c.Assert(fork.ID, check.Not(check.Equals), int64(0))
	c.Assert(fork.Name, check.Equals, "erikh/foobar2")
	c.Assert(fork.AutoCreated, check.Equals, true)

	gomock.InOrder(
		qs.getMock().GetRefs(sub.Fork, sub.HeadSHA).Return([]string{"heads/master"}, nil),
	)

	forkRef, err := processors.ManageRefs(qs.queueHandler, config.DefaultGithubClient, fork, sub.HeadSHA)
	c.Assert(err, check.IsNil)

	c.Assert(forkRef.SHA, check.Equals, sub.HeadSHA)
	c.Assert(forkRef.RefName, check.Equals, "heads/master")

	gomock.InOrder(
		qs.getMock().GetRefs(sub.Parent, sub.BaseSHA).Return([]string{"heads/master"}, nil),
	)

	parentRef, err := processors.ManageRefs(qs.queueHandler, config.DefaultGithubClient, parent, sub.BaseSHA)
	c.Assert(err, check.IsNil)

	c.Assert(parentRef.SHA, check.Equals, sub.BaseSHA)
	c.Assert(parentRef.RefName, check.Equals, "heads/master")

	gomock.InOrder(
		qs.getMock().GetDiffFiles(sub.Parent, sub.BaseSHA, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "untested/bar", "untested/task.yml"}, nil),
		qs.getMock().GetFileList(sub.Fork, sub.HeadSHA).Return([]string{"task.yml", "foo/task.yml", "foo/bar", "bar/task.yml", "bar/quux"}, nil),
	)

	rc := &types.RepoConfig{
		IgnoreDirs: []string{"untested"},
	}

	processMap, err := processors.PickTasks(config.DefaultGithubClient, sub, forkRef, parent, rc)
	c.Assert(err, check.IsNil)

	for _, key := range []string{".", "foo"} {
		_, ok := processMap[key]
		c.Assert(ok, check.Equals, true)
	}

	_, ok := processMap["bar"]
	c.Assert(ok, check.Equals, false)

	_, ok = processMap["untested"]
	c.Assert(ok, check.Equals, false)

	repoConfigBytes, e := ioutil.ReadFile("../../testdata/standard_repoconfig.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetRepository(sub.Parent).Return(&gh.Repository{FullName: gh.String(sub.Parent)}, nil)
	qs.getMock().GetFile(sub.Parent, "refs/heads/master", "tinyci.yml").Return(repoConfigBytes, nil)

	repoConfig, err := processors.GetRepoConfig(config.DefaultGithubClient, sub)
	c.Assert(err, check.IsNil)

	c.Assert(repoConfig.Queue, check.Equals, "repoconfig")
	c.Assert(repoConfig.WorkDir, check.Equals, "/sw")

	taskBytes, e := ioutil.ReadFile("../../testdata/standard_task.yml")
	c.Assert(e, check.IsNil)

	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "task.yml").Return(taskBytes, nil)
	qs.getMock().GetFile(sub.Fork, sub.HeadSHA, "foo/task.yml").Return(taskBytes, nil)

	parts := strings.SplitN(parent.Name, "/", 2)
	c.Assert(len(parts), check.Equals, 2)

	for _, name := range []string{"*root*", "foo"} {
		for x := 1; x <= 5; x++ {
			gomock.InOrder(
				qs.getMock().PendingStatus(parts[0], parts[1], fmt.Sprintf("%s:%d", name, x), sub.HeadSHA, "url"),
			)
		}
	}

	qs.getMock().ClearStates(sub.Parent, sub.HeadSHA).Return(nil)
	qis, err := processors.GenerateQueueItems(
		context.Background(),
		qs.queueHandler,
		config.DefaultGithubClient,
		&processors.InternalSubmission{
			Sub:        sub,
			ProcessMap: processMap,
			RepoConfig: repoConfig,
			ParentRepo: parent,
			Ref:        forkRef,
		})

	c.Assert(err, check.IsNil)
	c.Assert(len(qis), check.Equals, 10)
}
