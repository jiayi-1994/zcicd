package repository

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"gorm.io/gorm"
)

type ClusterRepository struct {
	db *gorm.DB
}

func NewClusterRepository(db *gorm.DB) *ClusterRepository {
	return &ClusterRepository{db: db}
}

func (r *ClusterRepository) Create(c *model.Cluster) error {
	return r.db.Create(c).Error
}

func (r *ClusterRepository) Get(id string) (*model.Cluster, error) {
	var c model.Cluster
	err := r.db.Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *ClusterRepository) Update(c *model.Cluster) error {
	return r.db.Save(c).Error
}

func (r *ClusterRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Cluster{}).Error
}

func (r *ClusterRepository) List() ([]model.Cluster, error) {
	var list []model.Cluster
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}
