package model

import "time"

type QualityGate struct {
	ID                  string   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID           string   `json:"project_id" gorm:"type:uuid;not null;uniqueIndex"`
	MinCoverage         *float64 `json:"min_coverage" gorm:"type:numeric(5,2);default:80.00"`
	MaxBugs             int      `json:"max_bugs" gorm:"default:0"`
	MaxVulnerabilities  int      `json:"max_vulnerabilities" gorm:"default:0"`
	MaxCodeSmells       int      `json:"max_code_smells" gorm:"default:50"`
	MaxDuplications     *float64 `json:"max_duplications" gorm:"type:numeric(5,2);default:5.00"`
	BlockDeploy         bool     `json:"block_deploy" gorm:"default:false"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (QualityGate) TableName() string { return "quality_gates" }
