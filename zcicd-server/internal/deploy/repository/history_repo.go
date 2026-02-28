package repository

import (
	"github.com/zcicd/zcicd-server/internal/deploy/model"
)

// DeployHistory CRUD (on same repo struct)

func (r *DeployRepository) CreateHistory(h *model.DeployHistory) error {
	return r.db.Create(h).Error
}

func (r *DeployRepository) GetHistory(id string) (*model.DeployHistory, error) {
	var h model.DeployHistory
	err := r.db.Preload("DeployConfig").Where("id = ?", id).First(&h).Error
	return &h, err
}

func (r *DeployRepository) UpdateHistory(h *model.DeployHistory) error {
	return r.db.Save(h).Error
}

func (r *DeployRepository) ListHistories(configID string, page, pageSize int) ([]model.DeployHistory, int64, error) {
	var histories []model.DeployHistory
	var total int64
	query := r.db.Where("deploy_config_id = ?", configID)
	query.Model(&model.DeployHistory{}).Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&histories).Error
	return histories, total, err
}

func (r *DeployRepository) GetLatestHistory(configID string) (*model.DeployHistory, error) {
	var h model.DeployHistory
	err := r.db.Where("deploy_config_id = ? AND status = 'succeeded'", configID).
		Order("created_at DESC").First(&h).Error
	return &h, err
}
