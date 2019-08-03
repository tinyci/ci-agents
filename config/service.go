package config

import (
	"net/url"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/tinyci/ci-agents/clients/asset"
	"github.com/tinyci/ci-agents/clients/data"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/clients/queue"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
)

// DefaultGithubClient if set, will override any requested github client.
var DefaultGithubClient github.Client

// TestClientConfig is a default test client configuration
var TestClientConfig = ClientConfig{
	Data:       "localhost:6000",
	Queue:      "localhost:6001",
	UI:         "http://localhost:6010",
	Asset:      "localhost:6002",
	Repository: "localhost:6003",
	Log:        "localhost:6005",
}

// ServiceConfig is the pre-normalized version of the config struct
type ServiceConfig map[string]interface{}

// UserConfig is the user-supplied configuration parsed from yaml.
type UserConfig struct {
	ServiceConfig ServiceConfig `yaml:"services"`
	ClientConfig  ClientConfig  `yaml:"clients"`

	OAuth          OAuthConfig `yaml:"oauth"`
	Auth           AuthConfig  `yaml:"auth"`
	HookURL        string      `yaml:"hook_url"`
	DSN            string      `yaml:"db"`
	TLS            CertConfig  `yaml:"tls"`
	Websockets     Websockets  `yaml:"websockets"`
	RequestLogging bool        `yaml:"log_requests"`
	Port           uint        `yaml:"port"`
	URL            string      `yaml:"url"`
	EnableTracing  bool        `yaml:"enable_tracing"`
}

// Service is the internal configuration for a service
type Service struct {
	UseSessions bool            `yaml:"-"`
	UseDB       bool            `yaml:"-"`
	NoTLSServer bool            `yaml:"-"`
	Formats     strfmt.Registry `yaml:"-"`

	Clients *Clients `yaml:"-"`

	Model *model.Model `yaml:"-"`
	Name  string       `yaml:"-"`
}

// Clients is a struct that encapsulates the various internal clients we use.
type Clients struct {
	Data  *data.Client
	Queue *queue.Client
	Asset *asset.Client
	Log   *log.SubLogger
}

// ClientConfig configures the clients
type ClientConfig struct {
	Data       string `yaml:"datasvc"`
	UI         string `yaml:"uisvc"`
	Queue      string `yaml:"queuesvc"`
	Asset      string `yaml:"assetsvc"`
	Log        string `yaml:"logsvc"`
	Repository string `yaml:"reposvc"`

	Cert CertConfig `yaml:"tls"`
}

func parseURL(svc, u string) *errors.Error {
	if strings.TrimSpace(u) == "" {
		// URLs do not need to be supplied for all services, so the services themselves
		// will validate this later. This is just for ensuring they're actual URLs.
		return nil
	}

	_, err := url.Parse(u)
	if err != nil {
		return errors.New(err).Wrapf("url for service %v is invalid: %q", svc, u)
	}

	return nil
}

// Validate validates the client configuration to ensure basic needs are met.
func (cc *ClientConfig) Validate() *errors.Error {
	urlmap := map[string]string{
		"datasvc":   cc.Data,
		"uisvc":     cc.UI,
		"queueusvc": cc.Queue,
		"logsvc":    cc.Log,
		"assetsvc":  cc.Asset,
	}

	for svc, u := range urlmap {
		if err := parseURL(svc, u); err != nil {
			return err
		}
	}

	return cc.Cert.Validate()
}

// CreateClients creates all the clients that are populated in the clients struct
func (cc *ClientConfig) CreateClients(uc UserConfig, service string) (*Clients, *errors.Error) {
	if err := cc.Validate(); err != nil {
		return nil, err
	}

	clientCert, err := cc.Cert.Load()
	if err != nil {
		return nil, err
	}

	clients := &Clients{}

	if cc.Log != "" {
		log.ConfigureRemote(cc.Log, clientCert, uc.EnableTracing)
	}

	clients.Log = log.New().WithService(service)

	if cc.Data != "" {
		dc, err := data.New(cc.Data, clientCert, uc.EnableTracing)
		if err != nil {
			return nil, err
		}

		clients.Data = dc
	}

	if cc.Queue != "" {
		qc, err := queue.New(cc.Queue, clientCert, uc.EnableTracing)
		if err != nil {
			return nil, err
		}

		clients.Queue = qc
	}

	if cc.Asset != "" {
		lc, err := asset.NewClient(cc.Asset, clientCert, uc.EnableTracing)
		if err != nil {
			return nil, err
		}

		clients.Asset = lc
	}

	return clients, nil
}

// CloseClients closes all clients.
func (c *Clients) CloseClients() {
	if c.Asset != nil {
		c.Asset.Close()
	}

	if log.RemoteClient != nil { // FIXME don't break the fourth wall
		log.RemoteClient.Close()
	}

	if c.Data != nil {
		c.Data.Close()
	}

	if c.Queue != nil {
		c.Queue.Close()
	}
}
