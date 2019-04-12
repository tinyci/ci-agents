package processors

import "github.com/tinyci/ci-agents/grpc/handler"

// DataServer is the handle into the GRPC subsystem for the datasvc.
type DataServer struct {
	H *handler.H
}
