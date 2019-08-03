package github

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/repository"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/grpc/codes"
)

// CommentError is for commenting on PRs when there is no better means of bubbling up an error.
func (rs *RepositoryServer) CommentError(ctx context.Context, cer *repository.CommentErrorRequest) (*empty.Empty, error) {
	owner, repo, retErr := utils.OwnerRepo(cer.RepoName)
	if retErr != nil {
		return nil, retErr.ToGRPC(codes.FailedPrecondition)
	}
	gh, err := rs.getClientForRepo(ctx, cer.RepoName)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	_, _, eerr := gh.Issues.CreateComment(ctx, owner, repo, int(cer.PrID), &github.IssueComment{
		Body: github.String(cer.Error),
	})

	if eerr != nil {
		return nil, eerr
	}

	return &empty.Empty{}, nil
}

func (rs *RepositoryServer) getStatusInfo(ctx context.Context, repoName string) (*github.Client, string, string, *errors.Error) {
	owner, repo, err := utils.OwnerRepo(repoName)
	if err != nil {
		return nil, "", "", err
	}

	gh, err := rs.getClientForRepo(ctx, repoName)
	if err != nil {
		return nil, "", "", err
	}

	return gh, owner, repo, nil
}

// PendingStatus updates the status for the sha for the given repo on github.
func (rs *RepositoryServer) PendingStatus(ctx context.Context, sr *repository.StatusRequest) (*empty.Empty, error) {
	gh, owner, repo, err := rs.getStatusInfo(ctx, sr.RepoName)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	_, _, eErr := gh.Repositories.CreateStatus(ctx, owner, repo, sr.Sha, &github.RepoStatus{
		TargetURL:   github.String(sr.Url),
		State:       github.String("pending"),
		Description: github.String("The run will be starting soon."),
		Context:     github.String(sr.RunName),
	})
	if eErr != nil {
		return nil, errors.New(eErr).Wrapf("creating status for %v/%v", owner, repo).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

// StartedStatus updates the status for the sha for the given repo on github.
func (rs *RepositoryServer) StartedStatus(ctx context.Context, sr *repository.StatusRequest) (*empty.Empty, error) {
	gh, owner, repo, err := rs.getStatusInfo(ctx, sr.RepoName)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	_, _, eErr := gh.Repositories.CreateStatus(ctx, owner, repo, sr.Sha, &github.RepoStatus{
		TargetURL:   github.String(sr.Url),
		State:       github.String("pending"),
		Description: github.String("The run has started!"),
		Context:     github.String(sr.RunName),
	})
	if eErr != nil {
		return nil, errors.New(eErr).Wrapf("creating status for %v/%v", owner, repo).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

func capStatus(str string) *string {
	if len(str) > 140 {
		return github.String(str[:140])
	}

	return github.String(str)
}

// ErrorStatus updates the status for the sha for the given repo on github.
func (rs *RepositoryServer) ErrorStatus(ctx context.Context, esr *repository.ErrorStatusRequest) (*empty.Empty, error) {
	gh, owner, repo, err := rs.getStatusInfo(ctx, esr.RepoName)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	_, _, eErr := gh.Repositories.CreateStatus(ctx, owner, repo, esr.Sha, &github.RepoStatus{
		TargetURL: github.String(esr.Url),
		State:     github.String("error"),
		// github statuses cap at 140c
		Description: capStatus(errors.New(esr.Error).Wrap("The run encountered an error").Error()),
		Context:     github.String(esr.RunName),
	})
	if eErr != nil {
		return nil, errors.New(eErr).Wrapf("creating status for %v/%v", owner, repo).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

// FinishedStatus updates the status for the sha for the given repo on github.
func (rs *RepositoryServer) FinishedStatus(ctx context.Context, fsr *repository.FinishedStatusRequest) (*empty.Empty, error) {
	gh, owner, repo, err := rs.getStatusInfo(ctx, fsr.RepoName)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	statusString := "failure"
	if fsr.Status {
		statusString = "success"
	}

	_, _, eErr := gh.Repositories.CreateStatus(ctx, owner, repo, fsr.Sha, &github.RepoStatus{
		TargetURL: github.String(fsr.Url),
		State:     github.String(statusString),
		// github statuses cap at 140c
		Description: capStatus(fmt.Sprintf("The run finished: %s! %v", statusString, fsr.Msg)),
		Context:     github.String(fsr.RunName),
	})
	if eErr != nil {
		return nil, errors.New(eErr).Wrapf("creating status for %v/%v", owner, repo).ToGRPC(codes.FailedPrecondition)
	}

	return &empty.Empty{}, nil
}

// ClearStates removes all status reports from a SHA in an attempt to restart
// the process.
func (rs *RepositoryServer) ClearStates(ctx context.Context, rsp *repository.RepoSHAPair) (*empty.Empty, error) {
	gh, owner, repo, err := rs.getStatusInfo(ctx, rsp.RepoName)
	if err != nil {
		return nil, err.ToGRPC(codes.FailedPrecondition)
	}

	statuses := []*github.RepoStatus{}
	contexts := map[string]struct{}{}

	var i int
	for {
		states, _, err := gh.Repositories.ListStatuses(ctx, owner, repo, rsp.Sha, &github.ListOptions{Page: i, PerPage: 200})
		if err != nil {
			return nil, errors.New(err).ToGRPC(codes.FailedPrecondition)
		}

		if len(states) == 0 {
			break
		}

		statuses = append(statuses, states...)
		i++
	}

	for _, status := range statuses {
		if _, ok := contexts[status.GetContext()]; ok {
			continue
		}

		contexts[status.GetContext()] = struct{}{}

		// XXX the context MUST be preserved for this to be overwritten. Do not
		// change it here.
		status.State = github.String("error")
		status.Description = github.String("The run that this test was a part of has been overridden by a new run. Pushing a new change will remove this error.")
		_, _, err := gh.Repositories.CreateStatus(ctx, owner, repo, rsp.Sha, status)
		if err != nil {
			return nil, errors.New(err).ToGRPC(codes.FailedPrecondition)
		}
	}

	return &empty.Empty{}, nil
}
