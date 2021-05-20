package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetSubmission returns a submission from the provided ID
func (ds *DataServer) GetSubmission(ctx context.Context, id *types.IntID) (*types.Submission, error) {
	s, err := ds.H.Model.GetSubmissionByID(ctx, id.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	ret, err := ds.C.ToProto(ctx, s)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	return ret.(*types.Submission), nil
}

// GetSubmissionRuns returns the runs associated with the provided submission.
func (ds *DataServer) GetSubmissionRuns(ctx context.Context, sub *data.SubmissionQuery) (*types.RunList, error) {
	ps, err := ds.C.FromProto(ctx, sub.Submission)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	protoSub := ps.(*models.Submission)

	runs, err := ds.H.Model.RunsForSubmission(ctx, protoSub.ID, sub.Page, sub.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	rl := &types.RunList{}

	for _, run := range runs {
		r, err := ds.C.ToProto(ctx, run)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		rl.List = append(rl.List, r.(*types.Run))
	}

	return rl, nil
}

// GetSubmissionTasks returns the tasks associated with the provided submission.
func (ds *DataServer) GetSubmissionTasks(ctx context.Context, sub *data.SubmissionQuery) (*types.TaskList, error) {
	tasks, err := ds.H.Model.TasksForSubmission(ctx, sub.Submission.Id, sub.Page, sub.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	tl := &types.TaskList{}

	for _, task := range tasks {
		t, err := ds.C.ToProto(ctx, task)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		tl.Tasks = append(tl.Tasks, t.(*types.Task))
	}

	return tl, nil
}

// PutSubmission registers a submission with the datasvc.
func (ds *DataServer) PutSubmission(ctx context.Context, sub *types.Submission) (*types.Submission, error) {
	s, err := ds.C.FromProto(ctx, sub)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.PutSubmission(ctx, s.(*models.Submission)); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	ret, err := ds.C.ToProto(ctx, s)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return ret.(*types.Submission), nil
}

// ListSubmissions lists the submissions with optional repository and ref filtering.
func (ds *DataServer) ListSubmissions(ctx context.Context, req *data.RepositoryFilterRequestWithPagination) (*types.SubmissionList, error) {
	list, err := ds.H.Model.SubmissionList(ctx, req.Page, req.PerPage, req.Repository, req.Sha)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	newList := &types.SubmissionList{Submissions: []*types.Submission{}}

	for _, sub := range list {
		s, err := ds.C.ToProto(ctx, sub)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		newList.Submissions = append(newList.Submissions, s.(*types.Submission))
	}

	return newList, nil
}

// CountSubmissions returns a count of all submissions that match the filter.
func (ds *DataServer) CountSubmissions(ctx context.Context, req *data.RepositoryFilterRequest) (*data.Count, error) {
	count, err := ds.H.Model.SubmissionCount(ctx, req.Repository, req.Sha)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &data.Count{Count: count}, nil
}

// CancelSubmission cancels a submission by ID.
func (ds *DataServer) CancelSubmission(ctx context.Context, id *types.IntID) (*empty.Empty, error) {
	empty := &empty.Empty{}

	if err := ds.H.Model.CancelSubmissionByID(ctx, id.ID); err != nil {
		return empty, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return empty, nil
}
