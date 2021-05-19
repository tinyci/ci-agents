package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// defaultOAuthExpiration is a constant for the default oauth state expiration time.
var defaultOAuthExpiration = 10 * time.Minute

func checkScopes(scopes []string) error {
	for _, scope := range scopes {
		if strings.Contains(scope, ",") {
			return errors.New("scope names cannot contain a comma")
		}
	}

	return nil
}

// OAuthRegisterState registers a state for a given set of scopes.
func (m *Model) OAuthRegisterState(ctx context.Context, state string, scopes []string) error {
	if err := checkScopes(scopes); err != nil {
		return err
	}

	timeout := m.config.OAuth.StateTimeout
	if timeout == 0 {
		timeout = defaultOAuthExpiration
	}

	oa := &models.OAuth{
		State:     state,
		ExpiresOn: time.Now().Add(timeout),
		Scopes:    strings.Join(scopes, ","),
	}

	return oa.Insert(ctx, m.db, boil.Infer())
}

// OAuthValidateState validates the oauth state and returns it, or it returns an error!
func (m *Model) OAuthValidateState(ctx context.Context, state string) ([]string, error) {
	oa, err := models.OAuths(qm.Where("state = ? and expires_on > now()", state)).One(ctx, m.db)
	if err != nil {
		return nil, err
	}

	return strings.Split(oa.Scopes, ","), nil
}
