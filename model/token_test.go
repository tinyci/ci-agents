package model

import (
	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/errors"
)

func (ms *modelSuite) TestBasicToken(c *check.C) {
	u, err := ms.model.CreateUser("erikh", testToken)
	c.Assert(err, check.IsNil)
	c.Assert(u, check.NotNil)

	token, err := ms.model.GetToken("erikh")
	c.Assert(err, check.IsNil)
	c.Assert(len(token), check.Not(check.Equals), "")

	u2, err := ms.model.FindUserByID(u.ID)
	c.Assert(err, check.IsNil)
	c.Assert(u2.LoginToken, check.Not(check.Equals), token)

	_, err = ms.model.GetToken("erikh")
	c.Assert(err, check.NotNil)

	c.Assert(ms.model.DeleteToken("erikh"), check.IsNil)

	token2, err := ms.model.GetToken("erikh")
	c.Assert(err, check.IsNil)
	c.Assert(len(token), check.Not(check.Equals), "")

	c.Assert(token, check.Not(check.Equals), token2)

	username, err := ms.model.ValidateToken(token2)
	c.Assert(err, check.IsNil)
	c.Assert(username, check.Equals, u.Username)
	_, err = ms.model.ValidateToken(token)
	c.Assert(err, check.Equals, errors.ErrInvalidAuth)
}
