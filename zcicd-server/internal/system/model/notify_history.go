package model

import "time"

type NotifyHistory struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChannelID    string    `json:"channel_id" gorm:"type:uuid;not null;index"`
	EventType    string    `json:"event_type" gorm:"size:64;not null"`
	Title        string    `json:"title" gorm:"size:256"`
	Content      string    `json:"content"`
	Status       string    `json:"status" gorm:"size:16;not null;default:'pending'"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func (NotifyHistory) TableName() string { return "notify_history" }
