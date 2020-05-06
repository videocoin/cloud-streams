package manager

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	billingv1 "github.com/videocoin/cloud-api/billing/private/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
)

func (m *Manager) startCheckStreamBalanceTask() {
	for range m.sbTicker.C {
		logger := m.logger.WithField("bg_task", "CheckStreamBalance")

		lockKey := "streams/tasks/check-stream-balance/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			logger.WithError(err).WithField("key", lockKey).Error("failed to obtain lock")
			return
		}
		if lock == nil {
			logger.WithField("key", lockKey).Error("failed to obtain lock")
			return
		}

		logger.Info("checking stream contract balance")

		ctx := context.Background()
		streams, err := m.ds.Stream.StatusReadyList(ctx)
		if err != nil {
			logger.WithError(err).Error("failed to get ready streams list")

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
			}

			continue
		}

		logger.Infof("ready streams count - %d", len(streams))

		emptyCtx := context.Background()
		for _, stream := range streams {
			logger := logger.WithField("stream_id", stream.ID)

			if stream.StreamContractAddress != "" && common.IsHexAddress(stream.StreamContractAddress) {
				addr := common.HexToAddress(stream.StreamContractAddress)

				logger := logger.
					WithField("user_id", stream.ID).
					WithField("to", stream.StreamContractAddress)

				logger.Info("get stream balance")

				toBalance, err := m.emitter.GetBalance(
					emptyCtx,
					&emitterv1.BalanceRequest{Address: addr.Bytes()},
				)
				if err != nil {
					logger.WithError(err).Error("failed to get balance")

					unlockErr := lock.Unlock()
					if unlockErr != nil {
						logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
					}

					continue
				}

				toBalanceValue := new(big.Int).SetBytes(toBalance.Value)
				toVID := new(big.Int).Div(toBalanceValue, big.NewInt(1000000000000000000))

				logger.Infof("balance is %d VID", toVID.Int64())
				logger = logger.WithField("to_balance", toVID.Int64())

				if toVID.Int64() <= int64(4) {
					logger.Info("deposit")

					_, err := m.emitter.Deposit(emptyCtx, &emitterv1.DepositRequest{
						StreamId: stream.ID,
						UserId:   stream.UserID,
						To:       addr.Bytes(),
						Value:    big.NewInt(1000000000000000000).Bytes(),
					})
					if err != nil {
						logger.WithError(err).Error("failed to deposit")

						unlockErr := lock.Unlock()
						if unlockErr != nil {
							logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
						}

						continue
					}
				}
			} else {
				logger.Warning("no stream contract address")
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
		}
	}
}

func (m *Manager) startCheckStreamAliveTask() {
	for range m.sbTicker.C {
		lockKey := "streams/tasks/check-stream-alive/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			m.logger.WithError(err).WithField("key", lockKey).Error("failed to obtain lock")
			return
		}
		if lock == nil {
			m.logger.WithField("key", lockKey).Error("failed to obtain lock")
			return
		}

		m.logger.Info("checking stream is alive")

		emptyCtx := context.Background()

		streams, err := m.ds.Stream.StatusReadyList(emptyCtx)
		if err != nil {
			m.logger.WithError(err).Error("failed to get ready streams")

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
			}

			continue
		}

		for _, stream := range streams {
			if stream.ReadyAt != nil {
				completedAt := stream.ReadyAt.Add(m.maxLiveStreamTime)
				if completedAt.Before(time.Now()) {
					logger := m.logger.WithField("id", stream.ID)
					logger.Info("completing stream")

					_, err := m.StopStream(emptyCtx, stream.ID, "", streamsv1.StreamStatusCompleted)
					if err != nil {
						logger.WithError(err).Errorf("failed to stop stream")
					} else {
						logger.Info("stream has been completed")
					}
				}
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
		}
	}
}

