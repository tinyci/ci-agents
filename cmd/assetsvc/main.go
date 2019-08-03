package main

import (
	"os"

	"github.com/tinyci/ci-agents/api/assetsvc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/asset"
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
		Name:           "assetsvc",
		Description:    "Asset & Log management for tinyCI",
		AppVersion:     Version,
		TinyCIVersion:  TinyCIVersion,
		DefaultService: config.DefaultServices.Asset,
		RegisterService: func(s *grpc.Server, h *handler.H) *errors.Error {
			asset.RegisterAssetServer(s, &assetsvc.AssetServer{H: h})
			return nil
		},
	}

	if err := s.Make().Run(os.Args); err != nil {
		errors.New(err).Exit()
	}
}
