package manager

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-pkg/dlock"
	ds "github.com/videocoin/cloud-streams/datastore"
)

type ManagerOpts struct {
	Ds      *ds.Datastore
	Logger  *logrus.Entry
	Emitter emitterv1.EmitterServiceClient
	DLock   *dlock.Locker
}

type Manager struct {
	logger    *logrus.Entry
	ds        *ds.Datastore
	emitter   emitterv1.EmitterServiceClient
	dlock     *dlock.Locker
	sbTicker  *time.Ticker
	sbTimeout time.Duration
}

func NewManager(opts *ManagerOpts) *Manager {
	sbTimeout := 10 * time.Second
	return &Manager{
		logger:    opts.Logger,
		ds:        opts.Ds,
		emitter:   opts.Emitter,
		dlock:     opts.DLock,
		sbTimeout: sbTimeout,
		sbTicker:  time.NewTicker(sbTimeout),
	}
}

func (m *Manager) StartBackgroundTasks() error {
	go m.startCheckStreamBalanceTask()
	return nil
}

func (m *Manager) StopBackgroundTasks() error {
	m.sbTicker.Stop()
	return nil
}

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

					if toVID.Int64() <= int64(1) {
						logger.Info("deposit")

						_, err := m.emitter.Deposit(emptyCtx, &emitterv1.DepositRequest{
							StreamId: stream.Id,
							UserId:   stream.UserId,
							To:       addr.Bytes(),
							Value:    big.NewInt(1000000000000000000).Bytes(),
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
