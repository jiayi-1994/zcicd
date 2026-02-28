package repository

import (
	"time"

	"github.com/zcicd/zcicd-server/internal/system/model"
	"gorm.io/gorm"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetTrends(days int) ([]model.DailyStat, error) {
	var stats []model.DailyStat
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.Where("stat_date >= ?", since).
		Order("stat_date ASC").Find(&stats).Error
	return stats, err
}

func (r *DashboardRepository) CountTable(table string) (int64, error) {
	var count int64
	err := r.db.Table(table).Count(&count).Error
	return count, err
}

func (r *DashboardRepository) CountTableWhere(table, where string, args ...interface{}) (int64, error) {
	var count int64
	err := r.db.Table(table).Where(where, args...).Count(&count).Error
	return count, err
}
