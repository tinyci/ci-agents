package model

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // pg support for gorm
	"github.com/tinyci/ci-agents/errors"
)

// Model is the outer layer of our internal database model, which will
// primarily be used by the data service.
type Model struct {
	*gorm.DB
}

// New returns the model structure after the db connection work has taken place.
func New(sqlURL string) (*Model, error) {
	db, err := gorm.Open("postgres", sqlURL)
	if err != nil {
		return nil, errors.New(err)
	}

	db = db.Set("gorm:auto_preload", true)

	if os.Getenv("SQL_DEBUG") != "" {
		db = db.Debug()
	} else {
		// this mutes it in test runs, where it's on by default I guess?! Very
		// noisy.
		db = db.LogMode(false)
	}

	return &Model{DB: db}, nil
}

// WrapError is a tail call for db transactions; it will return a wrapped and
// stack-annotated error with the msg if there is one; otherwise it will return
// nil. It also uses the errors package to normalize common errors returned
// from the DB.
func (m *Model) WrapError(call *gorm.DB, msg string) error {
	if call.Error == nil {
		return nil
	}

	return errors.MapError(call.Error).(errors.Error).Wrap(msg)
}

func doThing() error {
	return nil
}
