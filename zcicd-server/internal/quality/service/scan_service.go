package service

import (
	"time"

	"github.com/zcicd/zcicd-server/internal/quality/model"
	"github.com/zcicd/zcicd-server/internal/quality/repository"
)

type ScanService struct {
	repo *repository.ScanRepository
}

func NewScanService(repo *repository.ScanRepository) *ScanService {
	return &ScanService{repo: repo}
}

func (s *ScanService) CreateConfig(projectID string, req CreateScanConfigReq) (*model.ScanConfig, error) {
	c := &model.ScanConfig{
		ProjectID:       projectID,
		Name:            req.Name,
		ScanType:        req.ScanType,
		SonarProjectKey: req.SonarProjectKey,
		Enabled:         true,
	}
	if c.ScanType == "" {
		c.ScanType = "sonar"
	}
	return c, s.repo.CreateConfig(c)
}

func (s *ScanService) GetConfig(id string) (*model.ScanConfig, error) {
	return s.repo.GetConfig(id)
}

func (s *ScanService) UpdateConfig(id string, req UpdateScanConfigReq) (*model.ScanConfig, error) {
	c, err := s.repo.GetConfig(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.ScanType != "" {
		c.ScanType = req.ScanType
	}
	if req.SonarProjectKey != "" {
		c.SonarProjectKey = req.SonarProjectKey
	}
	if req.Enabled != nil {
		c.Enabled = *req.Enabled
	}
	return c, s.repo.UpdateConfig(c)
}

func (s *ScanService) DeleteConfig(id string) error {
	return s.repo.DeleteConfig(id)
}

func (s *ScanService) ListConfigs(projectID string, page, pageSize int) ([]model.ScanConfig, int64, error) {
	return s.repo.ListConfigs(projectID, page, pageSize)
}

func (s *ScanService) TriggerRun(configID string) (*model.ScanRun, error) {
	now := time.Now()
	run := &model.ScanRun{
		ScanConfigID: configID,
		Status:       "running",
		StartedAt:    &now,
	}
	return run, s.repo.CreateRun(run)
}

func (s *ScanService) GetRun(id string) (*model.ScanRun, error) {
	return s.repo.GetRun(id)
}

func (s *ScanService) ListRuns(configID string, page, pageSize int) ([]model.ScanRun, int64, error) {
	return s.repo.ListRuns(configID, page, pageSize)
}
