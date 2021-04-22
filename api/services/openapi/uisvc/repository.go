package uisvc

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
)

// GetRepositoriesSubscribed lists all subscribed repos as JSON.
func (h *H) GetRepositoriesSubscribed(ctx echo.Context, params uisvc.GetRepositoriesSubscribedParams) error {
	user, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	s := ""
	if params.Search != nil {
		s = *params.Search
	}

	repos, err := h.clients.Data.ListSubscriptions(ctx.Request().Context(), user.Username, s)
	if err != nil {
		return err
	}

	return ctx.JSON(200, repos)
}

// GetRepositoriesScan scans for owned and managed repositories for Add-to-CI operations.
func (h *H) GetRepositoriesScan(ctx echo.Context) error {
	user, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	if param, ok := h.config.ServiceConfig["last_scanned_wait"].(string); ok {
		dur, err := time.ParseDuration(param)
		if err != nil {
			return err
		}

		if user.LastScannedRepos != nil && time.Since(time.Time(*user.LastScannedRepos)) < dur {
			return ctx.NoContent(200)
		}
	}

	github, err := h.getClient(ctx)
	if err != nil {
		return err
	}

	scanCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	githubRepos, err := github.MyRepositories(scanCtx)
	if err != nil {
		return err
	}

	if err := h.clients.Data.PutRepositories(scanCtx, user.Username, githubRepos, true); err != nil {
		return err
	}

	return ctx.NoContent(200)
}

// GetRepositoriesMy lists the repositories the user can modify.
func (h *H) GetRepositoriesMy(ctx echo.Context, params uisvc.GetRepositoriesMyParams) error {
	user, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	repos, err := h.clients.Data.OwnedRepositories(ctx.Request().Context(), user.Username, params.Search)
	if err != nil {
		return err
	}

	return ctx.JSON(200, repos)
}

// GetRepositoriesVisible returns all the repos the user can see.
func (h *H) GetRepositoriesVisible(ctx echo.Context, params uisvc.GetRepositoriesVisibleParams) error {
	user, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	repos, err := h.clients.Data.AllRepositories(ctx.Request().Context(), user.Username, params.Search)
	if err != nil {
		return err
	}

	return ctx.JSON(200, repos)
}

// GetRepositoriesCiDelOwnerRepo removes the repository from CI. that's it.
func (h *H) GetRepositoriesCiDelOwnerRepo(ctx echo.Context, owner string, repository string) error {
	github, err := h.getClient(ctx)
	if err != nil {
		return err
	}

	repo, err := h.clients.Data.GetRepository(ctx.Request().Context(), path.Join(owner, repository))
	if err != nil {
		return err
	}

	if repo.Disabled {
		return errors.New("repo is not enabled")
	}

	if err := github.TeardownHook(context.Background(), owner, repository, h.config.HookURL); err != nil {
		return err
	}

	return ctx.NoContent(200)
}

// GetRepositoriesCiAddOwnerRepo adds the repository to CI and subscribes the user to it.
func (h *H) GetRepositoriesCiAddOwnerRepo(ctx echo.Context, owner string, repository string) error {
	user, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	github, err := h.getClient(ctx)
	if err != nil {
		return err
	}

	repoName := path.Join(owner, repository)
	if _, err := h.clients.Data.GetRepository(ctx.Request().Context(), repoName); err != nil {
		return err
	}

	if err := github.TeardownHook(context.Background(), owner, repository, h.config.HookURL); err != nil {
		return err
	}

	err = h.clients.Data.EnableRepository(context.Background(), user.Username, repoName)
	if err != nil {
		return err
	}

	postRepo, err := h.clients.Data.GetRepository(context.Background(), repoName)
	if err != nil {
		return err
	}

	if err := github.SetupHook(context.Background(), owner, repository, h.config.HookURL, postRepo.HookSecret); err != nil {
		if err := h.clients.Data.DisableRepository(context.Background(), user.Username, repoName); err != nil {
			return err
		}
		return err
	}

	err = h.clients.Data.AddSubscription(context.Background(), user.Username, repoName)
	if err != nil {
		return err
	}

	return ctx.JSON(200, postRepo)
}

// GetRepositoriesSubAddOwnerRepo adds a subscription for the user to the repo
func (h *H) GetRepositoriesSubAddOwnerRepo(ctx echo.Context, owner, repository string) error {
	user, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	if err := h.clients.Data.AddSubscription(context.Background(), user.Username, path.Join(owner, repository)); err != nil {
		return err
	}

	return ctx.NoContent(200)
}

// GetRepositoriesSubDelOwnerRepo removes the subscription to the repository from the user account.
func (h *H) GetRepositoriesSubDelOwnerRepo(ctx echo.Context, owner, repository string) error {
	user, err := h.getUser(ctx)
	if err != nil {
		return err
	}

	if err := h.clients.Data.DeleteSubscription(context.Background(), user.Username, path.Join(owner, repository)); err != nil {
		return err
	}

	return ctx.NoContent(200)
}
