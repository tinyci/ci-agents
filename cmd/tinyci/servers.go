package main

import (
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/api/services/grpc/assetsvc"
	"github.com/tinyci/ci-agents/api/services/grpc/auth/github"
	"github.com/tinyci/ci-agents/api/services/grpc/datasvc"
	"github.com/tinyci/ci-agents/api/services/grpc/logsvc"
	"github.com/tinyci/ci-agents/api/services/grpc/queuesvc"
	repoGithub "github.com/tinyci/ci-agents/api/services/grpc/repository/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/asset"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/auth"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/log"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/queue"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/cmdlib"
	"github.com/tinyci/ci-agents/config"
	"google.golang.org/grpc"
)

var servers = []*cmdlib.GRPCServer{
	{
		Name:           "assetsvc",
		Description:    "Asset & Log management for tinyCI",
		DefaultService: config.DefaultServices.Asset,
		RegisterService: func(s *grpc.Server, h *grpcHandler.H) error {
			asset.RegisterAssetServer(s, &assetsvc.AssetServer{H: h})
			return nil
		},
	},
	{
		Name:           "github-authsvc",
		Description:    "Github conduit for authentication in tinyCI",
		DefaultService: config.DefaultServices.Auth,
		RegisterService: func(s *grpc.Server, h *grpcHandler.H) error {
			auth.RegisterAuthServer(s, &github.AuthServer{H: h})
			return nil
		},
	},
	{
		Name:           "datasvc",
		Description:    "datasvc is the conduit for tinyCI to talk to a data store.",
		UseDB:          true,
		UseSessions:    true,
		DefaultService: config.DefaultServices.Data,
		RegisterService: func(s *grpc.Server, h *grpcHandler.H) error {
			data.RegisterDataServer(s, &datasvc.DataServer{H: h})
			return nil
		},
	},
	{
		Name:           "logsvc",
		Description:    "Centralized logging for tinyCI",
		DefaultService: config.DefaultServices.Log,
		RegisterService: func(s *grpc.Server, h *grpcHandler.H) error {
			log.RegisterLogServer(s, logsvc.New(nil))
			return nil
		},
	},
	{
		Name:           "queuesvc",
		Description:    "Queue & Run management for tinyCI",
		DefaultService: config.DefaultServices.Queue,
		RegisterService: func(s *grpc.Server, h *grpcHandler.H) error {
			queue.RegisterQueueServer(s, &queuesvc.QueueServer{H: h})
			return nil
		},
	},
	{
		Name:           "github-reposvc",
		Description:    "Github conduit for repository management in tinyCI",
		DefaultService: config.DefaultServices.Repository,
		RegisterService: func(s *grpc.Server, h *grpcHandler.H) error {
			repository.RegisterRepositoryServer(s, &repoGithub.RepositoryServer{H: h})
			return nil
		},
	},
}
