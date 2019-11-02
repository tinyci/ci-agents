package assetsvc

import (
	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/asset"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"google.golang.org/grpc"
)

// MakeAssetServer makes an instance of the assetsvc on port 6000. It returns a
// chan which can be closed to terminate it, and any boot-time errors.
func MakeAssetServer() (*handler.H, chan struct{}, error) {
	t, err := transport.Listen(nil, "tcp", config.DefaultServices.Asset.String())
	if err != nil {
		return nil, nil, errors.New(err)
	}

	h := &handler.H{
		Service: config.Service{Name: "assetsvc"},
		UserConfig: config.UserConfig{
			Auth: config.AuthConfig{
				TokenCryptKey: "1431d583a48a00243cc3d3d596ed362d77c50be4848dbf0d2f52bab841f072f9",
			},
		},
	}

	srv := grpc.NewServer()
	asset.RegisterAssetServer(srv, &AssetServer{H: h})

	doneChan, err := h.Boot(t, srv, make(chan struct{}))
	return h, doneChan, errors.New(err)
}
