package manager

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	billingv1 "github.com/videocoin/cloud-api/billing/private/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/dlock"
	"github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
)

var ErrHitBalanceLimitation = errors.New("account has insufficient funds to start streaming")
var ErrEndStreamNotAllowed = errors.New("not allowed to end stream")

type Opts struct {
	Ds                *datastore.Datastore
	Logger            *logrus.Entry
	Emitter           emitterv1.EmitterServiceClient
	Accounts          accountsv1.AccountServiceClient
	Users             usersv1.UserServiceClient
	Billing           billingv1.BillingServiceClient
	DLock             *dlock.Locker
	EB                *eventbus.EventBus
	MaxLiveStreamTime time.Duration
}

type Manager struct {
	logger            *logrus.Entry
	ds                *datastore.Datastore
	emitter           emitterv1.EmitterServiceClient
	accounts          accountsv1.AccountServiceClient
	users             usersv1.UserServiceClient
	billing           billingv1.BillingServiceClient
	dlock             *dlock.Locker
	eb                *eventbus.EventBus
	sbTicker          *time.Ticker
	sbTimeout         time.Duration
	maxLiveStreamTime time.Duration
}

func NewManager(opts *Opts) *Manager {
	sbTimeout := 10 * time.Second
	m := &Manager{
		logger:            opts.Logger,
		ds:                opts.Ds,
		emitter:           opts.Emitter,
		accounts:          opts.Accounts,
		users:             opts.Users,
		billing:           opts.Billing,
		dlock:             opts.DLock,
		eb:                opts.EB,
		sbTimeout:         sbTimeout,
		sbTicker:          time.NewTicker(sbTimeout),
		maxLiveStreamTime: opts.MaxLiveStreamTime,
	}

	m.eb.StreamStatusHandler = m.handleStreamStatus

	return m
}

func (m *Manager) StartBackgroundTasks() {
	go m.startCheckStreamBalanceTask()
	go m.startCheckStreamAliveTask()
	go m.startRemoveCompletedTask()
}

func (m *Manager) StopBackgroundTasks() error {
	m.sbTicker.Stop()
	return nil
}
