package model

import (
	"strings"
	"time"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/utils"
)

// OAuthExpiration is a constant for the oauth state expiration time. It really
// should be in the configuration.
var OAuthExpiration = 10 * time.Minute

// OAuth schema is for checking state return values from github.
type OAuth struct {
	State     string    `gorm:"primary key" json:"state"`
	Scopes    string    `json:"scopes"`
	ExpiresOn time.Time `json:"expires_on"`
}

// ToProto returns a protobuf representation of the model oauth request
func (o *OAuth) ToProto() *data.OAuthState {
	return &data.OAuthState{
		State:  o.State,
		Scopes: o.GetScopesList(),
	}
}

// SetScopes takes a []string of scope names, and marshals them to the struct field.
func (o *OAuth) SetScopes(scopes []string) {
	o.Scopes = strings.Join(scopes, ",")
}

// GetScopesList returns the split list of scopes
func (o *OAuth) GetScopesList() []string {
	return strings.Split(o.Scopes, ",")
}

// GetScopes returns the split list of scopes, mapped for easy comparison.
func (o *OAuth) GetScopes() map[string]struct{} {
	ret := map[string]struct{}{}
	for _, scope := range o.GetScopesList() {
		ret[scope] = struct{}{}
	}

	return ret
}

// OAuthRegisterState registers a state code in a uniqueness table that tracks
// it.
func (m *Model) OAuthRegisterState(state string, scopes []string) error {
	oa := &OAuth{
		State:     state,
		ExpiresOn: time.Now().Add(OAuthExpiration),
	}

	oa.SetScopes(scopes)
	return m.WrapError(m.Save(oa), "registering oauth state")
}

// OAuthValidateState validates that the state we sent actually exists and is
// ready to be consumed. In the event it is not, it returns error.
func (m *Model) OAuthValidateState(state string) (*OAuth, error) {
	oa := &OAuth{}
	err := m.WrapError(m.Where("state = ? and expires_on > now()", state).Find(&oa), "validating oauth state")
	if err != nil {
		return nil, utils.WrapError(utils.ErrNotFound, "%v", err)
	}

	defer m.Delete(&OAuth{State: state})

	return oa, nil
}
