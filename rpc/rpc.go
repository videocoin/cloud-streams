package rpc

import (
	"context"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/opentracing/opentracing-go"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-streams/datastore"
)

func (s *RpcServer) Create(ctx context.Context, req *v1.CreateStreamRequest) (*v1.StreamProfile, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("name", req.Name)
	span.SetTag("profile_id", req.ProfileId)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	userID, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger := s.logger.WithField("user_id", userID)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	stream, err := s.manager.CreateStream(
		ctx,
		req.Name,
		req.ProfileId,
		userID,
		s.baseInputURL,
		s.baseOutputURL,
		s.rtmpURL,
	)
	if err != nil {
		logFailedTo(logger, "create stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitCreateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.Id)

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfile, nil
}

func (s *RpcServer) Delete(ctx context.Context, req *v1.StreamRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, _, err := s.authenticate(ctx)
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

		logFailedTo(logger, "delete user stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitDeleteStream(
		opentracing.ContextWithSpan(ctx, span),
		req.Id)

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) Get(ctx context.Context, req *v1.StreamRequest) (*v1.StreamProfile, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, _, err := s.authenticate(ctx)
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

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfile, nil
}

func (s *RpcServer) List(ctx context.Context, req *protoempty.Empty) (*v1.StreamProfiles, error) {
	span := opentracing.SpanFromContext(ctx)

	userID, _, err := s.authenticate(ctx)
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

	streamProfiles, err := toStreamProfiles(streams)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfiles, nil
}

func (s *RpcServer) Update(ctx context.Context, req *v1.UpdateStreamRequest) (*v1.StreamProfile, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Error(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	userID, _, err := s.authenticate(ctx)
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

	err = s.manager.UpdateStream(
		ctx,
		stream,
		map[string]interface{}{"name": req.Name},
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.Id)

	return streamProfile, nil
}

func (s *RpcServer) UpdateStatus(ctx context.Context, req *v1.UpdateStreamRequest) (*protoempty.Empty, error) {
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

	err = s.manager.UpdateStream(
		ctx,
		stream,
		updates,
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.Id)

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) Run(ctx context.Context, req *v1.StreamRequest) (*v1.StreamProfile, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, _, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	span.SetTag("user_id", userID)
	logger = logger.WithField("user_id", userID)

	// account, err := s.accounts.GetByOwner(ctx, &accountsv1.AccountRequest{OwnerId: userID})
	// if err != nil {
	// 	logFailedTo(logger, "get account", err)
	// 	return nil, rpc.ErrRpcInternal
	// }

	// // temp balance limit force on api level
	// if account.Balance < 20 || account.Balance > 50 {
	// 	return nil, rpc.NewRpcPermissionError("hit balance limitation")
	// }

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
		map[string]interface{}{"status": v1.StreamStatusPreparing},
	)
	if err != nil {
		logFailedTo(logger, "update stream", err)
		return nil, rpc.ErrRpcInternal
	}

	_, err = s.emitter.InitStream(ctx, &emitterv1.InitStreamRequest{
		StreamId:         stream.Id,
		UserId:           userID,
		StreamContractId: stream.StreamContractId,
		ProfilesIds:      []string{stream.ProfileId},
	})
	if err != nil {
		logFailedTo(logger, "init stream", err)
		return nil, rpc.ErrRpcInternal
	}

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.Id)

	return streamProfile, nil
}

func (s *RpcServer) Stop(ctx context.Context, req *v1.StreamRequest) (*v1.StreamProfile, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)
	logger := s.logger.WithField("id", req.Id)

	userID, _, err := s.authenticate(ctx)
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

	if stream.Status == v1.StreamStatusCompleted {
		// nothing to do since it was already stopped
		streamProfile, err := toStreamProfile(stream)
		if err != nil {
			logFailedTo(logger, "", err)
			return nil, rpc.ErrRpcInternal
		}

		return streamProfile, nil
	}

	_, err = s.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
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

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		logFailedTo(logger, "", err)
		return nil, rpc.ErrRpcInternal
	}

	go s.eb.EmitUpdateStream(
		opentracing.ContextWithSpan(ctx, span),
		stream.Id)

	return streamProfile, nil
}
