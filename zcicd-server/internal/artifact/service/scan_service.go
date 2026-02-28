package service

import (
	"github.com/zcicd/zcicd-server/internal/artifact/model"
	"github.com/zcicd/zcicd-server/internal/artifact/repository"
)

type ScanService struct {
	repo *repository.ScanRepository
}

func NewScanService(repo *repository.ScanRepository) *ScanService {
	return &ScanService{repo: repo}
}

func (s *ScanService) TriggerScan(imageName string, req TriggerScanReq) (*model.ImageScan, error) {
	scan := &model.ImageScan{
		RegistryID: req.RegistryID,
		ImageName:  imageName,
		Tag:        req.Tag,
		Status:     "scanning",
	}
	return scan, s.repo.Create(scan)
}

func (s *ScanService) GetScan(id string) (*model.ImageScan, error) {
	return s.repo.Get(id)
}

func (s *ScanService) ListByImage(registryID, imageName string) ([]model.ImageScan, error) {
	return s.repo.ListByImage(registryID, imageName)
}
