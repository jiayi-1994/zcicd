package service

import (
	"github.com/zcicd/zcicd-server/internal/artifact/model"
	"github.com/zcicd/zcicd-server/internal/artifact/repository"
)

type ChartService struct {
	repo *repository.ChartRepository
}

func NewChartService(repo *repository.ChartRepository) *ChartService {
	return &ChartService{repo: repo}
}

func (s *ChartService) Create(req CreateChartReq) (*model.HelmChart, error) {
	c := &model.HelmChart{
		Name:    req.Name,
		RepoURL: req.RepoURL,
	}
	return c, s.repo.Create(c)
}

func (s *ChartService) Get(id string) (*model.HelmChart, error) {
	return s.repo.Get(id)
}

func (s *ChartService) List() ([]model.HelmChart, error) {
	return s.repo.List()
}

func (s *ChartService) Delete(id string) error {
	return s.repo.Delete(id)
}
