package model

import (
	"strings"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
)

func (ms *modelSuite) TestOAuthType(c *check.C) {
	scopes := []string{"test", "test2", "test3"}
	o := OAuth{State: "abc", Scopes: strings.Join(scopes, ",")}

	c.Assert(o.GetScopesList(), check.DeepEquals, scopes)
	scopeMap := map[string]struct{}{}
	for _, scope := range scopes {
		scopeMap[scope] = struct{}{}
	}
	c.Assert(o.GetScopes(), check.DeepEquals, scopeMap)
	c.Assert(o.ToProto(), check.FitsTypeOf, &data.OAuthState{})
	c.Assert(o.ToProto().Scopes, check.DeepEquals, o.GetScopesList())
	c.Assert(o.ToProto().State, check.Equals, o.State)
}

func (ms *modelSuite) TestOAuthStates(c *check.C) {
	scopes := []string{"test", "test2", "test3"}

	_, err := ms.model.OAuthValidateState("abc")
	c.Assert(err, check.NotNil)
	c.Assert(ms.model.OAuthRegisterState("abc", scopes), check.IsNil)
	o, err := ms.model.OAuthValidateState("abc")
	c.Assert(err, check.IsNil)
	c.Assert(o.GetScopesList(), check.DeepEquals, scopes)

	exp := OAuthExpiration
	defer func() { OAuthExpiration = exp }()

	OAuthExpiration = time.Second
	time.Sleep(time.Second)
	_, err = ms.model.OAuthValidateState("abc")
	c.Assert(err, check.NotNil)
}
