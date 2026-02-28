package model

import "time"

type EnvVariable struct {
	ID            string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EnvironmentID string    `json:"environment_id" gorm:"type:uuid;not null;index"`
	VarKey        string    `json:"var_key" gorm:"size:256;not null"`
	VarValue      string    `json:"var_value" gorm:"type:text"`
	IsSecret      bool      `json:"is_secret" gorm:"default:false"`
	Description   string    `json:"description" gorm:"size:512"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (EnvVariable) TableName() string { return "env_variables" }

type EnvResourceQuota struct {
	ID            string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EnvironmentID string    `json:"environment_id" gorm:"type:uuid;not null;uniqueIndex"`
	CPURequest    string    `json:"cpu_request" gorm:"size:32"`
	CPULimit      string    `json:"cpu_limit" gorm:"size:32"`
	MemoryRequest string    `json:"memory_request" gorm:"size:32"`
	MemoryLimit   string    `json:"memory_limit" gorm:"size:32"`
	PodLimit      int       `json:"pod_limit"`
	StorageLimit  string    `json:"storage_limit" gorm:"size:32"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (EnvResourceQuota) TableName() string { return "env_resource_quotas" }
