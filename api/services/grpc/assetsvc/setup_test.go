package assetsvc

import (
	"os"
	"testing"
	"time"

	check "github.com/erikh/check"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	client "github.com/tinyci/ci-agents/clients/asset"
	"github.com/tinyci/ci-agents/config"
)

type assetsvcSuite struct {
	assetsvcHandler  *grpcHandler.H
	assetsvcDoneChan chan struct{}
	assetClient      *client.Client
}

var _ = check.Suite(&assetsvcSuite{})

func TestAssetSvc(t *testing.T) {
	check.TestingT(t)
}

func (as *assetsvcSuite) SetUpTest(c *check.C) {
	os.RemoveAll("/var/tinyci/logs")
	var err error
	as.assetsvcHandler, as.assetsvcDoneChan, err = MakeAssetServer()
	c.Assert(err, check.IsNil)

	as.assetClient, err = client.NewClient(config.DefaultServices.Asset.String(), nil, false)
	c.Assert(err, check.IsNil)
}

func (as *assetsvcSuite) TearDownTest(c *check.C) {
	close(as.assetsvcDoneChan)
	time.Sleep(100 * time.Millisecond)
}
