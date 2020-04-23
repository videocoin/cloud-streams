package service

import (
	"context"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-pkg/dlock"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"github.com/videocoin/cloud-streams/rpc"
)

type Service struct {
	cfg        *Config
	rpc        *rpc.RPCServer
	privateRPC *rpc.PrivateRPCServer
	eb         *eventbus.EventBus
	dm         *manager.Manager
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

	ebConfig := &eventbus.Config{
		URI:     cfg.MQURI,
		Name:    cfg.Name,
		Logger:  grpclogrus.Extract(ctx).WithField("system", "eventbus"),
		Emitter: sc.Emitter,
		Users:   sc.Users,
	}
	eb, err := eventbus.New(ebConfig)
	if err != nil {
		return nil, err
	}

	managerOpts := &manager.Opts{
		Logger:            grpclogrus.Extract(ctx).WithField("system", "manager"),
		Ds:                ds,
		Emitter:           sc.Emitter,
		Accounts:          sc.Accounts,
		Users:             sc.Users,
		Billing:           sc.Billing,
		EB:                eb,
		DLock:             dlock,
		MaxLiveStreamTime: cfg.MaxLiveStreamTime,
	}
	manager := manager.NewManager(managerOpts)

	rpcConfig := &rpc.RPCServerOpts{
		Logger:          grpclogrus.Extract(ctx).WithField("system", "rpc"),
		Addr:            cfg.RPCAddr,
		Ds:              ds,
		Manager:         manager,
		Users:           sc.Users,
		Accounts:        sc.Accounts,
		Profiles:        sc.Profiles,
		BaseInputURL:    cfg.BaseInputURL,
		BaseOutputURL:   cfg.BaseOutputURL,
		RTMPURL:         cfg.RTMPURL,
		Emitter:         sc.Emitter,
		AuthTokenSecret: cfg.AuthTokenSecret,
		EventBus:        eb,
	}

	publicRPC, err := rpc.NewRPCServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	privateRPCConfig := &rpc.PrivateRPCServerOpts{
		Addr:     cfg.PrivateRPCAddr,
		Logger:   grpclogrus.Extract(ctx).WithField("system", "privaterpc"),
		Manager:  manager,
		Emitter:  sc.Emitter,
		Profiles: sc.Profiles,
		EventBus: eb,
	}

	privateRPC, err := rpc.NewPrivateRPCServer(privateRPCConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:        cfg,
		rpc:        publicRPC,
		privateRPC: privateRPC,
		eb:         eb,
		dm:         manager,
	}

	return svc, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		errCh <- s.rpc.Start()
	}()

	go func() {
		errCh <- s.privateRPC.Start()
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
