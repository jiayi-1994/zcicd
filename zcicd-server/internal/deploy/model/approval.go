package model

import "time"

type ApprovalRecord struct {
	ID              string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DeployHistoryID string     `json:"deploy_history_id" gorm:"type:uuid;not null;index"`
	EnvironmentID   string     `json:"environment_id" gorm:"type:uuid;not null"`
	RequestedBy     string     `json:"requested_by" gorm:"type:uuid;not null"`
	ApproverID      *string    `json:"approver_id" gorm:"type:uuid"`
	Status          string     `json:"status" gorm:"size:32;not null;default:'pending'"`
	Comment         string     `json:"comment" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at"`
	DecidedAt       *time.Time `json:"decided_at"`

	DeployHistory DeployHistory `json:"deploy_history,omitempty" gorm:"foreignKey:DeployHistoryID"`
}

func (ApprovalRecord) TableName() string { return "approval_records" }
