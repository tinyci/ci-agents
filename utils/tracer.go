package utils

import (
	"io"

	"github.com/tinyci/ci-agents/errors"
	"github.com/uber/jaeger-client-go/config"
)

// CreateTracer creates an opentracing-compatible jaegertracing client.
func CreateTracer(serviceName string) (io.Closer, *errors.Error) {
	// FIXME Taken from jaeger/opentracing examples; needs tunables.
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, errors.New(err)
	}

	cfg.Sampler = &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	closer, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return nil, errors.New(err)
	}

	return closer, nil
}
