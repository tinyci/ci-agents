package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	check "github.com/erikh/check"
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
		w.WriteHeader(500)
	}
}

type utilsSuite struct{}

var _ = check.Suite(&utilsSuite{})

func TestUtils(t *testing.T) {
	check.TestingT(t)
}

func (us *utilsSuite) TestRetryRequest(c *check.C) {
	errChan := make(chan error, 1)
	bh := &backoffHandler{}
	s := httptest.NewServer(bh)
	defer s.Close()

	req, err := http.NewRequest("GET", s.URL, nil)
	c.Assert(err, check.IsNil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = RetryRequest(ctx, s.Client(), req, nil)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.Equals, context.DeadlineExceeded)
	c.Assert(bh.count, check.Equals, 4)

	bh.count = 0

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = RetryRequest(ctx, s.Client(), req, errChan)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.Equals, context.DeadlineExceeded)
	c.Assert(bh.count, check.Equals, 4)
	c.Assert(<-errChan, check.Equals, context.DeadlineExceeded)

	bh.count = 0

	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	_, err = RetryRequest(ctx, s.Client(), req, nil)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.Equals, context.Canceled)
	c.Assert(bh.count, check.Equals, 3)

	bh.count = 0
	bh.succeedNext = true

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = RetryRequest(ctx, s.Client(), req, errChan)
	c.Assert(err, check.IsNil)
	c.Assert(bh.count, check.Equals, 1)
	c.Assert(<-errChan, check.IsNil) // empty response for async calls

	bh.count = 0
	go func() {
		time.Sleep(300 * time.Millisecond)
		bh.mutex.Lock()
		bh.succeedNext = true
		bh.mutex.Unlock()
	}()

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err = RetryRequest(ctx, s.Client(), req, errChan)
	c.Assert(err, check.IsNil)
	c.Assert(bh.count, check.Equals, 1)
	c.Assert(<-errChan, check.IsNil) // empty response for async calls
}
