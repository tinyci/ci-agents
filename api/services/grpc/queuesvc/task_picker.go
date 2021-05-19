package queuesvc

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	grpcHandler "github.com/tinyci/ci-agents/api/handlers/grpc"
	"github.com/tinyci/ci-agents/ci-gen/grpc/types"
	"github.com/tinyci/ci-agents/clients/log"
	topTypes "github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type taskPicker struct {
	handler *grpcHandler.H
	logger  *log.SubLogger
}

func (sp *submissionProcessor) newTaskPicker() *taskPicker {
	return &taskPicker{handler: sp.handler, logger: sp.logger}
}

func (tp *taskPicker) pick(ctx context.Context, sub *topTypes.Submission, repoInfo *repoInfo) ([]*types.QueueItem, error) {
	process := map[string]struct{}{}

	dirs, taskdirs, err := tp.toProcess(ctx, repoInfo)
	if err != nil {
		return nil, utils.WrapError(err, "determining what to process")
	}

	if (sub.All && sub.Manual) || (repoInfo.forkRef.Repository.Id == repoInfo.parent.Id && repoInfo.parentRef.RefName == repoInfo.mainBranch()) {
		for _, dir := range taskdirs {
			process[dir] = struct{}{}
		}
	} else {
		process = tp.selectTasks(dirs, taskdirs)
	}

	if err := tp.cancelPreviousRuns(ctx, repoInfo); err != nil {
		return nil, utils.WrapError(err, "while canceling the previous runs")
	}

	subRecord, err := tp.handler.Clients.Data.PutSubmission(ctx, &types.Submission{TicketID: repoInfo.ticketID, User: repoInfo.user, HeadRef: repoInfo.forkRef, BaseRef: repoInfo.parentRef})
	if err != nil {
		return nil, utils.WrapError(err, "couldn't convert submission")
	}

	// XXX reusing taskdirs here because it serves the same purpose, albeit
	// slightly different form here.
	tasks, taskdirs, err := tp.makeTaskDirs(ctx, process, subRecord, repoInfo)
	if err != nil {
		return nil, utils.WrapError(err, "computing task directories")
	}

	queueCreateTime := time.Now()
	tp.logger.Info(ctx, "Generating Queue Items")
	qis := []*types.QueueItem{}

	for _, dir := range taskdirs {
		task := tasks[dir]

		if len(task.Settings.Runs) > 0 {

			tmpQIs, err := tp.generateQueueItems(ctx, dir, task, repoInfo)
			if err != nil {
				return nil, utils.WrapError(err, "generating queue items")
			}

			qis = append(qis, tmpQIs...)
		}
	}
	tp.logger.Infof(ctx, "Computing queue items took %v", time.Since(queueCreateTime))

	return qis, nil
}

func (tp *taskPicker) getDiffFiles(ctx context.Context, repoInfo *repoInfo) (map[string]struct{}, []string, error) {
	client, err := repoInfo.client(tp.handler)
	if err != nil {
		return nil, nil, utils.WrapError(err, "obtaining parent owner's client")
	}

	diffFiles, err := client.GetDiffFiles(ctx, repoInfo.parent.Name, repoInfo.parentRef.Sha, repoInfo.forkRef.Sha)
	if err != nil {
		return nil, nil, utils.WrapError(err, "getting file list for diff")
	}

	dirs := map[string]struct{}{}

	for _, file := range diffFiles {
		dirs[path.Dir(file)] = struct{}{}
	}

	allFiles, err := client.GetFileList(ctx, repoInfo.fork.Name, repoInfo.forkRef.Sha)
	if err != nil {
		return nil, nil, err
	}

	return dirs, allFiles, nil
}

