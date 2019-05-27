package processors

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

const defaultMasterBranch = "heads/master"

// InternalSubmission is a transformed struct with the types pulled from the db.
type InternalSubmission struct {
	Sub        *types.Submission
	ProcessMap map[string]bool
	RepoConfig *types.RepoConfig
	ParentRepo *model.Repository
	ParentRef  *model.Ref
	Ref        *model.Ref
}

func makeQueueItemsFromTask(h *handler.H, client github.Client, is *InternalSubmission, dir string, task *model.Task) ([]*model.QueueItem, *errors.Error) {
	qis := []*model.QueueItem{}
	for name, run := range task.TaskSettings.Runs {
		dirStr := dir

		if dir == "." || dir == "" {
			dirStr = "*root*"
		}

		run := &model.Run{
			Name:        strings.Join([]string{dirStr, name}, ":"),
			RunSettings: run,
			Task:        task,
			CreatedAt:   time.Now(),
		}

		qi := &model.QueueItem{
			Run:       run,
			QueueName: run.RunSettings.Queue,
		}

		qis = append(qis, qi)

		parts := strings.SplitN(is.ParentRepo.Name, "/", 2)
		if len(parts) != 2 {
			return qis, errors.New("invalid repo name")
		}

		go func() {
			if err := client.PendingStatus(parts[0], parts[1], run.Name, is.Ref.SHA, h.URL); err != nil {
				fmt.Println(err)
			}
		}()
	}

	return qis, nil
}

// GenerateQueueItems is the final stage in the process that generates the
// queue items that will be passed on to runners. It is assumed these queue
// items must still be posted to the data svc.
func GenerateQueueItems(h *handler.H, client github.Client, is *InternalSubmission) ([]*model.QueueItem, *errors.Error) {
	qis := []*model.QueueItem{}

	if err := h.Clients.Data.CancelRefByName(is.Ref.Repository.ID, is.Ref.RefName); err != nil {
		// FIXME we should log this, but not stop on it.
		fmt.Printf("Couldn't cancel ref %q repo %d; will continue anyway: %v\n", is.Ref.RefName, is.ParentRepo.ID, err)
	}

	if err := client.ClearStates(is.Sub.Parent, is.Sub.HeadSHA); err != nil {
		// FIXME we should log this, but not stop on it.
		fmt.Printf("Couldn't clear states for repo %q ref %q: %v", is.Sub.Parent, is.Sub.HeadSHA, err)
	}

	taskdirs := []string{}

	for dir := range is.ProcessMap {
		taskdirs = append(taskdirs, dir)
	}

	tasks := map[string]*model.Task{}

	for i := 0; i < len(taskdirs); i++ {
		dir := taskdirs[i]

		// FIXME move this string.
		content, err := client.GetFile(is.Sub.Fork, is.Sub.HeadSHA, path.Join(dir, "task.yml"))
		if err != nil {
			return nil, err
		}

		ts, err := types.NewTaskSettings(content, false, is.RepoConfig)
		if err != nil {
			if is.Sub.PullRequest != 0 {
				if cerr := client.CommentError(is.Sub.Parent, is.Sub.PullRequest, err.Wrap("tinyCI had an error processing your pull request")); cerr != nil {
					return nil, cerr
				}
			}

			return nil, err
		}

		task := &model.Task{
			Parent:        is.ParentRepo,
			BaseSHA:       is.Sub.BaseSHA,
			PullRequestID: is.Sub.PullRequest,
			Ref:           is.Ref,
			Path:          dir,
			TaskSettings:  ts,
			CreatedAt:     time.Now(),
		}

		tasks[dir] = task

		for _, dir := range ts.Dependencies {
			if _, ok := is.ProcessMap[dir]; !ok {
				is.ProcessMap[dir] = true
				taskdirs = append(taskdirs, dir)
			}
		}
	}

	for dir, task := range tasks {
		if err := task.Validate(); err != nil {
			// an error here merely means the task is invalid (probably because it
			// has no runs and is only dependencies). otherwise, we can continue
			// processing the test.
			continue
		}

		retTask, err := h.Clients.Data.PutTask(task)
		if err != nil {
			return qis, err
		}

		tmpQIs, err := makeQueueItemsFromTask(h, client, is, dir, retTask)
		if err != nil {
			return qis, err
		}

		qis = append(qis, tmpQIs...)
	}

	return qis, nil
}

