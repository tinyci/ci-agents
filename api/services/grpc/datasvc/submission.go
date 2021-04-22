package datasvc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetSubmission returns a submission from the provided ID
func (ds *DataServer) GetSubmission(ctx context.Context, id *types.IntID) (*types.Submission, error) {
	s, err := ds.H.Model.GetSubmissionByID(id.ID)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return s.ToProto(), nil
}

// GetSubmissionRuns returns the runs associated with the provided submission.
func (ds *DataServer) GetSubmissionRuns(ctx context.Context, sub *data.SubmissionQuery) (*types.RunList, error) {
	protoSub, err := model.NewSubmissionFromProto(sub.Submission)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	runs, err := ds.H.Model.RunsForSubmission(protoSub, sub.Page, sub.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	rl := &types.RunList{}

	for _, run := range runs {
		rl.List = append(rl.List, run.ToProto())
	}

	return rl, nil
}

// GetSubmissionTasks returns the tasks associated with the provided submission.
func (ds *DataServer) GetSubmissionTasks(ctx context.Context, sub *data.SubmissionQuery) (*types.TaskList, error) {
	protoSub, err := model.NewSubmissionFromProto(sub.Submission)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	tasks, err := ds.H.Model.TasksForSubmission(protoSub, sub.Page, sub.PerPage)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	tl := &types.TaskList{}

	for _, task := range tasks {
		tl.Tasks = append(tl.Tasks, task.ToProto())
	}

	return tl, nil
}

// PutSubmission registers a submission with the datasvc.
func (ds *DataServer) PutSubmission(ctx context.Context, sub *types.Submission) (*types.Submission, error) {
	s, err := model.NewSubmissionFromProto(sub)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	if err := ds.H.Model.Create(s).Error; err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return s.ToProto(), nil
}

// ListSubmissions lists the submissions with optional repository and ref filtering.
func (ds *DataServer) ListSubmissions(ctx context.Context, req *data.RepositoryFilterRequestWithPagination) (*types.SubmissionList, error) {
	list, err := ds.H.Model.SubmissionList(req.Page, req.PerPage, req.Repository, req.Sha)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	newList := &types.SubmissionList{}

	for _, sub := range list {
		newList.Submissions = append(newList.Submissions, sub.ToProto())
	}

	return newList, nil
}

// CountSubmissions returns a count of all submissions that match the filter.
func (ds *DataServer) CountSubmissions(ctx context.Context, req *data.RepositoryFilterRequest) (*data.Count, error) {
	count, err := ds.H.Model.SubmissionCount(req.Repository, req.Sha)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return &data.Count{Count: count}, nil
}

// CancelSubmission cancels a submission by ID.
func (ds *DataServer) CancelSubmission(ctx context.Context, id *types.IntID) (*empty.Empty, error) {
	empty := &empty.Empty{}

	if err := ds.H.Model.CancelSubmissionByID(id.ID, ds.H.URL, nil); err != nil {
		return empty, status.Errorf(codes.FailedPrecondition, "%v", err)
	}

	return empty, nil
}
