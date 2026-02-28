package model

import (
	"time"

	"gorm.io/datatypes"
)

type ScanConfig struct {
	ID              string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID       string         `json:"project_id" gorm:"type:uuid;not null;index"`
	Name            string         `json:"name" gorm:"size:128;not null"`
	ScanType        string         `json:"scan_type" gorm:"size:32;not null;default:'sonar'"`
	SonarProjectKey string         `json:"sonar_project_key" gorm:"size:256"`
	Config          datatypes.JSON `json:"config" gorm:"default:'{}'"`
	Enabled         bool           `json:"enabled" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (ScanConfig) TableName() string { return "scan_configs" }
