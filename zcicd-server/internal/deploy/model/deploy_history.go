package model

import "time"

type DeployHistory struct {
	ID             string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DeployConfigID string     `json:"deploy_config_id" gorm:"type:uuid;not null;index"`
	Revision       string     `json:"revision" gorm:"size:128"`
	Status         string     `json:"status" gorm:"size:32;not null;default:'pending'"`
	SyncStatus     string     `json:"sync_status" gorm:"size:32"`
	HealthStatus   string     `json:"health_status" gorm:"size:32"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	Duration       int        `json:"duration"`
	TriggeredBy    string     `json:"triggered_by" gorm:"type:uuid"`
	RollbackFrom   *string    `json:"rollback_from" gorm:"type:uuid"`
	GitopsCommit   string     `json:"gitops_commit" gorm:"size:64"`
	DiffContent    string     `json:"diff_content" gorm:"type:text"`
	ErrorMessage   string     `json:"error_message" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at"`

	DeployConfig DeployConfig `json:"deploy_config,omitempty" gorm:"foreignKey:DeployConfigID"`
}

func (DeployHistory) TableName() string { return "deploy_histories" }
