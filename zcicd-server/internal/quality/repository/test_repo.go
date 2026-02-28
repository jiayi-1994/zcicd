package repository

import (
	"github.com/zcicd/zcicd-server/internal/quality/model"
	"gorm.io/gorm"
)

type TestRepository struct {
	db *gorm.DB
}

func NewTestRepository(db *gorm.DB) *TestRepository {
	return &TestRepository{db: db}
}

func (r *TestRepository) CreateConfig(c *model.TestConfig) error {
	return r.db.Create(c).Error
}

func (r *TestRepository) GetConfig(id string) (*model.TestConfig, error) {
	var c model.TestConfig
	err := r.db.Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *TestRepository) UpdateConfig(c *model.TestConfig) error {
	return r.db.Save(c).Error
}

func (r *TestRepository) DeleteConfig(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.TestConfig{}).Error
}

func (r *TestRepository) ListConfigs(projectID string, page, pageSize int) ([]model.TestConfig, int64, error) {
	var list []model.TestConfig
	var total int64
	q := r.db.Where("project_id = ?", projectID)
	q.Model(&model.TestConfig{}).Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

func (r *TestRepository) CreateRun(run *model.TestRun) error {
	return r.db.Create(run).Error
}

func (r *TestRepository) GetRun(id string) (*model.TestRun, error) {
	var run model.TestRun
	err := r.db.Where("id = ?", id).First(&run).Error
	return &run, err
}

func (r *TestRepository) UpdateRun(run *model.TestRun) error {
	return r.db.Save(run).Error
}

func (r *TestRepository) ListRuns(configID string, page, pageSize int) ([]model.TestRun, int64, error) {
	var list []model.TestRun
	var total int64
	q := r.db.Where("test_config_id = ?", configID)
	q.Model(&model.TestRun{}).Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	return list, total, err
}
