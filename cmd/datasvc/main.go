package main

import (
	"os"

	"github.com/tinyci/ci-agents/api/datasvc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
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
		Name:           "datasvc",
		Description:    "datasvc is the conduit for tinyCI to talk to a data store.",
		AppVersion:     Version,
		TinyCIVersion:  TinyCIVersion,
		UseDB:          true,
		UseSessions:    true,
		DefaultService: config.DefaultServices.Data,
		RegisterService: func(s *grpc.Server, h *handler.H) error {
			data.RegisterDataServer(s, &datasvc.DataServer{H: h})
			return nil
		},
	}

	if err := s.Make().Run(os.Args); err != nil {
		errors.New(err).Exit()
	}
}
