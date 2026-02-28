package service

import (
	"github.com/zcicd/zcicd-server/internal/deploy/model"
	"github.com/zcicd/zcicd-server/internal/deploy/repository"
)

type EnvService struct {
	envRepo *repository.EnvRepository
}

func NewEnvService(envRepo *repository.EnvRepository) *EnvService {
	return &EnvService{envRepo: envRepo}
}

// Variables

func (s *EnvService) ListVariables(envID string) ([]model.EnvVariable, error) {
	return s.envRepo.ListVariables(envID)
}

func (s *EnvService) CreateVariable(envID string, req EnvVariableReq) (*model.EnvVariable, error) {
	v := &model.EnvVariable{
		EnvironmentID: envID,
		VarKey:        req.VarKey,
		VarValue:      req.VarValue,
		IsSecret:      req.IsSecret,
		Description:   req.Description,
	}
	if err := s.envRepo.CreateVariable(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *EnvService) UpdateVariable(id string, req EnvVariableReq) (*model.EnvVariable, error) {
	v, err := s.envRepo.GetVariable(id)
	if err != nil {
		return nil, err
	}
	v.VarKey = req.VarKey
	v.VarValue = req.VarValue
	v.IsSecret = req.IsSecret
	v.Description = req.Description
	if err := s.envRepo.UpdateVariable(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *EnvService) DeleteVariable(id string) error {
	return s.envRepo.DeleteVariable(id)
}

func (s *EnvService) BatchUpsertVariables(envID string, req BatchEnvVariablesReq) error {
	vars := make([]model.EnvVariable, len(req.Variables))
	for i, v := range req.Variables {
		vars[i] = model.EnvVariable{
			VarKey:      v.VarKey,
			VarValue:    v.VarValue,
			IsSecret:    v.IsSecret,
			Description: v.Description,
		}
	}
	return s.envRepo.BatchUpsertVariables(envID, vars)
}

// Resource Quotas

func (s *EnvService) GetQuota(envID string) (*model.EnvResourceQuota, error) {
	return s.envRepo.GetQuota(envID)
}

func (s *EnvService) UpsertQuota(envID string, req EnvResourceQuotaReq) (*model.EnvResourceQuota, error) {
	q := &model.EnvResourceQuota{
		EnvironmentID: envID,
		CPURequest:    req.CPURequest,
		CPULimit:      req.CPULimit,
		MemoryRequest: req.MemoryRequest,
		MemoryLimit:   req.MemoryLimit,
		PodLimit:      req.PodLimit,
		StorageLimit:  req.StorageLimit,
	}
	if err := s.envRepo.UpsertQuota(q); err != nil {
		return nil, err
	}
	return q, nil
}
