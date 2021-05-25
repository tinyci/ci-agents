package uisvc

import (
	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
)

func (h *H) convertUserErrors(ctx echo.Context, list *types.UserErrors) ([]*uisvc.UserError, error) {
	ret := []*uisvc.UserError{}

	for _, ue := range list.Errors {
		u, err := h.C.FromProto(ctx.Request().Context(), ue)
		if err != nil {
			return nil, err
		}
		ret = append(ret, u.(*uisvc.UserError))
	}

	return ret, nil
}

func (h *H) convertSubmissions(ctx echo.Context, list *types.SubmissionList) ([]*uisvc.ModelSubmission, error) {
	subs := []*uisvc.ModelSubmission{}

	for _, sub := range list.Submissions {
		s, err := h.C.FromProto(ctx.Request().Context(), sub)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s.(*uisvc.ModelSubmission))
	}

	return subs, nil
}

func (h *H) convertTasks(ctx echo.Context, list *types.TaskList) ([]*uisvc.Task, error) {
	tasks := []*uisvc.Task{}

	for _, task := range list.Tasks {
		t, err := h.C.FromProto(ctx.Request().Context(), task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t.(*uisvc.Task))
	}

	return tasks, nil
}

func (h *H) convertRuns(ctx echo.Context, list *types.RunList) ([]*uisvc.Run, error) {
	runs := []*uisvc.Run{}

	for _, run := range list.List {
		r, err := h.C.FromProto(ctx.Request().Context(), run)
		if err != nil {
			return nil, err
		}
		runs = append(runs, r.(*uisvc.Run))
	}

	return runs, nil
}

func (h *H) convertRepositories(ctx echo.Context, list *types.RepositoryList) ([]*uisvc.Repository, error) {
	repositories := []*uisvc.Repository{}

	for _, repository := range list.List {
		r, err := h.C.FromProto(ctx.Request().Context(), repository)
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, r.(*uisvc.Repository))
	}

	return repositories, nil
}
