package db

import (
	"context"
	"errors"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/types"
)

func userReadValidateHook(ctx context.Context, db boil.ContextExecutor, u *models.User) error {
	if u.Username == "" {
		return errors.New("username is empty")
	}

	return nil
}

func userWriteValidateHook(ctx context.Context, db boil.ContextExecutor, u *models.User) error {
	if err := userReadValidateHook(ctx, db, u); err != nil {
		return err
	}

	if len(u.Token) == 0 {
		return errors.New("token is empty (nil)")
	}

	return nil
}

// SetGithubToken allows setting of the github token. Key must already be parsed. (config.Auth.Validate)
func (m *Model) SetGithubToken(ctx context.Context, u *models.User, tok *types.OAuthToken) error {
	var err error
	u.Token, err = types.EncryptToken(config.TokenCryptKey, tok)
	return err
}

// GetGithubToken allows setting of the github token. Key must already be parsed. (config.Auth.Validate)
func (m *Model) GetGithubToken(ctx context.Context, u *models.User) (*types.OAuthToken, error) {
	return types.DecryptToken(config.TokenCryptKey, u.Token)
}

// CreateUser initializes a user struct and writes it to the db.
func (m *Model) CreateUser(ctx context.Context, username string, token []byte) (*models.User, error) {
	u := &models.User{Username: username, Token: token}

	return u, u.Insert(ctx, m.db, boil.Infer())
}

// UpdateUser updates a user.
func (m *Model) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := user.Update(ctx, m.db, boil.Blacklist("username"))
	return err
}

// ListUsers lists all users.
func (m *Model) ListUsers(ctx context.Context) ([]*models.User, error) {
	return models.Users().All(ctx, m.db) // yarly
}

// FindUserByName finds a user by unique key username.
func (m *Model) FindUserByName(ctx context.Context, username string) (*models.User, error) {
	return models.Users(qm.Where("username = ?", username)).One(ctx, m.db)
}

// FindUserByID finds a user by primary key id.
func (m *Model) FindUserByID(ctx context.Context, id int64) (*models.User, error) {
	return models.FindUser(ctx, m.db, id)
}

// DeleteError deletes a given error for a user.
func (m *Model) DeleteError(ctx context.Context, u *models.User, id int64) error {
	_, err := u.UserErrors(models.UserErrorWhere.ID.EQ(id)).DeleteAll(ctx, m.db)
	return err
}

// AddError adds an error to the error list.
func (m *Model) AddError(ctx context.Context, u *models.User, err error) error {
	ue := &models.UserError{Error: err.Error()}
	if err := ue.SetUser(ctx, m.db, false, u); err != nil {
		return err
	}

	return ue.Insert(ctx, m.db, boil.Infer())
}

// GetErrors adds an error to the error list.
func (m *Model) GetErrors(ctx context.Context, u *models.User) (models.UserErrorSlice, error) {
	return u.UserErrors().All(ctx, m.db)
}
