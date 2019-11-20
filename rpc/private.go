package rpc

import (
	"context"
	"net"

	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-api/rpc"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-streams/datastore"
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

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		streamResponse.ID)

	return streamResponse, nil
}

func (s *PrivateRPCServer) PublishDone(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
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

	if stream.Status == v1.StreamStatusCompleted {
		streamResponse, err := toStreamResponsePrivate(stream)
		if err != nil {
			logFailedTo(logger, "", err)
			return nil, rpc.ErrRpcInternal
		}

		return streamResponse, nil
	}

	_, err = s.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
		UserId:                stream.UserId,
		StreamContractId:      stream.StreamContractId,
		StreamContractAddress: stream.StreamContractAddress,
	})
	if err != nil {
		logFailedTo(logger, "end stream", err)
	}

	err = s.manager.UpdateStream(
		ctx,
		stream,
		map[string]interface{}{"status": v1.StreamStatusCompleted},
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		streamResponse.ID)

	return streamResponse, nil
}

func (s *PrivateRPCServer) Run(ctx context.Context, req *privatev1.StreamRequest) (*privatev1.StreamResponse, error) {
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
		map[string]interface{}{"status": v1.StreamStatusPreparing},
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	_, err = s.emitter.InitStream(ctx, &emitterv1.InitStreamRequest{
		StreamId:         stream.Id,
		UserId:           stream.UserId,
		StreamContractId: stream.StreamContractId,
		ProfilesIds:      []string{stream.ProfileId},
	})
	if err != nil {
		logFailedTo(logger, "init stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.Id)

	return streamResponse, nil
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

	if stream.Status < v1.StreamStatusPrepared {
		return nil, rpc.ErrRpcBadRequest
	}

	if stream.Status == v1.StreamStatusCompleted {
		// nothing to do since it was already stopped
		streamResponse, err := toStreamResponsePrivate(stream)
		if err != nil {
			logFailedTo(logger, "", err)
			return nil, rpc.ErrRpcInternal
		}

		return streamResponse, nil
	}

	_, err = s.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
		UserId:                stream.UserId,
		StreamContractId:      stream.StreamContractId,
		StreamContractAddress: stream.StreamContractAddress,
	})

	if err != nil {
		logFailedTo(logger, "end stream", err)
	}

	err = s.manager.UpdateStream(
		ctx,
		stream,
		map[string]interface{}{"status": v1.StreamStatusCompleted},
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponsePrivate(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.Id)

	return streamResponse, nil
}

func toStreamResponsePrivate(stream *v1.Stream) (*privatev1.StreamResponse, error) {
	resp := new(privatev1.StreamResponse)
	if err := copier.Copy(resp, stream); err != nil {
		return nil, err
	}

	resp.ID = stream.Id
	resp.InputURL = stream.InputUrl
	resp.OutputURL = stream.OutputUrl
	resp.ProfileID = stream.ProfileId
	resp.StreamContractID = stream.StreamContractId
	resp.StreamContractAddress = stream.StreamContractAddress

	return resp, nil
}
