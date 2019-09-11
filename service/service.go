package service

import (
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/manager"
	"github.com/videocoin/cloud-streams/rpc"
	"google.golang.org/grpc"
)

type Service struct {
	cfg *Config
	rpc *rpc.RpcServer
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

	rpcConfig := &rpc.RpcServerOpts{
		Addr:            cfg.RPCAddr,
		Ds:              ds,
		Manager:         manager,
		Accounts:        accounts,
		BaseInputURL:    cfg.BaseInputURL,
		BaseOutputURL:   cfg.BaseOutputURL,
		Emitter:         emitter,
		AuthTokenSecret: cfg.AuthTokenSecret,
		Logger:          cfg.Logger,
	}

	rpc, err := rpc.NewRpcServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg: cfg,
		rpc: rpc,
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.rpc.Start()
	return nil
}

func (s *Service) Stop() error {
	return nil
}