// ManageRefs gathers or creates the refs necessary to relationally work with this task.
func ManageRefs(h *handler.H, client github.Client, repo *model.Repository, sha string) (*model.Ref, *errors.Error) {
	refs, err := client.GetRefs(repo.Github.GetFullName(), sha)
	if err != nil {
		return nil, err
	}

	var refName string

	if len(refs) > 0 {
		sort.Strings(refs)
		refName = refs[0]
	} else {
		refName = sha
	}

	if _, _, err := repo.OwnerRepo(); err != nil {
		return nil, err
	}

	ref, err := h.Clients.Data.GetRefByNameAndSHA(repo.Name, sha)
	if err != nil {
		if err.Contains(errors.ErrNotFound) {
			ref = &model.Ref{Repository: repo, RefName: refName, SHA: sha}

			id, err := h.Clients.Data.PutRef(ref)
			if err != nil {
				return nil, err
			}

			ref.ID = id
		} else {
			return nil, err
		}
	}

	return ref, nil
}

// ManageRepositories returns the parent and fork repo after confirming with github.
func ManageRepositories(h *handler.H, sub *types.Submission) (*model.Repository, *model.Repository, github.Client, *errors.Error) {
	if _, _, err := utils.OwnerRepo(sub.Parent); err != nil {
		return nil, nil, nil, err
	}

	parentIntRepo, eErr := h.Clients.Data.GetRepository(sub.Parent)
	if eErr != nil {
		return nil, nil, nil, eErr
	}

	if parentIntRepo.Disabled {
		return nil, nil, nil, errors.New("repository is not enabled")
	}

	repoOwner := parentIntRepo.Owner
	if repoOwner == nil {
		return nil, nil, nil, errors.New("No owner for target repository")
	}
	client := h.OAuth.GithubClient(repoOwner.Token)

	forkRepo, err := client.GetRepository(sub.Fork)
	if err != nil {
		return nil, nil, nil, err
	}

	// just make sure we still have access to the parent, we won't modify it here.
	if _, err := client.GetRepository(sub.Parent); err != nil {
		return nil, nil, nil, err
	}

	if _, _, err := utils.OwnerRepo(forkRepo.GetFullName()); err != nil {
		return nil, nil, nil, err
	}

retry:
	forkIntRepo, err := h.Clients.Data.GetRepository(forkRepo.GetFullName())
	if err != nil {
		if !err.Contains(errors.ErrNotFound) {
			return nil, nil, nil, err
		}

		if err := h.Clients.Data.PutRepositories(repoOwner.Username, []*gh.Repository{forkRepo}, true); err != nil {
			return nil, nil, nil, err
		}

		goto retry
	}

	return parentIntRepo, forkIntRepo, client, nil
}

// GetFileLists obtains file lists for the fork between the base and head shas.
func GetFileLists(client github.Client, sub *types.Submission) (map[string]interface{}, []string, *errors.Error) {
	diffFiles, err := client.GetDiffFiles(sub.Parent, sub.BaseSHA, sub.HeadSHA)
	if err != nil {
		return nil, nil, err
	}

	dirs := map[string]interface{}{}

	for _, file := range diffFiles {
		dirs[path.Dir(file)] = true
	}

	allFiles, err := client.GetFileList(sub.Fork, sub.HeadSHA)
	if err != nil {
		return nil, nil, err
	}

	return dirs, allFiles, nil
}

func getTaskDirs(client github.Client, sub *types.Submission, config *types.RepoConfig) (map[string]interface{}, []string, *errors.Error) {
	dirs, allFiles, err := GetFileLists(client, sub)
	if err != nil {
		return nil, nil, err
	}

	taskdirs := []string{}

	for _, file := range allFiles {
		if path.Base(file) == "task.yml" { // FIXME put this string somewhere else
			var skip bool
			for _, dir := range config.IgnoreDirs {
				if strings.HasPrefix(file, dir) {
					skip = true
					break
				}
			}

			if !skip {
				taskdirs = append(taskdirs, path.Dir(file))
			}
		}
	}

	sort.Strings(taskdirs)

	return dirs, taskdirs, nil
}

// PickTasks isolates the task dirs that need testing.
func PickTasks(client github.Client, sub *types.Submission, ref *model.Ref, parent *model.Repository, config *types.RepoConfig) (map[string]bool, *errors.Error) {
	process := []string{}

	mb := parent.Github.GetMasterBranch()
	if mb == "" {
		mb = defaultMasterBranch
	}

	dirs, taskdirs, err := getTaskDirs(client, sub, config)
	if err != nil {
		return nil, err
	}

	if (sub.All && sub.Manual) || (ref.Repository.ID == parent.ID && ref.RefName == mb) {
		process = taskdirs
	} else {
		for i := len(taskdirs) - 1; i >= 0; i-- {
			if _, ok := dirs[taskdirs[i]]; ok {
				process = append(process, taskdirs[i])
			} else {
				delete(dirs, taskdirs[i])
			}
		}

		for dir := range dirs {
			// the longest dirs will be at the end
			for i := len(taskdirs) - 1; i >= 0; i-- {
				if strings.HasPrefix(dir, taskdirs[i]) {
					process = append(process, taskdirs[i])
					break
				}
			}
		}

		process = append(process, ".")
	}

	processMap := map[string]bool{}

	for _, dir := range process {
		processMap[dir] = true
	}

	return processMap, nil
}

