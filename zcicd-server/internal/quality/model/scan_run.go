package model

import "time"

type ScanRun struct {
	ID             string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ScanConfigID   string     `json:"scan_config_id" gorm:"type:uuid;not null;index"`
	Status         string     `json:"status" gorm:"size:32;not null;default:'pending'"`
	Bugs           int        `json:"bugs" gorm:"default:0"`
	Vulnerabilities int       `json:"vulnerabilities" gorm:"default:0"`
	CodeSmells     int        `json:"code_smells" gorm:"default:0"`
	Coverage       *float64   `json:"coverage" gorm:"type:numeric(5,2)"`
	Duplications   *float64   `json:"duplications" gorm:"type:numeric(5,2)"`
	QualityRating  string     `json:"quality_rating" gorm:"size:1"`
	GateStatus     string     `json:"gate_status" gorm:"size:16"`
	ReportURL      string     `json:"report_url" gorm:"size:512"`
	ErrorMessage   string     `json:"error_message,omitempty"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (ScanRun) TableName() string { return "scan_runs" }
