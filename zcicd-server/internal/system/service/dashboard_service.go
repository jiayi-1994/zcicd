package service

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"github.com/zcicd/zcicd-server/internal/system/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

type OverviewStats struct {
	Projects     int64 `json:"projects"`
	Services     int64 `json:"services"`
	Clusters     int64 `json:"clusters"`
	Environments int64 `json:"environments"`
	Registries   int64 `json:"registries"`
	Integrations int64 `json:"integrations"`
}

func (s *DashboardService) GetOverview() (*OverviewStats, error) {
	stats := &OverviewStats{}
	var err error

	if stats.Projects, err = s.repo.CountTable("projects"); err != nil {
		return nil, err
	}
	if stats.Services, err = s.repo.CountTable("services"); err != nil {
		return nil, err
	}
	if stats.Clusters, err = s.repo.CountTable("clusters"); err != nil {
		return nil, err
	}
	if stats.Environments, err = s.repo.CountTable("environments"); err != nil {
		return nil, err
	}
	if stats.Registries, err = s.repo.CountTable("image_registries"); err != nil {
		return nil, err
	}
	if stats.Integrations, err = s.repo.CountTable("integrations"); err != nil {
		return nil, err
	}
	return stats, nil
}

func (s *DashboardService) GetTrends(days int) ([]model.DailyStat, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	return s.repo.GetTrends(days)
}
