package retry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/erikh/check"
)

type backoffHandler struct {
	mutex       sync.Mutex
	count       int
	succeedNext bool
}

func (bh *backoffHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bh.mutex.Lock()
	defer bh.mutex.Unlock()
	bh.count++
	if !bh.succeedNext {
		panic("intentionally failing to respond")
	}
}

func (us *retrySuite) TestHTTP(c *check.C) {
	bh := &backoffHandler{}
	s := httptest.NewServer(bh)
	defer s.Close()

	req, reqErr := http.NewRequest("GET", "http://localhost:1234", nil)
	c.Assert(reqErr, check.IsNil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := NewHTTP(s.Client()).Do(ctx, req)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.Equals, context.DeadlineExceeded)

	bh.count = 0

	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	_, err = NewHTTP(s.Client()).Do(ctx, req)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.Equals, context.Canceled)

	bh.count = 0
	bh.succeedNext = true

	req, reqErr = http.NewRequest("GET", s.URL, nil)
	c.Assert(reqErr, check.IsNil)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = NewHTTP(s.Client()).Do(ctx, req)
	c.Assert(err, check.IsNil)
	c.Assert(bh.count, check.Equals, 1)

	bh.count = 0
	bh.succeedNext = false
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	_, err = NewHTTPWithInterval(s.Client(), 100*time.Millisecond).Do(ctx, req)
	c.Assert(err, check.Equals, context.Canceled)
	c.Assert(bh.count, check.Equals, 11)

	bh.count = 0
	bh.succeedNext = false
	go func() {
		time.Sleep(300 * time.Millisecond)
		bh.mutex.Lock()
		bh.succeedNext = true
		bh.mutex.Unlock()
	}()

	_, err = NewHTTP(s.Client()).Do(context.Background(), req)
	c.Assert(err, check.IsNil)
	c.Assert(bh.count, check.Equals, 2)
}
