package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	v1 "github.com/videocoin/cloud-api/users/v1"
)

var (
	ErrUserNotFound = errors.New("user is not found")
)

type UserDatastore struct {
	db *gorm.DB
}

func NewUserDatastore(db *gorm.DB) (*UserDatastore, error) {
	db.AutoMigrate(&v1.User{})
	return &UserDatastore{db: db}, nil
}

func (ds *UserDatastore) Get(ctx context.Context, id string) (*v1.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()

	user := &v1.User{}

	if err := ds.db.Where("id = ?", id).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %s", err.Error())
	}

	return user, nil
}
