package repository

import (
	"github.com/zcicd/zcicd-server/internal/artifact/model"
	"gorm.io/gorm"
)

type ScanRepository struct {
	db *gorm.DB
}

func NewScanRepository(db *gorm.DB) *ScanRepository {
	return &ScanRepository{db: db}
}

func (r *ScanRepository) Create(scan *model.ImageScan) error {
	return r.db.Create(scan).Error
}

func (r *ScanRepository) Get(id string) (*model.ImageScan, error) {
	var s model.ImageScan
	err := r.db.Where("id = ?", id).First(&s).Error
	return &s, err
}

func (r *ScanRepository) ListByImage(registryID, imageName string) ([]model.ImageScan, error) {
	var list []model.ImageScan
	err := r.db.Where("registry_id = ? AND image_name = ?", registryID, imageName).
		Order("created_at DESC").Find(&list).Error
	return list, err
}
