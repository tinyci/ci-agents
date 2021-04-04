package queuesvc

import (
	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/queue"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/testutil"
	"google.golang.org/grpc"
)

// MakeQueueServer makes an instance of the queuesvc on port 6001. It returns a
// chan which can be closed to terminate it, and any boot-time errors.
func MakeQueueServer() (*handler.H, chan struct{}, error) {
	h := &handler.H{
		Service: config.Service{
			Name: "queuesvc",
		},
		UserConfig: config.UserConfig{
			DSN:          testutil.TestDBConfig,
			ClientConfig: config.TestClientConfig,
			URL:          "url",
			Port:         6001,
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	t, err := transport.Listen(nil, "tcp", config.DefaultServices.Queue.String())
	if err != nil {
		return nil, nil, err
	}

	srv := grpc.NewServer()
	queue.RegisterQueueServer(srv, &QueueServer{H: h})

	doneChan, err := h.Boot(t, srv, make(chan struct{}))
	return h, doneChan, err
}
