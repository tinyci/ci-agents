package datasvc

import (
	transport "github.com/erikh/go-transport"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/db"
	"github.com/tinyci/ci-agents/db/protoconv"
	"github.com/tinyci/ci-agents/testutil"
	"google.golang.org/grpc"
)

// MakeDataServer makes an instance of the datasvc on port 6000. It returns a
// chan which can be closed to terminate it, and any boot-time errors.
func MakeDataServer() (*DataServer, chan struct{}, error) {
	h := &grpcHandler.H{
		Service: config.Service{
			UseDB: true,
			Name:  "datasvc",
		},
		UserConfig: config.UserConfig{
			ClientConfig: config.TestClientConfig,
			DSN:          testutil.TestDBConfig,
			Port:         6000,
			URL:          "url",
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	t, err := transport.Listen(nil, "tcp", config.DefaultServices.Data.String())
	if err != nil {
		return nil, nil, err
	}

	srv := grpc.NewServer()

	db, err := db.NewConn(h.UserConfig.DSN)
	if err != nil {
		return nil, nil, err
	}

	ds := &DataServer{H: h, C: protoconv.New(db)}
	data.RegisterDataServer(srv, ds)

	doneChan, err := h.Boot(t, srv, make(chan struct{}))
	return ds, doneChan, err
}
