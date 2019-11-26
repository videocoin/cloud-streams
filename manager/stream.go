package manager

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/opentracing/opentracing-go"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	tracer "github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/cloud-pkg/uuid4"
	"github.com/videocoin/cloud-streams/datastore"
)

var (
	ErrStreamCantBeDeleted = errors.New("stream can't be deleted")
)

func (m *Manager) CreateStream(
	ctx context.Context,
	name string,
	profileID string,
	userID string,
	inputURL string,
	outputURL string,
	rtmpURL string,
) (*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.CreateStream")
	defer span.Finish()

	span.SetTag("name", name)
	span.SetTag("user_id", userID)
	span.SetTag("profile_id", profileID)

	id, err := uuid4.New()
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	span.SetTag("id", id)

	rand.Seed(time.Now().UTC().UnixNano())
	streamContractID := big.NewInt(int64(rand.Intn(math.MaxInt64)))
	stream, err := m.ds.Stream.Create(ctx, &v1.Stream{
		Id:               id,
		UserId:           userID,
		Name:             name,
		ProfileId:        profileID,
		InputUrl:         fmt.Sprintf("%s/%s/index.m3u8", inputURL, id),
		OutputUrl:        fmt.Sprintf("%s/%s/index.m3u8", outputURL, id),
		RtmpUrl:          fmt.Sprintf("%s/%s", rtmpURL, id),
		StreamContractId: streamContractID.Uint64(),
		Status:           v1.StreamStatusNew,
	})
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return stream, nil
}

func (m *Manager) Delete(ctx context.Context, id string) error {
	if err := m.ds.Stream.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetStreamByID(ctx context.Context, id string) (*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetStreamByID")
	defer span.Finish()

	stream, err := m.ds.Stream.Get(ctx, id)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return stream, nil
}

func (m *Manager) GetStreamListByUserID(ctx context.Context, userID string) ([]*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetStreamListByUserID")
	defer span.Finish()

	span.SetTag("user_id", userID)

	streams, err := m.ds.Stream.List(ctx, userID)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return streams, nil
}

func (m *Manager) UpdateStream(ctx context.Context, stream *v1.Stream, updates map[string]interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.UpdateStream")
	defer span.Finish()

	if value, ok := updates["name"]; ok {
		stream.Name = value.(string)
	}

	if value, ok := updates["stream_contract_address"]; ok {
		stream.StreamContractAddress = value.(string)
	}

	if value, ok := updates["status"]; ok {
		stream.Status = value.(v1.StreamStatus)
	}

	if value, ok := updates["input_status"]; ok {
		stream.InputStatus = value.(v1.InputStatus)
	}

	if value, ok := updates["ready_at"]; ok {
		stream.ReadyAt = pointer.ToTime(value.(time.Time))
	}

	if value, ok := updates["completed_at"]; ok {
		stream.CompletedAt = pointer.ToTime(value.(time.Time))
	}

	if stream.Status == v1.StreamStatusCompleted ||
		stream.Status == v1.StreamStatusFailed ||
		stream.Status == v1.StreamStatusCancelled {
		stream.InputStatus = v1.InputStatusNone
		updates["input_status"] = stream.InputStatus
	}

	if stream.Status == v1.StreamStatusCompleted && stream.CompletedAt == nil {
		stream.CompletedAt = pointer.ToTime(time.Now())
		updates["completed_at"] = stream.CompletedAt
	}

	if stream.Status == v1.StreamStatusReady && stream.ReadyAt == nil {
		stream.ReadyAt = pointer.ToTime(time.Now())
		updates["ready_at"] = stream.ReadyAt
	}

	if err := m.ds.Stream.Update(ctx, stream, updates); err != nil {
		tracer.SpanLogError(span, err)
		return err
	}

	return nil
}

func (m *Manager) GetUserStream(ctx context.Context, userID string, streamID string) (*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetUserStream")
	defer span.Finish()

	span.SetTag("stream_id", streamID)
	span.SetTag("user_id", userID)

	stream, err := m.ds.Stream.GetByUserID(ctx, userID, streamID)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return stream, nil
}

func (m *Manager) DeleteUserStream(ctx context.Context, userID string, streamID string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.DeleteUserStream")
	defer span.Finish()

	span.SetTag("stream_id", streamID)
	span.SetTag("user_id", userID)

	stream, err := m.ds.Stream.GetByUserID(ctx, userID, streamID)
	if err != nil {
		tracer.SpanLogError(span, err)
		return err
	}

	if !isRemovable(stream) {
		return ErrStreamCantBeDeleted
	}

	if err := m.ds.Stream.Delete(ctx, stream.Id); err != nil {
		tracer.SpanLogError(span, err)
		return err
	}

	return nil
}

func (m *Manager) RunStream(ctx context.Context, streamID string, userID string) (*v1.Stream, error) {
	logger := m.logger.WithField("id", streamID)

	if userID != "" {
		logger = logger.WithField("user_id", userID)

		err := m.checkBalance(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	var (
		stream *v1.Stream
		err    error
	)

	if userID != "" {
		stream, err = m.GetUserStream(ctx, userID, streamID)
	} else {
		stream, err = m.GetStreamByID(ctx, streamID)
	}

	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get user stream: %s", err)
	}

	if userID == "" && stream.UserId != "" {
		logger = logger.WithField("user_id", stream.UserId)

		err := m.checkBalance(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	updates := map[string]interface{}{"status": v1.StreamStatusPreparing}
	err = m.UpdateStream(ctx, stream, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update stream: %s", err)
	}

	_, err = m.emitter.InitStream(ctx, &emitterv1.InitStreamRequest{
		StreamId:         stream.Id,
		UserId:           userID,
		StreamContractId: stream.StreamContractId,
		ProfilesIds:      []string{stream.ProfileId},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init stream: %s", err)
	}

	go m.eb.EmitUpdateStream(ctx, streamID)

	return stream, nil
}

func (m *Manager) StopStream(ctx context.Context, streamID string, userID string) (*v1.Stream, error) {
	var (
		stream *v1.Stream
		err    error
	)

	if userID != "" {
		stream, err = m.GetUserStream(ctx, userID, streamID)
	} else {
		stream, err = m.GetStreamByID(ctx, streamID)
	}

	if err != nil {
		if err == datastore.ErrStreamNotFound {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get user stream: %s", err)
	}

	if stream.Status < v1.StreamStatusPrepared {
		return nil, ErrEndStreamNotAllowed
	}

	if stream.Status == v1.StreamStatusCompleted {
		return stream, nil
	}

	updates := map[string]interface{}{"status": v1.StreamStatusCompleted}
	err = m.UpdateStream(ctx, stream, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update stream: %s", err)
	}

	return stream, nil
}
