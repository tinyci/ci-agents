package queuesvc

import (
	"context"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
)

type taskPicker struct {
	handler *handler.H
	logger  *log.SubLogger
}

func (sp *submissionProcessor) newTaskPicker() *taskPicker {
	return &taskPicker{handler: sp.handler, logger: sp.logger}
}

func (tp *taskPicker) pick(ctx context.Context, sub *types.Submission, repoInfo *repoInfo) ([]*model.QueueItem, *errors.Error) {
	process := map[string]struct{}{}

	mb := repoInfo.parent.Github.GetMasterBranch()
	if mb == "" {
		mb = defaultMasterBranch
	}

	dirs, taskdirs, err := tp.toProcess(ctx, repoInfo)
	if err != nil {
		return nil, err.Wrap("determining what to process")
	}

	if (sub.All && sub.Manual) || (repoInfo.forkRef.Repository.ID == repoInfo.parent.ID && repoInfo.parentRef.RefName == mb) {
		for _, dir := range taskdirs {
			process[dir] = struct{}{}
		}
	} else {
		process = tp.selectTasks(dirs, taskdirs)
	}

	if err := tp.cancelPreviousRuns(ctx, repoInfo); err != nil {
		return nil, err.Wrap("while canceling the previous runs")
	}

	subRecord, err := tp.handler.Clients.Data.PutSubmission(ctx, &model.Submission{TicketID: repoInfo.ticketID, User: repoInfo.user, HeadRef: repoInfo.forkRef, BaseRef: repoInfo.parentRef})
	if err != nil {
		return nil, err.Wrap("couldn't convert submission")
	}

	// XXX reusing taskdirs here because it serves the same purpose, albeit
	// slightly different form here.
	tasks, taskdirs, err := tp.makeTaskDirs(ctx, process, subRecord, repoInfo)
	if err != nil {
		return nil, err.Wrap("computing task directories")
	}

	queueCreateTime := time.Now()
	tp.logger.Info(ctx, "Generating Queue Items")
	qis := []*model.QueueItem{}

	for _, dir := range taskdirs {
		task := tasks[dir]
		if err := task.Validate(); err != nil {
			// an error here merely means the task is invalid (probably because it
			// has no runs and is only dependencies). otherwise, we can continue
			// processing the test.
			continue
		}

		tmpQIs, err := tp.generateQueueItems(ctx, dir, task, repoInfo)
		if err != nil {
			return nil, err.Wrap("generating queue items")
		}

		qis = append(qis, tmpQIs...)
	}
	tp.logger.Infof(ctx, "Computing queue items took %v", time.Since(queueCreateTime))

	return qis, nil
}

func (tp *taskPicker) getDiffFiles(ctx context.Context, repoInfo *repoInfo) (map[string]struct{}, []string, *errors.Error) {
	client, err := repoInfo.client(tp.handler)
	if err != nil {
		return nil, nil, err.Wrap("obtaining parent owner's client")
	}

	diffFiles, err := client.GetDiffFiles(ctx, repoInfo.parent.Name, repoInfo.parentRef.SHA, repoInfo.forkRef.SHA)
	if err != nil {
		return nil, nil, err.Wrap("getting file list for diff")
	}

	dirs := map[string]struct{}{}

	for _, file := range diffFiles {
		dirs[path.Dir(file)] = struct{}{}
	}

	allFiles, err := client.GetFileList(ctx, repoInfo.fork.Name, repoInfo.forkRef.SHA)
	if err != nil {
		return nil, nil, err
	}

	return dirs, allFiles, nil
}

func (tp *taskPicker) toProcess(ctx context.Context, repoInfo *repoInfo) (map[string]struct{}, []string, *errors.Error) {
	dirs, allFiles, err := tp.getDiffFiles(ctx, repoInfo)
	if err != nil {
		return nil, nil, err
	}

	dirMap := map[string]struct{}{}
	taskdirs := []string{}

	for _, file := range allFiles {
		if path.Base(file) == taskConfigFilename {
			var skip bool
			for _, dir := range repoInfo.repoConfig.IgnoreDirs {
				if strings.HasPrefix(file, dir) {
					skip = true
					break
				}
			}

			if !skip {
				dir := path.Dir(file)

				if _, ok := dirMap[dir]; !ok {
					dirMap[dir] = struct{}{}
					taskdirs = append(taskdirs, dir)
				}
			}
		}
	}

	return dirs, taskdirs, nil
}

func (tp *taskPicker) selectTasks(dirs map[string]struct{}, taskdirs []string) map[string]struct{} {
	process := map[string]struct{}{}

	for i := len(taskdirs) - 1; i >= 0; i-- {
		if _, ok := dirs[taskdirs[i]]; ok {
			process[taskdirs[i]] = struct{}{}
		} else {
			delete(dirs, taskdirs[i])
		}
	}

	for dir := range dirs {
		// the longest dirs will be at the end
		for i := len(taskdirs) - 1; i >= 0; i-- {
			if strings.HasPrefix(dir, taskdirs[i]) {
				process[taskdirs[i]] = struct{}{}
				break
			}
		}
	}

	process["."] = struct{}{} // . is always tested.

	return process
}

