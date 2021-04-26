package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/utils"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetSubscriptionsForUser obtains all the subscriptions for a user
func (m *Model) GetSubscriptionsForUser(ctx context.Context, uid int64, search *string, page, perPage int64) ([]*models.Repository, error) {
	pg, ppg, err := utils.ScopePaginationInt(&page, &perPage)
	if err != nil {
		return nil, err
	}

	mods := []qm.QueryMod{
		qm.InnerJoin("users on subscriptions.user_id = users.id"),
		qm.Where("user_id = ?", uid),
		qm.Limit(ppg),
		qm.Offset(pg * ppg),
	}

	if search != nil {
		mods = append(mods, []qm.QueryMod{
			qm.InnerJoin("repositories on repositories.id = subscriptions.repository_id"),
			qm.Where("repositories.name LIKE ?", fmt.Sprintf("%%%s%%", strings.Replace(*search, "%", "%%", -1))),
		}...)
	}

	subs, err := models.Subscriptions(mods...).All(ctx, m.db)
	if err != nil {
		return nil, err
	}

	ids := []int64{}

	for _, sub := range subs {
		ids = append(ids, sub.RepositoryID)
	}

	return models.Repositories(models.RepositoryWhere.ID.IN(ids)).All(ctx, m.db)
}

// AddSubscriptionsForUser adds a series of repositories to a user's subscriptions list.
func (m *Model) AddSubscriptionsForUser(ctx context.Context, uid int64, repos []*models.Repository) error {
	u, err := m.FindUserByID(ctx, uid)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		sub := &models.Subscription{UserID: u.ID, RepositoryID: repo.ID}
		sub.Insert(ctx, m.db, boil.Infer()) // deliberately unchecked; this is to workaround a bug in sqlboiler
	}

	return nil
}

// RemoveSubscriptionForUser removes N subscriptions for the user
func (m *Model) RemoveSubscriptionForUser(ctx context.Context, uid int64, repos []*models.Repository) error {
	u, err := m.FindUserByID(ctx, uid)
	if err != nil {
		return err
	}

	ids := []int64{}

	for _, repo := range repos {
		ids = append(ids, repo.ID)
	}

	_, err = u.Subscriptions(models.SubscriptionWhere.RepositoryID.IN(ids)).DeleteAll(ctx, m.db)
	return err
}