// GetRepoConfig gathers the repo configuration from the parent fork.
func GetRepoConfig(client github.Client, sub *types.Submission) (*types.RepoConfig, *errors.Error) {
	repo, err := client.GetRepository(sub.Parent)
	if err != nil {
		return nil, err
	}

	masterBranch := repo.GetMasterBranch()
	if masterBranch == "" {
		masterBranch = defaultMasterBranch
	}

	// FIXME move this string.
	content, err := client.GetFile(sub.Parent, fmt.Sprintf("refs/%s", masterBranch), "tinyci.yml")
	if err != nil {
		return nil, err
	}

	return types.NewRepoConfig(content)
}

func resolveParentInfo(h *handler.H, sub *types.Submission) (*types.Submission, *errors.Error) {
	// to do this properly, we take the submitted by argument in the case of a
	// manual submission. In the uisvc, this is taken from session data -- never
	// from foreign input so unless a foreign agent can submit directly to the
	// queuesvc this should not be an issue.
	user, eErr := h.Clients.Data.GetUser(sub.SubmittedBy)
	if eErr != nil {
		return nil, eErr
	}

	token := &types.OAuthToken{}
	if err := utils.JSONIO(user.Token, token); err != nil {
		return nil, err
	}

	client := h.OAuth.GithubClient(token)
	repo, err := client.GetRepository(sub.Fork)
	if err != nil {
		return nil, err
	}

	if repo.GetFork() {
		sub.Parent = repo.GetParent().GetFullName()
	} else {
		sub.Parent = sub.Fork
	}

	if _, _, err := utils.OwnerRepo(sub.Parent); err != nil {
		return nil, err
	}

	ciRepo, err := h.Clients.Data.GetRepository(sub.Parent)
	if err != nil {
		return nil, err
	}

	if ciRepo.Disabled {
		return nil, errors.New("repository is disabled")
	}

	if len(sub.HeadSHA) != 40 {
		sub.HeadSHA, err = client.GetSHA(sub.Fork, sub.HeadSHA)
		if err != nil {
			return nil, err
		}
	}

	sub.BaseSHA, err = client.GetSHA(sub.Parent, "heads/master")
	return sub, err
}

// Process handles the overall processing of the submission. All other calls in this package originate here.
func Process(h *handler.H, sub *types.Submission) (retQI []*model.QueueItem, retErr *errors.Error) {
	var is *InternalSubmission
	since := time.Now()

	defer func() {
		h.Clients.Log.Infof("Processing Submission took %v", time.Since(since))

		if retErr != nil && is != nil && is.ParentRepo.Owner != nil {
			client := h.OAuth.GithubClient(is.ParentRepo.Owner.Token)
			owner, repo, err := is.ParentRepo.OwnerRepo()
			if err != nil {
				h.Clients.Log.Error(err.Wrapf("%s/%s", owner, repo))
				return
			}

			if err := client.FinishedStatus(owner, repo, "*global*", is.Ref.SHA, h.URL, false, fmt.Sprintf("failed to start job: %v", retErr)); err != nil {
				h.Clients.Log.Error(err)
			}
		}
	}()

	if err := sub.Validate(); err != nil {
		return nil, err
	}

	if sub.Manual {
		var err *errors.Error
		if sub, err = resolveParentInfo(h, sub); err != nil {
			return nil, err
		}
	}

	parentRepo, forkRepo, client, err := ManageRepositories(h, sub)
	if err != nil {
		return nil, err
	}

	modelRef, err := ManageRefs(h, client, forkRepo, sub.HeadSHA)
	if err != nil {
		return nil, err
	}

	parentRef, err := ManageRefs(h, client, parentRepo, sub.BaseSHA)
	if err != nil {
		return nil, err
	}

	repoConfig, err := GetRepoConfig(client, sub)
	if err != nil {
		return nil, err
	}

	// fork ref, parent repo. if it's a push parent and fork will be the same so it works out.
	processMap, err := PickTasks(client, sub, modelRef, parentRepo, repoConfig)
	if err != nil {
		return nil, err
	}

	is = &InternalSubmission{
		Sub:        sub,
		ProcessMap: processMap,
		RepoConfig: repoConfig,
		ParentRepo: parentRepo,
		Ref:        modelRef,
		ParentRef:  parentRef,
	}

	return GenerateQueueItems(h, client, is)
}
