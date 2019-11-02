package model

import (
	"time"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
)

// Session corresponds to the `sessions` table and encapsulates a web session.
type Session struct {
	Key       string    `gorm:"primary_key" json:"key"`
	Values    string    `json:"values"`
	ExpiresOn time.Time `json:"expires_on"`
}

// NewSessionFromProto returns a session from a protobuf representation.
func NewSessionFromProto(s *types.Session) *Session {
	return &Session{
		Key:       s.Key,
		Values:    s.Values,
		ExpiresOn: *MakeTime(s.ExpiresOn, false),
	}
}

// ToProto converts the session to protobuf representation.
func (s *Session) ToProto() *types.Session {
	return &types.Session{
		Key:       s.Key,
		Values:    s.Values,
		ExpiresOn: MakeTimestamp(&s.ExpiresOn),
	}
}

// LoadSession loads a session based on the key and returns it to the client
func (m *Model) LoadSession(id string) (*Session, error) {
	s := &Session{}
	return s, m.WrapError(m.Limit(1).Where("key = ? and expires_on > now()", id).First(s), "loading session")
}

// SaveSession does the opposite of LoadSession
func (m *Model) SaveSession(session *Session) error {
	return m.WrapError(m.Save(session), "saving session")
}
