package datasvc

import (
	"testing"
	"time"

	check "github.com/erikh/check"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/api/services/grpc/logsvc"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/testutil"
	"github.com/tinyci/ci-agents/testutil/testclients"
)

type datasvcSuite struct {
	model        *model.Model
	dataDoneChan chan struct{}
	dataServer   *grpcHandler.H
	logDoneChan  chan struct{}
	logServer    *grpcHandler.H
	client       *testclients.DataClient
}

var _ = check.Suite(&datasvcSuite{})

func TestDataSvc(t *testing.T) {
	check.TestingT(t)
}

func (ds *datasvcSuite) SetUpTest(c *check.C) {
	testutil.WipeDB(c)

	var err error
	ds.model, err = model.New(testutil.TestDBConfig)
	c.Assert(err, check.IsNil)

	ds.dataServer, ds.dataDoneChan, err = MakeDataServer()
	c.Assert(err, check.IsNil)

	ds.logServer, ds.logDoneChan, _, err = logsvc.MakeLogServer()
	c.Assert(err, check.IsNil)

	ds.client, err = testclients.NewDataClient()
	c.Assert(err, check.IsNil)
}

func (ds *datasvcSuite) TearDownTest(c *check.C) {
	close(ds.logDoneChan)
	close(ds.dataDoneChan)
	time.Sleep(100 * time.Millisecond)
}
