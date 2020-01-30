package manager

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	profilesv1 "github.com/videocoin/cloud-api/profiles/manager/v1"
	"github.com/videocoin/cloud-pkg/dlock"
)

func (m *Manager) startCheckStreamBalanceTask() error {
	for {
		select {
		case <-m.sbTicker.C:
			lockKey := "streams/tasks/check-stream-balance/lock"
			lock, err := m.dlock.Obtain(lockKey)
			if err != nil {
				return err
			}
			if lock == nil {
				m.logger.Errorf("failed to obtain lock %s", lockKey)
				return dlock.ErrObtainLock
			}

			defer lock.Unlock()

			m.logger.Info("checking stream contract balance")

			ctx := context.Background()
			streams, err := m.ds.Stream.StatusReadyList(ctx)
			if err != nil {
				m.logger.Error(err)
				continue
			}

			emptyCtx := context.Background()
			for _, stream := range streams {
				if stream.StreamContractAddress != "" && common.IsHexAddress(stream.StreamContractAddress) {
					addr := common.HexToAddress(stream.StreamContractAddress)

					logger := m.logger.
						WithField("user_id", stream.Id).
						WithField("to", stream.StreamContractAddress)

					toBalance, err := m.emitter.GetBalance(
						emptyCtx,
						&emitterv1.BalanceRequest{Address: addr.Bytes()},
					)
					if err != nil {
						logger.Error(err)
						continue
					}

					toBalanceValue := new(big.Int).SetBytes(toBalance.Value)
					toVID := new(big.Int).Div(toBalanceValue, big.NewInt(1000000000000000000))

					logger.Infof("balance is %d VID", toVID.Int64())
					logger = logger.WithField("to_balance", toVID.Int64())

					profile, err := m.profiles.Get(ctx, &profilesv1.ProfileRequest{
						Id: stream.ProfileId,
					})
					if err != nil {
						logger.Error(err)
						continue
					}

					i, e := big.NewInt(10), big.NewInt(19)
					deposit := i.Exp(i, e, nil).Bytes()
					if profile.Deposit != "" {
						d := new(big.Int)
						d.SetString(profile.Deposit, 10)
						if d == nil {
							logger.Errorf("Can't parse deposit %s", profile.Deposit)
							continue
						}
						deposit = d.Bytes()
					}
					reward := big.NewInt(10000000000000000)
					if profile.Reward != "" {
						r := new(big.Int)
						r.SetString(profile.Reward, 10)
						if r == nil {
							logger.Errorf("Can't parse reward %s", profile.Reward)
							continue
						}
						reward = r

					}
					cmp := toBalanceValue.Cmp(reward)
					if cmp <= 0 {
						logger.Info("deposit")

						_, err := m.emitter.Deposit(emptyCtx, &emitterv1.DepositRequest{
							StreamId: stream.Id,
							UserId:   stream.UserId,
							To:       addr.Bytes(),
							Value:    deposit,
						})
						if err != nil {
							logger.Error(err)
							continue
						}
					}
				}
			}
		}
	}

	return nil
}

func (m *Manager) startCheckStreamAliveTask() error {
	for {
		select {
		case <-m.sbTicker.C:
			lockKey := "streams/tasks/check-stream-alive/lock"
			lock, err := m.dlock.Obtain(lockKey)
			if err != nil {
				return err
			}
			if lock == nil {
				m.logger.Errorf("failed to obtain lock %s", lockKey)
				return dlock.ErrObtainLock
			}

			defer lock.Unlock()

			m.logger.Info("checking stream is alive")

			emptyCtx := context.Background()

			streams, err := m.ds.Stream.StatusReadyList(emptyCtx)
			if err != nil {
				m.logger.Error(err)
				continue
			}

			for _, stream := range streams {
				if stream.ReadyAt != nil {
					completedAt := stream.ReadyAt.Add(m.maxLiveStreamTime)
					if completedAt.Before(time.Now()) {
						logger := m.logger.WithField("id", stream.Id)
						logger.Info("completing stream")

						_, err := m.StopStream(emptyCtx, stream.Id, "")
						if err != nil {
							logger.Errorf("failed to complete stream: %s", err)
						}

						logger.Info("stream has been completed")
					}
				}
			}
		}
	}

	return nil
}
