package model

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       string         `json:"user_id" gorm:"type:uuid"`
	Username     string         `json:"username" gorm:"size:64"`
	Action       string         `json:"action" gorm:"size:64;not null"`
	ResourceType string         `json:"resource_type" gorm:"size:64;not null"`
	ResourceID   string         `json:"resource_id" gorm:"type:uuid"`
	ResourceName string         `json:"resource_name" gorm:"size:256"`
	ProjectID    string         `json:"project_id" gorm:"type:uuid"`
	Detail       datatypes.JSON `json:"detail"`
	IPAddress    string         `json:"ip_address" gorm:"size:45"`
	UserAgent    string         `json:"user_agent" gorm:"size:512"`
	RequestID    string         `json:"request_id" gorm:"size:64"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
