package manager

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"math/rand"

	profilesv1 "github.com/videocoin/cloud-api/profiles/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/uuid4"
)

func (m *Manager) Create(ctx context.Context, name, userId, inputURL, outputURL string, profileID profilesv1.ProfileId) (*v1.Stream, error) {
	id, err := uuid4.New()
	if err != nil {
		return nil, err
	}

	streamContractID := big.NewInt(int64(rand.Intn(math.MaxInt64)))
	stream, err := m.ds.Stream.Create(ctx, &v1.Stream{
		Id:               id,
		UserId:           userId,
		Name:             name,
		ProfileId:        profileID,
		InputUrl:         fmt.Sprintf("%s/%s", inputURL, id),
		OutputUrl:        fmt.Sprintf("%s/%s/index.m3u8", outputURL, id),
		StreamContractId: streamContractID.Uint64(),
		Status:           v1.StreamStatusNew,
	})
	if err != nil {
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

func (m *Manager) Get(ctx context.Context, id string) (*v1.Stream, error) {
	stream, err := m.ds.Stream.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (m *Manager) List(ctx context.Context, userId string) ([]*v1.Stream, error) {
	streams, err := m.ds.Stream.List(ctx, userId)
	if err != nil {
		return nil, err
	}

	return streams, nil
}

func (m *Manager) Update(ctx context.Context, stream *v1.Stream, updates map[string]interface{}) (*v1.Stream, error) {
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
		return nil, err
	}

	return stream, nil
}
