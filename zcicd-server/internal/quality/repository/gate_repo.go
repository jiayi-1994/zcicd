package repository

import (
	"github.com/zcicd/zcicd-server/internal/quality/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QualityGateRepository struct {
	db *gorm.DB
}

func NewQualityGateRepository(db *gorm.DB) *QualityGateRepository {
	return &QualityGateRepository{db: db}
}

func (r *QualityGateRepository) Get(projectID string) (*model.QualityGate, error) {
	var g model.QualityGate
	err := r.db.Where("project_id = ?", projectID).First(&g).Error
	return &g, err
}

func (r *QualityGateRepository) Upsert(g *model.QualityGate) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "project_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"min_coverage", "max_bugs", "max_vulnerabilities", "max_code_smells", "max_duplications", "block_deploy", "updated_at"}),
	}).Create(g).Error
}
