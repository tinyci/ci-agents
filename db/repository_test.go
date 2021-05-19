package db

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

var testRepository = &github.Repository{
	FullName: github.String("erikh/barbara"),
}

func TestRepositoryAssign(t *testing.T) {
	m := testInit(t)

	repo, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	users, err := m.CreateTestUsers(ctx, 2)
	assert.NilError(t, err)

	for _, user := range users {
		assert.NilError(t, m.AssignRepository(ctx, repo, user))

		repo, err = m.GetRepositoryByName(ctx, repo.Name)
		assert.NilError(t, err)

		owner, err := repo.Owner().One(ctx, m.db)
		assert.NilError(t, err)

		assert.Assert(t, cmp.Equal(owner.ID, user.ID))
		assert.Assert(t, cmp.Equal(owner.Username, user.Username))
	}
}

func TestRepositoryValidate(t *testing.T) {
	m := testInit(t)

	users, err := m.CreateTestUsers(ctx, 1)
	assert.NilError(t, err)

	failures := []struct {
		name   string
		github *github.Repository
		user   *models.User
	}{
		{"", testRepository, users[0]},
		{"erikh/barbara", &github.Repository{FullName: github.String("something/else")}, nil},
		{"erikh/barbara", &github.Repository{FullName: github.String("something/else")}, users[0]},
	}

	for i, failure := range failures {
		marshaled, err := json.Marshal(failure.github)
		assert.NilError(t, err)

		var userID int64
		if failure.user != nil {
			userID = failure.user.ID
		}

		r := &models.Repository{
			Name:    failure.name,
			OwnerID: userID,
			Github:  marshaled,
		}

		assert.Assert(t, r.Insert(ctx, m.db, boil.Infer()) != nil, fmt.Sprintf("iteration %d", i))
	}

	marshaled, err := json.Marshal(testRepository)
	assert.NilError(t, err)
	r := &models.Repository{
		Name:    "erikh/barbara",
		Github:  marshaled,
		OwnerID: users[0].ID,
	}

	assert.NilError(t, r.Insert(ctx, m.db, boil.Infer()))
	_, err = r.Update(ctx, m.db, boil.Infer())
	assert.NilError(t, err)

	r2, err := m.GetRepositoryByName(ctx, r.Name)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(r2.Name, r.Name))

	owner, err := r2.Owner().One(ctx, m.db)
	assert.NilError(t, err)

	assert.Assert(t, cmp.Equal(owner.ID, users[0].ID))
	assert.Assert(t, r2.Github != nil)

	gh := &github.Repository{}

	assert.NilError(t, json.Unmarshal(r2.Github, gh))

	assert.Assert(t, cmp.Equal(gh.GetFullName(), r.Name))
}

func TestRepositoryOwners(t *testing.T) {
	m := testInit(t)

	r1, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	r2, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	r1.Private = true
	_, err = r1.Update(ctx, m.db, boil.Infer())
	assert.NilError(t, err)

	// r2's owners should not be able to see r1.
	_, err = m.GetRepositoryByNameForUser(ctx, r1.Name, r2.OwnerID)
	assert.Assert(t, err != nil)

	// but r1's owners can see r2, because it's not private
	tmp, err := m.GetRepositoryByNameForUser(ctx, r2.Name, r1.OwnerID)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(tmp.Name, r2.Name))

	// and of course r1's owners can see r1...
	tmp, err = m.GetRepositoryByNameForUser(ctx, r1.Name, r1.OwnerID)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(tmp.Name, r1.Name))

	list, err := m.GetAllPublicRepos(ctx, nil)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(list), 1))
	assert.Assert(t, cmp.Equal(list[0].Name, r2.Name))

	searched, err := m.GetAllPublicRepos(ctx, stringp("not a match"))
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(searched), 0))

	searched, err = m.GetAllPublicRepos(ctx, &r2.Name)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(searched), 1))

	list, err = m.GetPrivateReposForUser(ctx, r1.OwnerID, nil)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(list), 1))
	assert.Assert(t, cmp.Equal(list[0].Name, r1.Name))

	searched, err = m.GetPrivateReposForUser(ctx, r1.OwnerID, stringp("not a match"))
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(searched), 0))

	searched, err = m.GetPrivateReposForUser(ctx, r1.OwnerID, &r1.Name)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(searched), 1))
	assert.Assert(t, cmp.Equal(searched[0].Name, r1.Name))

	list, err = m.GetPrivateReposForUser(ctx, r2.OwnerID, nil)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(list), 0))

	list, err = m.GetPrivateReposForUser(ctx, r2.OwnerID, &r2.Name)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(list), 0))

	list, err = m.GetVisibleReposForUser(ctx, r2.OwnerID, &r2.Name)
	assert.Assert(t, err)
	assert.Assert(t, cmp.Equal(len(list), 1))
}

func TestAddEnableRepository(t *testing.T) {
	m := testInit(t)

	owners, err := m.CreateTestUsers(ctx, 1)
	assert.NilError(t, err)

	err = m.SaveRepositories(ctx, []*github.Repository{
		{FullName: github.String("erikh/barbara")},
	}, owners[0].Username, false)
	assert.NilError(t, err)

	repo, err := m.GetRepositoryByNameForUser(ctx, "erikh/barbara", owners[0].ID)
	assert.NilError(t, err)

	assert.NilError(t, m.EnableRepository(ctx, repo, owners[0].ID))

	assert.Assert(t, repo.ID != 0)
	assert.Assert(t, cmp.Equal(repo.Name, "erikh/barbara"))
	assert.Assert(t, repo.HookSecret != "")
	assert.Assert(t, cmp.Equal(repo.Disabled.Bool, false))

	assert.NilError(t, m.DisableRepository(ctx, repo))
	assert.Assert(t, cmp.Equal(repo.Disabled.Bool, true))

	tmp, err := m.GetRepositoryByName(ctx, "erikh/barbara")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(tmp.Disabled.Bool, true))

	assert.NilError(t, m.EnableRepository(ctx, repo, owners[0].ID))
	assert.Assert(t, cmp.Equal(repo.Disabled.Bool, false))

	tmp, err = m.GetRepositoryByName(ctx, "erikh/barbara")
	assert.Assert(t, err)
	assert.Assert(t, cmp.Equal(tmp.Disabled.Bool, false))
}
