package logsvc

import "github.com/sirupsen/logrus"

// LogServer is the handle into the logging grpc service
type LogServer struct {
	DispatchTable DispatchTable
	Level         logrus.Level
}

// New creates a new LogServer.
func New(table DispatchTable, level logrus.Level) *LogServer {
	if table == nil {
		table = logLevelDispatch
	}

	logrus.SetLevel(level)

	return &LogServer{DispatchTable: table, Level: level}
}

func (ls *LogServer) changeLevel(level logrus.Level) {
	logrus.SetLevel(level)
	ls.Level = level
}
