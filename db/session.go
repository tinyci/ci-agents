package db

import (
	"context"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// LoadSession loads a session based on the key and returns it to the client
func (m *Model) LoadSession(ctx context.Context, key string) (*models.Session, error) {
	return models.Sessions(qm.Where("key = ? and expires_on > now()", key)).One(ctx, m.db)
}

// SaveSession does the opposite of LoadSession
func (m *Model) SaveSession(ctx context.Context, session *models.Session) error {
	return session.Insert(ctx, m.db, boil.Infer())
}
