package test

import (
	"testing"
	"time"

	check "github.com/erikh/check"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/grpc/handler"
	"github.com/tinyci/ci-agents/testutil/testservers"
)

type logsvcSuite struct {
	logsvcHandler  *handler.H
	logsvcDoneChan chan struct{}
	journal        *testservers.LogJournal
}

var _ = check.Suite(&logsvcSuite{})

func TestLogSvc(t *testing.T) {
	check.TestingT(t)
}

func (ls *logsvcSuite) SetUpTest(c *check.C) {
	var err error
	ls.logsvcHandler, ls.logsvcDoneChan, ls.journal, err = testservers.MakeLogServer()
	c.Assert(err, check.IsNil)

	log.ConfigureRemote("localhost:6005", nil)
}

func (ls *logsvcSuite) TearDownTest(c *check.C) {
	close(ls.logsvcDoneChan)
	time.Sleep(100 * time.Millisecond)
}
