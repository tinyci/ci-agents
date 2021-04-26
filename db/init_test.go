package db

import (
	"context"
	"testing"

	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/testutil"
)

var ctx = context.Background() // save typing

func testInit(t *testing.T) *Model {
	t.Helper()

	if err := testutil.WipeDB(); err != nil {
		t.Fatal(err)
	}

	conf := &config.UserConfig{
		DSN: testutil.TestDBConfig,
	}

	m, err := Open(conf)
	if err != nil {
		t.Fatal(err)
	}

	m.SetConnPoolSize(100)
	t.Cleanup(func() { m.db.Close() })

	return m
}
