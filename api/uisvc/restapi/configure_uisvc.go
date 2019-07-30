// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"github.com/tinyci/ci-agents/config"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// HandlerConfig is for on-disk configuration of the service. It will be parsed
// on boot and passed into each HandlerFunc.
type HandlerConfig struct{}

// MakeHandlerConfig takes a handlers.ServiceConfig and spits out a HandlerConfig.
func MakeHandlerConfig(sc config.ServiceConfig) *HandlerConfig {
	return &HandlerConfig{}
}

// Validate allows you to perform your own custom validations on the configuration.
func (hc HandlerConfig) Validate(h *handlers.H) *errors.Error {
	if h.ClientConfig.Data == "" {
		return errors.New("no datasvc url specified")
	}

	if h.ClientConfig.Log == "" {
		return errors.New("no logsvc url specified")
	}

	if err := h.Auth.Validate(true); err != nil {
		return err
	}

	return nil
}

// CustomInit allows you to perform any final magic before boot.
func (hc HandlerConfig) CustomInit(h *handlers.H) *errors.Error {
	h.NoTLSServer = true
	h.UseSessions = true
	h.Port = 6010
	h.Name = "uisvc"

	return nil
}

// DBConfigure configures the database if necessary.
func (hc HandlerConfig) DBConfigure(h *handlers.H) *errors.Error {
	return nil
}

// Configure allows you to configure the routes, in particular. Setting the
// processing functions here will be a big part of your day job.
func (hc HandlerConfig) Configure(router handlers.Routes) *errors.Error {
	router.SetProcessor("/errors", "get", Errors)
	router.SetProcessor("/login", "get", Login)
	router.SetProcessor("/login/upgrade", "get", Upgrade)
	router.SetProcessor("/loggedin", "get", LoggedIn)
	router.SetProcessor("/logout", "get", Logout)
	router.SetProcessor("/user/properties", "get", GetUserProperties)

	router.SetProcessor("/submit", "get", Submit)

	router.SetProcessor("/tasks", "get", ListTasks)
	router.SetProcessor("/tasks/count", "get", CountTasks)
	router.SetProcessor("/tasks/runs/{id}", "get", GetRunsForTask)
	router.SetProcessor("/tasks/runs/{id}/count", "get", CountRunsForTask)
	router.SetProcessor("/tasks/subscribed", "get", ListSubscribedTasksForUser)
	router.SetProcessor("/tasks/cancel/{id}", "post", CancelTask)

	router.SetProcessor("/submissions", "get", ListSubmissions)
	router.SetProcessor("/submissions/count", "get", CountSubmissions)
	router.SetProcessor("/submission/{id}", "get", GetSubmission)
	router.SetProcessor("/submission/{id}/tasks", "get", GetSubmissionTasks)

	router.SetProcessor("/runs", "get", ListRuns)
	router.SetProcessor("/runs/count", "get", CountRuns)
	router.SetProcessor("/run/{run_id}", "get", GetRun)

	router.SetProcessor("/cancel/{run_id}", "post", CancelRun)

	router.SetProcessor("/repositories/scan", "get", ScanRepositories)
	router.SetProcessor("/repositories/subscribed", "get", ListRepositoriesSubscribed)
	router.SetProcessor("/repositories/my", "get", ListRepositoriesMy)
	router.SetProcessor("/repositories/visible", "get", ListRepositoriesVisible)
	router.SetProcessor("/repositories/ci/add/{owner}/{repo}", "get", AddRepositoryToCI)
	router.SetProcessor("/repositories/ci/del/{owner}/{repo}", "get", DeleteRepositoryFromCI)
	router.SetProcessor("/repositories/sub/add/{owner}/{repo}", "get", AddRepositorySubscription)
	router.SetProcessor("/repositories/sub/del/{owner}/{repo}", "get", DeleteRepositorySubscription)

	router.SetProcessor("/token", "get", GetToken)
	router.SetProcessor("/token", "delete", DeleteToken)
	router.SetWebsocketProcessor("/log/attach/{id}", LogAttach)

	router.SetProcessor("/capabilities/{username}/{capability}", "post", AddCapability)
	router.SetProcessor("/capabilities/{username}/{capability}", "delete", RemoveCapability)

	return nil
}
