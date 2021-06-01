package logsvc

import (
	transport "github.com/erikh/go-transport"
	"github.com/sirupsen/logrus"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/log"
	client "github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/config"
	"google.golang.org/grpc"
)

// MakeLogServer makes a logsvc.
func MakeLogServer() (*grpcHandler.H, *LogServer, chan struct{}, *LogJournal, error) {
	journal := &LogJournal{Journal: map[string][]*log.LogMessage{}}

	logDispatch := DispatchTable{
		client.LevelDebug: func(wf Dispatcher, msg *log.LogMessage) {
			journal.Append(client.LevelDebug, msg)
		},
		client.LevelError: func(wf Dispatcher, msg *log.LogMessage) {
			journal.Append(client.LevelError, msg)
		},
		client.LevelInfo: func(wf Dispatcher, msg *log.LogMessage) {
			journal.Append(client.LevelInfo, msg)
		},
	}

	h := &grpcHandler.H{
		Service: config.Service{
			Name: "logsvc",
		},
		UserConfig: config.UserConfig{
			Port: 6005,
			// FIXME this is really dumb and should be unnecessary
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	t, err := transport.Listen(nil, "tcp", config.DefaultServices.Log.String())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	srv := grpc.NewServer()
	service := New(logDispatch, logrus.DebugLevel)
	log.RegisterLogServer(srv, service)

	doneChan, err := h.Boot(t, srv, make(chan struct{}))
	return h, service, doneChan, journal, err
}
