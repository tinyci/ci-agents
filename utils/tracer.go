package utils

import (
	"io"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/tinyci/ci-agents/errors"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
)

func getConfig() (*config.Configuration, *errors.Error) {
	// FIXME Taken from jaeger/opentracing examples; needs tunables.
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, errors.New(err)
	}

	cfg.Sampler = &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	return cfg, nil
}

// CreateTracer creates an opentracing-compatible jaegertracing client.
func CreateTracer(serviceName string) (io.Closer, *errors.Error) {
	cfg, eErr := getConfig()
	if eErr != nil {
		return nil, eErr
	}

	closer, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return nil, errors.New(err)
	}

	return closer, nil
}

// SetUpGRPCTracing configures grpc dial functions for tracing.
func SetUpGRPCTracing(client string) (io.Closer, []grpc.DialOption, *errors.Error) {
	cfg, eErr := getConfig()
	if eErr != nil {
		return nil, nil, eErr
	}

	tracer, closer, err := cfg.New(client)
	if err != nil {
		return nil, nil, errors.New(err)
	}

	return closer,
		[]grpc.DialOption{
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
			grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
		}, nil
}
