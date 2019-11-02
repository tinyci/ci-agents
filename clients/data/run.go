package data

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// RunCount returns the count of all items that match the repoName and sha.
func (c *Client) RunCount(ctx context.Context, repoName, sha string) (int64, error) {
	count, err := c.client.RunCount(ctx, &data.RefPair{RepoName: repoName, Sha: sha}, grpc.WaitForReady(true))
	if err != nil {
		return 0, errors.New(err)
	}

	return count.Count, nil
}

// ListRuns lists runs by repository name and sha
func (c *Client) ListRuns(ctx context.Context, repoName, sha string, page, perPage int64) ([]*model.Run, error) {
	list, err := c.client.RunList(ctx, &data.RunListRequest{Repository: repoName, Sha: sha, Page: page, PerPage: perPage}, grpc.WaitForReady(true))
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
func (c *Client) GetRun(ctx context.Context, id int64) (*model.Run, error) {
	run, err := c.client.GetRun(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewRunFromProto(run)
}

// GetRunUI retrieves a run by id.
func (c *Client) GetRunUI(ctx context.Context, id int64) (*model.Run, error) {
	run, err := c.client.GetRunUI(ctx, &types.IntID{ID: id}, grpc.WaitForReady(true))
	if err != nil {
		return nil, errors.New(err)
	}

	return model.NewRunFromProto(run)
}
