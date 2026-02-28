package repository

import (
	"github.com/zcicd/zcicd-server/internal/quality/model"
	"gorm.io/gorm"
)

type ScanRepository struct {
	db *gorm.DB
}

func NewScanRepository(db *gorm.DB) *ScanRepository {
	return &ScanRepository{db: db}
}

func (r *ScanRepository) CreateConfig(c *model.ScanConfig) error {
	return r.db.Create(c).Error
}

func (r *ScanRepository) GetConfig(id string) (*model.ScanConfig, error) {
	var c model.ScanConfig
	err := r.db.Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *ScanRepository) UpdateConfig(c *model.ScanConfig) error {
	return r.db.Save(c).Error
}

func (r *ScanRepository) DeleteConfig(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.ScanConfig{}).Error
}

func (r *ScanRepository) ListConfigs(projectID string, page, pageSize int) ([]model.ScanConfig, int64, error) {
	var list []model.ScanConfig
	var total int64
	q := r.db.Where("project_id = ?", projectID)
	q.Model(&model.ScanConfig{}).Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *ScanRepository) CreateRun(run *model.ScanRun) error {
	return r.db.Create(run).Error
}

func (r *ScanRepository) GetRun(id string) (*model.ScanRun, error) {
	var run model.ScanRun
	err := r.db.Where("id = ?", id).First(&run).Error
	return &run, err
}

func (r *ScanRepository) UpdateRun(run *model.ScanRun) error {
	return r.db.Save(run).Error
}

func (r *ScanRepository) ListRuns(configID string, page, pageSize int) ([]model.ScanRun, int64, error) {
	var list []model.ScanRun
	var total int64
	q := r.db.Where("scan_config_id = ?", configID)
	q.Model(&model.ScanRun{}).Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}
