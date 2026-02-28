package model

import "time"

type HelmChart struct {
	ID            string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `json:"name" gorm:"size:256;not null"`
	RepoURL       string    `json:"repo_url" gorm:"size:512;not null"`
	LatestVersion string    `json:"latest_version" gorm:"size:64"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (HelmChart) TableName() string { return "helm_charts" }
