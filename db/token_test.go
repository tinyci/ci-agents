package db

import (
	"testing"

	"github.com/tinyci/ci-agents/utils"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestBasicToken(t *testing.T) {
	m := testInit(t)
	u, err := m.CreateUser(ctx, "erikh", testToken)
	assert.NilError(t, err)
	assert.Assert(t, u != nil) //nolint:staticcheck

	token, err := m.GetToken(ctx, "erikh")
	assert.NilError(t, err)
	assert.Assert(t, len(token) != 0)

	u2, err := m.FindUserByID(ctx, u.ID) //nolint:staticcheck
	assert.NilError(t, err)
	assert.Assert(t, string(u2.LoginToken.Bytes) != token)

	_, err = m.GetToken(ctx, "erikh")
	assert.Assert(t, err != nil)

	assert.NilError(t, m.DeleteToken(ctx, "erikh"))

	token2, err := m.GetToken(ctx, "erikh")
	assert.NilError(t, err)
	assert.Assert(t, len(token) != 0)

	assert.Assert(t, token != token2)

	retUser, err := m.ValidateToken(ctx, token2)
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(retUser.Username, u.Username)) //nolint:staticcheck
	_, err = m.ValidateToken(ctx, token)
	assert.Error(t, err, utils.ErrInvalidAuth.Error())
}