func (m *Manager) startRemoveCompletedTask() {
	for range m.sbTicker.C {
		lockKey := "streams/tasks/remove-completed-task/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			m.logger.WithError(err).WithField("key", lockKey).Error("failed to obtain lock")
			return
		}
		if lock == nil {
			m.logger.WithField("key", lockKey).Error("failed to obtain lock")
			return
		}

		m.logger.Info("getting completed streams")

		emptyCtx := context.Background()
		streams, err := m.ds.Stream.ListForDeletion(emptyCtx)
		if err != nil {
			m.logger.WithError(err).Error("failed to list streams for deletion")

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
			}

			continue
		}

		for _, stream := range streams {
			m.logger.WithField("stream_id", stream.ID).Info("marking stream as deleted")
			err := m.ds.Stream.MarkAsDeleted(emptyCtx, stream)
			if err != nil {
				m.logger.WithError(err).Errorf("failed to mark streams as deleted")

				unlockErr := lock.Unlock()
				if unlockErr != nil {
					m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
				}

				continue
			}

			err = m.eb.EmitDeleteStreamContent(emptyCtx, stream.ID)
			if err != nil {
				m.logger.WithError(err).Errorf("failed to emit delete stream content")

				unlockErr := lock.Unlock()
				if unlockErr != nil {
					m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
				}

				continue
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
		}
	}
}

func (m *Manager) startStreamTotalCostTask() {
	for range m.stcTicker.C {
		lockKey := "streams/tasks/stream-total-cost/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			m.logger.WithError(err).WithField("key", lockKey).Error("failed to obtain lock")
			return
		}
		if lock == nil {
			m.logger.WithField("key", lockKey).Error("failed to obtain lock")
			return
		}

		m.logger.Info("getting charges")

		ctx := context.Background()
		charges, err := m.billing.GetCharges(ctx, &billingv1.ChargesRequest{})
		if err != nil {
			m.logger.WithError(err).Errorf("failed to get charges")
			unlockErr := lock.Unlock()
			if unlockErr != nil {
				m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
			}

			continue
		}

		for _, charge := range charges.Items {
			stream, err := m.ds.Stream.Get(ctx, charge.StreamID)
			if err != nil {
				continue
			}

			updates := map[string]interface{}{"total_cost": charge.TotalCost}
			err = m.ds.Stream.Update(ctx, stream, updates)
			if err != nil {
				continue
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
		}
	}
}

func (m *Manager) startCheckBalanceForLiveStreamsTask() {
	for range m.cblsTicker.C {
		lockKey := "streams/tasks/check-balance-task/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			m.logger.WithError(err).WithField("key", lockKey).Error("failed to obtain lock")
			return
		}
		if lock == nil {
			m.logger.WithField("key", lockKey).Error("failed to obtain lock")
			return
		}

		m.logger.Info("checking balance for live streams")

		emptyCtx := context.Background()
		streams, err := m.ds.Stream.StatusReadyList(emptyCtx)
		if err != nil {
			m.logger.WithError(err).Error("failed to get ready streams")

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
			}

			continue
		}

		lowBalance := map[string]bool{}

		for _, stream := range streams {
			if isLowBalance, ok := lowBalance[stream.UserID]; ok {
				if isLowBalance {
					_, err := m.StopStream(emptyCtx, stream.ID, stream.UserID, streamsv1.StreamStatusCompleted)
					if err != nil {
						m.logger.WithError(err).Errorf("failed to stop stream when low balance")
						continue
					}
				}
			}

			profile, err := m.billing.GetProfileByUserID(emptyCtx, &billingv1.ProfileRequest{UserID: stream.UserID})
			if err != nil {
				m.logger.WithError(err).Errorf("failed to get billing profile by user id")
				continue
			}

			if profile.Balance <= 0 {
				lowBalance[stream.UserID] = true
				_, err := m.StopStream(emptyCtx, stream.ID, stream.UserID, streamsv1.StreamStatusCompleted)
				if err != nil {
					m.logger.WithError(err).Error("failed to stop stream when low balance")
					continue
				}
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			m.logger.WithError(unlockErr).WithField("key", lockKey).Warning("failed to unlock")
		}
	}
}
