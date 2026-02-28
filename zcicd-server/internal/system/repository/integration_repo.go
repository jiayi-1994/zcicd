package repository

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"gorm.io/gorm"
)

type IntegrationRepository struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{db: db}
}

func (r *IntegrationRepository) Create(i *model.Integration) error {
	return r.db.Create(i).Error
}

func (r *IntegrationRepository) Get(id string) (*model.Integration, error) {
	var i model.Integration
	err := r.db.Where("id = ?", id).First(&i).Error
	return &i, err
}

func (r *IntegrationRepository) Update(i *model.Integration) error {
	return r.db.Save(i).Error
}

func (r *IntegrationRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Integration{}).Error
}

func (r *IntegrationRepository) List() ([]model.Integration, error) {
	var list []model.Integration
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}
