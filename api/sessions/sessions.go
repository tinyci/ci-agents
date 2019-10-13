package sessions

import (
	"context"
	"encoding/base32"
	"net/http"
	"strings"
	"time"

	sgin "github.com/gin-contrib/sessions"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/model"
)

// Session corresponds to the `sessions` table and encapsulates a web session.
type Session struct {
	Key       string    `json:"key"`
	Values    string    `json:"values"`
	ExpiresOn time.Time `json:"expires_on"`
}

// SessionManager manages many session objects.
type SessionManager struct {
	datasvc *data.Client
	options *sessions.Options
	codecs  []securecookie.Codec
}

// New creates a new session tracker for the db layer connection provided.
func New(datasvc *data.Client, options *sessions.Options, keyPairs ...[]byte) *SessionManager {
	if options == nil {
		options = &sessions.Options{
			MaxAge: 86400 * 30,
			Path:   "/",
		}
	}

	return &SessionManager{
		datasvc: datasvc,
		options: options,
		codecs:  securecookie.CodecsFromPairs(keyPairs...),
	}
}

// New creates a new session for the given name
func (sm *SessionManager) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(sm, name)
	if session == nil {
		return nil, nil
	}

	opts := *sm.options
	session.Options = &opts
	session.IsNew = true

	return session, sm.get(context.Background(), session, r, name)
}

func (sm *SessionManager) get(ctx context.Context, session *sessions.Session, r *http.Request, name string) error {
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err := securecookie.DecodeMulti(name, c.Value, &session.ID, sm.codecs...)
		if err == nil {
			if err := sm.LoadSession(ctx, session); err == nil {
				session.IsNew = false
			} else if errors.Cause(err).Error() != "record not found" {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// LoadSession loads a session from the database.
func (sm *SessionManager) LoadSession(ctx context.Context, session *sessions.Session) error {
	s, err := sm.datasvc.GetSession(ctx, session.ID)
	if err != nil {
		return err
	}

	if session.Options == nil {
		session.Options = &sessions.Options{Path: "/"}
	}

	session.Options.MaxAge = int(time.Until(time.Time(s.ExpiresOn)).Seconds())
	if session.Options.MaxAge < 1 {
		session.Options.MaxAge = -1
		return errors.New("cookie expired")
	}

	return securecookie.DecodeMulti(session.Name(), s.Values, &session.Values, sm.codecs...)
}

// Get gets the session for the given name
func (sm *SessionManager) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(sm, name)
}

// Save persists the session for the given name. Also requires the http writer.
func (sm *SessionManager) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	encoded, err := sm.SaveSession(s)
	if err != nil {
		return err
	}

	if s.IsNew {
		http.SetCookie(w, sessions.NewCookie(s.Name(), encoded, s.Options))
	}
	return nil
}

// Options sets the options for the session handler.
func (sm *SessionManager) Options(options sgin.Options) {
	sm.options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

// SaveSession saves the provided session with the codecs used.
func (sm *SessionManager) SaveSession(s *sessions.Session) (string, error) {
	if s.ID == "" {
		s.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(128)), "=")
	}

	if s.Options == nil {
		s.Options = &sessions.Options{
			MaxAge: 1440,
			Path:   "/",
		}
	}

	encoded, err := securecookie.EncodeMulti(s.Name(), s.ID, sm.codecs...)
	if err != nil {
		return "", err
	}

	values, err := securecookie.EncodeMulti(s.Name(), s.Values, sm.codecs...)
	if err != nil {
		return "", err
	}

	session := &model.Session{
		Key:       s.ID,
		Values:    values,
		ExpiresOn: time.Now().Add(time.Duration(s.Options.MaxAge) * time.Second),
	}

	if err := sm.datasvc.PutSession(context.Background(), session); err != nil {
		return "", err
	}

	return encoded, nil
}
