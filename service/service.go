package service

import (
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"github.com/videocoin/cloud-streams/rpc"
	"google.golang.org/grpc"
)

type Service struct {
	cfg        *Config
	rpc        *rpc.RpcServer
	privateRPC *rpc.PrivateRPCServer
	eb         *eventbus.EventBus
}

func NewService(cfg *Config) (*Service, error) {
	ds, err := ds.NewDatastore(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	manager := manager.NewManager(
		&manager.ManagerOpts{
			Ds:     ds,
			Logger: cfg.Logger.WithField("system", "manager"),
		})

	ulogger := cfg.Logger.WithField("system", "userscli")
	uGrpcDialOpts := grpcutil.ClientDialOptsWithRetry(ulogger)
	usersConn, err := grpc.Dial(cfg.UsersRPCAddr, uGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	users := usersv1.NewUserServiceClient(usersConn)

	alogger := cfg.Logger.WithField("system", "accountcli")
	aGrpcDialOpts := grpcutil.ClientDialOptsWithRetry(alogger)
	accountsConn, err := grpc.Dial(cfg.AccountsRPCAddr, aGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	accounts := accountsv1.NewAccountServiceClient(accountsConn)

	elogger := cfg.Logger.WithField("system", "emittercli")
	eGrpcDialOpts := grpcutil.ClientDialOptsWithRetry(elogger)
	emitterConn, err := grpc.Dial(cfg.EmitterRPCAddr, eGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	emitter := emitterv1.NewEmitterServiceClient(emitterConn)

	ebConfig := &eventbus.Config{
		URI:    cfg.MQURI,
		Name:   cfg.Name,
		Logger: cfg.Logger.WithField("system", "eventbus"),
		DM:     manager,
	}
	eb, err := eventbus.New(ebConfig)
	if err != nil {
		return nil, err
	}

	rpcConfig := &rpc.RpcServerOpts{
		Logger:          cfg.Logger,
		Addr:            cfg.RPCAddr,
		Ds:              ds,
		Manager:         manager,
		Users:           users,
		Accounts:        accounts,
		BaseInputURL:    cfg.BaseInputURL,
		BaseOutputURL:   cfg.BaseOutputURL,
		RTMPURL:         cfg.RTMPURL,
		Emitter:         emitter,
		AuthTokenSecret: cfg.AuthTokenSecret,
		EventBus:        eb,
	}

	publicRPC, err := rpc.NewRpcServer(rpcConfig)
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
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.rpc.Start()
	go s.privateRPC.Start()
	go s.eb.Start()
	return nil
}

func (s *Service) Stop() error {
	s.eb.Stop()
	return nil
}
