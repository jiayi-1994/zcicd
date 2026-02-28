package model

import "time"

type ImageScan struct {
	ID         string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RegistryID string     `json:"registry_id" gorm:"type:uuid;not null;index"`
	ImageName  string     `json:"image_name" gorm:"size:512;not null"`
	Tag        string     `json:"tag" gorm:"size:256;not null"`
	Status     string     `json:"status" gorm:"size:32;not null;default:'pending'"`
	Critical   int        `json:"critical" gorm:"default:0"`
	High       int        `json:"high" gorm:"default:0"`
	Medium     int        `json:"medium" gorm:"default:0"`
	Low        int        `json:"low" gorm:"default:0"`
	ReportURL  string     `json:"report_url" gorm:"size:512"`
	ScannedAt  *time.Time `json:"scanned_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (ImageScan) TableName() string { return "image_scans" }
