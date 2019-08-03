package main

import (
	"os"

	"github.com/tinyci/ci-agents/api/queuesvc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/queue"
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
		Name:           "queuesvc",
		Description:    "Queue & Run management for tinyCI",
		AppVersion:     Version,
		TinyCIVersion:  TinyCIVersion,
		DefaultService: config.DefaultServices.Queue,
		RegisterService: func(s *grpc.Server, h *handler.H) *errors.Error {
			queue.RegisterQueueServer(s, &queuesvc.QueueServer{H: h})
			return nil
		},
	}

	if err := s.Make().Run(os.Args); err != nil {
		errors.New(err).Exit()
	}
}
