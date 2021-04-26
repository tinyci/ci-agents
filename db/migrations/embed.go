package migrations

import (
	"database/sql"

	"github.com/rakyll/statik/fs"

	migrate "github.com/rubenv/sql-migrate"
)

//go:generate go run github.com/rakyll/statik -src=. -dest=.. -f -p migrations -include *.sql

// Migrate migrates the pre-connected database.
func Migrate(db *sql.DB) error {
	f, _ := fs.New()

	source := migrate.HttpFileSystemMigrationSource{
		FileSystem: f,
	}

	_, err := migrate.Exec(db, "postgres", source, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
