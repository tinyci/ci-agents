package main

import (
	"os"

	"github.com/tinyci/ci-agents/api/auth/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/cmdlib"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"google.golang.org/grpc"
)

// Version is the version of this service.
const Version = "1.0.0"

// TinyCIVersion is the version of tinyci supporting this service.
var TinyCIVersion = "" // to be changed by build processes

func main() {
	s := &cmdlib.GRPCServer{
		Name:           "github-authsvc",
		Description:    "Github conduit for authentication in tinyCI",
		AppVersion:     Version,
		TinyCIVersion:  TinyCIVersion,
		DefaultService: config.DefaultServices.Auth,
		RegisterService: func(s *grpc.Server, h *handler.H) error {
			auth.RegisterAuthServer(s, &github.AuthServer{H: h})
			return nil
		},
	}

	if err := s.Make().Run(os.Args); err != nil {
		errors.New(err).(errors.Error).Exit()
	}
}
