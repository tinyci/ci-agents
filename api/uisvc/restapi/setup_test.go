package restapi

import (
	"testing"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/api/assetsvc"
	"github.com/tinyci/ci-agents/api/datasvc"
	"github.com/tinyci/ci-agents/api/logsvc"
	"github.com/tinyci/ci-agents/api/queuesvc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/clients/asset"
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/testutil/testclients"
	"github.com/tinyci/ci-agents/testutil/testservers"
)

type uisvcSuite struct {
	datasvcClient  *testclients.DataClient
	queuesvcClient *testclients.QueueClient
	assetsvcClient *asset.Client
	logDoneChan    chan struct{}
	queueDoneChan  chan struct{}
	dataDoneChan   chan struct{}
	assetDoneChan  chan struct{}
	oauthDoneChan  chan struct{}
	logHandler     *handler.H
	dataHandler    *handler.H
	queueHandler   *handler.H
	assetHandler   *handler.H

	logJournal *logsvc.LogJournal
}

var _ = check.Suite(&uisvcSuite{})

func TestUISvc(t *testing.T) {
	config.DefaultEndpoint = testservers.TestEndpoint
	check.TestingT(t)
}

func (us *uisvcSuite) SetUpTest(c *check.C) {
	testutil.WipeDB(c)

	var err error
	us.dataHandler, us.dataDoneChan, err = datasvc.MakeDataServer()
	c.Assert(err, check.IsNil)

	us.queueHandler, us.queueDoneChan, err = queuesvc.MakeQueueServer()
	c.Assert(err, check.IsNil)

	us.assetHandler, us.assetDoneChan, err = assetsvc.MakeAssetServer()
	c.Assert(err, check.IsNil)

	us.logHandler, us.logDoneChan, us.logJournal, err = logsvc.MakeLogServer()
	c.Assert(err, check.IsNil)

	go us.logJournal.Tail()

	us.datasvcClient, err = testclients.NewDataClient()
	c.Assert(err, check.IsNil)

	us.queuesvcClient, err = testclients.NewQueueClient(us.datasvcClient)
	c.Assert(err, check.IsNil)

	us.assetsvcClient, err = asset.NewClient(config.DefaultServices.Asset.String(), nil, false)
	c.Assert(err, check.IsNil)

	us.oauthDoneChan, err = testservers.BootOAuthService()
	c.Assert(err, check.IsNil)
}

func (us *uisvcSuite) TearDownTest(c *check.C) {
	close(us.oauthDoneChan)
	close(us.dataDoneChan)
	close(us.queueDoneChan)
	close(us.logDoneChan)
	close(us.assetDoneChan)
	time.Sleep(100 * time.Millisecond)
}
