package rpc

import (
	"context"
	"net"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/auth"
	"github.com/videocoin/cloud-pkg/grpcutil"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RpcServerOpts struct {
	Addr            string
	AuthTokenSecret string
	BaseInputURL    string
	BaseOutputURL   string
	Manager         *manager.Manager
	Ds              *ds.Datastore
	Accounts        accountsv1.AccountServiceClient
	Emitter         emitterv1.EmitterServiceClient
	EventBus        *eventbus.EventBus
	Logger          *logrus.Entry
}

type RpcServer struct {
	addr            string
	authTokenSecret string
	baseInputURL    string
	baseOutputURL   string
	grpc            *grpc.Server
	listen          net.Listener
	ds              *ds.Datastore
	accounts        accountsv1.AccountServiceClient
	emitter         emitterv1.EmitterServiceClient
	manager         *manager.Manager
	logger          *logrus.Entry
	validator       *requestValidator
	eb              *eventbus.EventBus
}

func NewRpcServer(opts *RpcServerOpts) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		addr:            opts.Addr,
		authTokenSecret: opts.AuthTokenSecret,
		grpc:            grpcServer,
		listen:          listen,
		ds:              opts.Ds,
		accounts:        opts.Accounts,
		emitter:         opts.Emitter,
		manager:         opts.Manager,
		baseInputURL:    opts.BaseInputURL,
		baseOutputURL:   opts.BaseOutputURL,
		logger:          opts.Logger.WithField("system", "rpc"),
		validator:       newRequestValidator(),
		eb:              opts.EventBus,
	}

	v1.RegisterStreamServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RpcServer) Health(ctx context.Context, req *protoempty.Empty) (*rpc.HealthStatus, error) {
	return &rpc.HealthStatus{Status: "OK"}, nil
}

func (s *RpcServer) authenticate(ctx context.Context) (string, context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "authenticate")
	defer span.Finish()

	ctx = auth.NewContextWithSecretKey(ctx, s.authTokenSecret)
	ctx, err := auth.AuthFromContext(ctx)
	if err != nil {
		s.logger.Warningf("failed to auth from context: %s", err)
		return "", ctx, rpc.ErrRpcUnauthenticated
	}

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		s.logger.Warningf("failed to get user id from context: %s", err)
		return "", ctx, rpc.ErrRpcUnauthenticated
	}

	return userID, ctx, nil
}
