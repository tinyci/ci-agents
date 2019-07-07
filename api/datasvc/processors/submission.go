package processors

import (
	"context"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/model"
)

// PutSubmission registers a submission with the datasvc.
func (ds *DataServer) PutSubmission(ctx context.Context, sub *types.Submission) (*types.Submission, error) {
	s, err := model.NewSubmissionFromProto(sub)
	if err != nil {
		return nil, err
	}

	if err := ds.H.Model.Create(s).Error; err != nil {
		return nil, err
	}

	return s.ToProto(), nil
}
