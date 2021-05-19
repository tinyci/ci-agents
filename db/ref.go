package db

import (
	"context"
	"errors"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func refValidateHook(ctx context.Context, db boil.ContextExecutor, ref *models.Ref) error {
	if ref.Ref == "" {
		return errors.New("empty ref")
	}

	if ref.Sha == "" {
		return errors.New("empty SHA")
	}

	if len(ref.Sha) != 40 {
		return errors.New("invalid SHA")
	}

	return nil
}

// GetRefByNameAndSHA returns the ref matching the name/sha combination.
func (m *Model) GetRefByNameAndSHA(ctx context.Context, repoName, sha string) (*models.Ref, error) {
	return models.Refs(
		qm.InnerJoin("repositories on repositories.id = refs.repository_id"),
		models.RepositoryWhere.Name.EQ(repoName),
		models.RefWhere.Sha.EQ(sha),
	).One(ctx, m.db)
}

// PutRef adds a ref to the db
func (m *Model) PutRef(ctx context.Context, ref *models.Ref) error {
	return ref.Insert(ctx, m.db, boil.Infer())
}
