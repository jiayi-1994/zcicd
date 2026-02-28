package repository

import (
	"github.com/zcicd/zcicd-server/internal/artifact/model"
	"gorm.io/gorm"
)

type RegistryRepository struct {
	db *gorm.DB
}

func NewRegistryRepository(db *gorm.DB) *RegistryRepository {
	return &RegistryRepository{db: db}
}

func (r *RegistryRepository) Create(reg *model.ImageRegistry) error {
	return r.db.Create(reg).Error
}

func (r *RegistryRepository) Get(id string) (*model.ImageRegistry, error) {
	var reg model.ImageRegistry
	err := r.db.Where("id = ?", id).First(&reg).Error
	return &reg, err
}

func (r *RegistryRepository) Update(reg *model.ImageRegistry) error {
	return r.db.Save(reg).Error
}

func (r *RegistryRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.ImageRegistry{}).Error
}

func (r *RegistryRepository) List() ([]model.ImageRegistry, error) {
	var list []model.ImageRegistry
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}
