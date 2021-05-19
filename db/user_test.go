package db

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/types"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestCapabilityModification(t *testing.T) {
	m := testInit(t)

	caps := []types.Capability{types.CapabilityCancel, types.CapabilityModifyCI, types.CapabilityModifyUser, types.CapabilitySubmit}

	strCaps := []string{}
	for _, cap := range caps {
		strCaps = append(strCaps, string(cap))
	}

	fixedCaps := map[string][]string{
		"erikh2": strCaps,
	}
	u, err := m.CreateUser(ctx, "erikh", &types.OAuthToken{Token: "dummy"})
	assert.NilError(t, err)

	for _, cap := range caps {
		res, err := m.HasCapability(ctx, u, cap, fixedCaps)
		assert.NilError(t, err)
		assert.Assert(t, !res)
		assert.NilError(t, m.AddCapabilityToUser(ctx, u, cap))
		res, err = m.HasCapability(ctx, u, cap, fixedCaps)
		assert.NilError(t, err)
		assert.Assert(t, res)

		getCaps, err := m.GetCapabilities(ctx, u, fixedCaps)
		assert.NilError(t, err)
		listCaps := []string{}

		for _, cap := range getCaps {
			listCaps = append(listCaps, string(cap))
		}

		sort.Strings(listCaps)

		assert.Assert(t, cmp.DeepEqual([]string{string(cap)}, listCaps))

		assert.NilError(t, m.RemoveCapabilityFromUser(ctx, u, cap))
		res, err = m.HasCapability(ctx, u, cap, fixedCaps)
		assert.NilError(t, err)
		assert.Assert(t, !res)
	}

	u2, err := m.CreateUser(ctx, "erikh2", &types.OAuthToken{Token: "dummy"})
	assert.NilError(t, err)

	for _, cap := range caps {
		res, err := m.HasCapability(ctx, u2, cap, fixedCaps)
		assert.NilError(t, err)
		assert.Assert(t, res)
	}

	getCaps, err := m.GetCapabilities(ctx, u2, fixedCaps)
	assert.NilError(t, err)

	listCaps := []string{}

	for _, cap := range getCaps {
		listCaps = append(listCaps, string(cap))
	}

	sort.Strings(fixedCaps["erikh2"])
	sort.Strings(listCaps)

	assert.Assert(t, cmp.DeepEqual(listCaps, fixedCaps["erikh2"]))
}

func TestUserValidate(t *testing.T) {
	m := testInit(t)

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
		_, err := m.CreateUser(ctx, failure.username, failure.token)
		assert.Assert(t, err != nil)
	}

	u, err := m.CreateUser(ctx, "erikh", testToken)
	assert.NilError(t, err)
	assert.Assert(t, u.ID != 0)
	assert.Assert(t, cmp.Equal(u.Username, "erikh"))

	u2, err := m.FindUserByName(ctx, "erikh")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(u.ID, u2.ID))
	assert.Assert(t, cmp.Equal(u2.Username, "erikh"))

	token := &types.OAuthToken{}

	assert.NilError(t, u2.Token.Unmarshal(token))
	assert.Assert(t, cmp.DeepEqual(token, testToken))

	u2.Token = nil
	_, err = u2.Update(ctx, m.db, boil.Infer())
	assert.Assert(t, err != nil)

	u2.Token, err = json.Marshal(&types.OAuthToken{Token: "567890"})
	assert.NilError(t, err)
	_, err = u2.Update(ctx, m.db, boil.Infer())
	assert.NilError(t, err)

	u3, err := m.FindUserByName(ctx, "erikh")
	assert.NilError(t, err)

	token, token2 := &types.OAuthToken{}, &types.OAuthToken{}
	assert.NilError(t, u3.Token.Unmarshal(token))
	assert.NilError(t, u2.Token.Unmarshal(token2))

	assert.Assert(t, cmp.DeepEqual(token, token2))
}

