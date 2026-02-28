package model

import (
	"time"

	"gorm.io/datatypes"
)

type Environment struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID      string         `json:"project_id" gorm:"type:uuid;not null;index"`
	Name           string         `json:"name" gorm:"size:50;not null"`
	EnvType        string         `json:"env_type" gorm:"size:20;not null"`
	Namespace      string         `json:"namespace" gorm:"size:100"`
	ClusterID      string         `json:"cluster_id" gorm:"type:uuid"`
	IsProduction   bool           `json:"is_production" gorm:"default:false"`
	AutoDeploy     bool           `json:"auto_deploy" gorm:"default:false"`
	DeployStrategy datatypes.JSON `json:"deploy_strategy"`
	GlobalEnvVars  datatypes.JSON `json:"global_env_vars"`
	Status         string         `json:"status" gorm:"default:'active'"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (Environment) TableName() string { return "environments" }
