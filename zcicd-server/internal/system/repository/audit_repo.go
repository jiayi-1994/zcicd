package repository

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditRepository) List(page, pageSize int, filters map[string]string) ([]model.AuditLog, int64, error) {
	var list []model.AuditLog
	var total int64
	q := r.db.Model(&model.AuditLog{})
	if v, ok := filters["user_id"]; ok && v != "" {
		q = q.Where("user_id = ?", v)
	}
	if v, ok := filters["action"]; ok && v != "" {
		q = q.Where("action = ?", v)
	}
	if v, ok := filters["resource_type"]; ok && v != "" {
		q = q.Where("resource_type = ?", v)
	}
	if v, ok := filters["project_id"]; ok && v != "" {
		q = q.Where("project_id = ?", v)
	}
	q.Count(&total)
	err := q.Offset((page-1)*pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&list).Error
	return list, total, err
}
