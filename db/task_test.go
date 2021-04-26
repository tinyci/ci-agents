package db

import (
	"encoding/json"
	"testing"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestTaskValidate(t *testing.T) {
	m := testInit(t)

	parent, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	fork, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)
	baseref := &models.Ref{
		RepositoryID: parent.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, baseref.Insert(ctx, m.db, boil.Infer()))

	headref := &models.Ref{
		RepositoryID: fork.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, headref.Insert(ctx, m.db, boil.Infer()))

	sub := &models.Submission{
		BaseRefID: baseref.ID,
		HeadRefID: null.Int64From(headref.ID),
	}

	assert.NilError(t, sub.Insert(ctx, m.db, boil.Infer()))

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	failures := []struct {
		sub          *models.Submission
		TaskSettings *types.TaskSettings
	}{
		{nil, ts},
		{sub, nil},
	}

	for _, failure := range failures {
		var subID int64
		if failure.sub != nil {
			subID = failure.sub.ID
		}

		content, err := json.Marshal(failure.TaskSettings)
		assert.NilError(t, err)

		task := &models.Task{
			SubmissionID: subID,
			TaskSettings: content,
		}

		assert.Assert(t, task.Insert(ctx, m.db, boil.Infer()) != nil)
	}

	content, err := json.Marshal(ts)
	assert.NilError(t, err)

	task := &models.Task{
		SubmissionID: sub.ID,
		TaskSettings: content,
	}

	assert.NilError(t, task.Insert(ctx, m.db, boil.Infer()))

	t2, err := models.FindTask(ctx, m.db, task.ID)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(t2.ID, task.ID))

	ts = &types.TaskSettings{}
	assert.NilError(t, t2.TaskSettings.Unmarshal(&ts))

	assert.Assert(t, cmp.Equal(len(ts.Runs), 1))

	// relation check
	t2sub, err := t2.Submission().One(ctx, m.db)
	assert.NilError(t, err)

	t2ref, err := t2sub.BaseRef().One(ctx, m.db)
	assert.NilError(t, err)

	t2repo, err := t2ref.Repository().One(ctx, m.db)
	assert.NilError(t, err)

	assert.Assert(t, cmp.Equal(t2repo.Name, parent.Name))

	t2ref, err = t2sub.HeadRef().One(ctx, m.db)
	assert.NilError(t, err)

	t2repo, err = t2ref.Repository().One(ctx, m.db)
	assert.NilError(t, err)

	assert.Assert(t, cmp.Equal(t2repo.Name, fork.Name))
}

