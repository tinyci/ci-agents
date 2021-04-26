package datasvc

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"time"

	check "github.com/erikh/check"
	"github.com/golang/mock/gomock"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/mocks/github"
	"github.com/tinyci/ci-agents/testutil"
	topTypes "github.com/tinyci/ci-agents/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ctx = context.Background()

func (ds *datasvcSuite) TestBasicUser(c *check.C) {
	username := testutil.RandString(8)
	resp, err := ds.client.MakeUser(username)
	c.Assert(err, check.IsNil)
	c.Assert(resp.Id, check.Not(check.Equals), int64(0))
	c.Assert(resp.Username, check.Equals, username)

	user, err := ds.client.Client().GetUser(ctx, username)
	c.Assert(err, check.IsNil)
	c.Assert(resp.ProtoReflect(), check.DeepEquals, user.ProtoReflect())

	user.TokenJSON, err = json.Marshal(&topTypes.OAuthToken{Token: "this is in this test"})
	c.Assert(err, check.IsNil)
	c.Assert(ds.client.Client().PatchUser(ctx, user), check.IsNil)
	user, err = ds.client.Client().GetUser(ctx, username)
	c.Assert(err, check.IsNil)
	c.Assert(user.TokenJSON, check.NotNil)

	var token topTypes.OAuthToken

	c.Assert(json.Unmarshal(user.TokenJSON, &token), check.IsNil)

	c.Assert(token.Token, check.Equals, "this is in this test")

	user.TokenJSON, err = json.Marshal(testutil.DummyToken)
	c.Assert(err, check.IsNil)
	c.Assert(ds.client.Client().PatchUser(ctx, user), check.IsNil)

	now := timestamppb.Now()
	user.LastScannedRepos = now
	c.Assert(ds.client.Client().PatchUser(ctx, user), check.IsNil)
	user, err = ds.client.Client().GetUser(ctx, username)
	c.Assert(err, check.IsNil)
	c.Assert(user.LastScannedRepos.AsTime(), check.Not(check.DeepEquals), now.AsTime())

	user.Username = "notcool"
	c.Assert(ds.client.Client().PatchUser(ctx, user), check.NotNil)

	users := map[string]*types.User{ // preseed with already written record
		username: resp,
	}

	for i := 0; i < 10; i++ {
		username := testutil.RandString(8)
		resp, err := ds.client.MakeUser(username)
		c.Assert(err, check.IsNil)
		users[username] = resp
	}

	usersResp, err := ds.client.Client().ListUsers(ctx)
	c.Assert(err, check.IsNil)

	for _, item := range usersResp.Users {
		m, ok := users[item.Username]
		c.Assert(ok, check.Equals, true, check.Commentf("%s", item.Username))
		c.Assert(m.ProtoReflect(), check.DeepEquals, item.ProtoReflect(), check.Commentf("%s", item.Username))
	}
}

func (ds *datasvcSuite) TestUserLoginToken(c *check.C) {
	username := testutil.RandString(8)
	resp, err := ds.client.MakeUser(username)
	c.Assert(err, check.IsNil)
	c.Assert(resp.Id, check.Not(check.Equals), int64(0))
	c.Assert(resp.Username, check.Equals, username)

	_, err = ds.client.Client().GetToken(ctx, "quux")
	c.Assert(err, check.NotNil)

	accessToken, err := ds.client.Client().GetToken(ctx, username)
	c.Assert(err, check.IsNil)

	_, err = ds.client.Client().GetToken(ctx, username)
	c.Assert(err, check.NotNil)

	c.Assert(ds.client.Client().DeleteToken(ctx, username), check.IsNil)

	accessToken2, err := ds.client.Client().GetToken(ctx, username)
	c.Assert(err, check.IsNil)

	c.Assert(accessToken, check.Not(check.Equals), accessToken2)

	_, err = ds.client.Client().ValidateToken(ctx, accessToken)
	c.Assert(err, check.NotNil)
	u, err := ds.client.Client().ValidateToken(ctx, accessToken2)
	c.Assert(err, check.IsNil)
	c.Assert(u.Username, check.Equals, username)
}

func (ds *datasvcSuite) TestUserErrors(c *check.C) {
	username := testutil.RandString(8)
	_, err := ds.client.MakeUser(username)
	c.Assert(err, check.IsNil)
	messages := []string{}

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("%d", i)
		messages = append(messages, msg)

		c.Assert(
			ds.client.Client().AddError(ctx, msg, username),
			check.IsNil,
		)
	}

	resp, err := ds.client.Client().GetErrors(ctx, username)
	c.Assert(err, check.IsNil)
	c.Assert(len(resp.Errors), check.Equals, 10)
	respStr := []string{}

	for _, item := range resp.Errors {
		respStr = append(respStr, string(item.Error))
	}

	sort.Strings(respStr)
	c.Assert(messages, check.DeepEquals, respStr)
}

