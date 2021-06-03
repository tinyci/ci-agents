package db

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"errors"

	"github.com/gorilla/securecookie"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/utils"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (m *Model) createToken(name string) ([]byte, string) {
	key := securecookie.GenerateRandomKey(64)
	encoded := fmt.Sprintf("%s %s", name, string(key))
	return key, base64.StdEncoding.EncodeToString([]byte(encoded))
}

func (m *Model) unpackToken(token string) (string, []byte, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", nil, utils.WrapError(utils.ErrInvalidAuth, "%v", err.Error())
	}

	parts := strings.SplitN(string(b), " ", 2)
	if len(parts) != 2 {
		return "", nil, utils.WrapError(utils.ErrInvalidAuth, "%v", "invalid token")
	}

	return parts[0], []byte(parts[1]), nil
}

// ValidateToken checks that a token is valid for a given user.
func (m *Model) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	name, authToken, err := m.unpackToken(token)
	if err != nil {
		return nil, err // except for in this case, where it's done already above.
	}

	u, err := m.FindUserByName(ctx, name)
	if err != nil {
		return nil, utils.ErrInvalidAuth
	}

	if !u.LoginToken.Valid || len(u.LoginToken.Bytes) == 0 {
		return nil, utils.ErrInvalidAuth
	}

	if bytes.Equal(u.LoginToken.Bytes, authToken) {
		return u, nil
	}

	return nil, utils.ErrInvalidAuth
}

// GetToken retrieves a new token for logging in. If one exists, the
// DeleteToken method must be called first; otherwise this routine will throw
// an error.
func (m *Model) GetToken(ctx context.Context, name string) (string, error) {
	u, err := m.FindUserByName(ctx, name)
	if err != nil {
		return "", err
	}

	if u.LoginToken.Valid || len(u.LoginToken.Bytes) != 0 {
		return "", errors.New("Login token already exists, must delete the old one first")
	}

	key, token := m.createToken(name)
	u.LoginToken = null.BytesFrom(key)

	_, err = u.Update(ctx, m.db, boil.Infer())
	return token, err
}

// DeleteToken removes the existing token.
func (m *Model) DeleteToken(ctx context.Context, name string) error {
	u, err := m.FindUserByName(ctx, name)
	if err != nil {
		return err
	}

	u.LoginToken = null.BytesFromPtr(nil)
	_, err = u.Update(ctx, m.db, boil.Infer())
	return err
}
