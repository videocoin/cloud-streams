package rpc

import (
	"context"
	"net"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-api/rpc"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type PrivateRPCServerOpts struct {
	Addr    string
	Logger  *logrus.Entry
	Manager *manager.Manager
}

type PrivateRPCServer struct {
	addr    string
	logger  *logrus.Entry
	grpc    *grpc.Server
	listen  net.Listener
	manager *manager.Manager
}

func NewPrivateRPCServer(opts *PrivateRPCServerOpts) (*PrivateRPCServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)

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
	}

	privatev1.RegisterStreamsServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *PrivateRPCServer) Start() error {
	s.logger.Infof("starting private rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *PrivateRPCServer) Health(ctx context.Context, req *protoempty.Empty) (*rpc.HealthStatus, error) {
	return &rpc.HealthStatus{Status: "OK"}, nil
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

	streamResponse, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func toStreamResponse(stream *v1.Stream) (*privatev1.StreamResponse, error) {
	resp := new(privatev1.StreamResponse)
	if err := copier.Copy(resp, stream); err != nil {
		return nil, err
	}

	resp.ID = stream.Id
	resp.InputURL = stream.InputUrl
	resp.OutputURL = stream.OutputUrl
	resp.ProfileID = stream.ProfileId

	return resp, nil
}
