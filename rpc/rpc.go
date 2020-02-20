package rpc

import (
	"context"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/opentracing/opentracing-go"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	profilesv1 "github.com/videocoin/cloud-api/profiles/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/manager"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *RPCServer) Create(ctx context.Context, req *v1.CreateStreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("name", req.Name)
	span.SetTag("profile_id", req.ProfileId)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger := s.logger.WithField("user_id", userID)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	_, err = s.profiles.Get(ctx, &profilesv1.ProfileRequest{
		Id: req.ProfileId,
	})
	if err != nil {
		s.logger.Error(err)

		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.NotFound {
				respErr := &rpc.MultiValidationError{
					Errors: []*rpc.ValidationError{
						{
							Field:   "profile_id",
							Message: "profile id does not exist",
						},
					},
				}

				return nil, rpc.NewRpcValidationError(respErr)
			}
		}

		return nil, rpc.ErrRpcInternal
	}

	stream, err := s.manager.CreateStream(
		ctx,
		req.Name,
		req.ProfileId,
		userID,
		s.baseInputURL,
		s.baseOutputURL,
		s.rtmpURL,
		req.InputType,
		req.OutputType,
	)
	if err != nil {
		logFailedTo(logger, "create stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.logger.Error(s.eb.EmitCreateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.ID))

	streamResponse, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *RPCServer) Delete(ctx context.Context, req *v1.StreamRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger = logger.WithField("user_id", userID)

	err = s.manager.DeleteUserStream(ctx, userID, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrStreamCantBeDeleted {
			return nil, rpc.ErrRpcBadRequest
		}

		logFailedTo(logger, "delete user stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.logger.Error(s.eb.EmitDeleteStream(
		opentracing.ContextWithSpan(ctx, span),
		req.Id))

	return &protoempty.Empty{}, nil
}

func (s *RPCServer) Get(ctx context.Context, req *v1.StreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger = logger.WithField("user_id", userID)

	stream, err := s.manager.GetUserStream(ctx, userID, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get user stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *RPCServer) List(ctx context.Context, req *protoempty.Empty) (*v1.StreamListResponse, error) {
	span := opentracing.SpanFromContext(ctx)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger := s.logger.WithField("user_id", userID)

	streams, err := s.manager.GetStreamListByUserID(ctx, userID)
	if err != nil {
		logFailedTo(logger, "get stream list by user id", err)
		return nil, rpc.ErrRpcInternal
	}

	streamListResponse, err := toStreamListResponse(streams)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamListResponse, nil
}

func (s *RPCServer) Update(ctx context.Context, req *v1.UpdateStreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	span.SetTag("user_id", userID)
	logger = logger.WithField("user_id", userID)

	stream, err := s.manager.GetUserStream(ctx, userID, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get user stream", err)
		return nil, rpc.ErrRpcInternal
	}

	err = s.manager.UpdateStream(
		ctx,
		stream,
		map[string]interface{}{"name": req.Name},
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.logger.Error(s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.ID))

	return streamResponse, nil
}

func (s *RPCServer) UpdateStatus(ctx context.Context, req *v1.UpdateStreamRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	stream, err := s.manager.GetStreamByID(ctx, req.Id)
	if err != nil {
		logFailedTo(logger, "get stream by id", err)
		return nil, rpc.ErrRpcInternal
	}

	updates := make(map[string]interface{})

	if req.Status != v1.StreamStatusNew {
		updates["status"] = req.Status

		if req.Status == v1.StreamStatusPrepared {
			updates["input_status"] = v1.InputStatusPending
		}
	}

	if req.StreamContractAddress != "" {
		updates["stream_contract_address"] = req.StreamContractAddress
	}

	if req.InputStatus != v1.InputStatusNone {
		updates["input_status"] = req.InputStatus
	}

	if req.Status == v1.StreamStatusFailed {
		_, err = s.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
			UserId:                stream.UserID,
			StreamContractId:      stream.StreamContractID,
			StreamContractAddress: stream.StreamContractAddress,
		})

		if err != nil {
			logFailedTo(logger, "end stream", err)
		}
	}

	err = s.manager.UpdateStream(
		ctx,
		stream,
		updates,
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.logger.Error(s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.ID))

	return &protoempty.Empty{}, nil
}

func (s *RPCServer) Run(ctx context.Context, req *v1.StreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger = logger.WithField("user_id", userID)

	stream, err := s.manager.RunStream(ctx, req.Id, userID)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrHitBalanceLimitation {
			return nil, rpc.NewRpcPermissionError(err.Error())
		}

		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}

func (s *RPCServer) Stop(ctx context.Context, req *v1.StreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger = logger.WithField("user_id", userID)

	stream, err := s.manager.GetStreamByID(ctx, req.Id)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		logFailedTo(logger, "get stream by id", err)
		return nil, rpc.ErrRpcInternal
	}

	ss := v1.StreamStatusCancelled

	if stream.InputType == v1.InputTypeRTMP || stream.InputType == v1.InputTypeWebRTC {
		if stream.Status == v1.StreamStatusReady {
			ss = v1.StreamStatusCompleted
		}
	}

	stream, err = s.manager.StopStream(ctx, req.Id, userID, ss)
	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		if err == manager.ErrEndStreamNotAllowed {
			return nil, rpc.ErrRpcBadRequest
		}

		return nil, rpc.ErrRpcInternal
	}

	resp, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return resp, nil
}
