package repository

import (
	"github.com/zcicd/zcicd-server/internal/artifact/model"
	"gorm.io/gorm"
)

type ChartRepository struct {
	db *gorm.DB
}

func NewChartRepository(db *gorm.DB) *ChartRepository {
	return &ChartRepository{db: db}
}

func (r *ChartRepository) Create(c *model.HelmChart) error {
	return r.db.Create(c).Error
}

func (r *ChartRepository) Get(id string) (*model.HelmChart, error) {
	var c model.HelmChart
	err := r.db.Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *ChartRepository) List() ([]model.HelmChart, error) {
	var list []model.HelmChart
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *ChartRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.HelmChart{}).Error
}
