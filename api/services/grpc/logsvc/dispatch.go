package logsvc

import (
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/log"
	client "github.com/tinyci/ci-agents/clients/log"
)

// Dispatcher dispatches logs based on loglevel.
type Dispatcher interface {
	Debug(...interface{})
	Error(...interface{})
	Info(...interface{})
}

// DispatchTable is a level -> execution function map
type DispatchTable map[string]func(wf Dispatcher, msg *log.LogMessage)

var logLevelDispatch = DispatchTable{
	client.LevelDebug: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Debug(msg.Message)
	},
	client.LevelError: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Error(msg.Message)
	},
	client.LevelInfo: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Info(msg.Message)
	},
}
