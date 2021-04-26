package db

import (
	"testing"

	"github.com/tinyci/ci-agents/testutil"
	"gotest.tools/v3/assert"
)

func TestBasicConnection(t *testing.T) {
	m, err := NewConn("")
	assert.NilError(t, err) // no error on connect
	_, err = m.Exec("select 1")
	assert.Assert(t, err != nil)

	m, err = NewConn(testutil.TestDBConfig)
	assert.NilError(t, err)
	_, err = m.Exec("select 1")
	assert.NilError(t, err)

	m.Close()
}
