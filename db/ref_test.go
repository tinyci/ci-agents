package db

import (
	"testing"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	"github.com/tinyci/ci-agents/db/models"
)

func TestRefValidate(t *testing.T) {
	m := testInit(t)

	r, err := m.CreateTestRepository(ctx)
	assert.NilError(t, err)

	failures := []struct {
		repository *models.Repository
		refName    string
		sha        string
	}{
		{nil, "refs/heads/master", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{r, "", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{r, "refs/heads/master", ""},
		{r, "refs/heads/master", "abcdef"},
	}

	for i, failure := range failures {
		repoID := int64(0)

		if failure.repository != nil {
			repoID = failure.repository.ID
		}

		r := &models.Ref{
			RepositoryID: repoID,
			Ref:          failure.refName,
			Sha:          failure.sha,
		}

		assert.Assert(t, m.PutRef(ctx, r) != nil, "iteration %d", i)
		assert.Assert(t, r.Insert(ctx, m.db, boil.Infer()) != nil, "iteration %d", i)
	}

	ref := &models.Ref{
		RepositoryID: r.ID,
		Ref:          "refs/heads/master",
		Sha:          "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	assert.NilError(t, m.PutRef(ctx, ref))

	ref2, err := models.Refs(qm.Where("id = ?", ref.ID)).One(ctx, m.db)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(ref2.RepositoryID, ref.RepositoryID))

	repo, err := ref2.Repository().One(ctx, m.db)
	assert.NilError(t, err)

	_, err = m.GetRefByNameAndSHA(ctx, repo.Name, "baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	assert.Assert(t, err != nil)
	ref2, err = m.GetRefByNameAndSHA(ctx, repo.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(ref2.ID, ref.ID))
}
