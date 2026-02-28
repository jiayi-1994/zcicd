package service

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"github.com/zcicd/zcicd-server/internal/system/repository"
)

type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) List(page, pageSize int, filters map[string]string) ([]model.AuditLog, int64, error) {
	return s.repo.List(page, pageSize, filters)
}
