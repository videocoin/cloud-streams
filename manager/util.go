package manager

import (
	"context"
	"fmt"

	billingv1 "github.com/videocoin/cloud-api/billing/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	ds "github.com/videocoin/cloud-streams/datastore"
)

func (m *Manager) checkBalance(ctx context.Context, userID string) error {
	account, err := m.billing.GetProfileByUserID(ctx, &billingv1.ProfileRequest{UserID: userID})
	if err != nil {
		return fmt.Errorf("failed to get account: %s", err)
	}

	if account.Balance < 10 {
		return ErrHitBalanceLimitation
	}

	return nil
}

func isRemovable(stream *ds.Stream) bool {
	return stream.Status == v1.StreamStatusNew ||
		stream.Status == v1.StreamStatusCompleted ||
		stream.Status == v1.StreamStatusCancelled ||
		stream.Status == v1.StreamStatusFailed ||
		stream.Status == v1.StreamStatusDeleted
}
