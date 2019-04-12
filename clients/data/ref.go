package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/model"
)

// PutRef adds a ref to the database.
func (c *Client) PutRef(ref *model.Ref) (int64, *errors.Error) {
	id, err := c.client.PutRef(context.Background(), ref.ToProto())
	if err != nil {
		return 0, errors.New(err)
	}

	return id.Id, nil
}

// CancelRefByName cancels all jobs for a ref by name
func (c *Client) CancelRefByName(repoID int64, ref string) *errors.Error {
	_, err := c.client.CancelRefByName(context.Background(), &data.RepoRef{Repository: repoID, RefName: ref})
	return errors.New(err)
}

// GetRefByNameAndSHA retrieves a ref by it's repo name and SHA
func (c *Client) GetRefByNameAndSHA(repoName, sha string) (*model.Ref, *errors.Error) {
	ref, err := c.client.GetRefByNameAndSHA(context.Background(), &data.RefPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewRefFromProto(ref)
}
