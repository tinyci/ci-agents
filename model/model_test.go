package model

import (
	"testing"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/testutil"
)

type modelSuite struct {
	model *Model
}

var _ = check.Suite(&modelSuite{})

func TestModel(t *testing.T) {
	check.TestingT(t)
}

func (ms *modelSuite) SetUpTest(c *check.C) {
	testutil.WipeDB(c)

	var err error
	ms.model, err = New(testutil.TestDBConfig)
	c.Assert(err, check.IsNil)
}
