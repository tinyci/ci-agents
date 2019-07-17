package logsvc

import "github.com/sirupsen/logrus"

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

// LogServer is the handle into the logging grpc service
type LogServer struct {
	DispatchTable DispatchTable
}

// New creates a new LogServer.
func New(table DispatchTable) *LogServer {
	if table == nil {
		table = logLevelDispatch
	}

	return &LogServer{DispatchTable: table}
}