func (ds *datasvcSuite) TestRepositories(c *check.C) {
	stringp := func(s string) *string { return &s }
	username := testutil.RandString(8)
	_, err := ds.client.MakeUser(username)
	c.Assert(err, check.IsNil)

	repos := []string{}

	for i := 0; i < 10; i++ {
		owner := testutil.RandString(8) + "_" + strconv.Itoa(i)
		repo := testutil.RandString(8)
		fullRepo := path.Join(owner, repo)
		c.Assert(ds.client.MakeRepo(fullRepo, username, i%2 == 0, ""), check.IsNil)
		repos = append(repos, fullRepo)
	}

	all, err := ds.client.Client().AllRepositories(ctx, username, nil)
	c.Assert(err, check.IsNil)
	c.Assert(len(all.List), check.Equals, 10)

	all, err = ds.client.Client().AllRepositories(ctx, username, stringp("_2"))
	c.Assert(err, check.IsNil)
	c.Assert(len(all.List), check.Equals, 1)

	private, err := ds.client.Client().PrivateRepositories(ctx, username, "")
	c.Assert(err, check.IsNil)
	c.Assert(len(private.List), check.Equals, 5)

	private, err = ds.client.Client().PrivateRepositories(ctx, username, "_2")
	c.Assert(err, check.IsNil)
	c.Assert(len(private.List), check.Equals, 1)

	public, err := ds.client.Client().PublicRepositories(ctx, "")
	c.Assert(err, check.IsNil)
	c.Assert(len(public.List), check.Equals, 5)

	public, err = ds.client.Client().PublicRepositories(ctx, "_3")
	c.Assert(err, check.IsNil)
	c.Assert(len(public.List), check.Equals, 1)

	repo, err := ds.client.Client().GetRepository(ctx, repos[0])
	c.Assert(err, check.IsNil)
	c.Assert(repo.Name, check.Equals, repos[0])
	c.Assert(repo.Disabled, check.Equals, true)
	c.Assert(repo.HookSecret, check.Equals, "")

	c.Assert(ds.client.Client().EnableRepository(ctx, username, repos[0]), check.IsNil)
	repo, err = ds.client.Client().GetRepository(ctx, repos[0])
	c.Assert(err, check.IsNil)
	c.Assert(repo.HookSecret, check.Not(check.Equals), "")
	c.Assert(repo.Disabled, check.Equals, false)

	c.Assert(ds.client.Client().DisableRepository(ctx, username, repos[0]), check.IsNil)
	repo, err = ds.client.Client().GetRepository(ctx, repos[0])
	c.Assert(err, check.IsNil)
	c.Assert(repo.Disabled, check.Equals, true)
}

func (ds *datasvcSuite) TestSubscriptions(c *check.C) {
	username := testutil.RandString(8)
	_, err := ds.client.MakeUser(username)
	c.Assert(err, check.IsNil)

	repos := []string{}

	for i := 0; i < 10; i++ {
		owner := testutil.RandString(8) + "_" + strconv.Itoa(i)
		repo := testutil.RandString(8)
		fullRepo := path.Join(owner, repo)
		c.Assert(ds.client.MakeRepo(fullRepo, username, false, ""), check.IsNil)
		repos = append(repos, fullRepo)
	}

	username2 := testutil.RandString(8)
	_, err = ds.client.MakeUser(username2)
	c.Assert(err, check.IsNil)

	for _, repo := range repos {
		c.Assert(
			ds.client.Client().AddSubscription(ctx, username2, repo),
			check.IsNil,
		)
	}

	subs, err := ds.client.Client().ListSubscriptions(ctx, username2, "")
	c.Assert(err, check.IsNil)
	c.Assert(len(subs.List), check.Equals, 10)

	subs, err = ds.client.Client().ListSubscriptions(ctx, username2, "_1")
	c.Assert(err, check.IsNil)
	c.Assert(len(subs.List), check.Equals, 1)

	c.Assert(ds.client.MakeRepo("erikh/private", username, true, ""), check.IsNil)
	c.Assert(
		ds.client.Client().AddSubscription(ctx, username2, "erikh/private"),
		check.NotNil,
	)
}

