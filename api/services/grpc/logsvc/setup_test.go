package logsvc

import (
	"testing"
	"time"

	check "github.com/erikh/check"
	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	client "github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/config"
)

type logsvcSuite struct {
	service        *LogServer
	logsvcHandler  *grpcHandler.H
	logsvcDoneChan chan struct{}
	journal        *LogJournal
}

var _ = check.Suite(&logsvcSuite{})

func TestLogSvc(t *testing.T) {
	check.TestingT(t)
}

func (ls *logsvcSuite) SetUpTest(c *check.C) {
	var err error
	ls.logsvcHandler, ls.service, ls.logsvcDoneChan, ls.journal, err = MakeLogServer()
	c.Assert(err, check.IsNil)

	c.Assert(client.ConfigureRemote(config.DefaultServices.Log.String(), nil, false), check.IsNil)
}

func (ls *logsvcSuite) TearDownTest(c *check.C) {
	close(ls.logsvcDoneChan)
	time.Sleep(100 * time.Millisecond)
}
