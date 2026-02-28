package model

import (
	"time"

	"gorm.io/datatypes"
)

type DeployConfig struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID      string         `json:"project_id" gorm:"type:uuid;not null;index"`
	ServiceID      string         `json:"service_id" gorm:"type:uuid;not null;index"`
	EnvironmentID  string         `json:"environment_id" gorm:"type:uuid;not null;index"`
	Name           string         `json:"name" gorm:"size:128;not null"`
	DeployType     string         `json:"deploy_type" gorm:"size:32;not null;default:'helm'"`
	RepoURL        string         `json:"repo_url" gorm:"size:512;not null"`
	TargetRevision string         `json:"target_revision" gorm:"size:128;default:'main'"`
	ChartPath      string         `json:"chart_path" gorm:"size:256"`
	ValuesOverride datatypes.JSON `json:"values_override" gorm:"default:'{}'"`
	SyncPolicy     string         `json:"sync_policy" gorm:"size:32;default:'manual'"`
	AutoSync       bool           `json:"auto_sync" gorm:"default:false"`
	SelfHeal       bool           `json:"self_heal" gorm:"default:false"`
	Prune          bool           `json:"prune" gorm:"default:false"`
	ArgoAppName    string         `json:"argo_app_name" gorm:"size:256;index"`
	Namespace      string         `json:"namespace" gorm:"size:128"`
	Status         string         `json:"status" gorm:"size:32;default:'active'"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (DeployConfig) TableName() string { return "deploy_configs" }
