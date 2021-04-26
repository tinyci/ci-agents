package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/db/protoconv"
	topTypes "github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestQueueValidate(t *testing.T) {
	m := testInit(t)
	converter := protoconv.New(m.db)

	for iter := 0; iter < 100; iter++ {
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

		ref := &models.Ref{
			RepositoryID: fork.ID,
			Ref:          "refs/heads/master",
			Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		}

		assert.NilError(t, ref.Insert(ctx, m.db, boil.Infer()))

		sub := &models.Submission{
			BaseRefID: baseref.ID,
			HeadRefID: null.Int64From(ref.ID),
		}
		assert.NilError(t, sub.Insert(ctx, m.db, boil.Infer()))

		ts := &topTypes.TaskSettings{
			Mountpoint: "/tmp",
			Runs: map[string]*topTypes.RunSettings{
				"foobar": {
					Queue:   "default",
					Image:   "foo",
					Command: []string{"run", "me"},
				},
			},
		}

		content, err := json.Marshal(ts)
		assert.NilError(t, err)

		task := &models.Task{
			TaskSettings: content,
			SubmissionID: sub.ID,
		}

		assert.NilError(t, task.Insert(ctx, m.db, boil.Infer()))

		content, err = json.Marshal(ts.Runs["foobar"])
		assert.NilError(t, err)

		run := &models.Run{
			Name:        "foobar",
			RunSettings: content,
			TaskID:      task.ID,
		}

		assert.NilError(t, run.Insert(ctx, m.db, boil.Infer()))

		failures := []struct {
			queueName string
			run       *models.Run
		}{
			{"", run},
			{"default", nil},
		}

		for _, failure := range failures {
			var runID int64
			if failure.run != nil {
				runID = failure.run.ID
			}

			qi := &models.QueueItem{
				RunID:     runID,
				QueueName: failure.queueName,
			}

			assert.Assert(t, qi.Insert(ctx, m.db, boil.Infer()) != nil)
		}

		qi := &models.QueueItem{
			QueueName: "default",
			RunID:     run.ID,
		}

		assert.NilError(t, qi.Insert(ctx, m.db, boil.Infer()))

		qi2, err := models.FindQueueItem(ctx, m.db, qi.ID)
		assert.NilError(t, err)
		assert.Assert(t, qi.ID != 0)
		assert.Assert(t, cmp.Equal(qi.ID, qi2.ID))
		assert.Assert(t, cmp.Equal(qi2.QueueName, "default"))

		_, err = qi.Run().One(ctx, m.db)
		assert.NilError(t, err)

		qis := []*models.QueueItem{}

		for i := 0; i < 10; i++ {
			run, err = m.CreateTestRun(ctx)
			assert.NilError(t, err)

			qi := &models.QueueItem{
				QueueName: "default",
				RunID:     run.ID,
			}
			assert.NilError(t, err)
			assert.NilError(t, qi.Insert(ctx, m.db, boil.Infer()))
			qis = append(qis, qi)

			run, err = m.CreateTestRun(ctx)
			assert.NilError(t, err)

			qi = &models.QueueItem{
				QueueName: "default",
				RunID:     run.ID,
			}
			assert.NilError(t, qi.Insert(ctx, m.db, boil.Infer()))

			qis = append(qis, qi)
		}

		i, err := m.QueueTotalCount(ctx)
		assert.NilError(t, err)
		assert.Assert(t, cmp.Equal(i, int64(21*(iter+1)))) // relative to test iteration

		for _, qi := range qis {
			ro := "test"
			qi.RunningOn = null.StringFrom(ro)
			tmp, err := converter.ToProto(ctx, qi)
			assert.NilError(t, err)
			tmp, err = converter.FromProto(ctx, tmp.(*types.QueueItem))
			assert.NilError(t, err)

			qi2 := tmp.(*models.QueueItem)

			assert.Assert(t, cmp.Equal(qi2.RunID, qi.RunID))
			assert.Assert(t, cmp.Equal(qi2.ID, qi.ID))
			assert.Assert(t, cmp.Equal(qi2.QueueName, qi.QueueName))
			assert.Assert(t, cmp.Equal(qi2.Running, qi.Running))
			assert.Assert(t, cmp.Equal(qi2.RunningOn.String, "test"))
			assert.Assert(t, cmp.Equal(qi2.StartedAt, qi.StartedAt))
		}
	}
}

