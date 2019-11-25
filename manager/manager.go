package manager

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/dlock"
	"github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
)

var ErrHitBalanceLimitation = errors.New("hit balance limitation")
var ErrEndStreamNotAllowed = errors.New("not allowed to end stream")

type ManagerOpts struct {
	Ds                *datastore.Datastore
	Logger            *logrus.Entry
	Emitter           emitterv1.EmitterServiceClient
	Accounts          accountsv1.AccountServiceClient
	Users             usersv1.UserServiceClient
	DLock             *dlock.Locker
	EB                *eventbus.EventBus
	MaxLiveStreamTime int64
}

type Manager struct {
	logger            *logrus.Entry
	ds                *datastore.Datastore
	emitter           emitterv1.EmitterServiceClient
	accounts          accountsv1.AccountServiceClient
	users             usersv1.UserServiceClient
	dlock             *dlock.Locker
	eb                *eventbus.EventBus
	sbTicker          *time.Ticker
	sbTimeout         time.Duration
	maxLiveStreamTime int64
}

func NewManager(opts *ManagerOpts) *Manager {
	sbTimeout := 10 * time.Second
	m := &Manager{
		logger:            opts.Logger,
		ds:                opts.Ds,
		emitter:           opts.Emitter,
		accounts:          opts.Accounts,
		dlock:             opts.DLock,
		eb:                opts.EB,
		sbTimeout:         sbTimeout,
		sbTicker:          time.NewTicker(sbTimeout),
		maxLiveStreamTime: opts.MaxLiveStreamTime,
	}

	m.eb.StreamStatusHandler = m.handleStreamStatus

	return m
}

func (m *Manager) StartBackgroundTasks() error {
	go m.startCheckStreamBalanceTask()
	return nil
}

func (m *Manager) StopBackgroundTasks() error {
	m.sbTicker.Stop()
	return nil
}
