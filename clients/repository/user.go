package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
)

// MyLogin returns the username calling out to the API with its key. Can either
// be seeded by OAuth or Personal Token.
func (c *Client) MyLogin(ctx context.Context, token string) (string, error) {
	res, err := c.client.MyLogin(ctx, &repository.String{Name: token})
	if err != nil {
		return "", err
	}

	return res.Name, nil
}
