package restapi

import (
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/handlers"
)

// ListRepositoriesSubscribed lists all subscribed repos as JSON.
func ListRepositoriesSubscribed(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	search, _ := ctx.GetQuery("search")
	repos, err := h.Clients.Data.ListSubscriptions(user.Username, search)
	return repos, 200, err
}

// ScanRepositories scans for owned and managed repositories for Add-to-CI operations.
func ScanRepositories(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	if param, ok := h.ServiceConfig["last_scanned_wait"]; ok {
		dur, err := time.ParseDuration(param.(string))
		if err != nil {
			return nil, 500, errors.New(err)
		}

		if user.LastScannedRepos != nil && time.Since(time.Time(*user.LastScannedRepos)) < dur {
			return nil, 200, nil
		}
	}

	github, err := h.GetClient(ctx)
	if err != nil {
		return nil, 500, err
	}

	githubRepos, err := github.MyRepositories(ctx)
	if err != nil {
		return nil, 500, err
	}

	if err := h.Clients.Data.PutRepositories(user.Username, githubRepos, true); err != nil {
		return nil, 500, err
	}

	return nil, 200, nil
}

// ListRepositoriesMy lists the repositories the user can modify.
func ListRepositoriesMy(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	repos, err := h.Clients.Data.OwnedRepositories(user.Username, ctx.GetString("search"))
	if err != nil {
		return nil, 500, err
	}

	return repos, 200, nil
}

// ListRepositoriesVisible returns all the repos the user can see.
func ListRepositoriesVisible(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	repos, err := h.Clients.Data.AllRepositories(user.Username, ctx.GetString("search"))
	return repos, 200, err
}

// DeleteRepositoryFromCI removes the repository from CI. that's it.
func DeleteRepositoryFromCI(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	github, err := h.GetClient(ctx)
	if err != nil {
		return nil, 500, err
	}

	repo, err := h.Clients.Data.GetRepository(path.Join(ctx.GetString("owner"), ctx.GetString("repo")))
	if err != nil {
		return nil, 500, err
	}

	if repo.Disabled {
		return nil, 500, errors.New("repo is not enabled")
	}

	if err := github.TeardownHook(ctx, ctx.GetString("owner"), ctx.GetString("repo"), h.HookURL); err != nil {
		return nil, 500, err
	}

	return nil, 200, h.Clients.Data.DisableRepository(user.Username, path.Join(ctx.GetString("owner"), ctx.GetString("repo")))
}

// AddRepositoryToCI adds the repository to CI and subscribes the user to it.
func AddRepositoryToCI(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	github, err := h.GetClient(ctx)
	if err != nil {
		return nil, 500, err
	}

	repoName := path.Join(ctx.GetString("owner"), ctx.GetString("repo"))
	if _, err := h.Clients.Data.GetRepository(repoName); err != nil {
		return nil, 500, err
	}

	if err := github.TeardownHook(ctx, ctx.GetString("owner"), ctx.GetString("repo"), h.HookURL); err != nil {
		return nil, 500, err
	}

	err = h.Clients.Data.EnableRepository(user.Username, repoName)
	if err != nil {
		return nil, 500, err
	}

	postRepo, err := h.Clients.Data.GetRepository(repoName)
	if err != nil {
		return nil, 500, err
	}

	if err := github.SetupHook(ctx, ctx.GetString("owner"), ctx.GetString("repo"), h.HookURL, postRepo.HookSecret); err != nil {
		if err := h.Clients.Data.DisableRepository(user.Username, repoName); err != nil {
			return nil, 500, err
		}
		return nil, 500, err
	}

	err = h.Clients.Data.AddSubscription(user.Username, repoName)
	if err != nil {
		return nil, 500, err
	}

	return postRepo, 200, nil
}

// AddRepositorySubscription adds a subscription for the user to the repo
func AddRepositorySubscription(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	return nil, 200, h.Clients.Data.AddSubscription(user.Username, path.Join(ctx.GetString("owner"), ctx.GetString("repo")))
}

// DeleteRepositorySubscription removes the subscription to the repository from the user account.
func DeleteRepositorySubscription(h *handlers.H, ctx *gin.Context) (interface{}, int, *errors.Error) {
	user, err := h.GetUser(ctx)
	if err != nil {
		return nil, 500, err
	}

	return nil, 200, h.Clients.Data.DeleteSubscription(user.Username, path.Join(ctx.GetString("owner"), ctx.GetString("repo")))
}