func TestQueueManipulation(t *testing.T) {
	m := testInit(t)

	var (
		firstID, lastRunID int64
	)

	fillstart := time.Now()
	for i := 1; i <= 1000; i++ {
		run, err := m.CreateTestRun(ctx)
		assert.NilError(t, err)

		qi := &models.QueueItem{
			RunID:     run.ID,
			QueueName: "default",
		}

		assert.NilError(t, qi.Insert(ctx, m.db, boil.Infer()))
		if firstID == 0 {
			firstID = qi.ID
		}
		lastRunID = run.ID
	}
	fmt.Println("Filling queue took", time.Since(fillstart))

	count, err := m.QueueTotalCount(ctx)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(count, int64(1000)))

	_, err = m.QueueList(ctx, -1, 100)
	assert.Assert(t, err != nil)

	for i := 0; i < 10; i++ {
		list, err := m.QueueList(ctx, int64(i), 100)
		assert.NilError(t, err)
		assert.Assert(t, cmp.Equal(len(list), 100))
		for _, qi := range list {
			assert.Assert(t, !qi.Running)
			assert.Assert(t, !qi.RunningOn.Valid)

			run, err := qi.Run().One(ctx, m.db)
			assert.NilError(t, err)

			task, err := run.Task().One(ctx, m.db)
			assert.NilError(t, err)

			sub, err := task.Submission().One(ctx, m.db)
			assert.NilError(t, err)

			baseRef, err := sub.BaseRef().One(ctx, m.db)
			assert.NilError(t, err)

			repo, err := baseRef.Repository().One(ctx, m.db)
			assert.NilError(t, err)

			count, err := m.QueueTotalCountForRepository(ctx, repo.ID)
			assert.NilError(t, err)
			assert.Assert(t, cmp.Equal(count, int64(1))) // repo names are uniq'd

			_, err = m.QueueListForRepository(ctx, repo.ID, -1, 100)
			assert.Assert(t, err != nil)

			tmp, err := m.QueueListForRepository(ctx, repo.ID, 0, 100)
			assert.NilError(t, err)
			assert.Assert(t, cmp.DeepEqual(tmp[0], qi)) // repo names are uniq'd
		}
	}

	list, err := m.QueueList(ctx, 0, 100)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(list), 100))

	list2, err := m.QueueList(ctx, 1, 100)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(list), 100))

	// validate no overlap.
	for _, qi := range list {
		for _, qi2 := range list2 {
			assert.Assert(t, qi.ID != qi2.ID)
		}
	}

	start := time.Now()

	for i := lastRunID - 999; i < lastRunID; i++ {
		qi, err := m.NextQueueItem(ctx, "hostname", "") // testing empty string handling
		assert.NilError(t, err)

		run, err := qi.Run().One(ctx, m.db)
		assert.NilError(t, err)

		assert.Assert(t, cmp.Equal(qi.ID, firstID), fmt.Sprintf("%d", lastRunID-i))
		assert.Assert(t, cmp.Equal(qi.RunID, int64(i))) // ensures same order
		assert.Assert(t, run.Name != "")
		assert.Assert(t, qi.Running)
		assert.Assert(t, cmp.Equal(qi.RunningOn.String, "hostname"))
		assert.Assert(t, qi.StartedAt.Valid)

		task, err := run.Task().One(ctx, m.db)
		assert.NilError(t, err)

		sub, err := task.Submission().One(ctx, m.db)
		assert.NilError(t, err)

		baseRef, err := sub.BaseRef().One(ctx, m.db)
		assert.NilError(t, err)

		repo, err := baseRef.Repository().One(ctx, m.db)
		assert.NilError(t, err)

		assert.Assert(t, repo != nil) // checking the ORM works
		assert.Assert(t, cmp.Equal(run.RanOn.String, "hostname"))
		firstID++
	}

	fmt.Println("Iterating queue took", time.Since(start))
}

func TestQueueNamed(t *testing.T) {
	m := testInit(t)

	names := []string{
		"default", // testing default functionality
		"foo",
		"bar",
		"quux",
	}

	var firstID int64

	fillstart := time.Now()
	for i := 1; i <= 1000; i++ {
		for _, name := range names {
			run, err := m.CreateTestRun(ctx)
			assert.NilError(t, err)

			qi := &models.QueueItem{
				RunID:     run.ID,
				QueueName: name,
			}

			assert.NilError(t, qi.Insert(ctx, m.db, boil.Infer()))
			if firstID == 0 {
				firstID = run.ID
			}
		}
	}
	fmt.Println("Filling queue took", time.Since(fillstart))

	start := time.Now()

	for i := firstID; i < firstID+int64(1000*len(names)); i++ {
		qi, err := m.NextQueueItem(ctx, "hostname", names[(int64(i)-firstID)%int64(len(names))])
		assert.NilError(t, err)

		run, err := qi.Run().One(ctx, m.db)
		assert.NilError(t, err)

		assert.Assert(t, cmp.Equal(run.ID, i))
		assert.Assert(t, run.Name != "")
		assert.Assert(t, qi.Running)
		assert.Assert(t, cmp.Equal(qi.RunningOn.String, "hostname"))
		assert.Assert(t, qi.StartedAt.Valid)

		// checking relationships
		task, err := run.Task().One(ctx, m.db)
		assert.NilError(t, err)

		sub, err := task.Submission().One(ctx, m.db)
		assert.NilError(t, err)

		ref, err := sub.BaseRef().One(ctx, m.db)
		assert.NilError(t, err)

		_, err = ref.Repository().One(ctx, m.db)
		assert.NilError(t, err)

		assert.Assert(t, cmp.Equal(run.RanOn.String, "hostname"))
	}

	fmt.Println("Iterating queue took", time.Since(start))
}

