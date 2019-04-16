// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"github.com/tinyci/ci-agents/api/uisvc/processors"
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
	router.SetProcessor("/errors", "get", processors.Errors)
	router.SetProcessor("/login", "get", processors.Login)
	router.SetProcessor("/loggedin", "get", processors.LoggedIn)
	router.SetProcessor("/logout", "get", processors.Logout)
	router.SetProcessor("/user/properties", "get", processors.GetUserProperties)

	router.SetProcessor("/submit", "get", processors.Submit)

	router.SetProcessor("/tasks", "get", processors.ListTasks)
	router.SetProcessor("/tasks/count", "get", processors.CountTasks)
	router.SetProcessor("/tasks/runs/{id}", "get", processors.GetRunsForTask)
	router.SetProcessor("/tasks/runs/{id}/count", "get", processors.CountRunsForTask)
	router.SetProcessor("/tasks/subscribed", "get", processors.ListSubscribedTasksForUser)

	router.SetProcessor("/runs", "get", processors.ListRuns)
	router.SetProcessor("/runs/count", "get", processors.CountRuns)
	router.SetProcessor("/run/{run_id}", "get", processors.GetRun)

	router.SetProcessor("/cancel/{run_id}", "post", processors.CancelRun)

	router.SetProcessor("/repositories/subscribed", "get", processors.ListRepositoriesSubscribed)
	router.SetProcessor("/repositories/my", "get", processors.ListRepositoriesMy)
	router.SetProcessor("/repositories/visible", "get", processors.ListRepositoriesVisible)
	router.SetProcessor("/repositories/ci/add/{owner}/{repo}", "get", processors.AddRepositoryToCI)
	router.SetProcessor("/repositories/ci/del/{owner}/{repo}", "get", processors.DeleteRepositoryFromCI)
	router.SetProcessor("/repositories/sub/add/{owner}/{repo}", "get", processors.AddRepositorySubscription)
	router.SetProcessor("/repositories/sub/del/{owner}/{repo}", "get", processors.DeleteRepositorySubscription)

	router.SetProcessor("/token", "get", processors.GetToken)
	router.SetProcessor("/token", "delete", processors.DeleteToken)
	router.SetWebsocketProcessor("/log/attach/{id}", processors.LogAttach)

	router.SetProcessor("/capabilities/{username}/{capability}", "post", processors.AddCapability)
	router.SetProcessor("/capabilities/{username}/{capability}", "delete", processors.RemoveCapability)

	return nil
}
