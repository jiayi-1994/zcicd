package model

import (
	"time"

	"gorm.io/datatypes"
)

type NotifyChannel struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"size:128;not null"`
	ChannelType string         `json:"channel_type" gorm:"size:32;not null"`
	Config      datatypes.JSON `json:"config" gorm:"default:'{}'"`
	Enabled     bool           `json:"enabled" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (NotifyChannel) TableName() string { return "notify_channels" }
