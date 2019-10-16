package config

import (
	"net/url"
	"strings"

	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/types"
	"golang.org/x/oauth2"
	ghoauth "golang.org/x/oauth2/github"
)

// OAuthRepositoryScope determines how to get enough information out of github
// to manipulate repositories. The default is to not ask for any real
// permissions; but for this behavior we need more than that.
var OAuthRepositoryScope = []string{"repo"}

// DefaultEndpoint is the default endpoint for oauth2 operations.
var DefaultEndpoint = ghoauth.Endpoint

// OAuthConfig configures the oauth end of the uiservice handler, specifically
// focusing around the application credentials and login process.
type OAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

// Validate validates the oauth configuration
func (oc OAuthConfig) Validate() *errors.Error {
	if strings.TrimSpace(oc.ClientID) == "" {
		return errors.New("oauth2 client_id was missing")
	}

	if strings.TrimSpace(oc.ClientSecret) == "" {
		return errors.New("oauth2 client_secret was missing")
	}

	if strings.TrimSpace(oc.RedirectURL) == "" {
		return errors.New("oauth2 redirect_url was missing")
	}

	_, err := url.Parse(oc.RedirectURL)
	if err != nil {
		return errors.New(err).Wrap("parsing oauth2 redirect_url")
	}

	return nil
}

// SetDefaultGithubClient sets the default github client which is necessary for
// many testing scenarios. Not to be used in typical code.
func SetDefaultGithubClient(client github.Client) {
	githubClientMutex.Lock()
	defer githubClientMutex.Unlock()
	defaultGithubClient = client
}

// DefaultGithubClient returns the default github client set by SetDefaultGithubClient, if any.
func DefaultGithubClient() github.Client {
	githubClientMutex.RLock()
	defer githubClientMutex.RUnlock()
	return defaultGithubClient
}

// GithubClient either returns the client for the token, or if NoAuth is set
// returns the default client.
func (oc OAuthConfig) GithubClient(token *types.OAuthToken) github.Client {
	client := DefaultGithubClient()
	if client != nil {
		return client
	}

	return github.NewClientFromAccessToken(token.Token)
}

// Config returns the oauth configuration if one was provided.
func (oc OAuthConfig) Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     oc.ClientID,
		ClientSecret: oc.ClientSecret,
		RedirectURL:  oc.RedirectURL,
		Endpoint:     DefaultEndpoint,
		Scopes:       scopes,
	}
}
