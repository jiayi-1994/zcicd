package model

import "time"

type NotifyRule struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	EventType   string    `json:"event_type" gorm:"size:64;not null;index"`
	ChannelID   string    `json:"channel_id" gorm:"type:uuid;not null;index"`
	ProjectID   string    `json:"project_id" gorm:"type:uuid"`
	Severity    string    `json:"severity" gorm:"size:16;default:'all'"`
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (NotifyRule) TableName() string { return "notify_rules" }
