package processors

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/sirupsen/logrus"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/log"
	"github.com/tinyci/ci-agents/errors"
	"google.golang.org/grpc/codes"
)

// Put submits a log message to the service, which in our current case, echoes it to stdout by way of sirupsen/logrus.
func (ls *LogServer) Put(ctx context.Context, lm *log.LogMessage) (*empty.Empty, error) {
	dispatcher, ok := ls.DispatchTable[lm.GetLevel()]
	if !ok {
		return &empty.Empty{}, errors.Errorf("Invalid log level %q", lm.GetLevel()).ToGRPC(codes.FailedPrecondition)
	}

	fields := map[string]interface{}{}

	for key, value := range lm.Fields.Fields {
		switch kind := value.GetKind().(type) {
		case *_struct.Value_StringValue:
			fields[key] = kind.StringValue
		default:
			return &empty.Empty{}, errors.Errorf("%q must be a string value", key).ToGRPC(codes.FailedPrecondition)
		}
	}

	dispatcher(logrus.WithFields(fields), lm)
	return &empty.Empty{}, nil
}
