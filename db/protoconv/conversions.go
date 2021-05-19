package protoconv

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tinyci/ci-agents/ci-gen/grpc/services/data"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/db/models"
	topTypes "github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

// ErrConversionInvalidType is returned when you pass the wrong type.
var ErrConversionInvalidType = errors.New("could convert: invalid type")

func userErrorFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	e, ok := i.(*types.UserError)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	return &models.UserError{
		ID:     e.Id,
		UserID: e.UserID,
		Error:  e.Error,
	}, nil
}

func userErrorToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	e, ok := i.(*models.UserError)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	return &types.UserError{
		Id:     e.ID,
		UserID: e.UserID,
		Error:  e.Error,
	}, nil
}

func queueItemFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	qi, ok := i.(*types.QueueItem)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	return &models.QueueItem{
		ID:        qi.Id,
		RunID:     qi.Run.Id,
		Running:   qi.Running,
		RunningOn: null.StringFrom(qi.RunningOn),
		StartedAt: null.TimeFromPtr(timeFromPB(qi.StartedAt)),
		QueueName: qi.QueueName,
	}, nil
}

func queueItemToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	qi, ok := i.(*models.QueueItem)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	run, err := qi.Run().One(ctx, db)
	if err != nil {
		return nil, err
	}

	r, err := runToProto(ctx, db, run)
	if err != nil {
		return nil, err
	}

	return &types.QueueItem{
		Id:        qi.ID,
		Running:   qi.Running,
		RunningOn: qi.RunningOn.String,
		StartedAt: timeToPB(qi.StartedAt),
		QueueName: qi.QueueName,
		Run:       r.(*types.Run),
	}, nil
}

func runFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	run, ok := i.(*types.Run)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	tmp, err := taskFromProto(ctx, db, run.Task)
	if err != nil {
		return nil, err
	}

	task := tmp.(*models.Task)

	var ranOn *string
	if run.RanOnSet {
		ranOn = &run.RanOn
	}

	content, err := json.Marshal(run.Settings)
	if err != nil {
		return nil, err
	}

	return &models.Run{
		ID:          run.Id,
		Name:        run.Name,
		CreatedAt:   run.CreatedAt.AsTime(),
		StartedAt:   null.TimeFromPtr(timeFromPB(run.StartedAt)),
		FinishedAt:  null.TimeFromPtr(timeFromPB(run.FinishedAt)),
		Status:      makeStatus(run.Status, run.StatusSet),
		TaskID:      task.ID,
		RunSettings: content,
		RanOn:       null.StringFromPtr(ranOn),
	}, nil
}

func runToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	r, ok := i.(*models.Run)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	var status, set bool
	if r.Status.Valid {
		status = r.Status.Bool
		set = true
	}

	var (
		ranOn    string
		ranOnSet bool
	)

	if r.RanOn.Valid {
		ranOn = r.RanOn.String
		ranOnSet = true
	}

	task, err := r.Task().One(ctx, db)
	if err != nil {
		return nil, err
	}

	taskProto, err := taskToProto(ctx, db, task)
	if err != nil {
		return nil, err
	}

	rs := &topTypes.RunSettings{}
	if err := json.Unmarshal(r.RunSettings, rs); err != nil {
		return nil, err
	}

	return &types.Run{
		Id:         r.ID,
		Name:       r.Name,
		CreatedAt:  timestamppb.New(r.CreatedAt),
		StartedAt:  timeToPB(r.StartedAt),
		FinishedAt: timeToPB(r.FinishedAt),
		Status:     status,
		StatusSet:  set,
		Task:       taskProto.(*types.Task),
		Settings:   rs.ToProto(),
		RanOn:      ranOn,
		RanOnSet:   ranOnSet,
	}, nil
}

func taskFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	task, ok := i.(*types.Task)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	var sub *models.Submission
	if task.Submission != nil {
		tmp, err := subFromProto(ctx, db, task.Submission)
		if err != nil {
			return nil, err
		}

		sub = tmp.(*models.Submission)
	}

	if err := sub.Upsert(ctx, db, true, nil, boil.Infer(), boil.Infer()); err != nil {
		return nil, err
	}

	content, err := json.Marshal(task.Settings)
	if err != nil {
		return nil, err
	}

	return &models.Task{
		ID:           task.Id,
		Path:         task.Path,
		Canceled:     task.Canceled,
		FinishedAt:   null.TimeFromPtr(timeFromPB(task.FinishedAt)),
		StartedAt:    null.TimeFromPtr(timeFromPB(task.StartedAt)),
		CreatedAt:    task.CreatedAt.AsTime(),
		Status:       makeStatus(task.Status, task.StatusSet),
		TaskSettings: content,
		SubmissionID: sub.ID,
	}, nil
}

func taskToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	t, ok := i.(*models.Task)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	tmp, err := t.Submission().One(ctx, db)
	if err != nil {
		return nil, err
	}

	sub, err := subToProto(ctx, db, tmp)
	if err != nil {
		return nil, err
	}

	ts := &topTypes.TaskSettings{}
	if err := t.TaskSettings.Unmarshal(ts); err != nil {
		return nil, err
	}

	runCount, err := t.Runs().Count(ctx, db)
	if err != nil {
		return nil, err
	}

	return &types.Task{
		Id:         t.ID,
		Path:       t.Path,
		Canceled:   t.Canceled,
		FinishedAt: timeToPB(t.FinishedAt),
		StartedAt:  timeToPB(t.StartedAt),
		CreatedAt:  timestamppb.New(t.CreatedAt),
		Status:     t.Status.Bool,
		StatusSet:  t.Status.Valid,
		Settings:   ts.ToProto(),
		Runs:       runCount,
		Submission: sub.(*types.Submission),
	}, nil
}

func oauthToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	o, ok := i.(*models.OAuth)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	return &data.OAuthState{
		State:  o.State,
		Scopes: strings.Split(o.Scopes, ","),
	}, nil
}

func refFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	r, ok := i.(*types.Ref)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	return &models.Ref{
		ID:           r.Id,
		RepositoryID: r.Repository.Id,
		Ref:          r.RefName,
		Sha:          r.Sha,
	}, nil
}

func refToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	r, ok := i.(*models.Ref)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	repo, err := r.Repository().One(ctx, db)
	if err != nil {
		return nil, err
	}

	rp, err := repoToProto(ctx, db, repo)
	if err != nil {
		return nil, err
	}

	return &types.Ref{
		Id:         r.ID,
		Repository: rp.(*types.Repository),
		RefName:    r.Ref,
		Sha:        r.Sha,
	}, nil
}

func repoFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	r, ok := i.(*types.Repository)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	var owner *models.User

	if r.Owner != nil {
		var err error
		u, err := userFromProto(ctx, db, r.Owner)
		if err != nil {
			return nil, err
		}

		owner = u.(*models.User)
	}

	return &models.Repository{
		ID:          r.Id,
		Name:        r.Name,
		Private:     r.Private,
		Disabled:    null.BoolFrom(r.Disabled),
		OwnerID:     owner.ID,
		AutoCreated: r.AutoCreated,
		HookSecret:  r.HookSecret,
		Github:      r.Github,
	}, nil
}

func repoToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	r, ok := i.(*models.Repository)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	var owner *models.User

	owner, err := r.Owner().One(ctx, db)
	if err != nil {
		return nil, err
	}

	u, err := userToProto(ctx, db, owner)
	if err != nil {
		return nil, err
	}

	retOwner := u.(*types.User)

	return &types.Repository{
		Id:          r.ID,
		Name:        r.Name,
		Private:     r.Private,
		Disabled:    r.Disabled.Bool,
		Owner:       retOwner,
		AutoCreated: r.AutoCreated,
		HookSecret:  r.HookSecret,
		Github:      r.Github,
	}, nil
}

func userFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	u, ok := i.(*types.User)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	return &models.User{
		ID:               u.Id,
		Username:         u.Username,
		LastScannedRepos: null.TimeFrom(u.LastScannedRepos.AsTime()),
		Token:            u.TokenJSON,
	}, nil
}

func userToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	u, ok := i.(*models.User)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	return &types.User{
		Id:               u.ID,
		Username:         u.Username,
		LastScannedRepos: timestamppb.New(u.LastScannedRepos.Time),
		TokenJSON:        u.Token,
	}, nil
}

func subToProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	s, ok := i.(*models.Submission)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	var user *types.User

	pu, err := s.User().One(ctx, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err == nil {
		u, err := userToProto(ctx, db, pu)
		if err != nil {
			return nil, err
		}

		user = u.(*types.User)
	}

	// var status bool
	// if s.Status != nil {
	// 	status = *s.Status
	// }

	br, err := s.BaseRef().One(ctx, db)
	if err != nil {
		return nil, err
	}

	baseRef, err := refToProto(ctx, db, br)
	if err != nil {
		return nil, err
	}

	hr, err := s.HeadRef().One(ctx, db)
	if err != nil {
		return nil, err
	}

	headRef, err := refToProto(ctx, db, hr)
	if err != nil {
		return nil, err
	}

	tasksCount, err := s.Tasks().Count(ctx, db)
	if err != nil {
		return nil, err
	}

	return &types.Submission{
		Id:         s.ID,
		BaseRef:    baseRef.(*types.Ref),
		HeadRef:    headRef.(*types.Ref),
		User:       user,
		TasksCount: tasksCount,
		// RunsCount:  s.RunsCount,
		CreatedAt: timestamppb.New(s.CreatedAt),
		// StartedAt:  timestamppb.New(s.StartedAt),
		// FinishedAt: timestamppb.New(s.FinishedAt),
		// StatusSet: s.Status != nil,
		// Status:   status,
		// Canceled: s.Canceled,
		// TicketID: s.TicketID,
	}, nil
}

func subFromProto(ctx context.Context, db *sql.DB, i interface{}) (interface{}, error) {
	gt, ok := i.(*types.Submission)
	if !ok {
		return nil, ErrConversionInvalidType
	}

	var (
		u       *models.User
		headref *models.Ref
		err     error
	)

	if gt.User != nil {
		uc, err := userFromProto(ctx, db, gt.User)
		if err != nil {
			return nil, utils.WrapError(err, "converting for use in submission")
		}

		u = uc.(*models.User)
	}

	if gt.HeadRef != nil {
		headrefc, err := refFromProto(ctx, db, gt.HeadRef)
		if err != nil {
			return nil, utils.WrapError(err, "converting for use in submission")
		}

		headref = headrefc.(*models.Ref)
	}

	baserefc, err := refFromProto(ctx, db, gt.BaseRef)
	if err != nil {
		return nil, utils.WrapError(err, "converting for use in submission")
	}

	baseref := baserefc.(*models.Ref)

	if headref == nil {
		headref = baseref
	}

	// var status *bool
	// if gt.StatusSet {
	// 	status = &gt.Status
	// }
	//
	created := gt.CreatedAt.AsTime()

	if created.IsZero() {
		// this is a new record and hasn't been updated. Bump the created_at time.
		t := time.Now()
		created = t
	}
	//
	// finished := MakeTime(gt.FinishedAt, true)
	// started := MakeTime(gt.StartedAt, true)

	var uid *int64

	if u != nil {
		uid = &u.ID
	}

	var headrefID *int64
	if headref != nil {
		headrefID = &headref.ID
	}

	return &models.Submission{
		ID:        gt.Id,
		UserID:    null.Int64FromPtr(uid),
		BaseRefID: baseref.ID,
		HeadRefID: null.Int64FromPtr(headrefID),
		// TasksCount: gt.TasksCount,
		// RunsCount:  gt.RunsCount,
		CreatedAt: created,
		// FinishedAt: finished,
		// StartedAt:  started,
		// Status:     status,
		// Canceled:   gt.Canceled,
		// TicketID:   gt.TicketID,
	}, nil
}
