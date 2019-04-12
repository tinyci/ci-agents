package model

import (
	check "github.com/erikh/check"
	"github.com/google/go-github/github"
	gh "github.com/google/go-github/github"
)

var testRepository = &gh.Repository{
	FullName: gh.String("erikh/barbara"),
}

func (ms *modelSuite) TestRepositoryValidate(c *check.C) {
	users, err := ms.CreateUsers(2)
	c.Assert(err, check.IsNil)

	failures := []struct {
		name   string
		github *gh.Repository
		users  []*User
	}{
		{"", testRepository, users},
		{"erikh/barbara", nil, users},
		{"erikh/barbara", testRepository, nil},
		{"erikh/barbara", testRepository, []*User{}},
		{"erikh/barbara", &gh.Repository{FullName: github.String("something/else")}, users},
	}

	for i, failure := range failures {
		r := &Repository{
			Name:   failure.name,
			Owners: failure.users,
			Github: failure.github,
		}
		c.Assert(ms.model.Create(r).Error, check.NotNil, check.Commentf("iteration %d", i))
		c.Assert(ms.model.Save(r).Error, check.NotNil, check.Commentf("iteration %d", i))
	}

	r := &Repository{
		Name:   "erikh/barbara",
		Github: testRepository,
		Owners: users,
	}

	c.Assert(ms.model.Create(r).Error, check.IsNil)
	c.Assert(ms.model.Save(r).Error, check.IsNil)

	r2, err := ms.model.GetRepositoryByName(r.Name)
	c.Assert(err, check.IsNil)
	c.Assert(r2.Name, check.Equals, r.Name)
	c.Assert(len(r2.Owners), check.Equals, len(r.Owners))
	for x, owner := range r.Owners {
		c.Assert(owner.ID, check.Equals, r2.Owners[x].ID)
	}
	c.Assert(r2.Github, check.NotNil)
	c.Assert(r2.Github.GetFullName(), check.Equals, r.Name)
}

func (ms *modelSuite) TestRepositoryOwners(c *check.C) {
	r1, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	r2, err := ms.CreateRepository()
	c.Assert(err, check.IsNil)

	r1.Private = true
	c.Assert(ms.model.Save(r1).Error, check.IsNil)

	// r2's owners should not be able to see r1.
	_, err = ms.model.GetRepositoryByNameForUser(r1.Name, r2.Owners[0])
	c.Assert(err, check.NotNil)

	// but r1's owners can see r2, because it's not private
	tmp, err := ms.model.GetRepositoryByNameForUser(r2.Name, r1.Owners[0])
	c.Assert(err, check.IsNil)
	c.Assert(tmp.Name, check.Equals, r2.Name)

	// and of course r1's owners can see r1...
	tmp, err = ms.model.GetRepositoryByNameForUser(r1.Name, r1.Owners[0])
	c.Assert(err, check.IsNil)
	c.Assert(tmp.Name, check.Equals, r1.Name)

	list, err := ms.model.GetAllPublicRepos()
	c.Assert(err, check.IsNil)
	c.Assert(len(list), check.Equals, 1)
	c.Assert(list[0].Name, check.Equals, r2.Name)

	list, err = ms.model.GetPrivateReposForUser(r1.Owners[0])
	c.Assert(err, check.IsNil)
	c.Assert(len(list), check.Equals, 1)
	c.Assert(list[0].Name, check.Equals, r1.Name)

	list, err = ms.model.GetPrivateReposForUser(r2.Owners[0])
	c.Assert(err, check.IsNil)
	c.Assert(len(list), check.Equals, 0)
}

func (ms *modelSuite) TestAddEnableRepository(c *check.C) {
	owners, err := ms.CreateUsers(1)
	c.Assert(err, check.IsNil)

	err = ms.model.SaveRepositories([]*gh.Repository{
		{FullName: github.String("erikh/barbara")},
	}, owners[0].Username, false)
	c.Assert(err, check.IsNil)

	repo, err := ms.model.GetRepositoryByNameForUser("erikh/barbara", owners[0])
	c.Assert(err, check.IsNil)

	c.Assert(ms.model.EnableRepository(repo), check.IsNil)

	c.Assert(repo.ID, check.Not(check.Equals), 0)
	c.Assert(repo.Name, check.Equals, "erikh/barbara")
	c.Assert(repo.HookSecret, check.Not(check.Equals), "")
	c.Assert(repo.Disabled, check.Equals, false)

	c.Assert(ms.model.DisableRepository(repo), check.IsNil)
	c.Assert(repo.Disabled, check.Equals, true)

	tmp, err := ms.model.GetRepositoryByName("erikh/barbara")
	c.Assert(err, check.IsNil)
	c.Assert(tmp.Disabled, check.Equals, true)

	c.Assert(ms.model.EnableRepository(repo), check.IsNil)
	c.Assert(repo.Disabled, check.Equals, false)

	tmp, err = ms.model.GetRepositoryByName("erikh/barbara")
	c.Assert(err, check.IsNil)
	c.Assert(tmp.Disabled, check.Equals, false)
}
