package data

import (
	"context"

	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/data"
	"github.com/tinyci/ci-agents/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// RunCount returns the count of all items that match the repoName and sha.
func (c *Client) RunCount(repoName, sha string) (int64, *errors.Error) {
	count, err := c.client.RunCount(context.Background(), &data.RefPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return 0, errors.New(err)
	}

	return count.Count, nil
}

// ListRuns lists runs by repository name and sha
func (c *Client) ListRuns(repoName, sha string, page, perPage int64) ([]*model.Run, *errors.Error) {
	list, err := c.client.RunList(context.Background(), &data.RunListRequest{Repository: repoName, Sha: sha, Page: page, PerPage: perPage})
	if err != nil {
		return nil, errors.New(err)
	}

	runList := []*model.Run{}

	for _, run := range list.List {
		pr, err := model.NewRunFromProto(run)
		if err != nil {
			return nil, err
		}

		runList = append(runList, pr)
	}

	return runList, nil
}

// GetRun retrieves a run by id.
func (c *Client) GetRun(id int64) (*model.Run, *errors.Error) {
	run, err := c.client.GetRun(context.Background(), &types.IntID{ID: id})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewRunFromProto(run)
}

// GetRunUI retrieves a run by id.
func (c *Client) GetRunUI(id int64) (*model.Run, *errors.Error) {
	run, err := c.client.GetRunUI(context.Background(), &types.IntID{ID: id})
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewRunFromProto(run)
}
