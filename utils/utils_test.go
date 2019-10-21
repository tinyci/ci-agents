package utils

import (
	"encoding/json"
	. "testing"

	"github.com/erikh/check"
	"github.com/gin-gonic/gin"
)

type utilsSuite struct{}

var _ = check.Suite(&utilsSuite{})

func TestUtils(t *T) {
	check.TestingT(t)
}

func (us *utilsSuite) TestJSONContext(c *check.C) {
	ctx := &gin.Context{}

	var status bool

	content, err := json.Marshal(true)
	c.Assert(err, check.IsNil)

	ctx.Set("ok", content)

	c.Assert(JSONContext(ctx, "ok", &status), check.IsNil)

	var notStatus string
	c.Assert(JSONContext(ctx, "not_ok", &status), check.NotNil)
	c.Assert(JSONContext(ctx, "ok", &notStatus), check.NotNil)
}

func (us *utilsSuite) TestJSONIO(c *check.C) {
	s1 := struct {
		Foo  string
		Bar  string
		Quux bool
	}{"foo", "bar", true}

	s2 := struct {
		Foo  string
		Bar  string
		Quux bool
	}{}

	c.Assert(JSONIO(&s1, &s2), check.IsNil)
	c.Assert(s2.Foo, check.Equals, "foo")
	c.Assert(s2.Bar, check.Equals, "bar")
	c.Assert(s2.Quux, check.Equals, true)

	s3 := struct {
		Foo string
		Bar string
	}{}

	c.Assert(JSONIO(&s1, &s3), check.IsNil)
	c.Assert(s3.Foo, check.Equals, "foo")
	c.Assert(s3.Bar, check.Equals, "bar")

	s4 := struct {
		Foo int
		Bar string
	}{}
	c.Assert(JSONIO(&s1, &s4), check.NotNil)
}

func (us *utilsSuite) TestScopePagination(c *check.C) {
	_, _, err := ScopePaginationInt(-1, 0)
	c.Assert(err, check.NotNil)

	pg, ppg, err := ScopePaginationInt(0, 0)
	c.Assert(err, check.IsNil)
	c.Assert(pg, check.Equals, int64(0))
	c.Assert(ppg, check.Equals, defaultPerPage)

	pg, ppg, err = ScopePaginationInt(0, MaxPerPage*2)
	c.Assert(err, check.IsNil)
	c.Assert(pg, check.Equals, int64(0))
	c.Assert(ppg, check.Equals, MaxPerPage)

	pg, ppg, err = ScopePagination("", "")
	c.Assert(err, check.IsNil)
	c.Assert(pg, check.Equals, int64(0))
	c.Assert(ppg, check.Equals, defaultPerPage)

	for _, failure := range []string{"asdf", "-1", "d34dbeef"} {
		_, _, err := ScopePagination(failure, "")
		c.Assert(err, check.NotNil, check.Commentf("%v", failure))
		_, _, err = ScopePagination("", failure)
		c.Assert(err, check.NotNil, check.Commentf("%v", failure))
	}
}

type orResult struct {
	owner string
	repo  string
	error bool
}

func (us *utilsSuite) TestOwnerRepo(c *check.C) {
	results := map[string]orResult{
		"owner/repo": {
			owner: "owner",
			repo:  "repo",
		},
		"../..":       {error: true},
		"":            {error: true},
		"./":          {error: true},
		"/.":          {error: true},
		"/":           {error: true},
		"/asdf":       {error: true},
		"asdf/":       {error: true},
		"%/asdf":      {error: true},
		"asdf</>asdf": {error: true},
		"asdf&/asdf":  {error: true},
	}

	for test, result := range results {
		ck := check.IsNil
		if result.error {
			ck = check.NotNil
		}
		owner, repo, err := OwnerRepo(test)
		c.Assert(err, ck, check.Commentf("%v should be %v", test, result))
		c.Assert(owner, check.Equals, result.owner)
		c.Assert(repo, check.Equals, result.repo)
		c.Assert(IsOwnerRepo(test), check.Equals, !result.error)
	}
}

func (us *utilsSuite) TestIsSHA(c *check.C) {
	shaResults := map[string]bool{
		"abcdef": false, // not a full sha
		"":       false, // empty string
		"tttttttttttttttttttttttttttttttttttttttt": false, // 40 chars not a digest
		"0000000000000000000000000000000000000000": true,  // this is an actual sha we use frequently
		"be3d26c478991039e951097f2c99f56b55396941": true,
		"be3d26c478991039e951097f2c99f56b5539694.": false, // invalid chars
		".e3d26c478991039e951097f2c99f56b55396941": false, // invalid chars
		"be3d26c478991039e951g97f2c99f56b55396941": false, // invalid chars (g in middle)
	}

	for sha, result := range shaResults {
		c.Assert(IsSHA(sha), check.Equals, result, check.Commentf("%v should be %v", sha, result))
	}
}

func (us *utilsSuite) TestQualifyBranch(c *check.C) {
	results := map[string]struct {
		branch string
		error  bool
	}{
		"master":            {branch: "heads/master"},
		"refs/heads/master": {branch: "heads/master"},
		"be3d26c478991039e951097f2c99f56b55396941": {
			branch: "be3d26c478991039e951097f2c99f56b55396941",
			error:  true,
		}, // is a SHA
		".":             {error: true},
		"..":            {error: true},
		"heads/..":      {error: true},
		"":              {error: true},
		"/":             {error: true},
		"/heads/master": {error: true},
	}

	for origBranch, result := range results {
		newBranch, err := QualifyBranch(origBranch)
		ck := check.IsNil
		if result.error {
			ck = check.NotNil
		}

		c.Assert(err, ck, check.Commentf("error for %v should be %v", origBranch, result.error))
		c.Assert(newBranch, check.Equals, result.branch)
	}
}
