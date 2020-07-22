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
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-api/rpc"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-streams/datastore"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type PrivateServerOpts struct {
	Addr     string
	Logger   *logrus.Entry
	Manager  *manager.Manager
	Emitter  emitterv1.EmitterServiceClient
	EventBus *eventbus.EventBus
}

type PrivateServer struct {
	addr    string
	logger  *logrus.Entry
	grpc    *grpc.Server
	listen  net.Listener
	manager *manager.Manager
	emitter emitterv1.EmitterServiceClient
	eb      *eventbus.EventBus
}

func NewPrivateServer(opts *PrivateServerOpts) (*PrivateServer, error) {
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
	grpcServer := grpc.NewServer(grpcOpts...)

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	rpcServer := &PrivateServer{
		addr:    opts.Addr,
		logger:  opts.Logger.WithField("system", "privaterpc"),
		grpc:    grpcServer,
		listen:  listen,
		manager: opts.Manager,
		emitter: opts.Emitter,
		eb:      opts.EventBus,
	}

	healthService := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthService)

	privatev1.RegisterStreamsServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *PrivateServer) Start() error {
	s.logger.Infof("starting private rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *PrivateServer) Get(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)
	stream, err := s.manager.GetStreamByID(otCtx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	profile, err := s.manager.GetProfileByID(otCtx, stream.ProfileID)
	if err != nil {
		logFailedTo(logger, "get profile", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream, profile)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *PrivateServer) Publish(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)
	stream, err := s.manager.GetStreamByID(otCtx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	profile, err := s.manager.GetProfileByID(otCtx, stream.ProfileID)
	if err != nil {
		logFailedTo(logger, "get profile", err)
		return nil, rpc.ErrRpcInternal
	}

	cost := profile.Spec.Cost / 60 * req.Duration

	logger.Infof("cost %f", cost)

	balance, err := s.manager.GetBalance(otCtx, stream.UserID)
	if err != nil {
		logFailedTo(logger, "get balance", err)
		return nil, rpc.ErrRpcInternal
	}

	logger.Infof("balance %f", balance)

	if cost > balance {
		return nil, rpc.ErrRpcBadRequest
	}

	err = s.manager.UpdateStream(
		ctx,
		stream,
		map[string]interface{}{
			"status":       v1.StreamStatusPending,
			"input_status": v1.InputStatusActive,
		},
	)
	if err != nil {
		logFailedTo(logger, "mark as publish", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream, nil)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go func() {
		err := s.eb.EmitUpdateStream(otCtx, streamResponse.ID)
		if err != nil {
			s.logger.Error(err)
		}
	}()

	return streamResponse, nil
}

func (s *PrivateServer) PublishDone(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)
	_, err := s.manager.GetStreamByID(otCtx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	stream, err := s.manager.StopStream(otCtx, req.Id, "", v1.StreamStatusCompleted)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrEndStreamNotAllowed {
			return nil, rpc.ErrRpcBadRequest
		}

		logger.Errorf("failed to stop stream: %s", err)

		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream, nil)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *PrivateServer) Complete(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)
	stream, err := s.manager.GetStreamByID(otCtx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	err = s.manager.CompleteStream(otCtx, stream)
	if err != nil {
		logger.Errorf("failed to complete stream: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream, nil)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *PrivateServer) Run(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)

	stream, err := s.manager.RunStream(otCtx, req.Id, "")
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrHitBalanceLimitation {
			return nil, rpc.NewRpcPermissionError(err.Error())
		}

		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponsePrivate(stream, nil)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}

func (s *PrivateServer) Stop(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)
	stream, err := s.manager.GetStreamByID(otCtx, req.Id)
	if err != nil {
		logFailedTo(logger, "get stream", err)

		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		return nil, rpc.ErrRpcInternal
	}

	ss := v1.StreamStatusCancelled

	if stream.InputType == v1.InputTypeRTMP || stream.InputType == v1.InputTypeWebRTC {
		if stream.Status == v1.StreamStatusReady {
			ss = v1.StreamStatusCompleted
		}
	}

	stream, err = s.manager.StopStream(otCtx, req.Id, "", ss)
	if err != nil {
		logFailedTo(logger, "stop stream", err)

		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrEndStreamNotAllowed {
			return nil, rpc.ErrRpcBadRequest
		}

		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponsePrivate(stream, nil)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}

func (s *PrivateServer) UpdateStatus(ctx context.Context, req *privatev1.UpdateStatusRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.ID)
	logger := s.logger.WithField("stream_id", req.ID)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)
	stream, err := s.manager.GetStreamByID(otCtx, req.ID)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	logger.WithField("status", req.Status).Info("upating status")

	updates := map[string]interface{}{"status": req.Status}
	err = s.manager.UpdateStream(otCtx, stream, updates)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponsePrivate(stream, nil)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}

func toStreamResponsePrivate(stream *ds.Stream, profile *ds.Profile) (*privatev1.StreamResponse, error) {
	resp := new(privatev1.StreamResponse)
	if err := copier.Copy(resp, stream); err != nil {
		return nil, err
	}

	resp.ID = stream.ID
	resp.UserID = stream.UserID
	resp.InputURL = stream.InputURL
	resp.OutputURL = stream.OutputURL
	resp.ProfileID = stream.ProfileID
	resp.StreamContractID = stream.StreamContractID
	resp.StreamContractAddress = stream.StreamContractAddress
	resp.InputType = stream.InputType
	resp.OutputType = stream.OutputType

	if profile != nil {
		resp.ProfileCost = profile.Spec.Cost
	}

	return resp, nil
}
