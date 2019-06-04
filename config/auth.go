package config

import (
	"fmt"
	"strings"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
)

// SessionErrorsKey is the key used to retrieve the errors from the sessions table.
const SessionErrorsKey = "errors"

// SessionKey is the name of the cookie where the session will be stored.
const SessionKey = "tinyci"

// AuthConfig is the configuration for auth and secrets in the case auth isn't
// used.
type AuthConfig struct {
	SessionCryptKey   string              `yaml:"session_crypt_key"`
	TokenCryptKey     string              `yaml:"token_crypt_key"`
	FixedCapabilities map[string][]string `yaml:"fixed_capabilities"`

	sessionCryptKey []byte
	tokenCryptKey   []byte
}

// CertConfig manages the configuration of client and server certs for handler
// services.
type CertConfig struct {
	CAFile   string `yaml:"ca"`
	CertFile string `yaml:"cert"`
	KeyFile  string `yaml:"key"`
}

// Validate the certificate configuration (if supplied)
func (cc *CertConfig) Validate() *errors.Error {
	ca := strings.TrimSpace(cc.CAFile)
	cert := strings.TrimSpace(cc.CertFile)
	key := strings.TrimSpace(cc.KeyFile)

	if ca == "" && cert == "" && key == "" {
		return nil // no certificate information supplied
	}

	if ca == "" {
		return errors.New("missing ca certificate in TLS configuration")
	}

	if cert == "" {
		return errors.New("missing certificate in TLS configuration")
	}

	if key == "" {
		return errors.New("missing key in TLS configuration")
	}

	return nil
}

// Validate ensures the auth configuration is sane.
func (ac *AuthConfig) Validate(parseCrypt bool) *errors.Error {
	if parseCrypt {
		ac.sessionCryptKey = types.DecodeKey(ac.SessionCryptKey)
		if err := validateAESKey(ac.sessionCryptKey); err != nil {
			return err
		}
	}

	if ac.FixedCapabilities == nil {
		ac.FixedCapabilities = map[string][]string{}
	}

	return nil
}

func validateAESKey(key []byte) *errors.Error {
	switch len(key) {
	case 8, 16, 32:
	default:
		return errors.New("AES keys must be a multiple of 8, 16, or 32 bytes")
	}

	return nil
}

// ParseTokenKey reads the key from the config, validates it, and assigns it to the appropriate variables
func (ac *AuthConfig) ParseTokenKey() *errors.Error {
	ac.tokenCryptKey = types.DecodeKey(ac.TokenCryptKey)
	if err := validateAESKey(ac.tokenCryptKey); err != nil {
		return err
	}

	model.TokenCryptKey = ac.tokenCryptKey
	return nil
}

// ParsedSessionCryptKey returns the parsed session crypt key
func (ac *AuthConfig) ParsedSessionCryptKey() []byte {
	return ac.sessionCryptKey
}

// Load loads the cert based on the provided config and returns it.
func (cc CertConfig) Load() (*transport.Cert, *errors.Error) {
	if cc.CAFile == "" || cc.CertFile == "" || cc.KeyFile == "" {
		fmt.Println("Some TLS parameters were missing; running insecure!")
		return nil, nil
	}

	cert, err := transport.LoadCert(cc.CAFile, cc.CertFile, cc.KeyFile, "")
	if err != nil {
		return nil, errors.New(err)
	}

	return cert, nil
}