/*
func TestTaskList(t *testing.T) {
	m := testInit(t)

	parent, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	fork, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	baseref := &models.Ref{
		RepositoryID: parent.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, baseref.Insert(ctx, m.db, boil.Infer()))

	headref := &models.Ref{
		RepositoryID: fork.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, headref.Insert(ctx, m.db, boil.Infer()))

	sub := &models.Submission{
		BaseRef: baseref.ID,
		HeadRef: null.Int64From(headref.ID),
	}

	assert.NilError(t, sub.Insert(ctx, m.db, boil.Infer()))

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	content, err := json.Marshal(ts)
	assert.NilError(t, err)

	task := &models.Task{
		SubmissionID: sub.ID,
		TaskSettings: content,
	}

	assert.NilError(t, task.Insert(ctx, m.db, boil.Infer()))

	tasks, err := m.ListTasks(ctx, fork.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", 0, 100)
	assert.NilError(t, err)

	sub, err := tasks[0].Submission().One(ctx, m.db)
	assert.NilError(t, err)

	ref, err := sub.BaseRef().One(ctx, m.db)
	assert.NilError(t, err)

	assert.Assert(t, ref.Sha != "")

	count, err := m.CountTasks(ctx, "", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(count, int64(1)))

	count, err = m.CountTasks(ctx, fork.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(count, int64(1)))

	us, err := m.CreateTestUsers(1)
	assert.NilError(t, err)
	tasks, err = m.ListSubscribedTasksForUser(ctx, us[0].ID, 0, 100)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(tasks), 0))

	assert.NilError(t, m.AddSubscriptionsForUser(ctx, us[0], []*models.Repository{parent}))

	tasks, err = ms.types.ListSubscribedTasksForUser(us[0].ID, 0, 100)
	assert.NilError(t, err)
	assert.Assert(len(tasks), check.Not(check.Equals), 0)
}

func TestTaskListSHAList(t *testing.T) {
	parent, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	fork, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	baseref := &Ref{
		Repository: parent,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.Assert(ms.types.Save(baseref).Error, check.IsNil)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	shas := map[string]int{}

	now := time.Now()
	fmt.Print("generating tasks... ")
	for i := 0; i < 1000; i++ {
		count := rand.Intn(50)
		sha := testutil.RandString(40)
		shas[sha] = count
		ref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		assert.Assert(ms.types.Save(ref).Error, check.IsNil)

		sub := &Submission{
			BaseRef: baseref,
			HeadRef: ref,
		}
		assert.Assert(ms.types.Save(sub).Error, check.IsNil)

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				Submission:   sub,
				TaskSettings: ts,
			}
			assert.Assert(ms.types.Save(t2).Error, check.IsNil)
		}
	}

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing single repo multi-SHA... ")
	var tasklen int
	for i := 0; i < 1000; i++ {
		t, err := ms.types.ListTasks(fork.Name, "", int64(i), 100)
		assert.NilError(t, err)
		if len(t) > 0 {
			tasklen += len(t)
		} else {
			break
		}

		var lastID int64
		for _, tsk := range t {
			if lastID != 0 {
				assert.Assert(tsk.ID < lastID, check.Equals, true)
				lastID = tsk.ID
			}
		}
	}

	var totalcount int

	for _, count := range shas {
		totalcount += count
	}

	assert.Assert(tasklen, check.Equals, totalcount)

	for sha, count := range shas {
		x, err := ms.types.CountTasks(fork.Name, sha)
		assert.NilError(t, err)
		assert.Assert(x, check.Equals, int64(count))
		tasks, err := ms.types.ListTasks(fork.Name, sha, 0, 100)
		assert.NilError(t, err)
		assert.Assert(len(tasks), check.Equals, count)
		for _, task := range tasks {
			assert.Assert(task.Submission.HeadRef.SHA, check.Equals, sha)
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))

	// totalcount is already calculated in this test, so re-use it.

	count, err := ms.types.CountTasks("", "")
	assert.NilError(t, err)
	assert.Assert(count, check.Equals, int64(totalcount))
}

func TestTaskListParents(t *testing.T) {
	fork, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	now := time.Now()
	fmt.Print("generating tasks... ")
	parents := map[string]int{}
	for i := 0; i < 1000; i++ {
		count := rand.Intn(100)
		parent, err := m.CreateTestRepository(ctx)
		assert.NilError(t, err)
		parents[parent.Name] = count
		sha := testutil.RandString(40)
		baseref := &Ref{
			Repository: parent,
			RefName:    "refs/heads/master",
			SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		}
		assert.Assert(ms.types.Save(baseref).Error, check.IsNil)

		headref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		assert.Assert(ms.types.Save(headref).Error, check.IsNil)

		sub := &Submission{BaseRef: baseref, HeadRef: headref}
		assert.Assert(ms.types.Save(sub).Error, check.IsNil)

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				Submission:   sub,
				TaskSettings: ts,
			}
			assert.Assert(ms.types.Save(t2).Error, check.IsNil)
		}
	}

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing multi parent any-SHA... ")
	for parentName, count := range parents {
		x, err := ms.types.CountTasks(parentName, "")
		assert.NilError(t, err)
		assert.Assert(x, check.Equals, int64(count))
		tasks, err := ms.types.ListTasks(parentName, "", 0, 100)
		assert.NilError(t, err)
		assert.Assert(len(tasks), check.Equals, count)
		var lastID int64
		for _, task := range tasks {
			assert.Assert(task.Submission.BaseRef.Repository.Name, check.Equals, parentName)
			if lastID != 0 {
				assert.Assert(task.ID < lastID, check.Equals, true)
				lastID = task.ID
			}
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))

	var totalcount int64
	for _, count := range parents {
		totalcount += int64(count)
	}

	count, err := ms.types.CountTasks("", "")
	assert.NilError(t, err)
	assert.Assert(count, check.Equals, totalcount)
}

func TestTaskListForks(t *testing.T) {
	parent, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	baseref := &Ref{
		Repository: parent,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.Assert(ms.types.Save(baseref).Error, check.IsNil)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"foobar": {
				Image:   "foo",
				Command: []string{"run", "me"},
			},
		},
	}

	now := time.Now()
	fmt.Print("generating tasks... ")
	forks := map[string]int{}
	for i := 0; i < 1000; i++ {
		count := rand.Intn(100)
		fork, err := m.CreateTestRepository(ctx)
		assert.NilError(t, err)
		forks[fork.Name] = count
		sha := testutil.RandString(40)
		headref := &Ref{
			Repository: fork,
			RefName:    "refs/heads/master",
			SHA:        sha,
		}
		assert.Assert(ms.types.Save(headref).Error, check.IsNil)

		sub := &Submission{HeadRef: headref, BaseRef: baseref}

		for x := count - 1; x >= 0; x-- {
			t2 := &Task{
				TaskSettings: ts,
				Submission:   sub,
			}
			assert.Assert(ms.types.Save(t2).Error, check.IsNil)
		}
	}

	fmt.Printf("duration: %v\n", time.Since(now))

	now = time.Now()
	fmt.Print("testing multi fork any-SHA... ")
	for forkName, count := range forks {
		x, err := ms.types.CountTasks(forkName, "")
		assert.NilError(t, err)
		assert.Assert(x, check.Equals, int64(count))
		tasks, err := ms.types.ListTasks(forkName, "", 0, 100)
		assert.NilError(t, err)
		assert.Assert(len(tasks), check.Equals, count)
		var lastID int64
		for _, task := range tasks {
			assert.Assert(task.Submission.HeadRef.Repository.Name, check.Equals, forkName)
			if lastID != 0 {
				assert.Assert(task.ID < lastID, check.Equals, true)
				lastID = task.ID
			}
		}
	}
	fmt.Printf("duration: %v\n", time.Since(now))

	var totalcount int64
	for _, count := range forks {
		totalcount += int64(count)
	}

	count, err := ms.types.CountTasks("", "")
	assert.NilError(t, err)
	assert.Assert(count, check.Equals, totalcount)
}
*/
