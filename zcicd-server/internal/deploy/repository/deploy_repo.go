package repository

import (
	"github.com/zcicd/zcicd-server/internal/deploy/model"
	"gorm.io/gorm"
)

type DeployRepository struct {
	db *gorm.DB
}

func NewDeployRepository(db *gorm.DB) *DeployRepository {
	return &DeployRepository{db: db}
}

// DeployConfig CRUD

func (r *DeployRepository) CreateConfig(config *model.DeployConfig) error {
	return r.db.Create(config).Error
}

func (r *DeployRepository) GetConfig(id string) (*model.DeployConfig, error) {
	var config model.DeployConfig
	err := r.db.Where("id = ?", id).First(&config).Error
	return &config, err
}

func (r *DeployRepository) GetConfigByArgoApp(argoAppName string) (*model.DeployConfig, error) {
	var config model.DeployConfig
	err := r.db.Where("argo_app_name = ?", argoAppName).First(&config).Error
	return &config, err
}

func (r *DeployRepository) UpdateConfig(config *model.DeployConfig) error {
	return r.db.Save(config).Error
}

func (r *DeployRepository) DeleteConfig(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.DeployConfig{}).Error
}

func (r *DeployRepository) ListConfigs(projectID string, page, pageSize int) ([]model.DeployConfig, int64, error) {
	var configs []model.DeployConfig
	var total int64
	query := r.db.Where("project_id = ?", projectID)
	query.Model(&model.DeployConfig{}).Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&configs).Error
	return configs, total, err
}

func (r *DeployRepository) ListConfigsByEnv(projectID, envID string) ([]model.DeployConfig, error) {
	var configs []model.DeployConfig
	err := r.db.Where("project_id = ? AND environment_id = ?", projectID, envID).
		Order("created_at DESC").Find(&configs).Error
	return configs, err
}
