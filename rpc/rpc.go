package rpc

import (
	"context"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-streams/datastore"
	"github.com/videocoin/cloud-streams/manager"
)

func (s *Server) Create(ctx context.Context, req *v1.CreateStreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("name", req.Name)
	span.SetTag("profile_id", req.ProfileId)
	span.SetTag("input_type", req.InputType.String())
	span.SetTag("output_type", req.OutputType.String())

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

	profile, err := s.manager.GetProfileByID(ctx, req.ProfileId)
	if err != nil {
		logFailedTo(logger, "get profile", err)
		return nil, rpc.ErrRpcBadRequest
	}

	stream, err := s.manager.CreateStream(
		ctx,
		req.Name,
		profile.ID,
		userID,
		s.baseInputURL,
		s.baseOutputURL,
		s.rtmpURL,
		req.InputType,
		req.OutputType,
		req.DrmXml,
	)
	if err != nil {
		logFailedTo(logger, "create stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go func() {
		err := s.eb.EmitCreateStream(opentracing.ContextWithSpan(ctx, span), stream.ID)
		if err != nil {
			logFailedTo(logger, "emit create stream", err)
		}
	}()

	streamResponse, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamResponse, nil
}

func (s *Server) Delete(ctx context.Context, req *v1.StreamRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

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

	go func() {
		err := s.eb.EmitDeleteStream(opentracing.ContextWithSpan(ctx, span), req.Id)
		if err != nil {
			logFailedTo(logger, "emit delete stream", err)
		}
	}()

	return &protoempty.Empty{}, nil
}

func (s *Server) Get(ctx context.Context, req *v1.StreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

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

func (s *Server) List(ctx context.Context, req *protoempty.Empty) (*v1.StreamListResponse, error) {
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

func (s *Server) Update(ctx context.Context, req *v1.UpdateStreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

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

	err = s.manager.UpdateStream(ctx, stream, map[string]interface{}{"name": req.Name})
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamResponse, err := toStreamResponse(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go func() {
		err := s.eb.EmitUpdateStream(opentracing.ContextWithSpan(ctx, span), stream.ID)
		if err != nil {
			logFailedTo(logger, "emit update stream", err)
		}
	}()

	return streamResponse, nil
}

func (s *Server) UpdateStatus(ctx context.Context, req *v1.UpdateStreamRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)

	otCtx := opentracing.ContextWithSpan(context.Background(), span)
	stream, err := s.manager.GetStreamByID(otCtx, req.Id)
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
		s.manager.EndStream(otCtx, stream)
	}

	err = s.manager.UpdateStream(otCtx, stream, updates)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go func() {
		err := s.eb.EmitUpdateStream(otCtx, stream.ID)
		if err != nil {
			logFailedTo(logger, "emit update stream", err)
		}
	}()

	return &protoempty.Empty{}, nil
}

func (s *Server) Run(ctx context.Context, req *v1.StreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)
	otCtx := opentracing.ContextWithSpan(context.Background(), span)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger = logger.WithField("user_id", userID)

	stream, err := s.manager.RunStream(otCtx, req.Id, userID)
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

func (s *Server) Stop(ctx context.Context, req *v1.StreamRequest) (*v1.StreamResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("stream_id", req.Id)
	logger := s.logger.WithField("stream_id", req.Id)
	otCtx := opentracing.ContextWithSpan(context.Background(), span)

	userID, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)

	logger = logger.WithField("user_id", userID)
	logger.Info("getting stream by id")

	stream, err := s.manager.GetStreamByID(otCtx, req.Id)
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

	logger.WithField("status", ss).Info("stopping stream")

	stream, err = s.manager.StopStream(otCtx, req.Id, userID, ss)
	if err != nil {
		logger.WithError(err).Errorf("failed to stop stream")

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

func (s *Server) GetProfileList(ctx context.Context, _ *protoempty.Empty) (*v1.ProfileListResponse, error) {
	profiles, err := s.manager.ListEnabledProfiles(ctx)
	if err != nil {
		logFailedTo(s.logger, "get profile list", err)
		return nil, rpc.ErrRpcInternal
	}

	resp := &v1.ProfileListResponse{}
	if err := copier.Copy(&resp.Items, &profiles); err != nil {
		return nil, err
	}

	return resp, nil
}
