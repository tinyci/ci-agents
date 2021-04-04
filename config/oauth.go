package config

import (
	"net/url"
	"strings"

	"errors"

	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
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
func (oc OAuthConfig) Validate() error {
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
		return utils.WrapError(err, "parsing oauth2 redirect_url")
	}

	return nil
}

// SetDefaultGithubClient sets the default github client which is necessary for
// many testing scenarios. Not to be used in typical code.
func SetDefaultGithubClient(client github.Client, username string) {
	if username == "" {
		username = "default"
	}

	githubClientMutex.Lock()
	defer githubClientMutex.Unlock()
	defaultGithubClientMap[username] = client
}

// DefaultGithubClient returns the default github client set by SetDefaultGithubClient, if any.
func DefaultGithubClient(username string) github.Client {
	githubClientMutex.RLock()
	defer githubClientMutex.RUnlock()

	if username == "" {
		username = "default"
	}

	return defaultGithubClientMap[username]
}

// GithubClient either returns the client for the token.
func (oc OAuthConfig) GithubClient(token *types.OAuthToken) github.Client {
	client := DefaultGithubClient(token.Username)
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
