package datastore

import (
	"database/sql"
	"time"

	v1 "github.com/videocoin/cloud-api/streams/v1"
)

type Stream struct {
	Id                    string          `gorm:"type:varchar(36);PRIMARY_KEY"`
	UserId                string          `gorm:"type:varchar(255)"`
	Name                  string          `gorm:"type:varchar(255)"`
	ProfileId             string          `gorm:"type:varchar(255)"`
	Status                v1.StreamStatus `gorm:"type:string"`
	InputStatus           v1.InputStatus  `gorm:"type:string"`
	StreamContractId      uint64          `gorm:"type:bigint(20) unsigned"`
	StreamContractAddress string          `gorm:"type:varchar(255);DEFAULT:null"`
	InputUrl              string          `gorm:"type:varchar(255)"`
	OutputUrl             string          `gorm:"type:varchar(255)"`
	RtmpUrl               string          `gorm:"type:varchar(255)"`
	Refunded              sql.NullString  `gorm:"type:tinyint(1);DEFAULT:null"`
	CreatedAt             *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	UpdatedAt             *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	ReadyAt               *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
	CompletedAt           *time.Time      `gorm:"type:timestamp NULL;DEFAULT:null"`
}
