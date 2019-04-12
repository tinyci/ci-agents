package types

import (
	check "github.com/erikh/check"
)

func (ts *typesSuite) TestSubmissionValidation(c *check.C) {
	submap := map[bool][]*Submission{
		true: {
			{
				Parent:  "foo/bar",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "foo/bar",
				Fork:    "bar/foo",
				BaseSHA: "refs/heads/master",
				HeadSHA: "heads/master",
			},
			{
				Parent:  "foo.bar/bar.foo",
				Fork:    "bar/foo",
				BaseSHA: "refs/heads/master",
				HeadSHA: "heads/master",
			},
			{
				Parent:  "foo.bar/bar.foo",
				Fork:    "bar/foo",
				BaseSHA: "refs/heads/master",
				HeadSHA: "heads/master",
				All:     true,
				Manual:  true,
			},
		},
		false: {
			{
				Parent:  "/",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
				All:     true,
			},
			{
				Parent:  "/",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "../",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "/..",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "../..",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "./.",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "/.",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "./",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "bar/foo/",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "bar/foo",
				Fork:    "/",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "bar/foo",
				Fork:    "bar/foo/",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "bar/foo",
				Fork:    "",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "",
				Fork:    "bar/foo",
				BaseSHA: "master",
				HeadSHA: "master",
			},
			{
				Parent:  "",
				Fork:    "bar/foo",
				HeadSHA: "master",
			},
			{
				Parent:  "",
				Fork:    "bar/foo",
				BaseSHA: "master",
			},
		},
	}

	for success, subs := range submap {
		for _, sub := range subs {
			ck := check.IsNil
			if !success {
				ck = check.NotNil
			}
			c.Assert(sub.Validate(), ck, check.Commentf("Error is supposed to be a %v value: %#v", !success, sub))
		}
	}
}
