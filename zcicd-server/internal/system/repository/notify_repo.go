package repository

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"gorm.io/gorm"
)

type NotifyRepository struct {
	db *gorm.DB
}

func NewNotifyRepository(db *gorm.DB) *NotifyRepository {
	return &NotifyRepository{db: db}
}

// Channels

func (r *NotifyRepository) CreateChannel(c *model.NotifyChannel) error {
	return r.db.Create(c).Error
}

func (r *NotifyRepository) GetChannel(id string) (*model.NotifyChannel, error) {
	var c model.NotifyChannel
	err := r.db.Where("id = ?", id).First(&c).Error
	return &c, err
}

func (r *NotifyRepository) UpdateChannel(c *model.NotifyChannel) error {
	return r.db.Save(c).Error
}

func (r *NotifyRepository) DeleteChannel(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.NotifyChannel{}).Error
}

func (r *NotifyRepository) ListChannels() ([]model.NotifyChannel, error) {
	var list []model.NotifyChannel
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}
