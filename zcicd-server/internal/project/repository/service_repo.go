package repository

import (
	"context"

	"github.com/zcicd/zcicd-server/internal/project/model"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(ctx context.Context, svc *model.Service) error {
	return r.db.WithContext(ctx).Create(svc).Error
}

func (r *ServiceRepository) FindByID(ctx context.Context, id string) (*model.Service, error) {
	var svc model.Service
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&svc).Error; err != nil {
		return nil, err
	}
	return &svc, nil
}

func (r *ServiceRepository) Update(ctx context.Context, svc *model.Service) error {
	return r.db.WithContext(ctx).Save(svc).Error
}

func (r *ServiceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Service{}).Error
}

func (r *ServiceRepository) ListByProject(ctx context.Context, projectID string, page, pageSize int) ([]model.Service, int64, error) {
	var services []model.Service
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Service{}).Where("project_id = ?", projectID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&services).Error; err != nil {
		return nil, 0, err
	}
	return services, total, nil
}

func (r *ServiceRepository) FindByProjectAndName(ctx context.Context, projectID, name string) (*model.Service, error) {
	var svc model.Service
	if err := r.db.WithContext(ctx).Where("project_id = ? AND name = ?", projectID, name).First(&svc).Error; err != nil {
		return nil, err
	}
	return &svc, nil
}
