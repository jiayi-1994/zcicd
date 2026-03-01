package service

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"github.com/zcicd/zcicd-server/internal/system/repository"
)

type IntegrationService struct {
	repo *repository.IntegrationRepository
}

func NewIntegrationService(repo *repository.IntegrationRepository) *IntegrationService {
	return &IntegrationService{repo: repo}
}

func (s *IntegrationService) Create(req CreateIntegrationReq) (*model.Integration, error) {
	i := &model.Integration{
		Name:      req.Name,
		Type:      req.Type,
		Provider:  req.Provider,
		ConfigEnc: []byte(req.Config),
		Status:    "active",
	}
	return i, s.repo.Create(i)
}

func (s *IntegrationService) Get(id string) (*model.Integration, error) {
	return s.repo.Get(id)
}

func (s *IntegrationService) Update(id string, req UpdateIntegrationReq) (*model.Integration, error) {
	i, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		i.Name = req.Name
	}
	if req.Status != "" {
		i.Status = req.Status
	}
	if req.Config != "" {
		i.ConfigEnc = []byte(req.Config)
	}
	return i, s.repo.Update(i)
}

func (s *IntegrationService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *IntegrationService) List() ([]model.Integration, error) {
	return s.repo.List()
}
