package assetsvc

import (
	"os"
	"testing"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	client "github.com/tinyci/ci-agents/clients/asset"
)

type assetsvcSuite struct {
	assetsvcHandler  *handler.H
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

	as.assetClient, err = client.NewClient("localhost:6002", nil, false)
	c.Assert(err, check.IsNil)
}

func (as *assetsvcSuite) TearDownTest(c *check.C) {
	close(as.assetsvcDoneChan)
	time.Sleep(100 * time.Millisecond)
}
