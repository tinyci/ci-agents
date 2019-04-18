package testutil

import (
	"crypto/rand"
	"fmt"
	"math/big"

	check "github.com/erikh/check"
	"github.com/jinzhu/gorm"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

var (
	// TestDBConfig is the default test database configuration.
	TestDBConfig = "host=localhost user=tinyci password=tinyci"

	// DummyToken is a fake oauth2 token.
	DummyToken = &types.OAuthToken{Token: "123456"}
)

// WipeDB wipes all tables with `truncate table`.
func WipeDB(c *check.C) {
	m, err := gorm.Open("postgres", TestDBConfig)
	c.Assert(err, check.IsNil)
	defer m.Close()

	rows, err := m.Raw("select tablename from pg_catalog.pg_tables where schemaname='public'").Rows()
	c.Assert(err, check.IsNil)

	tables := []string{}

	for rows.Next() {
		var tableName string
		c.Assert(rows.Scan(&tableName), check.IsNil)
		tables = append(tables, tableName)
	}
	rows.Close()

	for _, table := range tables {
		c.Assert(m.Exec(fmt.Sprintf("truncate table %q", table)).Error, check.IsNil)
	}
}

// RandString returns a random string
func RandString(len int) string {
	ret := ""
	charlen := 'z' - 'a'

	for i := 0; i < len; i++ {
		rnd, _ := rand.Int(rand.Reader, big.NewInt(int64(charlen)))
		ret += fmt.Sprintf("%c", 'a'+rune(rnd.Int64()))
	}

	return ret
}

// JSONIO uses json serialization to convert from one struct to another
// p.s. this is terrible
func JSONIO(from, to interface{}) *errors.Error {
	return utils.JSONIO(from, to)
}
