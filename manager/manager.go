package manager

import (
	"github.com/sirupsen/logrus"
	ds "github.com/videocoin/cloud-streams/datastore"
)

type ManagerOpts struct {
	Ds     *ds.Datastore
	Logger *logrus.Entry
}

type Manager struct {
	ds     *ds.Datastore
	logger *logrus.Entry
}

func NewManager(opts *ManagerOpts) *Manager {
	return &Manager{
		ds:     opts.Ds,
		logger: opts.Logger,
	}
}
