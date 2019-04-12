package jsonbuffer

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/erikh/check"
	"github.com/tinyci/ci-agents/errors"
)

type jsonbufferSuite struct{}

var _ = check.Suite(&jsonbufferSuite{})

func TestJSONBuffer(t *testing.T) {
	check.TestingT(t)
}

func (ws *jsonbufferSuite) TestBasic(c *check.C) {
	f, err := os.Open("testdata/chat1.json")
	c.Assert(err, check.IsNil)
	defer f.Close()

	w := NewWrapper(f)

	// first read should be a message
	msg, err := w.Recv()
	c.Assert(err, check.IsNil)
	c.Assert(string(msg), check.Equals, "hello, world!")

	// second should be an error
	_, err = w.Recv()
	c.Assert(err, check.ErrorMatches, "^this is a test error$")

	buf := bytes.NewBuffer(nil)
	w2 := NewWrapper(buf)
	c.Assert(w2.Send("hello, world!"), check.IsNil)
	c.Assert(w2.SendError(errors.New("this is a test error")), check.IsNil)

	_, err = f.Seek(0, 0)
	c.Assert(err, check.IsNil)

	content, err := ioutil.ReadAll(f)
	c.Assert(err, check.IsNil)

	c.Assert(string(content), check.Equals, buf.String())
}

func (ws *jsonbufferSuite) TestStream(c *check.C) {
	buf := bytes.NewBuffer(nil)

	w := NewWrapper(buf)
	pr, pw := io.Pipe()

	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		io.Copy(w, pr)
	}()

	f, err := os.Open("testdata/chat1.json")
	c.Assert(err, check.IsNil)
	defer f.Close()

	_, err = io.Copy(pw, f)
	pw.Close()
	<-doneChan
	c.Assert(err, check.IsNil)

	f.Seek(0, 0)
	content, err := ioutil.ReadAll(f)
	c.Assert(err, check.IsNil)

	content2, err := ioutil.ReadAll(w)
	c.Assert(err, check.IsNil)
	c.Assert(content, check.DeepEquals, content2)
}

func (ws *jsonbufferSuite) TestStream2(c *check.C) {
	buf := bytes.NewBuffer(nil)
	w := NewWrapper(buf)
	pr, pw := io.Pipe()

	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		io.Copy(w, pr)
	}()

	byt := bytes.Repeat([]byte{0, 1}, 1024*512)
	buf2 := bytes.NewBuffer(byt)

	n, err := io.Copy(pw, buf2)
	c.Assert(err, check.IsNil)
	c.Assert(n, check.Equals, int64(1024*1024))
	pw.Close()
	<-doneChan

	content, err := ioutil.ReadAll(w)
	c.Assert(err, check.IsNil)

	c.Assert(len(content), check.Equals, 1024*1024)
	c.Assert(content, check.DeepEquals, byt)
}
