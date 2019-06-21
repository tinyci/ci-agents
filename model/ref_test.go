package model

import (
	check "github.com/erikh/check"
)

func (ms *modelSuite) TestRefValidate(c *check.C) {
	r, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	failures := []struct {
		repository *Repository
		refName    string
		sha        string
	}{
		{nil, "refs/heads/master", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{r, "", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		{r, "refs/heads/master", ""},
		{r, "refs/heads/master", "abcdef"},
	}

	for i, failure := range failures {
		r := &Ref{
			Repository: failure.repository,
			RefName:    failure.refName,
			SHA:        failure.sha,
		}

		c.Assert(r.Validate(), check.NotNil)

		c.Assert(ms.model.PutRef(r), check.NotNil, check.Commentf("iteration %d", i))
		c.Assert(ms.model.Save(r).Error, check.NotNil, check.Commentf("iteration %d", i))
	}

	ref := &Ref{
		Repository: r,
		RefName:    "refs/heads/master",
		SHA:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	c.Assert(ref.Validate(), check.IsNil)

	c.Assert(ms.model.PutRef(ref), check.IsNil)
	c.Assert(ms.model.Save(ref).Error, check.IsNil)

	ref2 := &Ref{}
	c.Assert(ms.model.Where("id = ?", ref.ID).First(ref2).Error, check.IsNil)
	c.Assert(ref2.Repository.Name, check.Equals, ref.Repository.Name)
	c.Assert(ref2.Repository.ID, check.Equals, ref.Repository.ID)
	c.Assert(ref2.RefName, check.Equals, ref.RefName)
	c.Assert(ref2.SHA, check.Equals, ref.SHA)

	_, err = ms.model.GetRefByNameAndSHA(ref2.Repository.Name, "baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.NotNil)
	ref2, err = ms.model.GetRefByNameAndSHA(ref2.Repository.Name, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c.Assert(err, check.IsNil)
	c.Assert(ref2.ID, check.Equals, ref.ID)
	c.Assert(ref2.Validate(), check.IsNil)
}
