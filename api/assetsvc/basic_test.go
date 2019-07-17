package assetsvc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	check "github.com/erikh/check"
)

func (as *assetsvcSuite) TestReadValidation(c *check.C) {
	c.Assert(as.assetClient.Read(context.Background(), 1, ioutil.Discard), check.NotNil)
}

func (as *assetsvcSuite) TestLogUpDown(c *check.C) {
	f, err := os.Open("/proc/cpuinfo")
	c.Assert(err, check.IsNil)
	origBuf := bytes.NewBuffer(nil)
	r := io.TeeReader(f, origBuf)
	c.Assert(as.assetClient.Write(context.Background(), 1, r), check.IsNil)
	f.Close()

	time.Sleep(100 * time.Millisecond)

	buf := bytes.NewBuffer(nil)
	c.Assert(as.assetClient.Read(context.Background(), 1, buf), check.IsNil)

	// the log finished message will be at the end, so just compare prefixes
	c.Assert(strings.HasPrefix(buf.String(), origBuf.String()), check.Equals, true)
}

func (as *assetsvcSuite) TestLogStream(c *check.C) {
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()

	go func() {
		c.Assert(as.assetClient.Write(context.Background(), 1, pr1), check.IsNil)
	}()

	go func() {
		c.Assert(as.assetClient.Write(context.Background(), 2, pr2), check.IsNil)
	}()

	for i := 0; i < 1000; i++ {
		fmt.Fprintln(pw1, i)
	}

	pw1.Close()
	time.Sleep(100 * time.Millisecond)
	c.Assert(as.assetClient.Read(context.Background(), 1, pw2), check.IsNil)

	pw2.Close()

	buf1 := bytes.NewBuffer(nil)
	buf2 := bytes.NewBuffer(nil)

	c.Assert(as.assetClient.Read(context.Background(), 1, buf1), check.IsNil)
	c.Assert(as.assetClient.Read(context.Background(), 2, buf2), check.IsNil)

	c.Assert(strings.HasPrefix(buf2.String(), buf2.String()), check.Equals, true)
}
