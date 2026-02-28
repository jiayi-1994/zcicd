package model

import (
	"time"

	"gorm.io/datatypes"
)

type WorkflowRun struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkflowID   string         `json:"workflow_id" gorm:"type:uuid;not null;index"`
	RunNumber    int            `json:"run_number" gorm:"not null"`
	Status       string         `json:"status" gorm:"size:32;default:'pending'"` // pending, running, succeeded, failed, cancelled, waiting_approval
	TriggerType  string         `json:"trigger_type" gorm:"size:32;not null"`
	TriggeredBy  *string        `json:"triggered_by" gorm:"type:uuid"`
	InputParams  datatypes.JSON `json:"input_params"`
	StagesStatus datatypes.JSON `json:"stages_status"`
	TektonRefs   datatypes.JSON `json:"tekton_refs"`
	StartedAt    *time.Time     `json:"started_at"`
	FinishedAt   *time.Time     `json:"finished_at"`
	DurationSec  *int           `json:"duration_sec"`
	ErrorMessage string         `json:"error_message" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	Workflow     *Workflow      `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
}

func (WorkflowRun) TableName() string { return "workflow_runs" }
