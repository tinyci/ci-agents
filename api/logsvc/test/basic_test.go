package test

import (
	"net/http"
	"net/url"

	check "github.com/erikh/check"
	logsvc "github.com/tinyci/ci-agents/api/logsvc/processors"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/model"
)

func (ls *logsvcSuite) Test01Journal(c *check.C) {
	log.New().Info("test")
	j := ls.journal.Journal[logsvc.LevelInfo]
	c.Assert(len(j), check.Equals, 1)
	c.Assert(j[0].Message, check.Equals, "test")
}

func (ls *logsvcSuite) TestLevels(c *check.C) {
	l := log.New()

	l.Info("test")
	l.Infof("test %d", 2)

	j := ls.journal.Journal[logsvc.LevelInfo]
	c.Assert(len(j), check.Equals, 2)
	c.Assert(j[0].Message, check.Equals, "test")
	c.Assert(j[1].Message, check.Equals, "test 2")

	l.Error("test")
	l.Errorf("test %d", 2)

	j = ls.journal.Journal[logsvc.LevelError]
	c.Assert(len(j), check.Equals, 2)
	c.Assert(j[0].Message, check.Equals, "test")
	c.Assert(j[1].Message, check.Equals, "test 2")

	l.Debug("test")
	l.Debugf("test %d", 2)

	j = ls.journal.Journal[logsvc.LevelDebug]
	c.Assert(len(j), check.Equals, 2)
	c.Assert(j[0].Message, check.Equals, "test")
	c.Assert(j[1].Message, check.Equals, "test 2")
}

func (ls *logsvcSuite) TestFields(c *check.C) {
	l := log.New()

	wf := l.WithFields(log.FieldMap{"test": "one", "test2": "two"})

	wf.Info("test")
	wf.Infof("test %d", 2)
	wf.Debug("test")
	wf.Debugf("test %d", 2)
	wf.Error("test")
	wf.Errorf("test %d", 2)

	for _, messages := range ls.journal.Journal {
		c.Assert(len(messages), check.Equals, 2)
		c.Assert(messages[0].Message, check.Equals, "test")
		c.Assert(messages[1].Message, check.Equals, "test 2")

		for i := 0; i < 2; i++ {
			val, ok := messages[i].Fields.Fields["test"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "one")

			val, ok = messages[i].Fields.Fields["test2"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "two")
		}
	}
}

func (ls *logsvcSuite) TestService(c *check.C) {
	l := log.New()

	ws := l.WithService("test")

	ws.Info("test")
	ws.Infof("test %d", 2)
	ws.Debug("test")
	ws.Debugf("test %d", 2)
	ws.Error("test")
	ws.Errorf("test %d", 2)

	for _, messages := range ls.journal.Journal {
		c.Assert(len(messages), check.Equals, 2)
		c.Assert(messages[0].Message, check.Equals, "test")
		c.Assert(messages[1].Message, check.Equals, "test 2")

		for i := 0; i < 2; i++ {
			val, ok := messages[i].Fields.Fields["service"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "test")
		}
	}
}

func (ls *logsvcSuite) TestRequest(c *check.C) {
	l := log.New()

	u, err := url.Parse("http://dogfood.tinyci.org")
	c.Assert(err, check.IsNil)

	req := &http.Request{
		RemoteAddr: "127.0.0.1:1234",
		Method:     "GET",
		URL:        u,
	}

	wr := l.WithRequest(req)

	wr.Info("test")
	wr.Infof("test %d", 2)
	wr.Debug("test")
	wr.Debugf("test %d", 2)
	wr.Error("test")
	wr.Errorf("test %d", 2)

	for _, messages := range ls.journal.Journal {
		c.Assert(len(messages), check.Equals, 2)
		c.Assert(messages[0].Message, check.Equals, "test")
		c.Assert(messages[1].Message, check.Equals, "test 2")

		for i := 0; i < 2; i++ {
			val, ok := messages[i].Fields.Fields["remote_addr"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "127.0.0.1")

			val, ok = messages[i].Fields.Fields["request_method"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "GET")

			val, ok = messages[i].Fields.Fields["request_url"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "http://dogfood.tinyci.org")
		}
	}

	u, err = url.Parse("https://pensando.io")
	c.Assert(err, check.IsNil)

	req = &http.Request{
		Header:     http.Header{"X-Forwarded-For": []string{"1.2.3.4,5.6.7.8"}},
		RemoteAddr: "127.0.0.1:1234",
		Method:     "POST",
		URL:        u,
	}

	ls.journal.Reset()

	wr = l.WithRequest(req)

	wr.Info("test")
	wr.Infof("test %d", 2)
	wr.Debug("test")
	wr.Debugf("test %d", 2)
	wr.Error("test")
	wr.Errorf("test %d", 2)

	for _, messages := range ls.journal.Journal {
		c.Assert(len(messages), check.Equals, 2)
		c.Assert(messages[0].Message, check.Equals, "test")
		c.Assert(messages[1].Message, check.Equals, "test 2")

		for i := 0; i < 2; i++ {
			val, ok := messages[i].Fields.Fields["remote_addr"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "1.2.3.4")

			val, ok = messages[i].Fields.Fields["request_method"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "POST")

			val, ok = messages[i].Fields.Fields["request_url"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "https://pensando.io")
		}
	}
}

func (ls *logsvcSuite) TestUser(c *check.C) {
	l := log.New()

	user := &model.User{ID: 1, Username: "erikh"}
	wu := l.WithUser(user)

	wu.Info("test")
	wu.Infof("test %d", 2)
	wu.Debug("test")
	wu.Debugf("test %d", 2)
	wu.Error("test")
	wu.Errorf("test %d", 2)

	for _, messages := range ls.journal.Journal {
		c.Assert(len(messages), check.Equals, 2)
		c.Assert(messages[0].Message, check.Equals, "test")
		c.Assert(messages[1].Message, check.Equals, "test 2")

		for i := 0; i < 2; i++ {
			val, ok := messages[i].Fields.Fields["user_id"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "1")

			val, ok = messages[i].Fields.Fields["username"]
			c.Assert(ok, check.Equals, true)
			c.Assert(val.GetStringValue(), check.Equals, "erikh")
		}
	}
}
