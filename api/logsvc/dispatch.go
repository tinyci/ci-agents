package logsvc

import (
	"fmt"
	"time"

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
		wf.Debug(formatMsg(msg))
	},
	client.LevelError: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Error(formatMsg(msg))
	},
	client.LevelInfo: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Info(formatMsg(msg))
	},
}

func formatMsg(msg *log.LogMessage) string {
	return fmt.Sprintf("[%v][%s] %s", time.Unix(msg.At.GetSeconds(), int64(msg.At.GetNanos())), msg.Service, msg.Message)
}
