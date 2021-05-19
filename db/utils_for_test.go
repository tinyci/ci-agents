package db

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var testToken = &types.OAuthToken{
	Token: "123456",
}

func stringp(s string) *string {
	return &s
}

func (m *Model) CreateTestRef(ctx context.Context, r *models.Repository, name, sha string) (*models.Ref, error) {
	ref := &models.Ref{
		RepositoryID: r.ID,
		Sha:          sha,
		Ref:          name,
	}

	return ref, ref.Insert(ctx, m.db, boil.Infer())
}

// CreateUsers creates `count` random users.
func (m *Model) CreateTestUsers(ctx context.Context, count int) ([]*models.User, error) {
	users := []*models.User{}

	for i := 0; i < count; i++ {
		u, err := m.CreateUser(ctx, testutil.RandString(8), testToken)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

// CreateRepository creates a random repository.
func (m *Model) CreateTestRepository(ctx context.Context) (*models.Repository, error) {
	return m.CreateTestRepositoryWithName(ctx, path.Join(testutil.RandString(8), testutil.RandString(8)))
}

func (m *Model) CreateTestRepositoryWithName(ctx context.Context, name string) (*models.Repository, error) {
	owners, err := m.CreateTestUsers(ctx, 1)
	if err != nil {
		return nil, err
	}

	content, err := json.Marshal(github.Repository{FullName: github.String(name)})
	if err != nil {
		return nil, err
	}

	r := &models.Repository{
		Name:    name,
		Github:  content,
		OwnerID: owners[0].ID,
	}

	if err := r.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	return r, nil
}

func (m *Model) CreateTestTaskForSubmission(ctx context.Context, sub *models.Submission) (*models.Task, error) {
	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			"default": {
				Image:   "foo",
				Command: []string{"run", "me"},
				Queue:   "default",
			},
		},
	}

	content, err := json.Marshal(ts)
	if err != nil {
		return nil, err
	}

	task := &models.Task{
		TaskSettings: content,
		SubmissionID: sub.ID,
	}

	if err := task.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	content, err = json.Marshal(ts.Runs["default"])
	if err != nil {
		return nil, err
	}

	run := &models.Run{
		Name:        "default",
		RunSettings: content,
		TaskID:      task.ID,
	}

	if err := run.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	qi := &models.QueueItem{
		RunID:     run.ID,
		QueueName: "default",
	}

	if err := qi.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	return task, nil
}

func (m *Model) CreateTestRun(ctx context.Context) (*models.Run, error) {
	parent, err := m.CreateTestRepository(ctx)
	if err != nil {
		return nil, err
	}

	fork, err := m.CreateTestRepository(ctx)
	if err != nil {
		return nil, err
	}

	baseref := &models.Ref{
		RepositoryID: parent.ID,
		Ref:          testutil.RandString(8),
		Sha:          testutil.RandHexString(40),
	}

	if err := baseref.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	headref := &models.Ref{
		RepositoryID: fork.ID,
		Ref:          testutil.RandString(8),
		Sha:          testutil.RandHexString(40),
	}

	if err := headref.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	sub := &models.Submission{
		HeadRefID: null.Int64From(headref.ID),
		BaseRefID: baseref.ID,
	}

	if err := sub.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	runName := testutil.RandString(8)

	ts := &types.TaskSettings{
		Mountpoint: "/tmp",
		Runs: map[string]*types.RunSettings{
			runName: {
				Image:   "foo",
				Command: []string{"run", "me"},
				Queue:   "default",
			},
		},
	}

	content, err := json.Marshal(ts)
	if err != nil {
		return nil, err
	}

	task := &models.Task{
		TaskSettings: content,
		SubmissionID: sub.ID,
	}

	if err := task.Insert(ctx, m.db, boil.Infer()); err != nil {
		return nil, err
	}

	content, err = json.Marshal(ts.Runs[runName])
	if err != nil {
		return nil, err
	}

	run := &models.Run{
		Name:        runName,
		RunSettings: content,
		TaskID:      task.ID,
	}

	return run, run.Insert(ctx, m.db, boil.Infer())
}

func (m *Model) FillTestQueue(ctx context.Context, count int64) ([]*models.QueueItem, error) {
	fillstart := time.Now()
	qis := []*models.QueueItem{}

	for i := int64(1); i <= count; i++ {
		run, err := m.CreateTestRun(ctx)
		if err != nil {
			return nil, err
		}

		qi := &models.QueueItem{
			RunID:     run.ID,
			QueueName: "default",
		}

		if err := qi.Insert(ctx, m.db, boil.Infer()); err != nil {
			return nil, err
		}

		qis = append(qis, qi)
	}

	fmt.Println("Filling queue took", time.Since(fillstart))

	return qis, nil
}

func (m *Model) CreateTestSubmission(ctx context.Context, sub *types.Submission) (*models.Submission, error) {
	if sub.SubmittedBy != "" {
		if _, err := m.CreateUser(ctx, sub.SubmittedBy, testutil.DummyToken); err != nil {
			return nil, err
		}
	}

	if sub.Fork != "" {
		r, err := m.CreateTestRepositoryWithName(ctx, sub.Fork)
		if err != nil {
			r, err = m.GetRepositoryByName(ctx, sub.Fork)
			if err != nil {
				return nil, err
			}
		}
		if _, err := m.CreateTestRef(ctx, r, "dummy", sub.HeadSHA); err != nil {
			return nil, err
		}
	}

	r, err := m.CreateTestRepositoryWithName(ctx, sub.Parent)
	if err != nil {
		r, err = m.GetRepositoryByName(ctx, sub.Parent)
		if err != nil {
			return nil, err
		}
	}

	if _, err := m.CreateTestRef(ctx, r, "dummy", sub.BaseSHA); err != nil {
		return nil, err
	}

	s, err := m.NewSubmissionFromMessage(ctx, sub)
	if err != nil {
		return nil, err
	}

	return s, s.Insert(ctx, m.db, boil.Infer())
}
