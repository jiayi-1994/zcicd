package repository

import (
	"context"

	"github.com/zcicd/zcicd-server/internal/workflow/model"

	"gorm.io/gorm"
)

type BuildRepository struct {
	db *gorm.DB
}

func NewBuildRepository(db *gorm.DB) *BuildRepository {
	return &BuildRepository{db: db}
}

func (r *BuildRepository) CreateConfig(ctx context.Context, cfg *model.BuildConfig) error {
	return r.db.WithContext(ctx).Create(cfg).Error
}

func (r *BuildRepository) FindConfigByID(ctx context.Context, id string) (*model.BuildConfig, error) {
	var cfg model.BuildConfig
	err := r.db.WithContext(ctx).First(&cfg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *BuildRepository) UpdateConfig(ctx context.Context, cfg *model.BuildConfig) error {
	return r.db.WithContext(ctx).Save(cfg).Error
}

func (r *BuildRepository) DeleteConfig(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.BuildConfig{}, "id = ?", id).Error
}

func (r *BuildRepository) ListConfigsByProject(ctx context.Context, projectID string, page, pageSize int) ([]model.BuildConfig, int64, error) {
	var list []model.BuildConfig
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BuildConfig{}).Where("project_id = ?", projectID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *BuildRepository) ListConfigsByService(ctx context.Context, serviceID string) ([]model.BuildConfig, error) {
	var list []model.BuildConfig
	err := r.db.WithContext(ctx).
		Where("service_id = ?", serviceID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *BuildRepository) CreateRun(ctx context.Context, run *model.BuildRun) error {
	return r.db.WithContext(ctx).Create(run).Error
}

func (r *BuildRepository) FindRunByID(ctx context.Context, id string) (*model.BuildRun, error) {
	var run model.BuildRun
	err := r.db.WithContext(ctx).
		Preload("BuildConfig").
		First(&run, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *BuildRepository) UpdateRun(ctx context.Context, run *model.BuildRun) error {
	return r.db.WithContext(ctx).Save(run).Error
}

func (r *BuildRepository) ListRuns(ctx context.Context, buildConfigID string, page, pageSize int) ([]model.BuildRun, int64, error) {
	var list []model.BuildRun
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BuildRun{}).Where("build_config_id = ?", buildConfigID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *BuildRepository) ListAllRuns(ctx context.Context, projectID string, page, pageSize int) ([]model.BuildRun, int64, error) {
	var list []model.BuildRun
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BuildRun{}).
		Joins("JOIN build_configs ON build_configs.id = build_runs.build_config_id").
		Where("build_configs.project_id = ?", projectID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("build_runs.created_at DESC").Offset(offset).Limit(pageSize).
		Preload("BuildConfig").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *BuildRepository) GetNextRunNumber(ctx context.Context, buildConfigID string) (int, error) {
	var maxNumber int
	r.db.WithContext(ctx).Model(&model.BuildRun{}).
		Where("build_config_id = ?", buildConfigID).
		Select("COALESCE(MAX(run_number), 0)").
		Scan(&maxNumber)
	return maxNumber + 1, nil
}
