package service

import (
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	profilesv1 "github.com/videocoin/cloud-api/profiles/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/dlock"
	"github.com/videocoin/cloud-pkg/grpcutil"
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

func NewService(cfg *Config) (*Service, error) {
	ds, err := ds.NewDatastore(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	dlock, err := dlock.New(cfg.RedisURI)
	if err != nil {
		return nil, err
	}

	conn, err := grpcutil.Connect(cfg.UsersRPCAddr, cfg.Logger.WithField("system", "userscli"))
	if err != nil {
		return nil, err
	}
	users := usersv1.NewUserServiceClient(conn)

	conn, err = grpcutil.Connect(cfg.AccountsRPCAddr, cfg.Logger.WithField("system", "accountcli"))
	if err != nil {
		return nil, err
	}
	accounts := accountsv1.NewAccountServiceClient(conn)

	conn, err = grpcutil.Connect(cfg.EmitterRPCAddr, cfg.Logger.WithField("system", "emittercli"))
	if err != nil {
		return nil, err
	}
	emitter := emitterv1.NewEmitterServiceClient(conn)

	conn, err = grpcutil.Connect(cfg.ProfilesRPCAddr, cfg.Logger.WithField("system", "profilescli"))
	if err != nil {
		return nil, err
	}
	profiles := profilesv1.NewProfilesServiceClient(conn)

	ebConfig := &eventbus.Config{
		URI:     cfg.MQURI,
		Name:    cfg.Name,
		Logger:  cfg.Logger.WithField("system", "eventbus"),
		Emitter: emitter,
		Users:   users,
	}
	eb, err := eventbus.New(ebConfig)
	if err != nil {
		return nil, err
	}

	managerOpts := &manager.Opts{
		Logger:            cfg.Logger.WithField("system", "manager"),
		Ds:                ds,
		Emitter:           emitter,
		Accounts:          accounts,
		Users:             users,
		EB:                eb,
		DLock:             dlock,
		MaxLiveStreamTime: cfg.MaxLiveStreamTime,
	}
	manager := manager.NewManager(managerOpts)

	rpcConfig := &rpc.RPCServerOpts{
		Logger:          cfg.Logger,
		Addr:            cfg.RPCAddr,
		Ds:              ds,
		Manager:         manager,
		Users:           users,
		Accounts:        accounts,
		Profiles:        profiles,
		BaseInputURL:    cfg.BaseInputURL,
		BaseOutputURL:   cfg.BaseOutputURL,
		RTMPURL:         cfg.RTMPURL,
		Emitter:         emitter,
		AuthTokenSecret: cfg.AuthTokenSecret,
		EventBus:        eb,
	}

	publicRPC, err := rpc.NewRPCServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	privateRPCConfig := &rpc.PrivateRPCServerOpts{
		Addr:     cfg.PrivateRPCAddr,
		Logger:   cfg.Logger.WithField("system", "privaterpc"),
		Manager:  manager,
		Emitter:  emitter,
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