func (ds *datasvcSuite) TestRuns(c *check.C) {
	config.SetDefaultGithubClient(github.NewMockClient(gomock.NewController(c)), "")
	now := time.Now()
	qis := []*types.QueueItem{}
	for i := 0; i < 1000; i++ {
		qi, err := ds.client.MakeQueueItem()
		c.Assert(err, check.IsNil)
		qis = append(qis, qi)
	}

	fmt.Printf("Filling queue took %v\n", time.Since(now))

	count, err := ds.client.Client().RunCount(ctx, "", "")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1000))

	count, err = ds.client.Client().RunCount(ctx, qis[0].Run.Task.Submission.BaseRef.Repository.Name, "")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1))

	count, err = ds.client.Client().RunCount(ctx, qis[0].Run.Task.Submission.BaseRef.Repository.Name, "foo")
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(0))

	count, err = ds.client.Client().RunCount(ctx, qis[0].Run.Task.Submission.BaseRef.Repository.Name, qis[0].Run.Task.Submission.BaseRef.Sha)
	c.Assert(err, check.IsNil)
	c.Assert(count, check.Equals, int64(1))

	for i := 0; i < 10; i++ {
		runs, err := ds.client.Client().ListRuns(ctx, "", "", int64(i), 100)
		c.Assert(err, check.IsNil)
		c.Assert(len(runs.List), check.Equals, 100, check.Commentf("Loop: %d", i))
	}

	runs, err := ds.client.Client().ListRuns(ctx, "", "", 10, 100)
	c.Assert(err, check.IsNil)
	c.Assert(len(runs.List), check.Equals, 0)

	_, err = ds.client.Client().ListRuns(ctx, "", "", -1, 100)
	c.Assert(err, check.NotNil)
}

func (ds *datasvcSuite) TestQueue(c *check.C) {
	config.SetDefaultGithubClient(github.NewMockClient(gomock.NewController(c)), "")
	now := time.Now()

	count := 100
	goroutines := 10
	total := count * goroutines

	errChan := make(chan error, total)

	for i := 0; i < goroutines; i++ {
		go func() {
			for x := 0; x < count; x++ {
				_, err := ds.client.MakeQueueItem()
				errChan <- err
			}
		}()
	}

	for i := 0; i < total; i++ {
		err := <-errChan
		c.Assert(err, check.IsNil)
	}

	fmt.Printf("Filling queue took %v\n", time.Since(now))

	for i := 0; i < 1000; i++ {
		qi, err := ds.client.Client().NextQueueItem(ctx, "default", "hi")
		c.Assert(err, check.IsNil)
		c.Assert(qi.Running, check.Equals, true)
		c.Assert(qi.RunningOn, check.Equals, "hi")
		c.Assert(qi.StartedAt.IsValid(), check.Equals, true)
		c.Assert(qi.Run.StartedAt.IsValid(), check.Equals, true)
	}

	_, err := ds.client.Client().NextQueueItem(ctx, "default", "hi")
	c.Assert(err, check.NotNil)
}

func (ds *datasvcSuite) TestOAuth(c *check.C) {
	c.Assert(ds.client.Client().OAuthRegisterState(ctx, "asdf", []string{"repo"}), check.IsNil)
	res, err := ds.client.Client().OAuthValidateState(ctx, "asdf")
	c.Assert(err, check.IsNil)
	c.Assert(res, check.DeepEquals, []string{"repo"})
	_, err = ds.client.Client().OAuthValidateState(ctx, "asdf2")
	c.Assert(err, check.NotNil)
}

