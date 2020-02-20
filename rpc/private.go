package rpc

import (
	"context"
	"net"
	"time"

	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-api/rpc"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-streams/datastore"
	ds "github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/eventbus"
	"github.com/videocoin/cloud-streams/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type PrivateRPCServerOpts struct {
	Addr     string
	Logger   *logrus.Entry
	Manager  *manager.Manager
	Emitter  emitterv1.EmitterServiceClient
	EventBus *eventbus.EventBus
}

type PrivateRPCServer struct {
	addr    string
	logger  *logrus.Entry
	grpc    *grpc.Server
	listen  net.Listener
	manager *manager.Manager
	emitter emitterv1.EmitterServiceClient
	eb      *eventbus.EventBus
}

func NewPrivateRPCServer(opts *PrivateRPCServerOpts) (*PrivateRPCServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	rpcServer := &PrivateRPCServer{
		addr:    opts.Addr,
		logger:  opts.Logger.WithField("system", "privaterpc"),
		grpc:    grpcServer,
		listen:  listen,
		manager: opts.Manager,
		emitter: opts.Emitter,
		eb:      opts.EventBus,
	}

	privatev1.RegisterStreamsServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *PrivateRPCServer) Start() error {
	s.logger.Infof("starting private rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *PrivateRPCServer) Get(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	stream, err := s.manager.GetStreamByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *PrivateRPCServer) Publish(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	stream, err := s.manager.GetStreamByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
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

	streamResponse, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.logger.Error(s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		streamResponse.ID))

	return streamResponse, nil
}

func (s *PrivateRPCServer) PublishDone(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	_, err := s.manager.GetStreamByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	stream, err := s.manager.StopStream(ctx, req.Id, "", v1.StreamStatusCompleted)
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

	err = s.manager.EndStream(ctx, stream)
	if err != nil {
		logFailedTo(logger, "end stream", err)
		return nil, rpc.ErrRpcInternal
	}

	err = s.manager.EndStream(ctx, stream)
	if err != nil {
		logger.Errorf("failed to end stream: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *PrivateRPCServer) Complete(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	stream, err := s.manager.GetStreamByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	err = s.manager.CompleteStream(ctx, stream)
	if err != nil {
		logger.Errorf("failed to complete stream: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	go func(stream *ds.Stream) {
		logger.Info("waiting to end stream")

		time.Sleep(time.Second * 10)

		logger.Info("end stream")

		err = s.manager.EndStream(context.Background(), stream)
		if err != nil {
			logger.Errorf("failed to end stream: %s", err)
		}
	}(stream)

	streamResponse, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *PrivateRPCServer) Run(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	stream, err := s.manager.RunStream(ctx, req.Id, "")
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrHitBalanceLimitation {
			return nil, rpc.NewRpcPermissionError(err.Error())
		}

		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}

func (s *PrivateRPCServer) Stop(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	stream, err := s.manager.GetStreamByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	ss := v1.StreamStatusCancelled

	if stream.InputType == v1.InputTypeRTMP || stream.InputType == v1.InputTypeWebRTC {
		if stream.Status == v1.StreamStatusReady {
			ss = v1.StreamStatusCompleted
		}
	}

	stream, err = s.manager.StopStream(ctx, req.Id, "", ss)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrEndStreamNotAllowed {
			return nil, rpc.ErrRpcBadRequest
		}

		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}

func (s *PrivateRPCServer) UpdateStatus(ctx context.Context, req *privatev1.UpdateStatusRequest) (*privatev1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.ID)
	logger := s.logger.WithField("id", req.ID)

	stream, err := s.manager.GetStreamByID(ctx, req.ID)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream", err)
		return nil, rpc.ErrRpcInternal
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}

	err = s.manager.UpdateStream(ctx, stream, updates)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}

func toStreamResponsePrivate(stream *ds.Stream) (*privatev1.StreamResponse, error) {
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

	return resp, nil
}
