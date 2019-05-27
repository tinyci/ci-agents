package processors

import (
	"fmt"
	"time"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/log"
)

// Dispatcher dispatches logs based on loglevel.
type Dispatcher interface {
	Debug(...interface{})
	Error(...interface{})
	Info(...interface{})
}

// DispatchTable is a level -> execution function map
type DispatchTable map[string]func(wf Dispatcher, msg *log.LogMessage)

const (
	// LevelDebug is the debug loglevel
	LevelDebug = "DEBUG"
	// LevelError is the error loglevel
	LevelError = "ERROR"
	// LevelInfo is the info loglevel
	LevelInfo = "INFO"
)

var logLevelDispatch = DispatchTable{
	LevelDebug: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Debug(formatMsg(msg))
	},
	LevelError: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Error(formatMsg(msg))
	},
	LevelInfo: func(wf Dispatcher, msg *log.LogMessage) {
		wf.Info(formatMsg(msg))
	},
}

func formatMsg(msg *log.LogMessage) string {
	return fmt.Sprintf("[%v][%s] %s", time.Unix(msg.At.GetSeconds(), int64(msg.At.GetNanos())), msg.Service, msg.Message)
}
