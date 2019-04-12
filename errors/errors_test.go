package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/erikh/check"
)

type errorSuite struct{}

var _ = check.Suite(&errorSuite{})

func TestErrors(t *testing.T) {
	check.TestingT(t)
}

func (es *errorSuite) TestWrap(c *check.C) {
	err := New("hi")

	c.Assert(err.Wrapf("foobar: %v", errors.New("hey there")).Error(), check.Equals, "foobar: hey there: hi")
	c.Assert(err.Wrap("hey there").Error(), check.Equals, "hey there: hi")
	c.Assert(err.Wrap(errors.New("hey there")).Error(), check.Equals, "hey there: hi")
}

func (es *errorSuite) TestLog(c *check.C) {
	err := New("hi")
	err.SetLog(true)
	c.Assert(err.GetLog(), check.Equals, true)

	err.SetLog(false)
	c.Assert(err.GetLog(), check.Equals, false)

	err = NewNoLog("hi")
	c.Assert(err.GetLog(), check.Equals, false)
}

func (es *errorSuite) TestFormat(c *check.C) {
	err := New("hi")
	e2 := foo(err)

	buf := fmt.Sprintf("%v", e2)
	c.Assert(strings.TrimSpace(buf), check.Equals, "in foo: hi")
	buf = fmt.Sprintf("%+v", e2)
	res := strings.TrimSpace(buf)
	c.Assert(res, check.Matches, "(?s).*in foo:.*github.com/tinyci/ci-agents/errors.foo.*")
	c.Assert(res, check.Matches, "(?s).*hi:.*github.com/tinyci/ci-agents/errors.\\(\\*errorSuite\\).TestFormat.*")
}

func (es *errorSuite) TestContains(c *check.C) {
	err := New("hi")
	e2 := foo(err)
	c.Assert(e2.Contains(err), check.Equals, true)
	c.Assert(e2.Contains(New("nope")), check.Equals, false)
}

func foo(err *Error) *Error {
	return err.Wrap("in foo")
}
