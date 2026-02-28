package model

import (
	"time"

	"gorm.io/datatypes"
)

type BuildConfig struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID      string         `json:"project_id" gorm:"type:uuid;not null;index"`
	ServiceID      string         `json:"service_id" gorm:"type:uuid;not null;index"`
	Name           string         `json:"name" gorm:"size:128;not null"`
	TemplateID     *string        `json:"template_id" gorm:"type:uuid"`
	RepoURL        string         `json:"repo_url" gorm:"size:512;not null"`
	Branch         string         `json:"branch" gorm:"size:128;default:'main'"`
	BuildEnv       datatypes.JSON `json:"build_env"`
	BuildScript    string         `json:"build_script" gorm:"type:text"`
	DockerfilePath string         `json:"dockerfile_path" gorm:"size:256;default:'Dockerfile'"`
	DockerContext  string         `json:"docker_context" gorm:"size:256;default:'.'"`
	ImageRepo      string         `json:"image_repo" gorm:"size:256;not null"`
	TagStrategy    string         `json:"tag_strategy" gorm:"size:32;default:'branch-commit'"`
	CacheEnabled   bool           `json:"cache_enabled" gorm:"default:true"`
	Variables      datatypes.JSON `json:"variables"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type BuildRun struct {
	ID            string        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BuildConfigID string        `json:"build_config_id" gorm:"type:uuid;not null;index"`
	RunNumber     int           `json:"run_number" gorm:"not null"`
	Status        string        `json:"status" gorm:"size:32;default:'pending'"` // pending, running, succeeded, failed, cancelled
	Branch        string        `json:"branch" gorm:"size:128"`
	CommitSHA     string        `json:"commit_sha" gorm:"size:64"`
	CommitMessage string        `json:"commit_message" gorm:"type:text"`
	ImageTag      string        `json:"image_tag" gorm:"size:256"`
	ImageDigest   string        `json:"image_digest" gorm:"size:128"`
	TektonRef     string        `json:"tekton_ref" gorm:"size:256"`
	LogPath       string        `json:"log_path" gorm:"size:512"`
	TriggeredBy   *string       `json:"triggered_by" gorm:"type:uuid"`
	StartedAt     *time.Time    `json:"started_at"`
	FinishedAt    *time.Time    `json:"finished_at"`
	DurationSec   *int          `json:"duration_sec"`
	CreatedAt     time.Time     `json:"created_at"`
	BuildConfig   *BuildConfig  `json:"build_config,omitempty" gorm:"foreignKey:BuildConfigID"`
}

type BuildTemplate struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string         `json:"name" gorm:"size:128;not null"`
	Language      string         `json:"language" gorm:"size:32;not null;index"`
	Framework     string         `json:"framework" gorm:"size:64"`
	Description   string         `json:"description" gorm:"type:text"`
	BuildEnv      datatypes.JSON `json:"build_env"`
	BuildScript   string         `json:"build_script" gorm:"type:text;not null"`
	DockerfileTpl string         `json:"dockerfile_tpl" gorm:"type:text"`
	TektonTaskTpl string         `json:"tekton_task_tpl" gorm:"type:text;not null"`
	IsSystem      bool           `json:"is_system" gorm:"default:false"`
	CreatedBy     *string        `json:"created_by" gorm:"type:uuid"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (BuildConfig) TableName() string   { return "build_configs" }
func (BuildRun) TableName() string      { return "build_runs" }
func (BuildTemplate) TableName() string { return "build_templates" }
