package rpc

import (
	"context"
	"net"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpctracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcvalidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	usersv1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/auth"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type ServerOpts struct { //nolint
	Addr            string
	AuthTokenSecret string
	BaseInputURL    string
	BaseOutputURL   string
	RTMPURL         string
	Manager         *manager.Manager
	Ds              *ds.Datastore
	Users           usersv1.UserServiceClient
	Accounts        accountsv1.AccountServiceClient
	Emitter         emitterv1.EmitterServiceClient
	EventBus        *eventbus.EventBus
	Logger          *logrus.Entry
}

type Server struct { //nolint
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
	emitter         emitterv1.EmitterServiceClient
	manager         *manager.Manager
	logger          *logrus.Entry
	validator       *requestValidator
	eb              *eventbus.EventBus
}

func NewServer(opts *ServerOpts) (*Server, error) {
	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpctracing.UnaryServerInterceptor(
				grpctracing.WithTracer(opentracing.GlobalTracer()),
				grpctracing.WithFilterFunc(tracingFilter),
			),
			grpcprometheus.UnaryServerInterceptor,
			grpclogrus.UnaryServerInterceptor(opts.Logger, grpclogrus.WithDecider(logrusFilter)),
			grpcvalidator.UnaryServerInterceptor(),
		)),
	}
	grpcSrv := grpc.NewServer(grpcOpts...)

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	validator, err := newRequestValidator()
	if err != nil {
		return nil, err
	}

	srv := &Server{
		logger:          opts.Logger,
		addr:            opts.Addr,
		authTokenSecret: opts.AuthTokenSecret,
		grpc:            grpcSrv,
		listen:          listen,
		ds:              opts.Ds,
		users:           opts.Users,
		accounts:        opts.Accounts,
		emitter:         opts.Emitter,
		manager:         opts.Manager,
		baseInputURL:    opts.BaseInputURL,
		baseOutputURL:   opts.BaseOutputURL,
		rtmpURL:         opts.RTMPURL,
		validator:       validator,
		eb:              opts.EventBus,
	}

	healthService := health.NewServer()
	healthv1.RegisterHealthServer(grpcSrv, healthService)

	v1.RegisterStreamsServiceServer(grpcSrv, srv)
	reflection.Register(grpcSrv)

	return srv, nil
}

func (s *Server) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *Server) authenticate(ctx context.Context) (string, error) {
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
