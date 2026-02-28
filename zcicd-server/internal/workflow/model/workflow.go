package model

import (
	"time"

	"gorm.io/datatypes"
)

type Workflow struct {
	ID            string          `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID     string          `json:"project_id" gorm:"type:uuid;not null;index"`
	Name          string          `json:"name" gorm:"size:128;not null"`
	Description   string          `json:"description" gorm:"type:text"`
	TriggerType   string          `json:"trigger_type" gorm:"size:32;default:'manual'"` // manual, webhook, cron, api
	TriggerConfig datatypes.JSON  `json:"trigger_config"`
	Enabled       bool            `json:"enabled" gorm:"default:true"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Stages        []WorkflowStage `json:"stages,omitempty" gorm:"foreignKey:WorkflowID"`
}

type WorkflowStage struct {
	ID         string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkflowID string         `json:"workflow_id" gorm:"type:uuid;not null;index"`
	Name       string         `json:"name" gorm:"size:128;not null"`
	StageType  string         `json:"stage_type" gorm:"size:32"` // build, test, deploy, custom, approval
	SortOrder  int            `json:"sort_order" gorm:"default:0"`
	Config     datatypes.JSON `json:"config"`
	Timeout    int            `json:"timeout" gorm:"default:3600"`
	Parallel   bool           `json:"parallel" gorm:"default:false"`
	Enabled    bool           `json:"enabled" gorm:"default:true"`
	Jobs       []StageJob     `json:"jobs,omitempty" gorm:"foreignKey:StageID"`
}

type StageJob struct {
	ID         string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StageID    string         `json:"stage_id" gorm:"type:uuid;not null;index"`
	Name       string         `json:"name" gorm:"size:128;not null"`
	JobType    string         `json:"job_type" gorm:"size:32;not null"` // build, test, deploy, custom, approval
	SortOrder  int            `json:"sort_order" gorm:"default:0"`
	Config     datatypes.JSON `json:"config"`
	TimeoutSec int            `json:"timeout_sec" gorm:"default:3600"`
	Enabled    bool           `json:"enabled" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (Workflow) TableName() string      { return "workflows" }
func (WorkflowStage) TableName() string { return "workflow_stages" }
func (StageJob) TableName() string      { return "stage_jobs" }
