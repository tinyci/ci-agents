package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
)

// GetFileList finds all the files in the tree for the given repository
func (c *Client) GetFileList(ctx context.Context, repoName, sha string) ([]string, error) {
	list, err := c.client.GetFileList(ctx, &repository.RepoSHAPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return nil, errors.New(err)
	}

	return list.List, nil
}

// GetSHA retrieves the SHA for the provided ref and repository.
func (c *Client) GetSHA(ctx context.Context, repoName, ref string) (string, error) {
	sha, err := c.client.GetSHA(ctx, &repository.RepoRefPair{RepoName: repoName, RefName: ref})
	if err != nil {
		return "", errors.New(err)
	}

	return sha.Name, nil
}

// GetRefs retreives many refs that have the corresponding SHA.
func (c *Client) GetRefs(ctx context.Context, repoName, sha string) ([]string, error) {
	refs, err := c.client.GetRefs(ctx, &repository.RepoSHAPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return nil, errors.New(err)
	}

	return refs.List, nil
}

// GetFile retrieves an entire file by way of the repoName, sha, and filename.
func (c *Client) GetFile(ctx context.Context, repoName, sha, filename string) ([]byte, error) {
	byts, err := c.client.GetFile(ctx, &repository.FileRequest{RepoName: repoName, Sha: sha, Filename: filename})
	if err != nil {
		return nil, errors.New(err)
	}

	return byts.Value, nil
}

// GetDiffFiles retrieves the files present in the diff between the base and the head.
func (c *Client) GetDiffFiles(ctx context.Context, repoName, base, head string) ([]string, error) {
	files, err := c.client.GetDiffFiles(ctx, &repository.FileDiffRequest{RepoName: repoName, Base: base, Head: head})
	if err != nil {
		return nil, errors.New(err)
	}

	return files.List, nil
}
