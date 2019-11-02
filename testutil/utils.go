package testutil

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	check "github.com/erikh/check"
	"github.com/jinzhu/gorm"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

var (
	// TestDBConfig is the default test database configuration.
	TestDBConfig = "host=localhost user=tinyci password=tinyci"

	// DummyToken is a fake oauth2 token.
	DummyToken = &types.OAuthToken{Token: "123456"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

var hexChars = strings.Split("0123456789abcdef", "")

// RandHexString creates random hexadecimal strings
func RandHexString(l int) string {
	ret := ""

	for i := 0; i < l; i++ {
		ret += hexChars[rand.Intn(len(hexChars))]
	}

	return ret
}

// RandString returns a random string
func RandString(len int) string {
	ret := ""
	charlen := int('z' - 'a')

	for i := 0; i < len; i++ {
		ret += fmt.Sprintf("%c", 'a'+rune(rand.Intn(charlen)))
	}

	return ret
}

// JSONIO uses json serialization to convert from one struct to another
// p.s. this is terrible
func JSONIO(from, to interface{}) error {
	return utils.JSONIO(from, to)
}
