package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
)

// MyLogin returns the username calling out to the API with its key. Can either
// be seeded by OAuth or Personal Token.
func (c *Client) MyLogin(ctx context.Context, token string) (string, *errors.Error) {
	res, err := c.client.MyLogin(ctx, &repository.String{Name: token})
	if err != nil {
		return "", errors.New(err)
	}

	return res.Name, nil
}
