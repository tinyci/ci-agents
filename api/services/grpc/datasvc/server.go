package datasvc

import grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"

// DataServer is the handle into the GRPC subsystem for the datasvc.
type DataServer struct {
	H *grpcHandler.H
}
