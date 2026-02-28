package service

import (
	"time"

	"github.com/zcicd/zcicd-server/internal/quality/model"
	"github.com/zcicd/zcicd-server/internal/quality/repository"
)

type TestService struct {
	repo *repository.TestRepository
}

func NewTestService(repo *repository.TestRepository) *TestService {
	return &TestService{repo: repo}
}

func (s *TestService) CreateConfig(projectID string, req CreateTestConfigReq) (*model.TestConfig, error) {
	c := &model.TestConfig{
		ProjectID: projectID,
		Name:      req.Name,
		TestType:  req.TestType,
		Framework: req.Framework,
		Command:   req.Command,
		Timeout:   req.Timeout,
		Enabled:   true,
	}
	if c.TestType == "" {
		c.TestType = "unit"
	}
	if c.Timeout == 0 {
		c.Timeout = 3600
	}
	return c, s.repo.CreateConfig(c)
}

func (s *TestService) GetConfig(id string) (*model.TestConfig, error) {
	return s.repo.GetConfig(id)
}

func (s *TestService) UpdateConfig(id string, req UpdateTestConfigReq) (*model.TestConfig, error) {
	c, err := s.repo.GetConfig(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		c.Name = req.Name
	}
	if req.TestType != "" {
		c.TestType = req.TestType
	}
	if req.Framework != "" {
		c.Framework = req.Framework
	}
	if req.Command != "" {
		c.Command = req.Command
	}
	if req.Timeout > 0 {
		c.Timeout = req.Timeout
	}
	if req.Enabled != nil {
		c.Enabled = *req.Enabled
	}
	return c, s.repo.UpdateConfig(c)
}

func (s *TestService) DeleteConfig(id string) error {
	return s.repo.DeleteConfig(id)
}

func (s *TestService) ListConfigs(projectID string, page, pageSize int) ([]model.TestConfig, int64, error) {
	return s.repo.ListConfigs(projectID, page, pageSize)
}

func (s *TestService) TriggerRun(configID string) (*model.TestRun, error) {
	now := time.Now()
	run := &model.TestRun{
		TestConfigID: configID,
		Status:       "running",
		StartedAt:    &now,
	}
	return run, s.repo.CreateRun(run)
}

func (s *TestService) GetRun(id string) (*model.TestRun, error) {
	return s.repo.GetRun(id)
}

func (s *TestService) ListRuns(configID string, page, pageSize int) ([]model.TestRun, int64, error) {
	return s.repo.ListRuns(configID, page, pageSize)
}
