package datasvc

import (
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/db/protoconv"
)

// DataServer is the handle into the GRPC subsystem for the datasvc.
type DataServer struct {
	H *grpcHandler.H
	C *protoconv.Converter
}

// New asdf
func New() {
}
