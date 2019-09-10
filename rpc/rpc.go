package rpc

import (
	"context"
	"fmt"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/streams/v1"
)

func (s *RpcServer) Create(ctx context.Context, req *v1.CreateStreamRequest) (*v1.StreamProfile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateStream")
	defer span.Finish()

	span.SetTag("name", req.Name)
	span.SetTag("profile_id", req.ProfileId)

	userID, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if verr := s.validator.validate(req); verr != nil {
		s.logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	stream, err := s.manager.Create(ctx, req.Name, userID, s.baseInputURL, s.baseOutputURL, req.ProfileId)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfile, nil
}

func (s *RpcServer) Delete(ctx context.Context, req *v1.StreamRequest) (*protoempty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	span.SetTag("id", req.Id)

	_, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	err = s.manager.Delete(ctx, req.Id)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) Get(ctx context.Context, req *v1.StreamRequest) (*v1.StreamProfile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	span.SetTag("id", req.Id)

	_, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	stream, err := s.manager.Get(ctx, req.Id)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfile, nil
}

func (s *RpcServer) List(ctx context.Context, req *protoempty.Empty) (*v1.StreamProfiles, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "List")
	defer span.Finish()

	userID, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	streams, err := s.manager.List(ctx, userID)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	streamProfiles, err := toStreamProfiles(streams)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfiles, nil
}

func (s *RpcServer) Update(ctx context.Context, req *v1.UpdateStreamRequest) (*v1.StreamProfile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Update")
	defer span.Finish()

	span.SetTag("id", req.Id)

	_, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	stream, err := s.manager.Get(ctx, req.Id)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	stream, err = s.manager.Update(
		ctx,
		stream,
		map[string]interface{}{
			"name": req.Name,
		},
	)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfile, nil
}

func (s *RpcServer) Run(ctx context.Context, req *v1.StreamRequest) (*v1.StreamProfile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Run")
	defer span.Finish()

	span.SetTag("id", req.Id)

	userID, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	account, err := s.accounts.GetByOwner(ctx, &accountsv1.AccountRequest{OwnerId: userID})
	if err != nil {
		s.logger.WithFields(
			logrus.Fields{
				"userId": userID,
			}).Errorf("failed to get account: %s", err.Error())
		return nil, rpc.ErrRpcInternal
	}

	// temp balance limit force on api level
	if account.Balance < 20 || account.Balance > 50 {
		s.logger.Errorf("hit balance limitation")
		return nil, rpc.ErrRpcBadRequest
	}

	stream, err := s.manager.Get(ctx, req.Id)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	if stream.UserId != userID {
		return nil, rpc.ErrRpcPermissionDenied
	}

	stream, err = s.manager.Update(
		ctx,
		stream,
		map[string]interface{}{
			"status": v1.StreamStatusPreparing,
		},
	)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	profileName := fmt.Sprintf("%d", stream.ProfileId)
	_, _ = s.emitter.InitStream(ctx, &emitterv1.InitStreamRequest{
		UserId:           userID,
		StreamContractId: stream.StreamContractId,
		ProfileNames:     []string{profileName},
	})

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfile, nil
}

func (s *RpcServer) Stop(ctx context.Context, req *v1.StreamRequest) (*v1.StreamProfile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Stop")
	defer span.Finish()

	span.SetTag("id", req.Id)

	_, _, err := s.authenticate(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	stream, err := s.manager.Get(ctx, req.Id)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	if stream.Status == v1.StreamStatusCompleted {
		// nothing to do since it was already stopped
		streamProfile, err := toStreamProfile(stream)
		if err != nil {
			s.logger.Error(err)
			return nil, rpc.ErrRpcInternal
		}

		return streamProfile, nil
	}

	_, _ = s.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
		StreamContractId:      stream.StreamContractId,
		StreamContractAddress: stream.StreamContractAddress,
	})

	stream, err = s.manager.Update(
		ctx,
		stream,
		map[string]interface{}{
			"status": v1.StreamStatusCompleted,
		},
	)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	streamProfile, err := toStreamProfile(stream)
	if err != nil {
		s.logger.Error(err)
		return nil, rpc.ErrRpcInternal
	}

	return streamProfile, nil
}
