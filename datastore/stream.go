package datastore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/streams/v1"
)

var (
	ErrStreamNotFound = errors.New("stream is not found")
)

type Stream struct {
	ID                    string          `gorm:"type:varchar(36);PRIMARY_KEY"`
	UserID                string          `gorm:"type:varchar(255)"`
	Name                  string          `gorm:"type:varchar(255)"`
	ProfileID             string          `gorm:"type:varchar(255)"`
	Status                v1.StreamStatus `gorm:"type:string"`
	InputStatus           v1.InputStatus  `gorm:"type:string"`
	StreamContractID      uint64          `gorm:"type:bigint(20) unsigned"`
	StreamContractAddress string          `gorm:"type:varchar(255);DEFAULT:null"`
	InputURL              string          `gorm:"type:varchar(255)"`
	OutputURL             string          `gorm:"type:varchar(255)"`
	RtmpURL               string          `gorm:"type:varchar(255)"`
	Refunded              sql.NullString  `gorm:"type:tinyint(1);DEFAULT:null"`
	CreatedAt             *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	UpdatedAt             *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	ReadyAt               *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	CompletedAt           *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	DeletedAt             *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	InputType             v1.InputType    `gorm:"type:string"`
	OutputType            v1.OutputType   `gorm:"type:string"`
	TotalCost             float64         `gorm:"type:decimal(10,6);default:0"`
	DrmXml                string          `gorm:"type:text;DEFAULT:null"`
}

type StreamDatastore struct {
	db *gorm.DB
}

func NewStreamDatastore(db *gorm.DB) (*StreamDatastore, error) {
	return &StreamDatastore{db: db}, nil
}

func (ds *StreamDatastore) Create(ctx context.Context, stream *Stream) (*Stream, error) {
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

	stream := &Stream{
		ID: id,
	}

	if err := ds.db.Delete(stream).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrStreamNotFound
		}

		return fmt.Errorf("failed to get stream by id %s: %s", id, err.Error())
	}

	return nil
}

func (ds *StreamDatastore) Get(ctx context.Context, id string) (*Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	span.SetTag("id", id)

	stream := &Stream{}

	if err := ds.db.Where("id = ?", id).First(stream).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStreamNotFound
		}

		return nil, fmt.Errorf("failed to get stream by id %s: %s", id, err.Error())
	}

	return stream, nil
}

func (ds *StreamDatastore) GetByUserID(ctx context.Context, userID string, streamID string) (*Stream, error) {
	stream := &Stream{}

	if err := ds.db.Where("user_id = ? AND id = ?", userID, streamID).First(stream).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrStreamNotFound
		}

		return nil, fmt.Errorf("failed to get stream by user id: %s", err.Error())
	}

	return stream, nil
}

func (ds *StreamDatastore) List(ctx context.Context, userID string) ([]*Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "List")
	defer span.Finish()

	span.SetTag("user_id", userID)

	streams := []*Stream{}

	if err := ds.db.Where("user_id = ?", userID).Order("created_at desc").Find(&streams).Error; err != nil {
		return nil, fmt.Errorf("failed to list streams by user id %s: %s", userID, err)
	}

	return streams, nil
}

func (ds *StreamDatastore) StatusReadyList(ctx context.Context) ([]*Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "StatusReadyList")
	defer span.Finish()

	streams := []*Stream{}

	if err := ds.db.Where("status = ?", v1.StreamStatusReady).Find(&streams).Error; err != nil {
		return nil, fmt.Errorf("failed to list streams: %s", err)
	}

	return streams, nil
}

func (ds *StreamDatastore) ListForDeletion(ctx context.Context) ([]*Stream, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListForDeletion")
	defer span.Finish()

	streams := []*Stream{}

	if err := ds.db.
		Where(
			"status = ? AND completed_at <= ?",
			v1.StreamStatusCompleted,
			time.Now().Add(time.Hour*24*-1),
		).
		Find(&streams).Error; err != nil {
		return nil, fmt.Errorf("failed to list streams for deletion: %s", err)
	}

	return streams, nil
}

func (ds *StreamDatastore) Update(ctx context.Context, stream *Stream, updates map[string]interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Update")
	defer span.Finish()

	if err := ds.db.Model(stream).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (ds *StreamDatastore) MarkAsDeleted(ctx context.Context, stream *Stream) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "MarkAsDeleted")
	defer span.Finish()

	stream.Status = v1.StreamStatusDeleted
	updates := map[string]interface{}{"status": stream.Status}
	if err := ds.db.Model(stream).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
