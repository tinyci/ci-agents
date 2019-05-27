package processors

import "github.com/tinyci/ci-agents/ci-gen/grpc/handler"

// DataServer is the handle into the GRPC subsystem for the datasvc.
type DataServer struct {
	H *handler.H
}