func TestQueueConcurrent(t *testing.T) {
	m := testInit(t)
	count := int64(1000)
	goRoutines := 10

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queueChan := make(chan *models.QueueItem, goRoutines)
	errChan := make(chan error, goRoutines)

	var firstID int64

	fillstart := time.Now()
	qis := []*models.QueueItem{}
	for i := int64(1); i <= count; i++ {
		run, err := m.CreateTestRun(ctx)
		assert.NilError(t, err)

		qi := &models.QueueItem{
			RunID:     run.ID,
			QueueName: "default",
		}

		qis = append(qis, qi)

		if firstID == 0 {
			firstID = run.ID
		}
	}

	assert.NilError(t, m.QueuePipelineAdd(ctx, qis))

	fmt.Println("Filling queue took", time.Since(fillstart))

	start := time.Now()

	for i := 0; i < goRoutines; i++ {
		go func(i int) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				qi, err := m.NextQueueItem(ctx, "hostname", "default")
				if err != nil {
					if errors.Is(err, utils.ErrNotFound) {
						return
					}

					errChan <- err
					return
				}

				queueChan <- qi
			}
		}(i)
	}

	for i := firstID; i <= firstID+count-1; i++ {
		select {
		case err := <-errChan:
			assert.NilError(t, err)
		case qi := <-queueChan:
			run, err := qi.Run().One(ctx, m.db)
			assert.NilError(t, err)
			assert.Assert(t, cmp.Equal(qi.RunID, i), fmt.Sprintf("%d", i-firstID)) // ensures same order
			assert.Assert(t, run.Name != "")
			assert.Assert(t, qi.Running)
			assert.Assert(t, cmp.Equal(qi.RunningOn.String, "hostname"))
			assert.Assert(t, qi.StartedAt.Valid)

			// check dep chain
			task, err := run.Task().One(ctx, m.db)
			assert.NilError(t, err)

			sub, err := task.Submission().One(ctx, m.db)
			assert.NilError(t, err)

			ref, err := sub.BaseRef().One(ctx, m.db)
			assert.NilError(t, err)

			repo, err := ref.Repository().One(ctx, m.db)
			assert.NilError(t, err)

			assert.Assert(t, repo != nil)
			assert.Assert(t, cmp.Equal(run.RanOn.String, "hostname"))
		}
	}

	cancel()

	fmt.Println("Iterating queue took", time.Since(start))
}

func TestQueueNamedConcurrent(t *testing.T) {
	m := testInit(t)

	names := []string{
		"default",
		"foo",
		"bar",
		"quux",
	}

	multiplier := 2
	count := 1000

	queueChan := make(chan *models.QueueItem, len(names)*multiplier)
	errChan := make(chan error, len(names)*multiplier)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var y int

	fillstart := time.Now()
	qis := []*models.QueueItem{}

	for i := 1; i <= count; i++ {
		for _, name := range names {
			y++
			run, err := m.CreateTestRun(ctx)
			assert.NilError(t, err)

			qi := &models.QueueItem{
				RunID:     run.ID,
				QueueName: name,
			}

			qis = append(qis, qi)
		}
	}

	assert.NilError(t, m.QueuePipelineAdd(ctx, qis))

	fmt.Println("Filling queue took", time.Since(fillstart))

	start := time.Now()

	for _, name := range names {
		for i := 0; i < multiplier; i++ {
			go func(name string) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}

					qi, err := m.NextQueueItem(ctx, "hostname", name)
					if err != nil {
						errChan <- err
						return
					}

					queueChan <- qi
				}
			}(name)
		}
	}

	for i := 1; i <= count; i++ {
		select {
		case err := <-errChan:
			assert.NilError(t, err)
		case qi := <-queueChan:
			run, err := qi.Run().One(ctx, m.db)
			assert.NilError(t, err)
			assert.Assert(t, qi.RunID != 0, fmt.Sprintf("%d", i)) // ensures same order
			assert.Assert(t, run.Name != "")
			assert.Assert(t, qi.Running)
			assert.Assert(t, cmp.Equal(qi.RunningOn.String, "hostname"))
			assert.Assert(t, qi.StartedAt.Valid)

			// check dep chain
			task, err := run.Task().One(ctx, m.db)
			assert.NilError(t, err)

			sub, err := task.Submission().One(ctx, m.db)
			assert.NilError(t, err)

			ref, err := sub.BaseRef().One(ctx, m.db)
			assert.NilError(t, err)

			repo, err := ref.Repository().One(ctx, m.db)
			assert.NilError(t, err)

			assert.Assert(t, repo != nil)
			assert.Assert(t, cmp.Equal(run.RanOn.String, "hostname"))
		}
	}

	fmt.Println("Iterating queue took", time.Since(start))
}
