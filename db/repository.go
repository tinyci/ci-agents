package db

import (
	"context"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/gorilla/securecookie"
	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/utils"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func repoValidateHook(ctx context.Context, db boil.ContextExecutor, r *models.Repository) error {
	if r.Name == "" {
		return errors.New("name is empty")
	}

	_, _, err := utils.OwnerRepo(r.Name)
	if err != nil {
		return fmt.Errorf("invalid repository name: %w", err)
	}

	if r.OwnerID == 0 {
		return errors.New("record has no owner")
	}

	gh := &github.Repository{}

	if err := r.Github.Unmarshal(&gh); err != nil {
		return utils.WrapError(err, "while unmarshaling github parameters for repository %q", r.Name)
	}

	if gh.GetFullName() != r.Name {
		return fmt.Errorf("github repository (%q) and repository name (%q) are not aligned", gh.GetFullName(), r.Name)
	}

	return nil
}

// GetRepositoryByName retrieves the repository by its unique name.
func (m *Model) GetRepositoryByName(ctx context.Context, name string) (*models.Repository, error) {
	return models.Repositories(models.RepositoryWhere.Name.EQ(name)).One(ctx, m.db)
}

// GetRepositoryByNameForUser retrieves the repository by its unique name, scoped by the userID's view.
func (m *Model) GetRepositoryByNameForUser(ctx context.Context, name string, userID int64) (*models.Repository, error) {
	return models.Repositories(qm.Where("owner_id = ? or not private", userID), models.RepositoryWhere.Name.EQ(name)).One(ctx, m.db)
}

// AssignRepository assigns the repository to the user explicitly. The user must already be persisted.
func (m *Model) AssignRepository(ctx context.Context, repo *models.Repository, owner *models.User) error {
	return repo.SetOwner(ctx, m.db, false, owner)
}

func (m *Model) getRepoSearch(ctx context.Context, where string, search *string, args ...interface{}) ([]*models.Repository, error) {
	if search == nil {
		return models.Repositories(qm.Where(where, args...)).All(ctx, m.db)
	}

	searchStr := *search
	searchStr = strings.Replace(searchStr, "%", "\\%", -1)
	searchStr = strings.Replace(searchStr, "_", "\\_", -1)

	args = append(args, "%"+searchStr+"%")
	return models.Repositories(qm.Where(where+" and name like ? escape '\\'", args...)).All(ctx, m.db)
}

// GetOwnedRepos returns all repositories owned by the user
func (m *Model) GetOwnedRepos(ctx context.Context, uid int64, search *string) ([]*models.Repository, error) {
	return m.getRepoSearch(ctx, "owner_id = ?", search, uid)
}

// GetAllPublicRepos retrieves all repos that are not private
func (m *Model) GetAllPublicRepos(ctx context.Context, search *string) ([]*models.Repository, error) {
	return m.getRepoSearch(ctx, "not private", search)
}

// GetPrivateReposForUser retrieves all private repos that the user owns.
func (m *Model) GetPrivateReposForUser(ctx context.Context, uid int64, search *string) ([]*models.Repository, error) {
	return m.getRepoSearch(ctx, "owner_id = ? and private", search, uid)
}

// GetVisibleReposForUser retrieves all repos the user can "see" in the
// database.
func (m *Model) GetVisibleReposForUser(ctx context.Context, uid int64, search *string) ([]*models.Repository, error) {
	r, err := m.GetAllPublicRepos(ctx, search)
	if err != nil {
		return nil, err
	}

	r2, err := m.GetPrivateReposForUser(ctx, uid, search)
	if err != nil {
		return nil, err
	}

	// reverse order to prefer private repos at the top
	return append(r2, r...), nil
}

// SaveRepositories saves github repositories; it sets the *User provided to
// the owner of it.
func (m *Model) SaveRepositories(ctx context.Context, repos []*github.Repository, username string, autoCreated bool) error {
	owner, err := m.FindUserByName(ctx, username)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		_, err := m.GetRepositoryByName(ctx, repo.GetFullName())
		if err != nil {
			localRepo, err := m.mkRepositoryFromGithub(repo, owner, autoCreated)
			if err != nil {
				return err
			}

			if err := localRepo.Insert(ctx, m.db, boil.Infer()); err != nil {
				return err
			}
		}
	}

	owner.LastScannedRepos = null.TimeFrom(time.Now())
	_, err = owner.Update(ctx, m.db, boil.Infer())
	return err
}

func (m *Model) mkRepositoryFromGithub(repo *github.Repository, owner *models.User, autoCreated bool) (*models.Repository, error) {
	content, err := json.Marshal(repo)
	if err != nil {
		return nil, err
	}

	return &models.Repository{
		Name:        repo.GetFullName(),
		Private:     repo.GetPrivate(),
		Disabled:    null.BoolFrom(true), // created repos are disabled by default
		Github:      content,
		OwnerID:     owner.ID,
		AutoCreated: autoCreated,
	}, nil
}

// DisableRepository removes it from CI.
func (m *Model) DisableRepository(ctx context.Context, repo *models.Repository) error {
	if !repo.Disabled.Valid || repo.Disabled.Bool {
		return errors.New("repo is not enabled")
	}

	repo.Disabled = null.BoolFrom(true)
	_, err := repo.Update(ctx, m.db, boil.Infer())
	return err
}

// EnableRepository adds it to CI.
func (m *Model) EnableRepository(ctx context.Context, repo *models.Repository, ownerID int64) error {
	if repo.Disabled.Valid && !repo.Disabled.Bool {
		return errors.New("repo is already enabled")
	}

	repo.Disabled = null.BoolFrom(false)
	repo.HookSecret = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(24)), "=")
	repo.OwnerID = ownerID
	_, err := repo.Update(ctx, m.db, boil.Infer())
	return err
}