func (tp *taskPicker) toProcess(ctx context.Context, repoInfo *repoInfo) (map[string]struct{}, []string, error) {
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

func (tp *taskPicker) cancelPreviousRuns(ctx context.Context, repoInfo *repoInfo) error {
	if err := tp.handler.Clients.Data.CancelRefByName(ctx, repoInfo.forkRef.Repository.Id, repoInfo.forkRef.RefName); err != nil {
		tp.logger.Errorf(ctx, "Couldn't cancel ref %q repo %d; will continue anyway: %v\n", repoInfo.forkRef.RefName, repoInfo.parent.Id, err)
	}

	client, err := repoInfo.client(tp.handler)
	if err != nil {
		return utils.WrapError(err, "could not retrieve parent owner's github client")
	}

	if err := client.ClearStates(ctx, repoInfo.parent.Name, repoInfo.forkRef.Sha); err != nil {
		tp.logger.Errorf(ctx, "Couldn't clear states for repo %q ref %q: %v", repoInfo.parent.Name, repoInfo.forkRef.Sha, err)
	}

	return nil
}

func (tp *taskPicker) makeTask(ctx context.Context, subRecord *types.Submission, dir string, repoInfo *repoInfo) (*types.Task, error) {
	client, err := repoInfo.client(tp.handler)
	if err != nil {
		return nil, utils.WrapError(err, "obtaining client for parent owner")
	}

	content, err := client.GetFile(ctx, repoInfo.fork.Name, repoInfo.forkRef.Sha, path.Join(dir, taskConfigFilename))
	if err != nil {
		return nil, utils.WrapError(err, "obtaining task instructions for repo %q sha %q dir %q", repoInfo.fork.Name, repoInfo.forkRef.Sha, dir)
	}

	ts, err := topTypes.NewTaskSettings(content, false, *repoInfo.repoConfig)
	if err != nil {
		if repoInfo.ticketID != 0 {
			if cerr := client.CommentError(ctx, repoInfo.parent.Name, repoInfo.ticketID, utils.WrapError(err, "tinyCI had an error processing your pull request")); cerr != nil {
				return nil, utils.WrapError(cerr, "attempting to alert the user about the error in their pull request")
			}
		}

		return nil, utils.WrapError(err, "validating task settings for repo %q sha %q dir %q", repoInfo.fork.Name, repoInfo.forkRef.Sha, dir)
	}

	settings := &types.TaskSettings{}

	return &types.Task{
		Path:       dir,
		Settings:   settings,
		CreatedAt:  timestamppb.Now(),
		Submission: subRecord,
	}, utils.JSONIO(ts, settings)
}

func (tp *taskPicker) makeTaskDirs(ctx context.Context, process map[string]struct{}, subRecord *types.Submission, repoInfo *repoInfo) (map[string]*types.Task, []string, error) {
	tasks := map[string]*types.Task{}

	tp.logger.Info(ctx, "Computing task dirs")

	taskdirs := []string{}
	for dir := range process {
		taskdirs = append(taskdirs, dir)
	}

	for i := 0; i < len(taskdirs); i++ {
		task, err := tp.makeTask(ctx, subRecord, taskdirs[i], repoInfo)
		if err != nil {
			return nil, nil, utils.WrapError(err, "making task")
		}

		tasks[taskdirs[i]] = task
		for _, dir := range task.Settings.Dependencies {
			if _, ok := process[dir]; !ok {
				process[dir] = struct{}{}
				taskdirs = append(taskdirs, dir)
			}
		}
	}

	sort.Strings(taskdirs)

	return tasks, taskdirs, nil
}

func (tp *taskPicker) generateQueueItems(ctx context.Context, dir string, task *types.Task, repoInfo *repoInfo) ([]*types.QueueItem, error) {
	qis := []*types.QueueItem{}

	task, err := tp.handler.Clients.Data.PutTask(ctx, task)
	if err != nil {
		return nil, utils.WrapError(err, "Could not insert task")
	}

	names := []string{}

	for name := range task.Settings.Runs {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		qi, err := tp.makeRunQueue(ctx, name, dir, task, repoInfo)
		if err != nil {
			return nil, utils.WrapError(err, "constructing queue item")
		}
		qis = append(qis, qi)
	}

	return qis, nil
}

func (tp *taskPicker) makeRunQueue(ctx context.Context, name, dir string, task *types.Task, repoInfo *repoInfo) (*types.QueueItem, error) {
	rs := task.Settings.Runs[name]

	dirStr := dir

	if dir == "." || dir == "" {
		dirStr = "*root*"
	}

	run := &types.Run{
		Name:      strings.Join([]string{dirStr, name}, ":"),
		Settings:  rs,
		Task:      task,
		CreatedAt: timestamppb.Now(),
	}

	go tp.setPendingStatus(ctx, run, repoInfo)

	return &types.QueueItem{
		Run:       run,
		QueueName: run.Settings.Queue,
	}, nil
}

func (tp *taskPicker) setPendingStatus(ctx context.Context, run *types.Run, repoInfo *repoInfo) {
	parts := strings.SplitN(repoInfo.parent.Name, "/", 2)
	if len(parts) != 2 {
		tp.logger.Error(ctx, fmt.Errorf("invalid repo name %q", repoInfo.parent.Name))
	}

	client, err := repoInfo.client(tp.handler)
	if err != nil {
		tp.logger.Error(ctx, utils.WrapError(err, "could not obtain client for parent owner"))
	}

	if err := client.PendingStatus(ctx, parts[0], parts[1], run.Name, repoInfo.forkRef.Sha, tp.handler.URL); err != nil {
		tp.logger.Error(ctx, utils.WrapError(err, "could not set pending status"))
	}
}
