package rpc

import (
	"context"
	"net"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	profilesv1 "github.com/videocoin/cloud-api/profiles/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/auth"
	"github.com/videocoin/cloud-pkg/grpcutil"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type RpcServerOpts struct {
	Addr            string
	AuthTokenSecret string
	BaseInputURL    string
	BaseOutputURL   string
	RTMPURL         string
	Manager         *manager.Manager
	Ds              *ds.Datastore
	Users           usersv1.UserServiceClient
	Accounts        accountsv1.AccountServiceClient
	Profiles        profilesv1.ProfilesServiceClient
	Emitter         emitterv1.EmitterServiceClient
	EventBus        *eventbus.EventBus
	Logger          *logrus.Entry
}

type RpcServer struct {
	addr            string
	authTokenSecret string
	baseInputURL    string
	baseOutputURL   string
	rtmpURL         string
	grpc            *grpc.Server
	listen          net.Listener
	ds              *ds.Datastore
	users           usersv1.UserServiceClient
	accounts        accountsv1.AccountServiceClient
	profiles        profilesv1.ProfilesServiceClient
	emitter         emitterv1.EmitterServiceClient
	manager         *manager.Manager
	logger          *logrus.Entry
	validator       *requestValidator
	eb              *eventbus.EventBus
}

func NewRpcServer(opts *RpcServerOpts) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}
	validator, err := newRequestValidator()
	if err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		addr:            opts.Addr,
		authTokenSecret: opts.AuthTokenSecret,
		grpc:            grpcServer,
		listen:          listen,
		ds:              opts.Ds,
		users:           opts.Users,
		accounts:        opts.Accounts,
		profiles:        opts.Profiles,
		emitter:         opts.Emitter,
		manager:         opts.Manager,
		baseInputURL:    opts.BaseInputURL,
		baseOutputURL:   opts.BaseOutputURL,
		rtmpURL:         opts.RTMPURL,

		logger:    opts.Logger.WithField("system", "rpc"),
		validator: validator,
		eb:        opts.EventBus,
	}

	v1.RegisterStreamServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RpcServer) authenticate(ctx context.Context) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "authenticate")
	defer span.Finish()

	ctx = auth.NewContextWithSecretKey(ctx, s.authTokenSecret)
	ctx, jwtToken, err := auth.AuthFromContext(ctx)
	if err != nil {
		s.logger.Warningf("failed to auth from context: %s", err)
		return "", rpc.ErrRpcUnauthenticated
	}

	tokenType, ok := auth.TypeFromContext(ctx)
	if ok {
		if usersv1.TokenType(tokenType) == usersv1.TokenTypeAPI {
			_, err := s.users.GetApiToken(context.Background(), &usersv1.ApiTokenRequest{Token: jwtToken})
			if err != nil {
				s.logger.Errorf("failed to get api token: %s", err)
				return "", rpc.ErrRpcUnauthenticated
			}
		}
	}

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		s.logger.Warningf("failed to get user id from context: %s", err)
		return "", rpc.ErrRpcUnauthenticated
	}

	return userID, nil
}
