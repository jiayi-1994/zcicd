package model

import "time"

type User struct {
	ID           string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string     `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Email        string     `json:"email" gorm:"uniqueIndex;size:100;not null"`
	PasswordHash string     `json:"-" gorm:"not null"`
	DisplayName  string     `json:"display_name" gorm:"size:100"`
	AvatarURL    string     `json:"avatar_url"`
	Status       string     `json:"status" gorm:"default:'active'"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserRole struct {
	UserID    string    `json:"user_id" gorm:"type:uuid;primaryKey"`
	Role      string    `json:"role" gorm:"primaryKey;size:50"`
	ScopeType string    `json:"scope_type" gorm:"primaryKey;size:20;default:'system'"`
	ScopeID   string    `json:"scope_id" gorm:"primaryKey;type:uuid;default:'00000000-0000-0000-0000-000000000000'"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
