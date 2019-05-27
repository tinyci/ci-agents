package asset

import (
	"context"
	"io"

	transport "github.com/erikh/go-transport"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/asset"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"google.golang.org/grpc"
)

// Client is a handle into the asset client.
type Client struct {
	ac asset.AssetClient
}

// NewClient creates a new *Client for use.
func NewClient(cert *transport.Cert, addr string) (*Client, *errors.Error) {
	t, err := transport.GRPCDial(cert, addr)
	if err != nil {
		return nil, errors.New(err)
	}
	return &Client{ac: asset.NewAssetClient(t)}, nil
}

// Write writes a log at id with the supplied reader providing the content.
func (c *Client) Write(id int64, f io.Reader) *errors.Error {
	s, err := c.ac.PutLog(context.Background(), grpc.WaitForReady(true))
	if err != nil {
		return errors.New(err)
	}

	buf := make([]byte, 64)

	for {
		var done bool
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return errors.New(err)
		} else if err == io.EOF {
			done = true
		}

		ls := &asset.LogSend{
			ID:    id,
			Chunk: buf[:n],
		}

		if err := s.Send(ls); err != nil && err != io.EOF {
			return errors.New(err)
		} else if err == io.EOF {
			done = true
		}

		if done {
			if err := s.CloseSend(); err != nil {
				return errors.New(err)
			}

			return nil
		}
	}
}

func (c *Client) Read(id int64, w io.Writer) *errors.Error {
	as, err := c.ac.GetLog(context.Background(), &types.IntID{ID: id}, grpc.WaitForReady(false))
	if err != nil {
		return errors.New(err)
	}

	md, err := as.Header()
	if err != nil {
		return errors.New(err)
	}

	errs := md.Get("errors")
	if len(errs) > 0 {
		return errors.New(errs[0])
	}

	for {
		chunk, err := as.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return errors.New(err)
		}

		if _, err := w.Write(chunk.Chunk); err != nil {
			return errors.New(err)
		}
	}
}
