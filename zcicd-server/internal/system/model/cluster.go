package model

import "time"

type Cluster struct {
	ID            string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `json:"name" gorm:"size:128;uniqueIndex;not null"`
	DisplayName   string    `json:"display_name" gorm:"size:256"`
	Description   string    `json:"description"`
	Provider      string    `json:"provider" gorm:"size:32"`
	APIServerURL  string    `json:"api_server_url" gorm:"size:512;not null"`
	KubeConfigRef string    `json:"kube_config_ref" gorm:"size:256"`
	Status        string    `json:"status" gorm:"size:16;default:'connected'"`
	NodeCount     int       `json:"node_count" gorm:"default:0"`
	Version       string    `json:"version" gorm:"size:32"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Cluster) TableName() string { return "clusters" }
