package retry

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/tinyci/ci-agents/errors"
)

var defaultInterval = time.Second

// HTTP wraps a *http.Client with retry data.
type HTTP struct {
	interval time.Duration

	*http.Client
}

// NewHTTP constructs a new client that retries every defaultInterval.
func NewHTTP(client *http.Client) *HTTP {
	return &HTTP{Client: client, interval: defaultInterval}
}

// NewHTTPWithInterval is just like new only you can program the interval yourself.
func NewHTTPWithInterval(client *http.Client, interval time.Duration) *HTTP {
	return &HTTP{Client: client, interval: interval}
}

// Do executes the request; under certain response scenarios it will retry if
// there are errors.
func (c *HTTP) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		resp, err := c.Client.Do(req)
		if err != nil {
			switch err {
			case context.Canceled, context.DeadlineExceeded:
				return nil, err
			case io.EOF:
				goto sleep
			default:
				switch err := err.(type) {
				case *net.OpError, *url.Error:
					goto sleep
				default:
					return nil, errors.New(err)
				}
			}
		}

		return resp, nil

	sleep:
		time.Sleep(c.interval)
	}
}
