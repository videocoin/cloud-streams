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
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	tracer "github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/cloud-pkg/uuid4"
	"github.com/videocoin/cloud-streams/datastore"
	ds "github.com/videocoin/cloud-streams/datastore"
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
	inputType v1.InputType,
	outputType v1.OutputType,
) (*ds.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.CreateStream")
	defer span.Finish()

	span.SetTag("name", name)
	span.SetTag("user_id", userID)
	span.SetTag("profile_id", profileID)
	span.SetTag("input_type", inputType.String())
	span.SetTag("output_type", outputType.String())

	id, err := uuid4.New()
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	span.SetTag("id", id)

	rand.Seed(time.Now().UTC().UnixNano())
	streamContractID := big.NewInt(int64(rand.Intn(math.MaxInt64)))

	outputMediaURL := fmt.Sprintf("%s/%s/index.m3u8", outputURL, id)
	if outputType == v1.OutputTypeFile {
		outputMediaURL = fmt.Sprintf("%s/%s/index.mp4", outputURL, id)
	}

	stream, err := m.ds.Stream.Create(ctx, &ds.Stream{
		ID:               id,
		UserID:           userID,
		Name:             name,
		ProfileID:        profileID,
		InputURL:         fmt.Sprintf("%s/%s/index.m3u8", inputURL, id),
		OutputURL:        outputMediaURL,
		RtmpURL:          fmt.Sprintf("%s/%s", rtmpURL, id),
		StreamContractID: streamContractID.Uint64(),
		Status:           v1.StreamStatusNew,
		InputType:        inputType,
		OutputType:       outputType,
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

func (m *Manager) GetStreamByID(ctx context.Context, id string) (*ds.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.GetStreamByID")
	defer span.Finish()

	stream, err := m.ds.Stream.Get(ctx, id)
	if err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return stream, nil
}

func (m *Manager) GetStreamListByUserID(ctx context.Context, userID string) ([]*ds.Stream, error) {
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

func (m *Manager) UpdateStream(ctx context.Context, stream *ds.Stream, updates map[string]interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "manager.UpdateStream")
	defer span.Finish()

	if value, ok := updates["name"]; ok {
		stream.Name = value.(string)
	}

	if value, ok := updates["stream_contract_address"]; ok {
		stream.StreamContractAddress = value.(string)
	}

	if value, ok := updates["status"]; ok {
		if stream.Status != v1.StreamStatusFailed &&
			stream.Status != v1.StreamStatusCompleted &&
			stream.Status != v1.StreamStatusCancelled {
			stream.Status = value.(v1.StreamStatus)
		}
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

func (m *Manager) GetUserStream(ctx context.Context, userID string, streamID string) (*ds.Stream, error) {
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

	if err := m.ds.Stream.Delete(ctx, stream.ID); err != nil {
		tracer.SpanLogError(span, err)
		return err
	}

	return nil
}

func (m *Manager) RunStream(ctx context.Context, streamID string, userID string) (*ds.Stream, error) {
	logger := m.logger.WithField("id", streamID)

	if userID != "" {
		logger = logger.WithField("user_id", userID)

		err := m.checkBalance(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	var (
		stream *ds.Stream
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

	if userID == "" && stream.UserID != "" {
		logger.WithField("user_id", stream.UserID)

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

	go func() {
		initStream, err := m.emitter.InitStream(ctx, &emitterv1.InitStreamRequest{
			StreamId:         stream.ID,
			UserId:           userID,
			StreamContractId: stream.StreamContractID,
			ProfilesIds:      []string{stream.ProfileID},
		})
		if err != nil {
			logger.WithError(err).Error("failed to init stream")

			updates := map[string]interface{}{"status": v1.StreamStatusFailed}
			updateErr := m.UpdateStream(ctx, stream, updates)
			if updateErr != nil {
				logger.Errorf("failed to update stream: %s", updateErr)
			}

			return
		}

		updates = map[string]interface{}{
			"stream_contract_address": initStream.StreamContractAddress,
			"status":                  v1.StreamStatusPrepared,
			"input_status":            v1.InputStatusPending,
		}
		err = m.UpdateStream(ctx, stream, updates)
		if err != nil {
			logger.WithError(err).Error("failed to update stream")
			return
		}

		err = m.eb.EmitUpdateStream(ctx, streamID)
		if err != nil {
			logger.WithError(err).Error("failed to emit update stream")
			return
		}
	}()

	return stream, nil
}

func (m *Manager) StopStream(
	ctx context.Context,
	streamID string,
	userID string,
	streamStatus v1.StreamStatus,
) (*ds.Stream, error) {
	var (
		stream *ds.Stream
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

	if stream.Status == v1.StreamStatusPrepared ||
		stream.Status == v1.StreamStatusPending {
		m.EndStream(ctx, stream)
	}

	if stream.Status == v1.StreamStatusCompleted {
		return stream, nil
	}

	updates := map[string]interface{}{"status": streamStatus}
	err = m.UpdateStream(ctx, stream, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update stream: %s", err)
	}
	go func() {
		err := m.eb.EmitUpdateStream(ctx, streamID)
		if err != nil {
			m.logger.Error(err)
		}
	}()

	return stream, nil
}

func (m *Manager) CompleteStream(ctx context.Context, stream *ds.Stream) error {
	stream.Status = v1.StreamStatusCompleted

	updates := map[string]interface{}{"status": stream.Status}
	err := m.UpdateStream(ctx, stream, updates)
	if err != nil {
		return err
	}

	m.EndStream(ctx, stream)

	go func() {
		err := m.eb.EmitUpdateStream(ctx, stream.ID)
		if err != nil {
			m.logger.Error(err)
		}
	}()
	return nil
}

func (m *Manager) EndStream(ctx context.Context, stream *ds.Stream) {
	resp, err := m.emitter.EndStream(ctx, &emitterv1.EndStreamRequest{
		UserId:                stream.UserID,
		StreamContractId:      stream.StreamContractID,
		StreamContractAddress: stream.StreamContractAddress,
	})

	if err != nil {
		logger := m.logger.WithFields(logrus.Fields{
			"user_id":                 stream.UserID,
			"stream_contract_id":      stream.StreamContractID,
			"stream_contract_address": stream.StreamContractAddress,
		})

		if resp != nil {
			logger = logger.
				WithField("endstream_tx", resp.EndStreamTx).
				WithField("escrow_tx", resp.EscrowRefundTx)
		}

		logger.WithError(err).Error("failed to end stream")
	}
}
