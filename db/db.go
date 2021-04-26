package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq" // postgres db driver
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/db/models"
)

// NewConn creates a new DB connection. This is typically used to migrate the
// database and perform other administrative functions. You probably don't need
// this and want Open() instead.
func NewConn(dsn string) (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}

// Model is the handle into the DB subsystem.
type Model struct {
	db     *sql.DB
	config *config.UserConfig
}

// Open opens a handle into the database, exposing its functionality.
func Open(conf *config.UserConfig) (*Model, error) {
	sqlDB, err := NewConn(conf.DSN)
	if err != nil {
		return nil, err
	}

	if os.Getenv("SQL_DEBUG") != "" {
		boil.DebugMode = true
		boil.DebugWriter = os.Stderr
	}

	sqlDB.SetConnMaxIdleTime(-1)
	sqlDB.SetConnMaxLifetime(-1)

	registerHooks()

	return &Model{db: sqlDB, config: conf}, nil
}

// GetDB is used to bridge some gaps, mostly by the protoconv lib
func (m *Model) GetDB() *sql.DB {
	return m.db
}

func registerHooks() {
	models.AddRefHook(boil.AfterSelectHook, refValidateHook)
	models.AddRefHook(boil.BeforeInsertHook, refValidateHook)
	models.AddRefHook(boil.BeforeUpdateHook, refValidateHook)
	models.AddRefHook(boil.BeforeUpsertHook, refValidateHook)

	models.AddRepositoryHook(boil.AfterSelectHook, repoValidateHook)
	models.AddRepositoryHook(boil.BeforeInsertHook, repoValidateHook)
	models.AddRepositoryHook(boil.BeforeUpdateHook, repoValidateHook)
	models.AddRepositoryHook(boil.BeforeUpsertHook, repoValidateHook)

	models.AddRunHook(boil.AfterSelectHook, runValidateHook)
	models.AddRunHook(boil.BeforeInsertHook, runValidateHook)
	models.AddRunHook(boil.BeforeUpdateHook, runValidateHook)
	models.AddRunHook(boil.BeforeUpsertHook, runValidateHook)

	models.AddQueueItemHook(boil.AfterSelectHook, queueItemValidateHook)
	models.AddQueueItemHook(boil.BeforeInsertHook, queueItemValidateHook)
	models.AddQueueItemHook(boil.BeforeUpdateHook, queueItemValidateHook)
	models.AddQueueItemHook(boil.BeforeUpsertHook, queueItemValidateHook)

	models.AddTaskHook(boil.AfterSelectHook, taskValidateHook)
	models.AddTaskHook(boil.BeforeInsertHook, taskValidateHook)
	models.AddTaskHook(boil.BeforeUpdateHook, taskValidateHook)
	models.AddTaskHook(boil.BeforeUpsertHook, taskValidateHook)

	// select hook is different here.
	models.AddUserHook(boil.AfterSelectHook, userReadValidateHook)
	models.AddUserHook(boil.BeforeInsertHook, userWriteValidateHook)
	models.AddUserHook(boil.BeforeUpdateHook, userWriteValidateHook)
	models.AddUserHook(boil.BeforeUpsertHook, userWriteValidateHook)
}

// SetConnPoolSize sets the connection pool parameters.
func (m *Model) SetConnPoolSize(size int) {
	m.db.SetMaxIdleConns(size)
	m.db.SetMaxOpenConns(size)
}
