package datastore

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //nolint
)

type Datastore struct {
	Stream *StreamDatastore
}

func NewDatastore(uri string) (*Datastore, error) {
	ds := new(Datastore)

	db, err := gorm.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	db.LogMode(false)

	streamDs, err := NewStreamDatastore(db)
	if err != nil {
		return nil, err
	}

	ds.Stream = streamDs

	return ds, nil
}
