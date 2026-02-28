package model

import (
	"time"

	"gorm.io/datatypes"
)

type Service struct {
	ID              string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID       string         `json:"project_id" gorm:"type:uuid;not null;index"`
	Name            string         `json:"name" gorm:"size:100;not null"`
	ServiceType     string         `json:"service_type" gorm:"size:20"`
	Language        string         `json:"language" gorm:"size:30"`
	RepoURL         string         `json:"repo_url"`
	Branch          string         `json:"branch" gorm:"default:'main'"`
	DockerfilePath  string         `json:"dockerfile_path" gorm:"default:'Dockerfile'"`
	BuildContext    string         `json:"build_context" gorm:"default:'.'"`
	DeployType      string         `json:"deploy_type" gorm:"size:20"`
	HelmChartPath   string         `json:"helm_chart_path"`
	HelmValues      datatypes.JSON `json:"helm_values"`
	K8sManifests    string         `json:"k8s_manifests" gorm:"type:text"`
	HealthCheckPath string         `json:"health_check_path" gorm:"default:'/healthz'"`
	Ports           datatypes.JSON `json:"ports"`
	EnvVars         datatypes.JSON `json:"env_vars"`
	Resources       datatypes.JSON `json:"resources"`
	Status          string         `json:"status" gorm:"default:'active'"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Service) TableName() string { return "services" }
