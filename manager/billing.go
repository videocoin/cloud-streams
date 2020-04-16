package manager

import (
	"context"

	billingv1 "github.com/videocoin/cloud-api/billing/private/v1"
)

func (m *Manager) GetBalance(ctx context.Context, userID string) (float64, error) {
	profile, err := m.billing.GetProfileByUserID(ctx, &billingv1.ProfileRequest{UserID: userID})
	if err != nil {
		return 0, err
	}

	return profile.Balance, nil
}
