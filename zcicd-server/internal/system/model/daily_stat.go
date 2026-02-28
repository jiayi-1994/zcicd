package model

import "time"

type DailyStat struct {
	ID                string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	StatDate          time.Time `json:"stat_date" gorm:"type:date;uniqueIndex;not null"`
	TotalBuilds       int       `json:"total_builds" gorm:"default:0"`
	SuccessfulBuilds  int       `json:"successful_builds" gorm:"default:0"`
	FailedBuilds      int       `json:"failed_builds" gorm:"default:0"`
	TotalDeploys      int       `json:"total_deploys" gorm:"default:0"`
	SuccessfulDeploys int       `json:"successful_deploys" gorm:"default:0"`
	FailedDeploys     int       `json:"failed_deploys" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (DailyStat) TableName() string { return "daily_stats" }
