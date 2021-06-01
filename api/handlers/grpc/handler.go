package grpc

import (
	"context"
	"fmt"
	"io"
	"net"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/db"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc"
)

// H is the standard handler for GRPC. It contains similar information to the http version.
type H struct {
	config.UserConfig `yaml:",inline"`
	config.Service    `yaml:",inline"`
	Model             *db.Model
}

// CreateServer creates the grpc server
func (h *H) CreateServer(installLogger bool) (*grpc.Server, io.Closer, error) {
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

	if installLogger {
		return grpc.NewServer(grpc.StreamInterceptor(h.logStreamInterceptor), grpc.UnaryInterceptor(h.logUnaryInterceptor)), nil, nil
	}

	return grpc.NewServer(), nil, nil
}

func (h *H) logUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if h.Service.Clients.Log != nil {
		u := uuid.New()
		started := time.Now()
		go h.Service.Clients.Log.WithFields(log.FieldMap{
			"method":    path.Base(info.FullMethod),
			"service":   h.Service.Name,
			"uuid":      u.String(),
			"startedAt": fmt.Sprintf("%v", started),
		}).Debug(ctx, "")

		res, err := handler(ctx, req)
		if err != nil {
			go h.Service.Clients.Log.WithFields(log.FieldMap{
				"method":     path.Base(info.FullMethod),
				"service":    h.Service.Name,
				"uuid":       u.String(),
				"finishedAt": fmt.Sprintf("%v", time.Now()),
				"duration":   fmt.Sprintf("%v", time.Since(started)),
			}).Error(ctx, err)
		} else {
			go h.Service.Clients.Log.WithFields(log.FieldMap{
				"method":     path.Base(info.FullMethod),
				"service":    h.Service.Name,
				"uuid":       u.String(),
				"finishedAt": fmt.Sprintf("%v", time.Now()),
				"duration":   fmt.Sprintf("%v", time.Since(started)),
			}).Debug(ctx, "")
		}

		return res, err
	}

	return handler(ctx, req)
}

func (h *H) logStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if h.Service.Clients.Log != nil {
		u := uuid.New()
		started := time.Now()
		go h.Service.Clients.Log.WithFields(log.FieldMap{
			"method":    path.Base(info.FullMethod),
			"service":   h.Service.Name,
			"uuid":      u.String(),
			"startedAt": fmt.Sprintf("%v", started),
		}).Debug(ss.Context(), "")

		err := handler(srv, ss)

		go h.Service.Clients.Log.WithFields(log.FieldMap{
			"method":     path.Base(info.FullMethod),
			"service":    h.Service.Name,
			"uuid":       u.String(),
			"finishedAt": fmt.Sprintf("%v", time.Now()),
			"duration":   fmt.Sprintf("%v", time.Since(started)),
		}).Debug(ss.Context(), "")

		return err
	}

	return handler(srv, ss)
}

// Boot boots the service. It returns a done channel for closing and any errors.
func (h *H) Boot(t net.Listener, s *grpc.Server, finished chan struct{}) (chan struct{}, error) {
	if h.Service.UseDB {
		var err error
		h.Model, err = db.Open(&h.UserConfig)
		if err != nil {
			return nil, err
		}

		if size, ok := h.UserConfig.ServiceConfig["db_pool_size"].(int); ok {
			h.Model.SetConnPoolSize(size)
		}
	}

	if err := h.Auth.ParseTokenKey(); err != nil {
		return nil, err
	}

	var err error
	h.Service.Clients, err = h.UserConfig.ClientConfig.CreateClients(h.UserConfig, h.Name)
	if err != nil {
		return nil, err
	}

	doneChan := make(chan struct{})
	started := make(chan struct{})

	go func(t net.Listener, s *grpc.Server) {
		close(started)
		if err := s.Serve(t); err != nil {
			h.Clients.Log.Error(context.Background(), err)
		}
	}(t, s)

	go func(t net.Listener, s *grpc.Server) {
		<-started
		<-doneChan
		h.Service.Clients.CloseClients()
		s.GracefulStop()
		t.Close()
		close(finished)
	}(t, s)

	<-started
	return doneChan, nil
}