func (ds *datasvcSuite) TestRef(c *check.C) {
	username := testutil.RandString(8)
	_, err := ds.client.MakeUser(username)
	c.Assert(err, check.IsNil)

	ownerName, repoName := testutil.RandString(8), testutil.RandString(8)

	c.Assert(ds.client.MakeRepo(path.Join(ownerName, repoName), username, false, ""), check.IsNil)

	repo, err := ds.client.Client().GetRepository(ctx, path.Join(ownerName, repoName))
	c.Assert(err, check.IsNil)

	id, err := ds.client.Client().PutRef(ctx, &types.Ref{
		Repository: repo,
		RefName:    "heads/hi",
		Sha:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	})
	c.Assert(err, check.IsNil)
	c.Assert(id, check.Not(check.Equals), int64(0))

	ref, err := ds.client.Client().GetRefByNameAndSHA(ctx, path.Join(ownerName, repoName), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.IsNil)
	c.Assert(ref.Id, check.Equals, id)

	_, err = ds.client.Client().GetRefByNameAndSHA(ctx, path.Join(ownerName, repoName), "baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.NotNil)

	_, err = ds.client.Client().GetRefByNameAndSHA(ctx, path.Join(testutil.RandString(8), repoName), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.NotNil)

	_, err = ds.client.Client().GetRefByNameAndSHA(ctx, path.Join(ownerName, testutil.RandString(8)), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.NotNil)
}

func (ds *datasvcSuite) TestSubmissions(c *check.C) {
	var lastRepo string

	submissions := 25
	count := 100
	goroutines := 10
	total := count * goroutines

	for i := 0; i < submissions; i++ {
		username := testutil.RandString(8)
		user, err := ds.client.MakeUser(username)
		c.Assert(err, check.IsNil)

		ownerName, repoName := testutil.RandString(8), testutil.RandString(8)

		lastRepo = path.Join(ownerName, repoName) // for a later check

		c.Assert(ds.client.MakeRepo(path.Join(ownerName, repoName), username, false, ""), check.IsNil)

		repo, err := ds.client.Client().GetRepository(ctx, path.Join(ownerName, repoName))
		c.Assert(err, check.IsNil)

		id, err := ds.client.Client().PutRef(ctx, &types.Ref{
			Repository: repo,
			RefName:    "heads/hi",
			Sha:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		})
		c.Assert(err, check.IsNil)
		c.Assert(id, check.Not(check.Equals), int64(0))

		ref, err := ds.client.Client().GetRefByNameAndSHA(ctx, path.Join(ownerName, repoName), "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		c.Assert(err, check.IsNil)
		c.Assert(ref.Id, check.Equals, id)

		s, err := ds.client.Client().PutSubmission(ctx, &types.Submission{BaseRef: ref, HeadRef: ref, User: user})
		c.Assert(err, check.IsNil)
		c.Assert(s.Id, check.Not(check.Equals), int64(0))

		errChan := make(chan error, goroutines)
		taskChan := make(chan *types.Task, total)

		for i := 0; i < goroutines; i++ {
			go func() {
				for x := 0; x < count; x++ {
					runName := testutil.RandString(8)

					ts := &types.TaskSettings{
						Workdir:    "/tmp",
						Mountpoint: "/tmp",
						Runs: map[string]*types.RunSettings{
							runName: {
								Image:   "foo",
								Command: []string{"run", "me"},
								Queue:   "default",
							},
						},
					}

					task := &types.Task{
						Settings:   ts,
						Submission: s,
					}

					t, err := ds.client.Client().PutTask(ctx, task)
					if err != nil {
						fmt.Println(err)
						errChan <- err
						return
					}

					taskChan <- t
				}
			}()
		}

		tasks := []*types.Task{}
		for i := 0; i < total; i++ {
			select {
			case err := <-errChan:
				c.Assert(err, check.IsNil)
			case task := <-taskChan:
				tasks = append(tasks, task)
			}
		}

		sort.SliceStable(tasks, func(i, j int) bool { return tasks[i].Id < tasks[j].Id })

		s2, err := ds.client.Client().GetSubmissionByID(ctx, s.Id)
		c.Assert(err, check.IsNil)
		c.Assert(s2.Id, check.Equals, s.Id)
		c.Assert(s2.TasksCount, check.Equals, int64(1000))
		c.Assert(s2.CreatedAt.IsValid(), check.Equals, true)

		for x := int64(0); x < 10; x++ {
			tasks2, err := ds.client.Client().GetTasksForSubmission(ctx, s, x, 100)
			c.Assert(err, check.IsNil)

			for _, task := range tasks2.Tasks {
				c.Assert(task.Submission.CreatedAt.IsValid(), check.Equals, true)
			}

			sliceTasks := tasks[x*100 : (x+1)*100]

			for i := 0; i < 100; i++ {
				c.Assert(tasks2.Tasks[i].Id, check.Equals, sliceTasks[i].Id)
			}
		}
	}

	list, err := ds.client.Client().ListSubmissions(ctx, 0, 100, "", "")
	c.Assert(err, check.IsNil)
	c.Assert(len(list.Submissions), check.Equals, 25)

	list, err = ds.client.Client().ListSubmissions(ctx, 0, 100, lastRepo, "")
	c.Assert(err, check.IsNil)
	c.Assert(len(list.Submissions), check.Equals, 1)

	list, err = ds.client.Client().ListSubmissions(ctx, 0, 100, lastRepo, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.IsNil)
	c.Assert(len(list.Submissions), check.Equals, 1)

	_, err = ds.client.Client().ListSubmissions(ctx, 0, 100, lastRepo, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab")
	c.Assert(err, check.NotNil)

	_, err = ds.client.Client().ListSubmissions(ctx, 0, 100, "a/b", "")
	c.Assert(err, check.NotNil)

	subCount, err := ds.client.Client().CountSubmissions(ctx, "", "")
	c.Assert(err, check.IsNil)
	c.Assert(subCount, check.Equals, int64(25))

	subCount, err = ds.client.Client().CountSubmissions(ctx, lastRepo, "")
	c.Assert(err, check.IsNil)
	c.Assert(subCount, check.Equals, int64(1))

	subCount, err = ds.client.Client().CountSubmissions(ctx, lastRepo, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.IsNil)
	c.Assert(subCount, check.Equals, int64(1))

	_, err = ds.client.Client().CountSubmissions(ctx, lastRepo, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab")
	c.Assert(err, check.NotNil)

	_, err = ds.client.Client().CountSubmissions(ctx, "a/b", "")
	c.Assert(err, check.NotNil)
}
