package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"google.golang.org/grpc"
)

// PutRef adds a ref to the database.
func (c *Client) PutRef(ctx context.Context, ref *types.Ref) (int64, error) {
	id, err := c.client.PutRef(ctx, ref, grpc.WaitForReady(true))
	if err != nil {
		return 0, err
	}

	return id.Id, nil
}

// CancelRefByName cancels all jobs for a ref by name
func (c *Client) CancelRefByName(ctx context.Context, repoID int64, ref string) error {
	_, err := c.client.CancelRefByName(ctx, &data.RepoRef{Repository: repoID, RefName: ref}, grpc.WaitForReady(true))
	return err
}

// GetRefByNameAndSHA retrieves a ref by it's repo name and SHA
func (c *Client) GetRefByNameAndSHA(ctx context.Context, repoName, sha string) (*types.Ref, error) {
	return c.client.GetRefByNameAndSHA(ctx, &data.RefPair{RepoName: repoName, Sha: sha}, grpc.WaitForReady(true))
}
