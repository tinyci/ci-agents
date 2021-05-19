package testutil

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

var (
	// TestDBConfig is the default test database configuration.
	TestDBConfig = "host=localhost user=tinyci password=tinyci"

	// DummyToken is a fake oauth2 token.
	DummyToken = &types.OAuthToken{Token: "123456"}

	// WipeTables is a list of tables to be wiped by WipeDB. The order is
	// important, basically it's the root, forward looking, from a perspective of
	// foreign key relationships. This list is cleared in /reverse order/.
	WipeTables = []string{
		"o_auths",
		"sessions",
		"users",
		"repositories",
		"refs",
		"submissions",
		"tasks",
		"runs",
		"queue_items",
		"subscriptions",
		"user_capabilities",
		"user_errors",
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// WipeDB wipes all tables with `delete from`.
func WipeDB() error {
	m, err := sql.Open("postgres", TestDBConfig)
	if err != nil {
		return err
	}
	defer m.Close()

	m.SetMaxOpenConns(1)

	tx, err := m.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := len(WipeTables) - 1; i >= 0; i-- {
		_, err := tx.Exec(fmt.Sprintf("delete from %q", WipeTables[i]))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
