package repository

import (
	"context"

	"github.com/zcicd/zcicd-server/internal/project/model"
	"gorm.io/gorm"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (r *EnvironmentRepository) Create(ctx context.Context, env *model.Environment) error {
	return r.db.WithContext(ctx).Create(env).Error
}

func (r *EnvironmentRepository) FindByID(ctx context.Context, id string) (*model.Environment, error) {
	var env model.Environment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&env).Error; err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *EnvironmentRepository) Update(ctx context.Context, env *model.Environment) error {
	return r.db.WithContext(ctx).Save(env).Error
}

func (r *EnvironmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Environment{}).Error
}

func (r *EnvironmentRepository) ListByProject(ctx context.Context, projectID string) ([]model.Environment, error) {
	var envs []model.Environment
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("created_at ASC").Find(&envs).Error; err != nil {
		return nil, err
	}
	return envs, nil
}

func (r *EnvironmentRepository) FindByProjectAndName(ctx context.Context, projectID, name string) (*model.Environment, error) {
	var env model.Environment
	if err := r.db.WithContext(ctx).Where("project_id = ? AND name = ?", projectID, name).First(&env).Error; err != nil {
		return nil, err
	}
	return &env, nil
}
