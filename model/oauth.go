package model

import (
	"time"

	"github.com/tinyci/ci-agents/errors"
)

// OAuth schema is for checking state return values from github.
type OAuth struct {
	State     string    `gorm:"primary key" json:"state"`
	ExpiresOn time.Time `json:"expires_on"`
}

// OAuthRegisterState registers a state code in a uniqueness table that tracks
// it.
func (m *Model) OAuthRegisterState(state string) *errors.Error {
	return m.WrapError(
		m.Save(
			&OAuth{
				State:     state,
				ExpiresOn: time.Now().Add(10 * time.Minute),
			},
		), "registering oauth state")
}

// OAuthValidateState validates that the state we sent actually exists and is
// ready to be consumed. In the event it is not, it returns *errors.Error.
func (m *Model) OAuthValidateState(state string) *errors.Error {
	var count int64

	err := m.WrapError(
		m.
			Model(&OAuth{}).
			Where("state = ? and expires_on > now()", state).
			Count(&count),
		"validating oauth state")

	if err != nil {
		return errors.New(err)
	}

	defer m.Delete(&OAuth{State: state})

	if count > 0 {
		return nil
	}

	return errors.ErrNotFound
}
