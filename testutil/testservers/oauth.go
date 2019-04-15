package testservers

// a fair amount of this was cribbed and modified from openshift/osin docs and
// examples.

import (
	"net"
	"net/http"

	"github.com/openshift/osin"
	"github.com/tinyci/ci-agents/errors"
	"golang.org/x/oauth2"
)

// TestEndpoint is the endpoint to be consumed by config/auth.go's stuff in tests.
var TestEndpoint = oauth2.Endpoint{
	AuthURL:  "http://localhost:14000/authorize",
	TokenURL: "http://localhost:14000/token",
}

type storage struct {
	clients   map[string]osin.Client
	authorize map[string]*osin.AuthorizeData
	access    map[string]*osin.AccessData
	refresh   map[string]string
}

func newStorage() osin.Storage {
	r := &storage{
		clients:   make(map[string]osin.Client),
		authorize: make(map[string]*osin.AuthorizeData),
		access:    make(map[string]*osin.AccessData),
		refresh:   make(map[string]string),
	}

	r.clients["erikh"] = &osin.DefaultClient{
		Id:          "client id",
		Secret:      "client secret",
		RedirectUri: "http://localhost:6010/login",
	}

	return r
}

func (s *storage) Clone() osin.Storage {
	return s
}

func (s *storage) Close() {
}

func (s *storage) GetClient(id string) (osin.Client, error) {
	if c, ok := s.clients[id]; ok {
		return c, nil
	}
	return nil, osin.ErrNotFound
}

func (s *storage) SetClient(id string, client osin.Client) error {
	s.clients[id] = client
	return nil
}

func (s *storage) SaveAuthorize(data *osin.AuthorizeData) error {
	s.authorize[data.Code] = data
	return nil
}

func (s *storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	if d, ok := s.authorize[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *storage) RemoveAuthorize(code string) error {
	delete(s.authorize, code)
	return nil
}

func (s *storage) SaveAccess(data *osin.AccessData) error {
	s.access[data.AccessToken] = data
	if data.RefreshToken != "" {
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

func (s *storage) LoadAccess(code string) (*osin.AccessData, error) {
	if d, ok := s.access[code]; ok {
		return d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *storage) RemoveAccess(code string) error {
	delete(s.access, code)
	return nil
}

func (s *storage) LoadRefresh(code string) (*osin.AccessData, error) {
	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, osin.ErrNotFound
}

func (s *storage) RemoveRefresh(code string) error {
	delete(s.refresh, code)
	return nil
}

// BootOAuthService boots an oauth service
// FIXME make this cancelable
func BootOAuthService() (chan struct{}, *errors.Error) {
	server := osin.NewServer(osin.NewServerConfig(), newStorage())

	sm := &http.ServeMux{}

	sm.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
			ar.Authorized = true
			server.FinishAuthorizeRequest(resp, r, ar)
		}
		osin.OutputJSON(resp, w, r)
	})

	sm.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		resp := server.NewResponse()
		defer resp.Close()

		if ar := server.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			server.FinishAccessRequest(resp, r, ar)
		}
		osin.OutputJSON(resp, w, r)
	})

	s := &http.Server{Handler: sm}
	l, err := net.Listen("tcp", ":14000")
	if err != nil {
		return nil, errors.New(err)
	}

	doneChan := make(chan struct{})
	go func() {
		<-doneChan
		s.Close()
		l.Close()
	}()

	go s.Serve(l)

	return doneChan, nil
}
