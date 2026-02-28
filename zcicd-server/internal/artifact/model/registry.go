package model

import "time"

type ImageRegistry struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string    `json:"name" gorm:"size:128;not null"`
	RegistryType string    `json:"registry_type" gorm:"size:32;not null;default:'harbor'"`
	Endpoint     string    `json:"endpoint" gorm:"size:512;not null"`
	Username     string    `json:"username" gorm:"size:128"`
	PasswordEnc  []byte    `json:"-" gorm:"type:bytea"`
	IsDefault    bool      `json:"is_default" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ImageRegistry) TableName() string { return "image_registries" }
