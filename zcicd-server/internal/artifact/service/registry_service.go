package service

import (
	"github.com/zcicd/zcicd-server/internal/artifact/model"
	"github.com/zcicd/zcicd-server/internal/artifact/repository"
)

type RegistryService struct {
	repo *repository.RegistryRepository
}

func NewRegistryService(repo *repository.RegistryRepository) *RegistryService {
	return &RegistryService{repo: repo}
}

func (s *RegistryService) Create(req CreateRegistryReq) (*model.ImageRegistry, error) {
	reg := &model.ImageRegistry{
		Name:         req.Name,
		RegistryType: req.RegistryType,
		Endpoint:     req.Endpoint,
		Username:     req.Username,
		IsDefault:    req.IsDefault,
	}
	if reg.RegistryType == "" {
		reg.RegistryType = "harbor"
	}
	return reg, s.repo.Create(reg)
}

func (s *RegistryService) Get(id string) (*model.ImageRegistry, error) {
	return s.repo.Get(id)
}

func (s *RegistryService) Update(id string, req UpdateRegistryReq) (*model.ImageRegistry, error) {
	reg, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		reg.Name = req.Name
	}
	if req.Endpoint != "" {
		reg.Endpoint = req.Endpoint
	}
	if req.Username != "" {
		reg.Username = req.Username
	}
	if req.IsDefault != nil {
		reg.IsDefault = *req.IsDefault
	}
	return reg, s.repo.Update(reg)
}

func (s *RegistryService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *RegistryService) List() ([]model.ImageRegistry, error) {
	return s.repo.List()
}
