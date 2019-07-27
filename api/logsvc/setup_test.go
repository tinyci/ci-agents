package logsvc

import (
	"testing"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	client "github.com/tinyci/ci-agents/clients/log"
)

type logsvcSuite struct {
	logsvcHandler  *handler.H
	logsvcDoneChan chan struct{}
	journal        *LogJournal
}

var _ = check.Suite(&logsvcSuite{})

func TestLogSvc(t *testing.T) {
	check.TestingT(t)
}

func (ls *logsvcSuite) SetUpTest(c *check.C) {
	var err error
	ls.logsvcHandler, ls.logsvcDoneChan, ls.journal, err = MakeLogServer()
	c.Assert(err, check.IsNil)

	client.ConfigureRemote("localhost:6005", nil, false)
}

func (ls *logsvcSuite) TearDownTest(c *check.C) {
	close(ls.logsvcDoneChan)
	time.Sleep(100 * time.Millisecond)
}