func TestUserErrors(t *testing.T) {
	m := testInit(t)

	u, err := m.CreateUser(ctx, "erikh", testToken)
	assert.NilError(t, err)

	errs := []string{
		"hi there",
		"hello, world!",
		"greetings",
		"error message!",
	}

	for _, err := range errs {
		assert.NilError(t, m.AddError(ctx, u, errors.New(err)))
	}

	errs2, err := m.GetErrors(ctx, u)
	assert.NilError(t, err)

	mappedErrs := []string{}

	for _, err := range errs2 {
		mappedErrs = append(mappedErrs, err.Error)
	}

	assert.Assert(t, cmp.DeepEqual(mappedErrs, errs))

	u2, err := models.FindUser(ctx, m.db, u.ID)
	assert.NilError(t, err)

	errs2, err = m.GetErrors(ctx, u2)
	assert.NilError(t, err)

	mappedErrs = []string{}

	for _, err := range errs2 {
		mappedErrs = append(mappedErrs, err.Error)
	}

	assert.Assert(t, cmp.DeepEqual(mappedErrs, errs))

	for _, err := range errs2 {
		assert.NilError(t, m.DeleteError(ctx, u2, err.ID))
	}

	errs2, err = m.GetErrors(ctx, u2)
	assert.NilError(t, err)

	assert.Assert(t, len(errs2) == 0)
}

func TestUserSubscriptions(t *testing.T) {
	m := testInit(t)

	u, err := m.CreateUser(ctx, testutil.RandString(8), testToken)
	assert.NilError(t, err)

	subs, err := m.GetSubscriptionsForUser(ctx, u.ID, nil, 0, 20)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(len(subs), 0))

	repos := []*models.Repository{}

	for i := rand.Intn(10) + 1; i >= 0; i-- {
		repo, err := m.CreateTestRepository(ctx)
		assert.NilError(t, err)

		repos = append(repos, repo)
	}

	assert.NilError(t, m.AddSubscriptionsForUser(ctx, u.ID, repos))

	subs, err = m.GetSubscriptionsForUser(ctx, u.ID, nil, 0, 20)
	assert.NilError(t, err)

	assert.Assert(t, cmp.Equal(len(repos), len(subs)))

	sort.Slice(repos, func(i, j int) bool { return strings.Compare(repos[i].Name, repos[j].Name) < 0 })
	sort.Slice(subs, func(i, j int) bool { return strings.Compare(subs[i].Name, subs[j].Name) < 0 })

	for i := 0; i < len(repos); i++ {
		assert.Assert(t, cmp.Equal(repos[i].Name, subs[i].Name))
	}

	assert.NilError(t, m.RemoveSubscriptionForUser(ctx, u.ID, []*models.Repository{repos[0]}))
	subs, err = m.GetSubscriptionsForUser(ctx, u.ID, nil, 0, 20)
	assert.NilError(t, err)

	sort.Slice(subs, func(i, j int) bool { return strings.Compare(subs[i].Name, subs[j].Name) < 0 })

	assert.Assert(t, cmp.Equal(len(repos[1:]), len(subs)))

	for i := 0; i < len(repos[1:]); i++ {
		assert.Assert(t, cmp.Equal(repos[i+1].Name, subs[i].Name))
	}
}

func TestSaveRepositories(t *testing.T) {
	m := testInit(t)

	owners, err := m.CreateTestUsers(ctx, 1)
	assert.NilError(t, err)

	err = m.SaveRepositories(ctx, []*github.Repository{
		{
			FullName: github.String("erikh/barbara"),
		},
		{
			FullName: github.String("foo/bar"),
		},
	}, owners[0].Username, false)
	assert.NilError(t, err)

	list, err := m.GetAllPublicRepos(ctx, nil)
	assert.NilError(t, err)

	sort.Slice(list, func(i, j int) bool { return strings.Compare(list[i].Name, list[j].Name) < 0 })

	names := []string{}
	for _, item := range list {
		names = append(names, item.Name)
	}

	assert.Assert(t, cmp.DeepEqual(names, []string{"erikh/barbara", "foo/bar"}))
}
