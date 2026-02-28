package repository

import (
	"github.com/zcicd/zcicd-server/internal/deploy/model"
	"gorm.io/gorm"
)

type ApprovalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) *ApprovalRepository {
	return &ApprovalRepository{db: db}
}

func (r *ApprovalRepository) Create(record *model.ApprovalRecord) error {
	return r.db.Create(record).Error
}

func (r *ApprovalRepository) Get(id string) (*model.ApprovalRecord, error) {
	var record model.ApprovalRecord
	err := r.db.Preload("DeployHistory").Preload("DeployHistory.DeployConfig").
		Where("id = ?", id).First(&record).Error
	return &record, err
}

func (r *ApprovalRepository) GetByHistoryID(historyID string) (*model.ApprovalRecord, error) {
	var record model.ApprovalRecord
	err := r.db.Where("deploy_history_id = ?", historyID).First(&record).Error
	return &record, err
}

func (r *ApprovalRepository) Update(record *model.ApprovalRecord) error {
	return r.db.Save(record).Error
}

func (r *ApprovalRepository) ListPending(approverID string, page, pageSize int) ([]model.ApprovalRecord, int64, error) {
	var records []model.ApprovalRecord
	var total int64
	query := r.db.Where("status = 'pending'")
	query.Model(&model.ApprovalRecord{}).Count(&total)
	err := query.Preload("DeployHistory").Preload("DeployHistory.DeployConfig").
		Offset((page-1)*pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&records).Error
	return records, total, err
}