func (tp *taskPicker) cancelPreviousRuns(ctx context.Context, repoInfo *repoInfo) *errors.Error {
	if err := tp.handler.Clients.Data.CancelRefByName(ctx, repoInfo.forkRef.Repository.ID, repoInfo.forkRef.RefName); err != nil {
		tp.logger.Errorf(ctx, "Couldn't cancel ref %q repo %d; will continue anyway: %v\n", repoInfo.forkRef.RefName, repoInfo.parent.ID, err)
	}

	client, err := repoInfo.client(tp.handler)
	if err != nil {
		return err.Wrap("could not retrieve parent owner's github client")
	}

	if err := client.ClearStates(ctx, repoInfo.parent.Name, repoInfo.forkRef.SHA); err != nil {
		tp.logger.Errorf(ctx, "Couldn't clear states for repo %q ref %q: %v", repoInfo.parent.Name, repoInfo.forkRef.SHA, err)
	}

	return nil
}

func (tp *taskPicker) makeTask(ctx context.Context, subRecord *model.Submission, dir string, repoInfo *repoInfo) (*model.Task, *errors.Error) {
	client, err := repoInfo.client(tp.handler)
	if err != nil {
		return nil, err.Wrapf("obtaining client for parent owner")
	}

	content, err := client.GetFile(ctx, repoInfo.fork.Name, repoInfo.forkRef.SHA, path.Join(dir, taskConfigFilename))
	if err != nil {
		return nil, err.Wrapf("obtaining task instructions for repo %q sha %q dir %q", repoInfo.fork.Name, repoInfo.forkRef.SHA, dir)
	}

	ts, err := types.NewTaskSettings(content, false, repoInfo.repoConfig)
	if err != nil {
		if repoInfo.ticketID != 0 {
			if cerr := client.CommentError(ctx, repoInfo.parent.Name, repoInfo.ticketID, err.Wrap("tinyCI had an error processing your pull request")); cerr != nil {
				return nil, cerr.Wrap("attempting to alert the user about the error in their pull request")
			}
		}

		return nil, err.Wrapf("validating task settings for repo %q sha %q dir %q", repoInfo.fork.Name, repoInfo.forkRef.SHA, dir)
	}

	return &model.Task{
		Path:         dir,
		TaskSettings: ts,
		CreatedAt:    time.Now(),
		Submission:   subRecord,
	}, nil
}

func (tp *taskPicker) makeTaskDirs(ctx context.Context, process map[string]struct{}, subRecord *model.Submission, repoInfo *repoInfo) (map[string]*model.Task, []string, *errors.Error) {
	tasks := map[string]*model.Task{}

	tp.logger.Info(ctx, "Computing task dirs")

	taskdirs := []string{}
	for dir := range process {
		taskdirs = append(taskdirs, dir)
	}

	for i := 0; i < len(taskdirs); i++ {
		task, err := tp.makeTask(ctx, subRecord, taskdirs[i], repoInfo)
		if err != nil {
			return nil, nil, err.Wrap("making task")
		}

		tasks[taskdirs[i]] = task

		for _, dir := range task.TaskSettings.Dependencies {
			if _, ok := process[dir]; !ok {
				process[dir] = struct{}{}
				taskdirs = append(taskdirs, dir)
			}
		}
	}
	sort.Strings(taskdirs)

	return tasks, taskdirs, nil
}

func (tp *taskPicker) generateQueueItems(ctx context.Context, dir string, task *model.Task, repoInfo *repoInfo) ([]*model.QueueItem, *errors.Error) {
	qis := []*model.QueueItem{}

	task, err := tp.handler.Clients.Data.PutTask(ctx, task)
	if err != nil {
		return nil, err.Wrap("Could not insert task")
	}

	names := []string{}

	for name := range task.TaskSettings.Runs {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		qi, err := tp.makeRunQueue(ctx, name, dir, task, repoInfo)
		if err != nil {
			return nil, err.Wrap("constructing queue item")
		}
		qis = append(qis, qi)
	}

	return qis, nil
}

func (tp *taskPicker) makeRunQueue(ctx context.Context, name, dir string, task *model.Task, repoInfo *repoInfo) (*model.QueueItem, *errors.Error) {
	rs := task.TaskSettings.Runs[name]

	dirStr := dir

	if dir == "." || dir == "" {
		dirStr = "*root*"
	}

	run := &model.Run{
		Name:        strings.Join([]string{dirStr, name}, ":"),
		RunSettings: rs,
		Task:        task,
		CreatedAt:   time.Now(),
	}

	go tp.setPendingStatus(ctx, run, repoInfo)

	return &model.QueueItem{
		Run:       run,
		QueueName: run.RunSettings.Queue,
	}, nil
}

func (tp *taskPicker) setPendingStatus(ctx context.Context, run *model.Run, repoInfo *repoInfo) {
	parts := strings.SplitN(repoInfo.parent.Name, "/", 2)
	if len(parts) != 2 {
		tp.logger.Error(ctx, errors.Errorf("invalid repo name %q", repoInfo.parent.Name))
	}

	client, err := repoInfo.client(tp.handler)
	if err != nil {
		tp.logger.Error(ctx, err.Wrap("could not obtain client for parent owner"))
	}

	if err := client.PendingStatus(ctx, parts[0], parts[1], run.Name, repoInfo.forkRef.SHA, tp.handler.URL); err != nil {
		tp.logger.Error(ctx, err.Wrap("could not set pending status"))
	}
}
