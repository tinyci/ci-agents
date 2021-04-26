package db

import (
	"context"
	"strings"

	"github.com/tinyci/ci-agents/db/models"
	"github.com/tinyci/ci-agents/types"
)

// AddCapabilityToUser adds a capability to a user account.
func (m *Model) AddCapabilityToUser(ctx context.Context, u *models.User, cap types.Capability) error {
	return u.AddUserCapabilities(ctx, m.db, true, &models.UserCapability{Name: string(cap)})
}

// RemoveCapabilityFromUser removes a capability from a user account.
func (m *Model) RemoveCapabilityFromUser(ctx context.Context, u *models.User, cap types.Capability) error {
	_, err := u.UserCapabilities(models.UserCapabilityWhere.Name.EQ(string(cap))).DeleteAll(ctx, m.db)
	return err
}

// GetCapabilities returns the capabilities the supplied user account has.
func (m *Model) GetCapabilities(ctx context.Context, u *models.User, fixedCaps map[string][]string) ([]types.Capability, error) {
	caps := map[string]struct{}{}

	if fc, ok := fixedCaps[u.Username]; ok {
		for _, cap := range fc {
			caps[cap] = struct{}{}
		}
	}

	dbCaps, err := u.UserCapabilities().All(ctx, m.db)
	if err != nil {
		return nil, err
	}

	for _, cap := range dbCaps {
		caps[cap.Name] = struct{}{}
	}

	realCaps := []types.Capability{}

	for cap := range caps {
		realCaps = append(realCaps, types.Capability(cap))
	}

	return realCaps, nil
}

// HasCapability returns true if the user is capable of performing the operation.
func (m *Model) HasCapability(ctx context.Context, u *models.User, cap types.Capability, fixedCaps map[string][]string) (bool, error) {
	// if we have fixed caps, we consult that table only; these are overrides for
	// users that exist within the configuration file for the datasvc.
	if caps, ok := fixedCaps[u.Username]; ok {
		for _, thisCap := range caps {
			if strings.TrimSpace(thisCap) == strings.TrimSpace(string(cap)) {
				return true, nil
			}
		}

		return false, nil
	}

	return u.UserCapabilities(models.UserCapabilityWhere.Name.EQ(string(cap))).Exists(ctx, m.db)
}
