package model

import "time"

type TestConfig struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID string    `json:"project_id" gorm:"type:uuid;not null;index"`
	Name      string    `json:"name" gorm:"size:128;not null"`
	TestType  string    `json:"test_type" gorm:"size:32;not null;default:'unit'"`
	Framework string    `json:"framework" gorm:"size:64"`
	Command   string    `json:"command"`
	Timeout   int       `json:"timeout" gorm:"default:3600"`
	Enabled   bool      `json:"enabled" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TestConfig) TableName() string { return "test_configs" }
