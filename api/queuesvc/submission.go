package queuesvc

import (
	"context"
	"fmt"
	"sort"

	gh "github.com/google/go-github/github"
	"github.com/tinyci/ci-agents/ci-gen/grpc/handler"
	"github.com/tinyci/ci-agents/clients/github"
	"github.com/tinyci/ci-agents/clients/log"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/model"
	"github.com/tinyci/ci-agents/types"
	"github.com/tinyci/ci-agents/utils"
)

const (
	defaultMainBranch  = "heads/master"
	repoConfigFilename = "tinyci.yml"
	taskConfigFilename = "task.yml"
)

// small cache of repository information we need
type repoInfo struct {
	ghParent   *gh.Repository
	ghFork     *gh.Repository
	parent     *model.Repository
	fork       *model.Repository
	parentRef  *model.Ref
	forkRef    *model.Ref
	user       *model.User
	repoConfig *types.RepoConfig
	ticketID   int64
}

type submissionProcessor struct {
	handler  *handler.H
	logger   *log.SubLogger
	repoInfo *repoInfo
}

func getLogger(sub *types.Submission, h *handler.H) *log.SubLogger {
	if sub != nil {
		return h.Clients.Log.WithFields(log.FieldMap{
			"parent":       sub.Parent,
			"fork":         sub.Fork,
			"head":         sub.HeadSHA,
			"base":         sub.BaseSHA,
			"manual":       fmt.Sprintf("%v", sub.Manual),
			"submitted_by": sub.SubmittedBy,
			"all":          fmt.Sprintf("%v", sub.All),
		})
	}
	return h.Clients.Log
}

func (qs *QueueServer) newSubmissionProcessor() *submissionProcessor {
	return &submissionProcessor{repoInfo: &repoInfo{}, handler: qs.H}
}

func (sp *submissionProcessor) process(ctx context.Context, sub *types.Submission) ([]*model.QueueItem, *errors.Error) {
	sp.logger = getLogger(sub, sp.handler)
	if err := sp.configureRepositories(ctx, sub); err != nil {
		return nil, err.Wrap("configuring repositories for submission")
	}

	client, err := sp.repoInfo.client(sp.handler)
	if err != nil {
		return nil, err.Wrap("fetching client for parent repository")
	}

	sp.repoInfo.repoConfig, err = sp.getRepoConfig(ctx, client)
	if err != nil {
		return nil, err.Wrap("obtaining repository configuration")
	}

	tp := sp.newTaskPicker()

	return tp.pick(ctx, sub, sp.repoInfo)
}

func (sp *submissionProcessor) configureRepositories(ctx context.Context, sub *types.Submission) *errors.Error {
	if err := sub.Validate(); err != nil {
		return err.Wrap("validating submission")
	}

	// manual submissions must be resolvable by the submitter to avoid security
	// leaks, so this uses the user's account to look up the parent info and
	// returns it so that it can be added to the submission data.
	if sub.Manual {
		user, userClient, err := sp.getSubmittedUserClient(ctx, sub.SubmittedBy)
		if err != nil {
			return err.Wrap("getting submitting user account info")
		}

		repo, err := userClient.GetRepository(ctx, sub.Fork)
		if err != nil {
			return err.Wrap("obtaining fork repository for submission -- probably no access")
		}

		sub.Parent, err = sp.selectParentOrFork(ctx, userClient, repo)
		if err != nil {
			return err.Wrap("while deriving parent information from fork")
		}

		sp.repoInfo.user = user
	}

	parent, err := sp.parentRepository(ctx, sub.Parent)
	if err != nil {
		return err.Wrap("obtaining parent repository")
	}

	client, err := sp.repoInfo.client(sp.handler)
	if err != nil {
		return err.Wrap("obtaining github client for parent repo owner")
	}

	if parent.Disabled {
		return errors.New("repository is not enabled")
	}

	sp.repoInfo.ghParent, err = client.GetRepository(ctx, parent.Name)
	if err != nil {
		return err.Wrap("checking access to parent repository on github")
	}

	fork, err := sp.makeFork(ctx, client, parent, sub.Fork)
	if err != nil {
		return err.Wrap("locating or creating fork record")
	}

	sp.repoInfo.ticketID = sub.TicketID

	if len(sub.HeadSHA) != 40 { // FIXME could be trumped with long branch names
		sub.HeadSHA, err = client.GetSHA(ctx, sub.Fork, sub.HeadSHA)
		if err != nil {
			return err.Wrap("while obtaining the HEAD SHA for the head repo/branch")
		}
	}

	sub.BaseSHA, err = client.GetSHA(ctx, sub.Parent, sp.repoInfo.mainBranch())
	if err != nil {
		return err.Wrap("while selecting HEAD SHA for base repo/branch")
	}

	if sub.BaseSHA == "0000000000000000000000000000000000000000" {
		fmt.Println("here")
		if sub.Fork == sub.Parent {
			// new branch; set to head ref
			sub.BaseSHA = sub.HeadSHA
		} else {
			return errors.New("base SHA was blank but this was not a new branch")
		}
	}

	sp.repoInfo.forkRef, err = sp.manageRefs(ctx, client, fork, sub.HeadSHA)
	if err != nil {
		return err
	}

	sp.repoInfo.parentRef, err = sp.manageRefs(ctx, client, parent, sub.BaseSHA)
	if err != nil {
		return err
	}
	return nil
}

