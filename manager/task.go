package manager

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
)

func (m *Manager) startCheckStreamBalanceTask() {
	for range m.sbTicker.C {
		lockKey := "streams/tasks/check-stream-balance/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			m.logger.Error(err)
			return
		}
		if lock == nil {
			m.logger.Errorf("failed to obtain lock %s", lockKey)
			return
		}

		m.logger.Info("checking stream contract balance")

		ctx := context.Background()
		streams, err := m.ds.Stream.StatusReadyList(ctx)
		if err != nil {
			m.logger.Error(err)

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
			}

			continue
		}

		emptyCtx := context.Background()
		for _, stream := range streams {
			if stream.StreamContractAddress != "" && common.IsHexAddress(stream.StreamContractAddress) {
				addr := common.HexToAddress(stream.StreamContractAddress)

				logger := m.logger.
					WithField("user_id", stream.ID).
					WithField("to", stream.StreamContractAddress)

				toBalance, err := m.emitter.GetBalance(
					emptyCtx,
					&emitterv1.BalanceRequest{Address: addr.Bytes()},
				)
				if err != nil {
					logger.Error(err)

					unlockErr := lock.Unlock()
					if unlockErr != nil {
						logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
					}

					continue
				}

				toBalanceValue := new(big.Int).SetBytes(toBalance.Value)
				toVID := new(big.Int).Div(toBalanceValue, big.NewInt(1000000000000000000))

				logger.Infof("balance is %d VID", toVID.Int64())
				logger = logger.WithField("to_balance", toVID.Int64())

				if toVID.Int64() <= int64(1) {
					logger.Info("deposit")

					_, err := m.emitter.Deposit(emptyCtx, &emitterv1.DepositRequest{
						StreamId: stream.ID,
						UserId:   stream.UserID,
						To:       addr.Bytes(),
						Value:    big.NewInt(1000000000000000000).Bytes(),
					})
					if err != nil {
						logger.Error(err)

						unlockErr := lock.Unlock()
						if unlockErr != nil {
							logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
						}

						continue
					}
				}
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
		}
	}
}

func (m *Manager) startCheckStreamAliveTask() {
	for range m.sbTicker.C {
		lockKey := "streams/tasks/check-stream-alive/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			m.logger.Error(err)
			return
		}
		if lock == nil {
			m.logger.Errorf("failed to obtain lock %s", lockKey)
			return
		}

		m.logger.Info("checking stream is alive")

		emptyCtx := context.Background()

		streams, err := m.ds.Stream.StatusReadyList(emptyCtx)
		if err != nil {
			m.logger.Error(err)

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
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
						logger.Errorf("failed to complete stream: %s", err)
					}

					logger.Info("stream has been completed")
				}
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
		}
	}
}

func (m *Manager) startRemoveCompletedTask() {
	for range m.sbTicker.C {
		lockKey := "streams/tasks/remove-completed-task/lock"
		lock, err := m.dlock.Obtain(lockKey)
		if err != nil {
			m.logger.Error(err)
			return
		}
		if lock == nil {
			m.logger.Errorf("failed to obtain lock %s", lockKey)
			return
		}

		m.logger.Info("getting completed streams")

		emptyCtx := context.Background()
		streams, err := m.ds.Stream.ListForDeletion(emptyCtx)
		if err != nil {
			m.logger.Error(err)

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
			}

			continue
		}

		for _, stream := range streams {
			m.logger.WithField("stream_id", stream.ID).Info("marking stream as deleted")
			err := m.ds.Stream.MarkAsDeleted(emptyCtx, stream)
			if err != nil {
				m.logger.Errorf("failed to mark as deleted: %s", err)

				unlockErr := lock.Unlock()
				if unlockErr != nil {
					m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
				}

				continue
			}

			err = m.eb.EmitDeleteStreamContent(emptyCtx, stream.ID)
			if err != nil {
				m.logger.Errorf("failed to emit delete stream content: %s", err)

				unlockErr := lock.Unlock()
				if unlockErr != nil {
					m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
				}

				continue
			}
		}

		unlockErr := lock.Unlock()
		if unlockErr != nil {
			m.logger.Infof("failed to unlock %s: %s", lockKey, unlockErr)
		}
	}
}
