package model

import "time"

type Integration struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	Type        string    `json:"type" gorm:"size:32;not null"`
	Provider    string    `json:"provider" gorm:"size:32;not null"`
	ConfigEnc   []byte    `json:"-" gorm:"type:bytea;not null"`
	Status      string    `json:"status" gorm:"size:16;default:'active'"`
	LastCheckAt *time.Time `json:"last_check_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Integration) TableName() string { return "integrations" }
