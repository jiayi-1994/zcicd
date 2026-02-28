package repository

import (
	"context"

	"github.com/zcicd/zcicd-server/internal/project/model"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) FindByID(ctx context.Context, id string) (*model.Project, error) {
	var project model.Project
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindByIdentifier(ctx context.Context, identifier string) (*model.Project, error) {
	var project model.Project
	if err := r.db.WithContext(ctx).Where("identifier = ?", identifier).First(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Project{}).Error
}

func (r *ProjectRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Project{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR identifier LIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}
