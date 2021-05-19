package db

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestOAuthStates(t *testing.T) {
	scopes := []string{"test", "test2", "test3"}

	m := testInit(t)
	m.config.OAuth.StateTimeout = 2 * time.Second

	_, err := m.OAuthValidateState(ctx, "abc")
	assert.Assert(t, err != nil)
	assert.NilError(t, m.OAuthRegisterState(ctx, "abc", scopes))

	retScopes, err := m.OAuthValidateState(ctx, "abc")
	assert.NilError(t, err)

	assert.Assert(t, cmp.DeepEqual(retScopes, scopes))
}
