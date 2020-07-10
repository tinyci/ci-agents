package utils

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/tinyci/ci-agents/errors"
)

var (
	initialWait = 100 * time.Millisecond
	maxWait     = time.Minute
)

// RetryRequest attempts a request on the client repeatedly until it succeeds
// or the context is canceled. It will retry for time.Duration and back off
// slowly over several polls. If a channel is provided, it will return to the
// channel when the operation succeeds with any error if relevant, but only
// calling it sync will return the response.
func RetryRequest(ctx context.Context, client *http.Client, req *http.Request, returnChan chan error) (resp *http.Response, retErr error) {
	defer func() {
		if returnChan != nil {
			returnChan <- retErr
		}
	}()

	wait := initialWait

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if err := CheckStatus(resp); err != nil {
			// FIXME log this error
			time.Sleep(wait)
			if wait > maxWait {
				wait = maxWait
			} else {
				wait *= 2
			}
			continue
		}

		return resp, nil
	}
}

// Get a *url.URL
func Get(url *url.URL, client *http.Client) (*http.Response, error) {
	return client.Get(url.String())
}

// ParseJSON simplification of json.NewDecoder.
func ParseJSON(reader io.Reader, dst interface{}) error {
	return json.NewDecoder(reader).Decode(dst)
}

// CheckStatus returns error if the status is not 200 OK.
func CheckStatus(resp *http.Response) error {
	if resp.StatusCode != 200 {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(content))
	}

	return nil
}