func (sp *submissionProcessor) manageRefs(ctx context.Context, client github.Client, repo *model.Repository, sha string) (*model.Ref, *errors.Error) {
	refs, err := client.GetRefs(ctx, repo.Name, sha)
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

	ref, err := sp.handler.Clients.Data.GetRefByNameAndSHA(ctx, repo.Name, sha)
	if err != nil {
		if err.Contains(errors.ErrNotFound) {
			ref = &model.Ref{Repository: repo, RefName: refName, SHA: sha}

			id, err := sp.handler.Clients.Data.PutRef(ctx, ref)
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

func (sp *submissionProcessor) makeFork(ctx context.Context, client github.Client, parent *model.Repository, fork string) (*model.Repository, *errors.Error) {
	var err *errors.Error
	sp.repoInfo.ghFork, err = client.GetRepository(ctx, fork)
	if err != nil {
		return nil, err.Wrap("obtaining fork information from github")
	}

	if _, _, err := utils.OwnerRepo(sp.repoInfo.ghFork.GetFullName()); err != nil {
		return nil, err.Wrap("validating name of fork repository")
	}

retry:
	forkRepo, err := sp.forkRepository(ctx, sp.repoInfo.ghFork.GetFullName())
	if err != nil {
		if !err.Contains(errors.ErrNotFound) {
			return nil, err
		}

		if err := sp.handler.Clients.Data.PutRepositories(ctx, parent.Owner.Username, []*gh.Repository{sp.repoInfo.ghFork}, true); err != nil {
			return nil, err
		}
		goto retry
	}

	return forkRepo, nil
}

func (ri *repoInfo) client(h *handler.H) (github.Client, *errors.Error) {
	repoOwner := ri.parent.Owner
	if repoOwner == nil {
		return nil, errors.New("No owner for target repository")
	}

	return h.OAuth.GithubClient(repoOwner.Token), nil
}

func (sp *submissionProcessor) getSubmittedUserClient(ctx context.Context, submittedBy string) (*model.User, github.Client, *errors.Error) {
	if submittedBy == "" {
		return nil, nil, errors.New("invalid submission -- no `submitted by` field supplied")
	}

	user, err := sp.handler.Clients.Data.GetUser(ctx, submittedBy)
	if err != nil {
		return nil, nil, err.Wrap("obtaining user information for submitter")
	}

	token := &types.OAuthToken{}
	if err := utils.JSONIO(user.Token, token); err != nil {
		return nil, nil, err.Wrap("Decoding token from user account")
	}

	client := sp.handler.OAuth.GithubClient(token)

	return user, client, nil
}

func (sp *submissionProcessor) parentRepository(ctx context.Context, parent string) (*model.Repository, *errors.Error) {
	var err *errors.Error
	if sp.repoInfo.parent == nil {
		sp.repoInfo.parent, err = sp.handler.Clients.Data.GetRepository(ctx, parent)
	}

	return sp.repoInfo.parent, err
}

func (sp *submissionProcessor) forkRepository(ctx context.Context, fork string) (*model.Repository, *errors.Error) {
	var err *errors.Error
	if sp.repoInfo.fork == nil {
		sp.repoInfo.fork, err = sp.handler.Clients.Data.GetRepository(ctx, fork)
	}

	return sp.repoInfo.fork, err
}

func (sp *submissionProcessor) selectParentOrFork(ctx context.Context, client github.Client, fork *gh.Repository) (string, *errors.Error) {
	forkRepo, err := sp.forkRepository(ctx, fork.GetFullName())
	// this is ok; if modelRepo is nil then it's disabled.
	enabled := err == nil && !forkRepo.Disabled

	ret := fork.GetFullName()

	if !enabled && fork.GetFork() {
		ret = fork.GetParent().GetFullName()
		sp.logger.Info(ctx, "Selected parent of fork")
	} else {
		sp.logger.Info(ctx, "Selected fork; is directly enabled")
	}

	if _, _, err := utils.OwnerRepo(ret); err != nil {
		return "", err.Wrap("validating structure of parent repo name")
	}

	return ret, nil
}

func (ri *repoInfo) mainBranch() string {
	defaultBranch := ri.ghParent.GetDefaultBranch()

	if defaultBranch == "" {
		defaultBranch = defaultMainBranch
	} else {
		defaultBranch = "heads/" + defaultBranch
	}

	return defaultBranch
}

func (sp *submissionProcessor) getRepoConfig(ctx context.Context, client github.Client) (*types.RepoConfig, *errors.Error) {
	content, err := client.GetFile(ctx, sp.repoInfo.parent.Name, fmt.Sprintf("refs/%s", sp.repoInfo.mainBranch()), repoConfigFilename)
	if err != nil {
		return nil, err
	}

	return types.NewRepoConfig(content)
}
