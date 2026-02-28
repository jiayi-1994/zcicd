package repository

import (
	"context"

	"github.com/zcicd/zcicd-server/internal/workflow/model"

	"gorm.io/gorm"
)

type TemplateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db: db}
}

func (r *TemplateRepository) Create(ctx context.Context, tpl *model.BuildTemplate) error {
	return r.db.WithContext(ctx).Create(tpl).Error
}

func (r *TemplateRepository) FindByID(ctx context.Context, id string) (*model.BuildTemplate, error) {
	var tpl model.BuildTemplate
	err := r.db.WithContext(ctx).First(&tpl, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (r *TemplateRepository) Update(ctx context.Context, tpl *model.BuildTemplate) error {
	return r.db.WithContext(ctx).Save(tpl).Error
}

func (r *TemplateRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.BuildTemplate{}, "id = ?", id).Error
}

func (r *TemplateRepository) List(ctx context.Context, language string) ([]model.BuildTemplate, error) {
	var list []model.BuildTemplate
	query := r.db.WithContext(ctx).Model(&model.BuildTemplate{})
	if language != "" {
		query = query.Where("language = ?", language)
	}
	err := query.Order("language, name").Find(&list).Error
	return list, err
}

func (r *TemplateRepository) ListSystem(ctx context.Context) ([]model.BuildTemplate, error) {
	var list []model.BuildTemplate
	err := r.db.WithContext(ctx).
		Where("is_system = ?", true).
		Order("language, name").
		Find(&list).Error
	return list, err
}
