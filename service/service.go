package service

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-pkg/dlock"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"github.com/videocoin/cloud-streams/rpc"
)

type Service struct {
	cfg  *Config
	rpc  *rpc.Server
	prpc *rpc.PrivateServer
	eb   *eventbus.EventBus
	dm   *manager.Manager
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	ds, err := ds.NewDatastore(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	dlock, err := dlock.New(cfg.RedisURI)
	if err != nil {
		return nil, err
	}

	sc, err := clientv1.NewServiceClientFromEnvconfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	eb, err := eventbus.New(&eventbus.Config{
		URI:     cfg.MQURI,
		Name:    cfg.Name,
		Logger:  ctxlogrus.Extract(ctx).WithField("system", "eventbus"),
		Emitter: sc.Emitter,
		Users:   sc.Users,
	})
	if err != nil {
		return nil, err
	}

	manager := manager.NewManager(&manager.Opts{
		Logger:            ctxlogrus.Extract(ctx).WithField("system", "manager"),
		Ds:                ds,
		Emitter:           sc.Emitter,
		Accounts:          sc.Accounts,
		Users:             sc.Users,
		Billing:           sc.Billing,
		EB:                eb,
		DLock:             dlock,
		MaxLiveStreamTime: cfg.MaxLiveStreamTime,
	})

	srv, err := rpc.NewServer(&rpc.ServerOpts{
		Logger:          ctxlogrus.Extract(ctx).WithField("system", "rpc"),
		Addr:            cfg.RPCAddr,
		Ds:              ds,
		Manager:         manager,
		Users:           sc.Users,
		Accounts:        sc.Accounts,
		BaseInputURL:    cfg.BaseInputURL,
		BaseOutputURL:   cfg.BaseOutputURL,
		RTMPURL:         cfg.RTMPURL,
		Emitter:         sc.Emitter,
		AuthTokenSecret: cfg.AuthTokenSecret,
		EventBus:        eb,
	})
	if err != nil {
		return nil, err
	}

	psrv, err := rpc.NewPrivateServer(&rpc.PrivateServerOpts{
		Addr:     cfg.PrivateRPCAddr,
		Logger:   ctxlogrus.Extract(ctx).WithField("system", "private-rpc"),
		Manager:  manager,
		Emitter:  sc.Emitter,
		EventBus: eb,
	})
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:  cfg,
		rpc:  srv,
		prpc: psrv,
		eb:   eb,
		dm:   manager,
	}

	return svc, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		errCh <- s.rpc.Start()
	}()

	go func() {
		errCh <- s.prpc.Start()
	}()

	go func() {
		errCh <- s.eb.Start()
	}()

	s.dm.StartBackgroundTasks()
}

func (s *Service) Stop() error {
	err := s.eb.Stop()
	if err != nil {
		return err
	}
	err = s.dm.StopBackgroundTasks()
	return err
}

func (s *Service) LoadFixtures(presetsRoot string) error {
	var presetsFiles []string

	err := filepath.Walk(presetsRoot, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			presetsFiles = append(presetsFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	datastore, err := ds.NewDatastore(s.cfg.DBURI)
	if err != nil {
		return err
	}

	m := &runtime.JSONPb{OrigName: true, EmitDefaults: true, EnumsAsInts: false}
	ctx := context.Background()
	profileIds := []string{}

	for _, file := range presetsFiles {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		profile := new(ds.Profile)
		err = m.Unmarshal(data, &profile)
		if err != nil {
			return err
		}

		_, err = datastore.Profile.Get(ctx, profile.ID)
		if err != nil {
			if err == ds.ErrProfileNotFound {
				_, createErr := datastore.Profile.Create(ctx, profile)
				if createErr != nil {
					return createErr
				}

				profileIds = append(profileIds, profile.ID)
				continue
			}

			return err
		}

		err = datastore.Profile.Update(ctx, profile)
		if err != nil {
			return err
		}

		profileIds = append(profileIds, profile.ID)
	}

	if len(profileIds) > 0 {
		err = datastore.Profile.DeleteAllExceptIds(ctx, profileIds)
		if err != nil {
			return err
		}
	}

	return nil
}
