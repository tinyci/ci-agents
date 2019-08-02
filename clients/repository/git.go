package repository

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
)

// GetFileList finds all the files in the tree for the given repository
func (c *Client) GetFileList(repoName, sha string) ([]string, *errors.Error) {
	list, err := c.client.GetFileList(context.Background(), &repository.RepoSHAPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return nil, errors.New(err)
	}

	return list.List, nil
}

// GetSHA retrieves the SHA for the provided ref and repository.
func (c *Client) GetSHA(repoName, ref string) (string, *errors.Error) {
	sha, err := c.client.GetSHA(context.Background(), &repository.RepoRefPair{RepoName: repoName, RefName: ref})
	if err != nil {
		return "", errors.New(err)
	}

	return sha.Name, nil
}

// GetRefs retreives many refs that have the corresponding SHA.
func (c *Client) GetRefs(repoName, sha string) ([]string, *errors.Error) {
	refs, err := c.client.GetRefs(context.Background(), &repository.RepoSHAPair{RepoName: repoName, Sha: sha})
	if err != nil {
		return nil, errors.New(err)
	}

	return refs.List, nil
}

// GetFile retrieves an entire file by way of the repoName, sha, and filename.
func (c *Client) GetFile(repoName, sha, filename string) ([]byte, *errors.Error) {
	byts, err := c.client.GetFile(context.Background(), &repository.FileRequest{RepoName: repoName, Sha: sha, Filename: filename})
	if err != nil {
		return nil, errors.New(err)
	}

	return byts.Value, nil
}

// GetDiffFiles retrieves the files present in the diff between the base and the head.
func (c *Client) GetDiffFiles(repoName, base, head string) ([]string, *errors.Error) {
	files, err := c.client.GetDiffFiles(context.Background(), &repository.FileDiffRequest{RepoName: repoName, Base: base, Head: head})
	if err != nil {
		return nil, errors.New(err)
	}

	return files.List, nil
}
