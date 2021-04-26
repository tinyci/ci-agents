package uisvc

import "github.com/tinyci/ci-agents/ci-gen/grpc/types"

func sanitizeSubmissions(subs []*types.Submission) []*types.Submission {
	for _, sub := range subs {
		sanitizeRepository(sub.BaseRef.Repository)
		sanitizeRepository(sub.HeadRef.Repository)
	}

	return subs
}

func sanitizeUser(user *types.User) *types.User {
	user.TokenJSON = nil
	user.Errors = nil
	user.LastScannedRepos = nil
	return user
}

func sanitizeRuns(runs []*types.Run) []*types.Run {
	for _, run := range runs {
		sanitizeRun(run)
	}
	return runs
}

func sanitizeRun(run *types.Run) *types.Run {
	sanitizeTask(run.Task)
	return run
}

func sanitizeRepositories(repos []*types.Repository) []*types.Repository {
	for _, repo := range repos {
		sanitizeRepository(repo)
	}

	return repos
}

func sanitizeRepository(repo *types.Repository) *types.Repository {
	repo.HookSecret = ""
	repo.Github = nil
	sanitizeUser(repo.Owner)
	return repo
}

func sanitizeTask(task *types.Task) *types.Task {
	sanitizeRepository(task.Submission.BaseRef.Repository)
	sanitizeRepository(task.Submission.HeadRef.Repository)
	return task
}

func sanitizeTasks(tasks []*types.Task) []*types.Task {
	for _, task := range tasks {
		sanitizeTask(task)
	}

	return tasks
}
