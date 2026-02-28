package model

import "time"

type TestRun struct {
	ID           string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TestConfigID string     `json:"test_config_id" gorm:"type:uuid;not null;index"`
	BuildRunID   string     `json:"build_run_id" gorm:"type:uuid"`
	Status       string     `json:"status" gorm:"size:32;not null;default:'pending'"`
	Total        int        `json:"total" gorm:"default:0"`
	Passed       int        `json:"passed" gorm:"default:0"`
	Failed       int        `json:"failed" gorm:"default:0"`
	Skipped      int        `json:"skipped" gorm:"default:0"`
	Coverage     *float64   `json:"coverage" gorm:"type:numeric(5,2)"`
	Duration     *int       `json:"duration"`
	ReportURL    string     `json:"report_url" gorm:"size:512"`
	ErrorMessage string     `json:"error_message,omitempty"`
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (TestRun) TableName() string { return "test_runs" }
