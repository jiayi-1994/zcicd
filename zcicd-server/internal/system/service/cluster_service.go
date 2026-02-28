package service

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"github.com/zcicd/zcicd-server/internal/system/repository"
)

type ClusterService struct {
	repo *repository.ClusterRepository
}

func NewClusterService(repo *repository.ClusterRepository) *ClusterService {
	return &ClusterService{repo: repo}
}

func (s *ClusterService) Create(req CreateClusterReq) (*model.Cluster, error) {
	c := &model.Cluster{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Provider:      req.Provider,
		APIServerURL:  req.APIServerURL,
		KubeConfigRef: req.KubeConfigRef,
		Status:        "connected",
	}
	return c, s.repo.Create(c)
}

func (s *ClusterService) Get(id string) (*model.Cluster, error) {
	return s.repo.Get(id)
}

func (s *ClusterService) Update(id string, req UpdateClusterReq) (*model.Cluster, error) {
	c, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	if req.DisplayName != "" {
		c.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		c.Description = req.Description
	}
	if req.KubeConfigRef != "" {
		c.KubeConfigRef = req.KubeConfigRef
	}
	if req.Status != "" {
		c.Status = req.Status
	}
	return c, s.repo.Update(c)
}

func (s *ClusterService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ClusterService) List() ([]model.Cluster, error) {
	return s.repo.List()
}
