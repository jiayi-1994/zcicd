package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Project struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string         `json:"name" gorm:"size:100;not null"`
	Identifier    string         `json:"identifier" gorm:"uniqueIndex;size:50;not null"`
	Description   string         `json:"description"`
	OwnerID       string         `json:"owner_id" gorm:"type:uuid;not null"`
	RepoURL       string         `json:"repo_url"`
	DefaultBranch string         `json:"default_branch" gorm:"default:'main'"`
	Visibility    string         `json:"visibility" gorm:"default:'private'"`
	Status        string         `json:"status" gorm:"default:'active'"`
	Settings      datatypes.JSON `json:"settings"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Project) TableName() string { return "projects" }
