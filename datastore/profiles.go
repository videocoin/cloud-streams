package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/uuid4"
)

var (
	ErrProfileNotFound = errors.New("profile is not found")
)

type ProfileDatastore struct {
	db *gorm.DB
}

type Profile struct {
	ID          string           `gorm:"type:varchar(36);PRIMARY_KEY"`
	Name        string           `gorm:"type:varchar(255)"`
	Description string           `gorm:"type:varchar(255)"`
	IsEnabled   bool             `gorm:"type:tinyint(1);DEFAULT:0" json:"is_enabled"`
	Spec        v1.Spec          `gorm:"type:json;DEFAULT:null"`
	Rel         string           `gorm:"type:text"`
	Capacity    *v1.CapacityInfo `gorm:"type:json;DEFAULT:null"`
}

func (Profile) TableName() string {
	return "streams_profiles"
}

func NewProfileDatastore(db *gorm.DB) (*ProfileDatastore, error) {
	return &ProfileDatastore{db: db}, nil
}

func (ds *ProfileDatastore) Create(ctx context.Context, profile *Profile) (*Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Create")
	defer span.Finish()

	tx := ds.db.Begin()

	if profile.ID == "" {
		id, _ := uuid4.New()
		profile.ID = id
	}

	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return profile, nil
}

func (ds *ProfileDatastore) Delete(ctx context.Context, id string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()

	span.SetTag("id", id)

	profile := &Profile{
		ID: id,
	}

	if err := ds.db.Delete(profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrProfileNotFound
		}

		return fmt.Errorf("failed to get profile by id %s: %s", id, err.Error())
	}

	return nil
}

func (ds *ProfileDatastore) Get(ctx context.Context, id string) (*Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	span.SetTag("id", id)

	profile := &Profile{}

	if err := ds.db.Where("id = ?", id).First(profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrProfileNotFound
		}

		return nil, fmt.Errorf("failed to get profile by id %s: %s", id, err.Error())
	}

	return profile, nil
}

func (ds *ProfileDatastore) ListEnabled(ctx context.Context) ([]*Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListEnabled")
	defer span.Finish()

	profiles := []*Profile{}

	if err := ds.db.Where("is_enabled = ?", true).Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("failed to list enabled profiles: %s", err)
	}

	return profiles, nil
}

func (ds *ProfileDatastore) List(ctx context.Context) ([]*Profile, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "List")
	defer span.Finish()

	profiles := []*Profile{}

	if err := ds.db.Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("failed to list profiles: %s", err)
	}

	return profiles, nil
}

func (ds *ProfileDatastore) Update(ctx context.Context, profile *Profile) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Update")
	defer span.Finish()

	return ds.db.Save(profile).Error
}

func (ds *ProfileDatastore) DeleteAllExceptIds(ctx context.Context, ids []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "DeleteAllExceptIds")
	defer span.Finish()

	if len(ids) == 0 {
		err := ds.db.Delete(&Profile{}).Error
		if err != nil {
			return err
		}
	} else {
		err := ds.db.Where("id NOT IN (?)", ids).Delete(&Profile{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
