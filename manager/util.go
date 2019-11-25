package manager

import (
	"context"
	"fmt"
	"math/big"

	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
)

func (m *Manager) checkBalance(ctx context.Context, userID string) error {
	account, err := m.accounts.GetByOwner(ctx, &accountsv1.AccountRequest{OwnerId: userID})
	if err != nil {
		return fmt.Errorf("failed to get account: %s", err)
	}

	balance, ok := new(big.Int).SetString(account.Balance, 10)
	balanceVID := new(big.Int).Div(balance, big.NewInt(1000000000000000000))
	if !ok || balanceVID.Cmp(big.NewInt(20)) == -1 {
		return ErrHitBalanceLimitation
	}

	return nil
}

func isRemovable(stream *v1.Stream) bool {
	return stream.Status == v1.StreamStatusNew ||
		stream.Status == v1.StreamStatusCompleted ||
		stream.Status == v1.StreamStatusCancelled ||
		stream.Status == v1.StreamStatusFailed
}
