package restapi

import (
	"strings"

	"github.com/tinyci/ci-agents/handlers"

	"github.com/tinyci/ci-agents/gen/svc/uisvc/restapi/operations"
)

// SetRoutes sets the routes in the handler so gin can execute them.
func (hc HandlerConfig) SetRoutes(h *handlers.H) {
	h.Routes = handlers.Routes{}

	p := ""

	p = handlers.TransformSwaggerRoute("/capabilities/{username}/{capability}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("DELETE")] = &handlers.Route{
		ParamValidator: operations.DeleteCapabilitiesUsernameCapabilityValidateURLParams,
		Handler:        operations.DeleteCapabilitiesUsernameCapability,
		Method:         "DELETE",
		UseCORS:        true,
		UseAuth:        true,
		Capability:     "modify:user",
	}
	p = handlers.TransformSwaggerRoute("/token")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("DELETE")] = &handlers.Route{
		ParamValidator: operations.DeleteTokenValidateURLParams,
		Handler:        operations.DeleteToken,
		Method:         "DELETE",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/errors")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetErrorsValidateURLParams,
		Handler:        operations.GetErrors,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/log/attach/{id}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetLogAttachIDValidateURLParams,
		Handler:        operations.GetLogAttachID,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/loggedin")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetLoggedinValidateURLParams,
		Handler:        operations.GetLoggedin,
		Method:         "GET",
		UseCORS:        true,
	}
	p = handlers.TransformSwaggerRoute("/login")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetLoginValidateURLParams,
		Handler:        operations.GetLogin,
		Method:         "GET",
		UseCORS:        true,
	}
	p = handlers.TransformSwaggerRoute("/login/upgrade")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetLoginUpgradeValidateURLParams,
		Handler:        operations.GetLoginUpgrade,
		Method:         "GET",
		UseCORS:        true,
	}
	p = handlers.TransformSwaggerRoute("/logout")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetLogoutValidateURLParams,
		Handler:        operations.GetLogout,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/repositories/ci/add/{owner}/{repo}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRepositoriesCiAddOwnerRepoValidateURLParams,
		Handler:        operations.GetRepositoriesCiAddOwnerRepo,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
		Capability:     "modify:ci",
		TokenScope:     "repo",
	}
	p = handlers.TransformSwaggerRoute("/repositories/ci/del/{owner}/{repo}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRepositoriesCiDelOwnerRepoValidateURLParams,
		Handler:        operations.GetRepositoriesCiDelOwnerRepo,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
		Capability:     "modify:ci",
		TokenScope:     "repo",
	}
	p = handlers.TransformSwaggerRoute("/repositories/my")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRepositoriesMyValidateURLParams,
		Handler:        operations.GetRepositoriesMy,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/repositories/sub/add/{owner}/{repo}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRepositoriesSubAddOwnerRepoValidateURLParams,
		Handler:        operations.GetRepositoriesSubAddOwnerRepo,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/repositories/sub/del/{owner}/{repo}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRepositoriesSubDelOwnerRepoValidateURLParams,
		Handler:        operations.GetRepositoriesSubDelOwnerRepo,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/repositories/subscribed")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRepositoriesSubscribedValidateURLParams,
		Handler:        operations.GetRepositoriesSubscribed,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/repositories/visible")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRepositoriesVisibleValidateURLParams,
		Handler:        operations.GetRepositoriesVisible,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/run/{run_id}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRunRunIDValidateURLParams,
		Handler:        operations.GetRunRunID,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/runs")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRunsValidateURLParams,
		Handler:        operations.GetRuns,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/runs/count")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetRunsCountValidateURLParams,
		Handler:        operations.GetRunsCount,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/submit")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetSubmitValidateURLParams,
		Handler:        operations.GetSubmit,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
		Capability:     "submit",
	}
	p = handlers.TransformSwaggerRoute("/tasks")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetTasksValidateURLParams,
		Handler:        operations.GetTasks,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/tasks/count")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetTasksCountValidateURLParams,
		Handler:        operations.GetTasksCount,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/tasks/runs/{id}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetTasksRunsIDValidateURLParams,
		Handler:        operations.GetTasksRunsID,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/tasks/runs/{id}/count")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetTasksRunsIDCountValidateURLParams,
		Handler:        operations.GetTasksRunsIDCount,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/tasks/subscribed")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetTasksSubscribedValidateURLParams,
		Handler:        operations.GetTasksSubscribed,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/token")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetTokenValidateURLParams,
		Handler:        operations.GetToken,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/user/properties")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("GET")] = &handlers.Route{
		ParamValidator: operations.GetUserPropertiesValidateURLParams,
		Handler:        operations.GetUserProperties,
		Method:         "GET",
		UseCORS:        true,
		UseAuth:        true,
	}
	p = handlers.TransformSwaggerRoute("/cancel/{run_id}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("POST")] = &handlers.Route{
		ParamValidator: operations.PostCancelRunIDValidateURLParams,
		Handler:        operations.PostCancelRunID,
		Method:         "POST",
		UseCORS:        true,
		UseAuth:        true,
		Capability:     "cancel",
	}
	p = handlers.TransformSwaggerRoute("/capabilities/{username}/{capability}")

	if h.Routes[p] == nil {
		h.Routes[p] = map[string]*handlers.Route{}
	}

	h.Routes[p][strings.ToLower("POST")] = &handlers.Route{
		ParamValidator: operations.PostCapabilitiesUsernameCapabilityValidateURLParams,
		Handler:        operations.PostCapabilitiesUsernameCapability,
		Method:         "POST",
		UseCORS:        true,
		UseAuth:        true,
		Capability:     "modify:user",
	}
}
