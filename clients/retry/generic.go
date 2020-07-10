package retry

import (
	"io"
	"net"
	"net/url"
	"time"

	"github.com/tinyci/ci-agents/errors"
)

// Generic is a generic retry mechanism that can retry any block of code.
type Generic struct {
	interval time.Duration
}

// NewGeneric implements a generic repeat/retry mechanism.
func NewGeneric() *Generic {
	return &Generic{interval: defaultInterval}
}

// NewGenericWithInterval provides a generic with a selected interval instead of the default.
func NewGenericWithInterval(interval time.Duration) *Generic {
	return &Generic{interval: interval}
}

// Do provides a mechanism to retry the function until it succeeds. Any error
// returned will continue the loop.
func (g *Generic) Do(f func() error) *errors.Error {
	for {
		if err := f(); err != nil {
			switch err {
			case io.EOF:
				goto sleep
			default:
				switch err.(type) {
				case *net.OpError, *url.Error:
					goto sleep
				default:
					return errors.New(err)
				}
			}
		sleep:
			time.Sleep(g.interval)
			continue
		}

		return nil
	}
}
