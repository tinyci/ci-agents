package testservers

import (
	"github.com/tinyci/ci-agents/ci-gen/gen/svc/uisvc/restapi"
	d "github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/clients/tinyci"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
)

// MakeUIServer makes a uisvc.
func MakeUIServer(client github.Client) (*handlers.H, chan struct{}, *tinyci.Client, *tinyci.Client, *errors.Error) {
	h := &handlers.H{
		Config: restapi.HandlerConfig{},
		Service: config.Service{
			Name: "uisvc",
		},
		UserConfig: config.UserConfig{
			OAuth: config.OAuthConfig{
				ClientID:     "client id",
				ClientSecret: "client secret",
				RedirectURL:  "http://localhost:6010/login",
			},
			ClientConfig: config.TestClientConfig,
			URL:          "http://localhost",
			Auth: config.AuthConfig{
				SessionCryptKey: "0431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
				TokenCryptKey:   "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
			Websockets: config.Websockets{
				InsecureWebSockets: true,
			},
			RequestLogging: true,
			Port:           6010,
		},
	}

	d, err := d.New("localhost:6000", nil, false)
	if err != nil {
		return nil, nil, nil, nil, errors.New(err)
	}

	config.DefaultGithubClient = client
	finished := make(chan struct{})
	doneChan, err := handlers.Boot(nil, h, finished)
	if err != nil {
		return nil, nil, nil, nil, errors.New(err)
	}

	u, err := d.PutUser(&model.User{Username: "erikh", Token: &types.OAuthToken{Token: "dummy", Scopes: []string{"repo"}}})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for _, cap := range model.AllCapabilities {
		if err := d.AddCapability(u, cap); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	token, err := d.GetToken("erikh")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tc, err := tinyci.New("http://localhost:6010", token, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if _, err := tc.Errors(); err != nil { // connectivity check
		return nil, nil, nil, nil, err
	}

	_, err = d.PutUser(&model.User{Username: "erikh2", Token: &types.OAuthToken{Token: "dummy"}})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	token, err = d.GetToken("erikh2")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	utc, err := tinyci.New("http://localhost:6010", token, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return h, doneChan, tc, utc, nil
}
