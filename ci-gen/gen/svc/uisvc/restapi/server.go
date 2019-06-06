package restapi

import (
	"strings"

	"github.com/tinyci/ci-agents/handlers"

	"github.com/tinyci/ci-agents/ci-gen/gen/svc/uisvc/restapi/operations"
)

// SetRoutes sets the routes in the handler so gin can execute them.
func (hc HandlerConfig) SetRoutes(h *handlers.H) {
	h.Routes = handlers.Routes{}

	addRoute(
		h,
		handlers.TransformSwaggerRoute("/capabilities/{username}/{capability}"),
		"DELETE",
		&handlers.Route{
			ParamValidator: operations.DeleteCapabilitiesUsernameCapabilityValidateURLParams,
			Handler:        operations.DeleteCapabilitiesUsernameCapability,
			Method:         "DELETE",
			UseCORS:        true,
			UseAuth:        true,
			Capability:     "modify:user",
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/token"),
		"DELETE",
		&handlers.Route{
			ParamValidator: operations.DeleteTokenValidateURLParams,
			Handler:        operations.DeleteToken,
			Method:         "DELETE",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/errors"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetErrorsValidateURLParams,
			Handler:        operations.GetErrors,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/log/attach/{id}"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetLogAttachIDValidateURLParams,
			Handler:        operations.GetLogAttachID,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/loggedin"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetLoggedinValidateURLParams,
			Handler:        operations.GetLoggedin,
			Method:         "GET",
			UseCORS:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/login"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetLoginValidateURLParams,
			Handler:        operations.GetLogin,
			Method:         "GET",
			UseCORS:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/login/upgrade"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetLoginUpgradeValidateURLParams,
			Handler:        operations.GetLoginUpgrade,
			Method:         "GET",
			UseCORS:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/logout"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetLogoutValidateURLParams,
			Handler:        operations.GetLogout,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/repositories/ci/add/{owner}/{repo}"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRepositoriesCiAddOwnerRepoValidateURLParams,
			Handler:        operations.GetRepositoriesCiAddOwnerRepo,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
			Capability:     "modify:ci",
			TokenScope:     "repo",
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/repositories/ci/del/{owner}/{repo}"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRepositoriesCiDelOwnerRepoValidateURLParams,
			Handler:        operations.GetRepositoriesCiDelOwnerRepo,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
			Capability:     "modify:ci",
			TokenScope:     "repo",
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/repositories/my"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRepositoriesMyValidateURLParams,
			Handler:        operations.GetRepositoriesMy,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/repositories/sub/add/{owner}/{repo}"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRepositoriesSubAddOwnerRepoValidateURLParams,
			Handler:        operations.GetRepositoriesSubAddOwnerRepo,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/repositories/sub/del/{owner}/{repo}"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRepositoriesSubDelOwnerRepoValidateURLParams,
			Handler:        operations.GetRepositoriesSubDelOwnerRepo,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/repositories/subscribed"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRepositoriesSubscribedValidateURLParams,
			Handler:        operations.GetRepositoriesSubscribed,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/repositories/visible"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRepositoriesVisibleValidateURLParams,
			Handler:        operations.GetRepositoriesVisible,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/run/{run_id}"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRunRunIDValidateURLParams,
			Handler:        operations.GetRunRunID,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/runs"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRunsValidateURLParams,
			Handler:        operations.GetRuns,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/runs/count"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetRunsCountValidateURLParams,
			Handler:        operations.GetRunsCount,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/submit"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetSubmitValidateURLParams,
			Handler:        operations.GetSubmit,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
			Capability:     "submit",
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/tasks"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetTasksValidateURLParams,
			Handler:        operations.GetTasks,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/tasks/count"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetTasksCountValidateURLParams,
			Handler:        operations.GetTasksCount,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/tasks/runs/{id}"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetTasksRunsIDValidateURLParams,
			Handler:        operations.GetTasksRunsID,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/tasks/runs/{id}/count"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetTasksRunsIDCountValidateURLParams,
			Handler:        operations.GetTasksRunsIDCount,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/tasks/subscribed"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetTasksSubscribedValidateURLParams,
			Handler:        operations.GetTasksSubscribed,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/token"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetTokenValidateURLParams,
			Handler:        operations.GetToken,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/user/properties"),
		"GET",
		&handlers.Route{
			ParamValidator: operations.GetUserPropertiesValidateURLParams,
			Handler:        operations.GetUserProperties,
			Method:         "GET",
			UseCORS:        true,
			UseAuth:        true,
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/cancel/{run_id}"),
		"POST",
		&handlers.Route{
			ParamValidator: operations.PostCancelRunIDValidateURLParams,
			Handler:        operations.PostCancelRunID,
			Method:         "POST",
			UseCORS:        true,
			UseAuth:        true,
			Capability:     "cancel",
		},
	)
	addRoute(
		h,
		handlers.TransformSwaggerRoute("/capabilities/{username}/{capability}"),
		"POST",
		&handlers.Route{
			ParamValidator: operations.PostCapabilitiesUsernameCapabilityValidateURLParams,
			Handler:        operations.PostCapabilitiesUsernameCapability,
			Method:         "POST",
			UseCORS:        true,
			UseAuth:        true,
			Capability:     "modify:user",
		},
	)
}

func addRoute(h *handlers.H, path, method string, r *handlers.Route) {
	if h.Routes[path] == nil {
		h.Routes[path] = map[string]*handlers.Route{}
	}

	h.Routes[path][strings.ToLower(method)] = r
}
