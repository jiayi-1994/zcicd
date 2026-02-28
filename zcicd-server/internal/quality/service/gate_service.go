package service

import (
	"github.com/zcicd/zcicd-server/internal/quality/model"
	"github.com/zcicd/zcicd-server/internal/quality/repository"
)

type QualityGateService struct {
	repo *repository.QualityGateRepository
}

func NewQualityGateService(repo *repository.QualityGateRepository) *QualityGateService {
	return &QualityGateService{repo: repo}
}

func (s *QualityGateService) Get(projectID string) (*model.QualityGate, error) {
	return s.repo.Get(projectID)
}

func (s *QualityGateService) Upsert(projectID string, req QualityGateReq) (*model.QualityGate, error) {
	g := &model.QualityGate{ProjectID: projectID}
	if req.MinCoverage != nil {
		g.MinCoverage = req.MinCoverage
	}
	if req.MaxBugs != nil {
		g.MaxBugs = *req.MaxBugs
	}
	if req.MaxVulnerabilities != nil {
		g.MaxVulnerabilities = *req.MaxVulnerabilities
	}
	if req.MaxCodeSmells != nil {
		g.MaxCodeSmells = *req.MaxCodeSmells
	}
	if req.MaxDuplications != nil {
		g.MaxDuplications = req.MaxDuplications
	}
	if req.BlockDeploy != nil {
		g.BlockDeploy = *req.BlockDeploy
	}
	return g, s.repo.Upsert(g)
}
