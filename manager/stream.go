package manager

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"

	"github.com/opentracing/opentracing-go"
	profilesv1 "github.com/videocoin/cloud-api/profiles/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	tracer "github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/cloud-pkg/uuid4"
)

func (m *Manager) CreateStream(ctx context.Context, name, userID, inputURL, outputURL string, profileID profilesv1.ProfileId) (*v1.Stream, error) {
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

	streamContractID := big.NewInt(int64(rand.Intn(math.MaxInt64)))
	stream, err := m.ds.Stream.Create(ctx, &v1.Stream{
		Id:               id,
		UserId:           userID,
		Name:             name,
		ProfileId:        profileID,
		InputUrl:         fmt.Sprintf("%s/%s", inputURL, id),
		OutputUrl:        fmt.Sprintf("%s/%s/index.m3u8", outputURL, id),
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

func (m *Manager) UpdateStream(ctx context.Context, stream *v1.Stream, updates map[string]interface{}) (*v1.Stream, error) {
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

	if err := m.ds.Stream.Update(ctx, stream, updates); err != nil {
		tracer.SpanLogError(span, err)
		return nil, err
	}

	return stream, nil
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

	if err := m.ds.Stream.Delete(ctx, stream.Id); err != nil {
		tracer.SpanLogError(span, err)
		return err
	}

	return nil
}
