package repository

import (
	"github.com/zcicd/zcicd-server/internal/system/model"
	"gorm.io/gorm"
)

type RuleRepository struct {
	db *gorm.DB
}

func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) Create(rule *model.NotifyRule) error {
	return r.db.Create(rule).Error
}

func (r *RuleRepository) Get(id string) (*model.NotifyRule, error) {
	var rule model.NotifyRule
	err := r.db.Where("id = ?", id).First(&rule).Error
	return &rule, err
}

func (r *RuleRepository) Update(rule *model.NotifyRule) error {
	return r.db.Save(rule).Error
}

func (r *RuleRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.NotifyRule{}).Error
}

func (r *RuleRepository) List() ([]model.NotifyRule, error) {
	var list []model.NotifyRule
	err := r.db.Order("created_at DESC").Find(&list).Error
	return list, err
}
