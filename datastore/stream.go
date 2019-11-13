package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/streams/v1"

	"github.com/jinzhu/gorm"
)

var (
	ErrStreamNotFound = errors.New("stream is not found")
)

type StreamDatastore struct {
	db *gorm.DB
}

func NewStreamDatastore(db *gorm.DB) (*StreamDatastore, error) {
	db.AutoMigrate(&v1.Stream{})
	return &StreamDatastore{db: db}, nil
}

func (ds *StreamDatastore) Create(ctx context.Context, stream *v1.Stream) (*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	tx := ds.db.Begin()

	time, err := ptypes.Timestamp(ptypes.TimestampNow())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	stream.CreatedAt = &time

	if err = tx.Create(stream).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return stream, nil
}

func (ds *StreamDatastore) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	span.SetTag("id", id)

	stream := &v1.Stream{
		Id: id,
	}

	if err := ds.db.Delete(stream).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrStreamNotFound
		}

		return fmt.Errorf("failed to get stream by id %s: %s", id, err.Error())
	}

	return nil
}

func (ds *StreamDatastore) Get(ctx context.Context, id string) (*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	span.SetTag("id", id)

	stream := &v1.Stream{}

	if err := ds.db.Where("id = ?", id).First(stream).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStreamNotFound
		}

		return nil, fmt.Errorf("failed to get stream by id %s: %s", id, err.Error())
	}

	return stream, nil
}

func (ds *StreamDatastore) GetByUserID(ctx context.Context, userID string, streamID string) (*v1.Stream, error) {
	stream := &v1.Stream{}

	if err := ds.db.Where("user_id = ? AND id = ?", userID, streamID).First(stream).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStreamNotFound
		}

		return nil, fmt.Errorf("failed to get stream by user id: %s", err.Error())
	}

	return stream, nil
}

func (ds *StreamDatastore) List(ctx context.Context, userID string) ([]*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "List")
	defer span.Finish()

	span.SetTag("user_id", userID)

	streams := []*v1.Stream{}

	if err := ds.db.Where("user_id = ?", userID).Find(&streams).Error; err != nil {
		return nil, fmt.Errorf("failed to list streams by user id %s: %s", userID, err)
	}

	return streams, nil
}

func (ds *StreamDatastore) StatusReadyList(ctx context.Context) ([]*v1.Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "StatusReadyList")
	defer span.Finish()

	streams := []*v1.Stream{}

	if err := ds.db.Where("status = ?", v1.StreamStatusReady).Find(&streams).Error; err != nil {
		return nil, fmt.Errorf("failed to list streams: %s", err)
	}

	return streams, nil
}

func (ds *StreamDatastore) Update(ctx context.Context, stream *v1.Stream, updates map[string]interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Update")
	defer span.Finish()

	if err := ds.db.Model(stream).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
