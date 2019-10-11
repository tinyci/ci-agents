package model

import (
	"encoding/json"
	"math/rand"
	"sort"
	"strings"

	check "github.com/erikh/check"
	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
)

var testToken = &types.OAuthToken{
	Token: "123456",
}

 func (ms *modelSuite) TestUserTokens(c *check.C) {
	usr, err := ms.model.CreateUser("roz", testToken)
	c.Assert(err, check.IsNil)
	testTokenJSON, e := json.Marshal(testToken)
	if e != nil {
		c.Error("error while marshalling :" + e.Error())
	}

	usrJSON, e := json.Marshal(usr)
	if e != nil {
		c.Error("error while marshalling :" + e.Error())
	}

	if strings.Contains(string(usrJSON), string(testTokenJSON)) {
		c.Error("contains TokenJSON")
	}

	typesUser := usr.ToProto()
	usr2, err := NewUserFromProto(typesUser)
	c.Check(err, check.IsNil)
	usr2JSON, e := json.Marshal(usr2)
	if e != nil {
		c.Error("error while marshalling :" + e.Error())
	}

	c.Check(usr2JSON, check.DeepEquals, usrJSON)
	usrList := []*User{usr, usr2}
	usrListJSON, e := json.Marshal(usrList)
	if e != nil {
		c.Error("error while marshalling :" + e.Error())
	}

	usr2List, err := MakeUsers(MakeUserList(usrList))
	c.Check(err, check.IsNil)
	usr2ListJSON, e := json.Marshal(usr2List)
	if e != nil {
		c.Error("error while marshalling :" + e.Error())
	}
	c.Check(usrListJSON, check.DeepEquals, usr2ListJSON)
	if strings.Contains(string(usrListJSON), string(testTokenJSON)) {
		c.Error("contains TokenJSON")
	}
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
