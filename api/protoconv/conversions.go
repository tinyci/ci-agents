package protoconv

import (
	"context"
	"fmt"
	"time"

	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/db/protoconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ErrConversionInvalidType is returned when you pass the wrong type.
var ErrConversionInvalidType = protoconv.ErrConversionInvalidType

func timeToPtr(t *timestamppb.Timestamp) *time.Time {
	var ret *time.Time
	if t.IsValid() {
		t2 := t.AsTime()
		ret = &t2
	}

	return ret
}

func ueFromProto(ctx context.Context, i interface{}) (interface{}, error) {
	ue, ok := i.(*types.UserError)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	return &uisvc.UserError{
		Id:    &ue.Id,
		Error: &ue.Error,
	}, nil
}

func ueToProto(ctx context.Context, i interface{}) (interface{}, error) {
	return nil, nil
}

func subFromProto(ctx context.Context, i interface{}) (interface{}, error) {
	s, ok := i.(*types.Submission)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	baseRef, err := refFromProto(ctx, s.BaseRef)
	if err != nil {
		return nil, err
	}

	headRef, err := refFromProto(ctx, s.HeadRef)
	if err != nil {
		return nil, err
	}

	var status *bool
	if s.StatusSet {
		status = &s.Status
	}

	user, err := userFromProto(ctx, s.User)
	if err != nil {
		return nil, err
	}

	return &uisvc.ModelSubmission{
		Id:         &s.Id,
		CreatedAt:  timeToPtr(s.CreatedAt),
		FinishedAt: timeToPtr(s.FinishedAt),
		StartedAt:  timeToPtr(s.StartedAt),

		BaseRef: baseRef.(*uisvc.Ref),
		HeadRef: headRef.(*uisvc.Ref),

		Canceled: &s.Canceled,
		Status:   status,

		RunsCount:  &s.RunsCount,
		TasksCount: &s.TasksCount,
		TicketId:   &s.TicketID,

		User: user.(*uisvc.User),
	}, nil
}

func subToProto(ctx context.Context, i interface{}) (interface{}, error) {
	_, ok := i.(*uisvc.ModelSubmission)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	// FIXME finish

	return nil, nil
}

func refFromProto(ctx context.Context, i interface{}) (interface{}, error) {
	r, ok := i.(*types.Ref)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	repo, err := repoFromProto(ctx, r.Repository)
	if err != nil {
		return nil, err
	}

	return &uisvc.Ref{
		RefName:    &r.RefName,
		Repository: repo.(*uisvc.Repository),
		Sha:        &r.Sha,
	}, nil
}

func refToProto(ctx context.Context, i interface{}) (interface{}, error) {
	return nil, nil
}

func runFromProto(ctx context.Context, i interface{}) (interface{}, error) {
	r, ok := i.(*types.Run)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	var createdAt, finishedAt, startedAt *time.Time
	if r.CreatedAt.IsValid() {
		t := r.CreatedAt.AsTime()
		createdAt = &t
	}

	if r.FinishedAt.IsValid() {
		t := r.FinishedAt.AsTime()
		finishedAt = &t
	}

	if r.StartedAt.IsValid() {
		t := r.StartedAt.AsTime()
		startedAt = &t
	}

	var status *bool

	if r.StatusSet {
		status = &r.Status
	}

	task, err := taskFromProto(ctx, r.Task)
	if err != nil {
		return nil, err
	}

	return &uisvc.Run{
		CreatedAt:  createdAt,
		FinishedAt: finishedAt,
		StartedAt:  startedAt,
		Id:         &r.Id,
		Name:       &r.Name,
		RanOn:      &r.RanOn,
		Status:     status,
		// Settings   *RunSettings `json:"settings,omitempty"`
		Task: task.(*uisvc.Task),
	}, nil
}

func runToProto(ctx context.Context, i interface{}) (interface{}, error) {
	return nil, nil
}

func taskFromProto(ctx context.Context, i interface{}) (interface{}, error) {
	t, ok := i.(*types.Task)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	var createdAt, finishedAt, startedAt *time.Time
	if t.CreatedAt.IsValid() {
		tm := t.CreatedAt.AsTime()
		createdAt = &tm
	}

	if t.FinishedAt.IsValid() {
		tm := t.FinishedAt.AsTime()
		finishedAt = &tm
	}

	if t.StartedAt.IsValid() {
		tm := t.StartedAt.AsTime()
		startedAt = &tm
	}

	var status *bool

	if t.StatusSet {
		status = &t.Status
	}

	sub, err := subFromProto(ctx, t.Submission)
	if err != nil {
		return nil, err
	}

	return &uisvc.Task{
		Canceled:   &t.Canceled,
		CreatedAt:  createdAt,
		FinishedAt: finishedAt,
		StartedAt:  startedAt,
		Id:         &t.Id,
		Path:       &t.Path,
		Status:     status,
		Submission: sub.(*uisvc.ModelSubmission),
		Runs:       &t.Runs,
		/*
			Settings   *TaskSettings    `json:"settings,omitempty"`
		*/
	}, nil
}

func taskToProto(ctx context.Context, i interface{}) (interface{}, error) {
	return nil, nil
}

func repoFromProto(ctx context.Context, i interface{}) (interface{}, error) {
	r, ok := i.(*types.Repository)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	return &uisvc.Repository{
		AutoCreated: &r.AutoCreated,
		Disabled:    &r.Disabled,
		Id:          &r.Id,
		Name:        &r.Name,
		Private:     &r.Private,
	}, nil
}

func repoToProto(ctx context.Context, i interface{}) (interface{}, error) {
	return nil, nil
}

func userFromProto(ctx context.Context, i interface{}) (interface{}, error) {
	u, ok := i.(*types.User)
	if !ok {
		return nil, fmt.Errorf("%T: %w", i, ErrConversionInvalidType)
	}

	return &uisvc.User{
		// FIXME we might need to expand this later. If we need to include
		// sensitive data, this function should appropriately account for the fact
		// that a lot of data that goes through it should not be exposed to the
		// user.
		Id:       &u.Id,
		Username: &u.Username,
	}, nil
}

func userToProto(ctx context.Context, i interface{}) (interface{}, error) {
	return nil, nil
}
