package manager

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	ds "github.com/videocoin/cloud-streams/datastore"
)

type ManagerOpts struct {
	Ds      *ds.Datastore
	Logger  *logrus.Entry
	Emitter emitterv1.EmitterServiceClient
}

type Manager struct {
	ds        *ds.Datastore
	emitter   emitterv1.EmitterServiceClient
	logger    *logrus.Entry
	sbTicker  *time.Ticker
	sbTimeout time.Duration
}

func NewManager(opts *ManagerOpts) *Manager {
	sbTimeout := 5 * time.Second
	return &Manager{
		logger:    opts.Logger,
		ds:        opts.Ds,
		emitter:   opts.Emitter,
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
			m.logger.Info("check stream contract balance")
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
					balance, err := m.emitter.GetBalance(emptyCtx, &emitterv1.BalanceRequest{Address: addr.Bytes()})
					if err != nil {
						m.logger.Error(err)
						continue
					}

					v := new(big.Int).SetBytes(balance.Value)
					vid := new(big.Int).Div(v, big.NewInt(1000000000000000000))
					if vid.Cmp(big.NewInt(1)) <= 1 {
						m.logger.
							WithField("user_id", stream.Id).
							WithField("to", stream.StreamContractAddress).
							Info("deposit")
						_, err := m.emitter.Deposit(emptyCtx, &emitterv1.DepositRequest{
							UserId: stream.UserId,
							To:     addr.Bytes(),
							Value:  big.NewInt(1000000000000000000).Bytes(),
						})
						if err != nil {
							m.logger.Error(err)
							continue
						}
					} else {
						m.logger.
							WithField("addr", stream.StreamContractAddress).
							Debugf("balance is %d", vid)
					}
				}
			}
		}
	}

	return nil
}
