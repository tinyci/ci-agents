package handler

import (
	"context"
	"io"
	"net"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc"
)

// H is the standard handler for GRPC. It contains similar information to the http version.
type H struct {
	config.UserConfig `yaml:",inline"`
	config.Service    `yaml:",inline"`
}

// CreateServer creates the grpc server
func (h *H) CreateServer() (*grpc.Server, io.Closer, error) {
	if h.EnableTracing {
		closer, err := utils.CreateTracer(h.Name)
		if err != nil {
			return nil, nil, err
		}
		s := grpc.NewServer(
			grpc.UnaryInterceptor(
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())),
			grpc.StreamInterceptor(
				otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer())))
		return s, closer, nil
	}

	return grpc.NewServer(), nil, nil
}

// Boot boots the service. It returns a done channel for closing and any errors.
func (h *H) Boot(t net.Listener, s *grpc.Server, finished chan struct{}) (chan struct{}, error) {
	if h.Service.UseDB {
		var err error
		h.Model, err = model.New(h.UserConfig.DSN)
		if err != nil {
			return nil, err
		}
	}

	if err := h.Auth.ParseTokenKey(); err != nil {
		return nil, err
	}

	var err error
	h.Clients, err = h.UserConfig.ClientConfig.CreateClients(h.UserConfig, h.Name)
	if err != nil {
		return nil, err
	}

	doneChan := make(chan struct{})

	go func(t net.Listener, s *grpc.Server) {
		if err := s.Serve(t); err != nil {
			h.Clients.Log.Error(context.Background(), err)
		}
	}(t, s)

	go func(t net.Listener, s *grpc.Server) {
		<-doneChan
		s.GracefulStop()
		t.Close()
		h.Clients.CloseClients()
		close(finished)
	}(t, s)

	return doneChan, nil
}
