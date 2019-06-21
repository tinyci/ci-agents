package model

import (
	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/testutil"
)

func (ms *modelSuite) TestBasicConnection(c *check.C) {
	_, err := New("")
	c.Assert(err, check.NotNil)

	m, err := New(testutil.TestDBConfig)
	c.Assert(err, check.IsNil)
	c.Assert(m.Exec("select 1").Error, check.IsNil)
}
