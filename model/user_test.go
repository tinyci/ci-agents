package model

import (
	"encoding/json"
	"math/rand"
	"sort"
	"strings"

	"errors"

	check "github.com/erikh/check"
	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

var testToken = &types.OAuthToken{
	Token: "123456",
}

func (ms *modelSuite) TestUserTokens(c *check.C) {
	usr, err := ms.model.CreateUser("roz", testToken)
	c.Assert(err, check.IsNil)
	testTokenJSON, err2 := json.Marshal(testToken)
	c.Assert(err2, check.IsNil)

	usrJSON, err2 := json.Marshal(usr)
	c.Assert(err2, check.IsNil)

	c.Assert(strings.Contains(string(usrJSON), string(testTokenJSON)), check.Equals, false)

	typesUser := usr.ToProto()
	usr2, err := NewUserFromProto(typesUser)

	c.Assert(err, check.IsNil)
	usr2JSON, err2 := json.Marshal(usr2)

	c.Assert(err2, check.IsNil)
	c.Assert(usr2JSON, check.DeepEquals, usrJSON)

	usrList := []*User{usr, usr2}
	usrListJSON, err2 := json.Marshal(usrList)

	c.Assert(err2, check.IsNil)

	usr2List, err := MakeUsers(MakeUserList(usrList))
	c.Assert(err, check.IsNil)
	usr2ListJSON, err2 := json.Marshal(usr2List)

	c.Assert(err2, check.IsNil)
	c.Assert(usrListJSON, check.DeepEquals, usr2ListJSON)
	c.Assert(strings.Contains(string(usrListJSON), string(testTokenJSON)), check.Equals, false)
}

func (ms *modelSuite) TestCapabilityModification(c *check.C) {
	caps := []Capability{CapabilityCancel, CapabilityModifyCI, CapabilityModifyUser, CapabilitySubmit}

	strCaps := []string{}
	for _, cap := range caps {
		strCaps = append(strCaps, string(cap))
	}

	fixedCaps := map[string][]string{
		"erikh2": strCaps,
	}
	u, err := ms.model.CreateUser("erikh", &types.OAuthToken{Token: "dummy"})
	c.Assert(err, check.IsNil)

	for _, cap := range caps {
		res, err := ms.model.HasCapability(u, cap, fixedCaps)
		c.Assert(err, check.IsNil)
		c.Assert(res, check.Equals, false)
		c.Assert(ms.model.AddCapabilityToUser(u, cap), check.IsNil)
		res, err = ms.model.HasCapability(u, cap, fixedCaps)
		c.Assert(err, check.IsNil)
		c.Assert(res, check.Equals, true)

		getCaps, err := ms.model.GetCapabilities(u, fixedCaps)
		c.Assert(err, check.IsNil)
		listCaps := []string{}

		for _, cap := range getCaps {
			listCaps = append(listCaps, string(cap))
		}

		sort.Strings(listCaps)

		c.Assert([]string{string(cap)}, check.DeepEquals, listCaps)

		c.Assert(ms.model.RemoveCapabilityFromUser(u, cap), check.IsNil)
		res, err = ms.model.HasCapability(u, cap, fixedCaps)
		c.Assert(err, check.IsNil)
		c.Assert(res, check.Equals, false)
	}

	u2, err := ms.model.CreateUser("erikh2", &types.OAuthToken{Token: "dummy"})
	c.Assert(err, check.IsNil)

	for _, cap := range caps {
		res, err := ms.model.HasCapability(u2, cap, fixedCaps)
		c.Assert(err, check.IsNil)
		c.Assert(res, check.Equals, true)
	}

	getCaps, err := ms.model.GetCapabilities(u2, fixedCaps)
	c.Assert(err, check.IsNil)

	listCaps := []string{}

	for _, cap := range getCaps {
		listCaps = append(listCaps, string(cap))
	}

	sort.Strings(fixedCaps["erikh2"])
	sort.Strings(listCaps)

	c.Assert(listCaps, check.DeepEquals, fixedCaps["erikh2"])
}

func (ms *modelSuite) TestUserValidate(c *check.C) {
	failcases := []struct {
		username string
		token    *types.OAuthToken
	}{
		{"", nil},
		{"", testToken},
		{"erikh", nil},
		{"erikh", &types.OAuthToken{}},
	}

	for _, failure := range failcases {
		_, err := ms.model.CreateUser(failure.username, failure.token)
		c.Assert(err, check.NotNil)
	}

	u, err := ms.model.CreateUser("erikh", testToken)
	c.Assert(err, check.IsNil)
	c.Assert(u.ID, check.Not(check.Equals), 0)
	c.Assert(u.Username, check.Equals, "erikh")

	u2, err := ms.model.FindUserByName("erikh")
	c.Assert(err, check.IsNil)
	c.Assert(u.ID, check.Equals, u2.ID)
	c.Assert(u2.Username, check.Equals, "erikh")
	c.Assert(u2.Token, check.DeepEquals, testToken)

	u2.Token = nil
	c.Assert(ms.model.Save(u2).Error, check.NotNil)

	u2.Token = &types.OAuthToken{Token: "567890"}
	c.Assert(ms.model.Save(u2).Error, check.IsNil)

	u3, err := ms.model.FindUserByName("erikh")
	c.Assert(err, check.IsNil)
	c.Assert(u3.Token, check.DeepEquals, u2.Token)
}

func (ms *modelSuite) TestUserErrors(c *check.C) {
	u, err := ms.model.CreateUser("erikh", testToken)
	c.Assert(err, check.IsNil)

	errs := []string{
		"hi there",
		"hello, world!",
		"greetings",
		"error message!",
	}

	for _, err := range errs {
		u.AddError(errors.New(err))
	}

	errs2 := []string{}

	for _, err := range u.Errors {
		errs2 = append(errs2, err.Error)
	}

	c.Assert(errs2, check.DeepEquals, errs)
	c.Assert(ms.model.Save(u).Error, check.IsNil)

	u2 := &User{ID: u.ID}
	c.Assert(ms.model.Find(&u2).Error, check.IsNil)

	errs2 = []string{}

	for _, err := range u2.Errors {
		errs2 = append(errs2, err.Error)
	}

	c.Assert(errs2, check.DeepEquals, errs)

	for _, err := range u2.Errors {
		c.Assert(ms.model.DeleteError(u2, err.ID), check.IsNil)
	}

	u2 = &User{ID: u.ID}
	c.Assert(ms.model.Find(&u2).Error, check.IsNil)

	c.Assert(u2.Errors, check.DeepEquals, []UserError{})
}

func (ms *modelSuite) TestUserSubscriptions(c *check.C) {
	u, err := ms.model.CreateUser(testutil.RandString(8), testToken)
	c.Assert(err, check.IsNil)
	c.Assert(len(u.Subscribed), check.Equals, 0)

	repos := RepositoryList{}

	for i := rand.Intn(10) + 1; i >= 0; i-- {
		repo, err := ms.CreateRepository()
		c.Assert(err, check.IsNil)

		repos = append(repos, repo)
	}

	c.Assert(ms.model.AddSubscriptionsForUser(u, repos), check.IsNil)
	c.Assert(ms.model.Save(u).Error, check.IsNil)

	u2 := &User{}
	c.Assert(ms.model.Preload("Subscribed").Where("id = ?", u.ID).First(u2).Error, check.IsNil)
	c.Assert(u2.ID, check.Equals, u.ID)
	c.Assert(len(repos), check.Equals, len(u2.Subscribed))

	sort.Stable(repos)
	list := RepositoryList(u2.Subscribed)
	sort.Stable(list)

	for i := 0; i < len(repos); i++ {
		c.Assert(repos[i].Name, check.Equals, list[i].Name)
	}

	c.Assert(ms.model.RemoveSubscriptionForUser(u, repos[0]), check.IsNil)

	u2 = &User{}
	c.Assert(ms.model.Preload("Subscribed").Where("id = ?", u.ID).First(u2).Error, check.IsNil)
	c.Assert(u2.ID, check.Equals, u.ID)

	list = RepositoryList(u2.Subscribed)
	sort.Stable(list)

	c.Assert(len(repos[1:]), check.Equals, len(list))

	for i := 0; i < len(repos[1:]); i++ {
		c.Assert(repos[i+1].Name, check.Equals, list[i].Name)
	}
}

func (ms *modelSuite) TestSaveRepositories(c *check.C) {
	owners, err := ms.CreateUsers(1)
	c.Assert(err, check.IsNil)

	err = ms.model.SaveRepositories([]*github.Repository{
		{
			FullName: github.String("erikh/barbara"),
		},
		{
			FullName: github.String("foo/bar"),
		},
	}, owners[0].Username, false)
	c.Assert(err, check.IsNil)

	list, err := ms.model.GetAllPublicRepos("")
	c.Assert(err, check.IsNil)

	sort.Stable(list)
	names := []string{}

	for _, item := range list {
		names = append(names, item.Name)
	}

	c.Assert(names, check.DeepEquals, []string{"erikh/barbara", "foo/bar"})
}
