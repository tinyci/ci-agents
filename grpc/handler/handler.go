package handler

import (
	"net"

	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// H is the standard handler for GRPC. It contains similar information to the http version.
type H struct {
	config.UserConfig `yaml:",inline"`
	config.Service    `yaml:",inline"`
}

// Boot boots the service. It returns a done channel for closing and any errors.
func (h *H) Boot(t net.Listener, s *grpc.Server) (chan struct{}, *errors.Error) {
	if h.Service.UseDB {
		var err *errors.Error
		h.Model, err = model.New(h.UserConfig.DSN)
		if err != nil {
			return nil, err
		}
	}

	if err := h.Auth.ParseTokenKey(); err != nil {
		return nil, err
	}

	var err *errors.Error
	h.Clients, err = h.UserConfig.ClientConfig.CreateClients(h.Name)
	if err != nil {
		return nil, err
	}

	doneChan := make(chan struct{})

	go func(t net.Listener, s *grpc.Server) {
		go s.Serve(t)

		<-doneChan
		t.Close()
	}(t, s)

	return doneChan, nil
}
