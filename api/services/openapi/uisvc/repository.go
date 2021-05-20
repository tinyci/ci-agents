package uisvc

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tinyci/ci-agents/ci-gen/openapi/services/uisvc"
	"github.com/tinyci/ci-agents/utils"
)

// GetRepositoriesSubscribed lists all subscribed repos as JSON.
func (h *H) GetRepositoriesSubscribed(ctx echo.Context, params uisvc.GetRepositoriesSubscribedParams) error {
	user, ok := h.getUser(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	s := ""
	if params.Search != nil {
		s = *params.Search
	}

	repos, err := h.clients.Data.ListSubscriptions(ctx.Request().Context(), user.Username, s)
	if err != nil {
		return err
	}

	r, err := h.convertRepositories(ctx, repos)
	if err != nil {
		return err
	}

	return ctx.JSON(200, r)
}

// GetRepositoriesScan scans for owned and managed repositories for Add-to-CI operations.
func (h *H) GetRepositoriesScan(ctx echo.Context) error {
	user, ok := h.getUser(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	if param, ok := h.Config.ServiceConfig["last_scanned_wait"].(string); ok {
		dur, err := time.ParseDuration(param)
		if err != nil {
			return err
		}

		if user.LastScannedRepos.IsValid() && time.Since(user.LastScannedRepos.AsTime()) < dur {
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
	user, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	repos, err := h.clients.Data.OwnedRepositories(ctx.Request().Context(), user, params.Search)
	if err != nil {
		return err
	}

	r, err := h.convertRepositories(ctx, repos)
	if err != nil {
		return err
	}

	return ctx.JSON(200, r)
}

// GetRepositoriesVisible returns all the repos the user can see.
func (h *H) GetRepositoriesVisible(ctx echo.Context, params uisvc.GetRepositoriesVisibleParams) error {
	user, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	repos, err := h.clients.Data.AllRepositories(ctx.Request().Context(), user, params.Search)
	if err != nil {
		return err
	}

	r, err := h.convertRepositories(ctx, repos)
	if err != nil {
		return err
	}

	return ctx.JSON(200, r)
}

// GetRepositoriesCiDelOwnerRepo removes the repository from CI. that's it.
func (h *H) GetRepositoriesCiDelOwnerRepo(ctx echo.Context, owner string, repository string) error {
	username, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

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

	if err := github.TeardownHook(context.Background(), owner, repository, h.Config.HookURL); err != nil {
		return err
	}

	if err := h.clients.Data.DisableRepository(ctx.Request().Context(), username, path.Join(owner, repository)); err != nil {
		return err
	}

	return ctx.NoContent(200)
}

// GetRepositoriesCiAddOwnerRepo adds the repository to CI and subscribes the user to it.
func (h *H) GetRepositoriesCiAddOwnerRepo(ctx echo.Context, owner string, repository string) error {
	user, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	github, err := h.getClient(ctx)
	if err != nil {
		return err
	}

	repoName := path.Join(owner, repository)
	if _, err := h.clients.Data.GetRepository(ctx.Request().Context(), repoName); err != nil {
		return err
	}

	if err := github.TeardownHook(context.Background(), owner, repository, h.Config.HookURL); err != nil {
		return err
	}

	err = h.clients.Data.EnableRepository(context.Background(), user, repoName)
	if err != nil {
		return err
	}

	postRepo, err := h.clients.Data.GetRepository(context.Background(), repoName)
	if err != nil {
		return err
	}

	if err := github.SetupHook(context.Background(), owner, repository, h.Config.HookURL, postRepo.HookSecret); err != nil {
		if err := h.clients.Data.DisableRepository(context.Background(), user, repoName); err != nil {
			return err
		}
		return err
	}

	err = h.clients.Data.AddSubscription(context.Background(), user, repoName)
	if err != nil {
		return err
	}

	r, err := h.C.FromProto(ctx.Request().Context(), postRepo)
	if err != nil {
		return err
	}

	return ctx.JSON(200, r)
}

// GetRepositoriesSubAddOwnerRepo adds a subscription for the user to the repo
func (h *H) GetRepositoriesSubAddOwnerRepo(ctx echo.Context, owner, repository string) error {
	user, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	if err := h.clients.Data.AddSubscription(context.Background(), user, path.Join(owner, repository)); err != nil {
		return err
	}

	return ctx.NoContent(200)
}

// GetRepositoriesSubDelOwnerRepo removes the subscription to the repository from the user account.
func (h *H) GetRepositoriesSubDelOwnerRepo(ctx echo.Context, owner, repository string) error {
	user, ok := h.getUsername(ctx)
	if !ok {
		return utils.ErrInvalidAuth
	}

	if err := h.clients.Data.DeleteSubscription(context.Background(), user, path.Join(owner, repository)); err != nil {
		return err
	}

	return ctx.NoContent(200)
}
