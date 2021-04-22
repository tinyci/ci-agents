package logsvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/sirupsen/logrus"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Put submits a log message to the service, which in our current case, echoes it to stdout by way of sirupsen/logrus.
func (ls *LogServer) Put(ctx context.Context, lm *log.LogMessage) (*empty.Empty, error) {
	dispatcher, ok := ls.DispatchTable[lm.GetLevel()]
	if !ok {
		return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "Invalid log level %q", lm.GetLevel())
	}

	fields := map[string]interface{}{}

	for key, value := range lm.Fields.Fields {
		switch kind := value.GetKind().(type) {
		case *_struct.Value_StringValue:
			fields[key] = kind.StringValue
		default:
			return &empty.Empty{}, status.Errorf(codes.FailedPrecondition, "%q must be a string value", key)
		}
	}

	dispatcher(logrus.WithFields(fields), lm)
	return &empty.Empty{}, nil
}
