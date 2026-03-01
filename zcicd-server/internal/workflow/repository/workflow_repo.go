package repository

import (
	"context"

	"github.com/zcicd/zcicd-server/internal/workflow/model"

	"gorm.io/gorm"
)

type WorkflowRepository struct {
	db *gorm.DB
}

func NewWorkflowRepository(db *gorm.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) Create(ctx context.Context, wf *model.Workflow) error {
	return r.db.WithContext(ctx).Create(wf).Error
}

func (r *WorkflowRepository) FindByID(ctx context.Context, id string) (*model.Workflow, error) {
	var wf model.Workflow
	err := r.db.WithContext(ctx).
		Preload("Stages").
		Preload("Stages.Jobs").
		First(&wf, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wf, nil
}

func (r *WorkflowRepository) Update(ctx context.Context, wf *model.Workflow) error {
	return r.db.WithContext(ctx).Save(wf).Error
}

func (r *WorkflowRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("workflow_id = ?", id).Delete(&model.StageJob{}).Error; err != nil {
			return err
		}
		if err := tx.Where("workflow_id = ?", id).Delete(&model.WorkflowStage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Workflow{}, "id = ?", id).Error
	})
}

func (r *WorkflowRepository) ListByProject(ctx context.Context, projectID string, page, pageSize int) ([]model.Workflow, int64, error) {
	var list []model.Workflow
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Workflow{}).Where("project_id = ?", projectID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *WorkflowRepository) CreateRun(ctx context.Context, run *model.WorkflowRun) error {
	return r.db.WithContext(ctx).Create(run).Error
}

func (r *WorkflowRepository) FindRunByID(ctx context.Context, id string) (*model.WorkflowRun, error) {
	var run model.WorkflowRun
	err := r.db.WithContext(ctx).
		Preload("Workflow").
		Preload("Workflow.Stages").
		First(&run, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *WorkflowRepository) UpdateRun(ctx context.Context, run *model.WorkflowRun) error {
	return r.db.WithContext(ctx).Save(run).Error
}

func (r *WorkflowRepository) ListRuns(ctx context.Context, workflowID string, page, pageSize int) ([]model.WorkflowRun, int64, error) {
	var list []model.WorkflowRun
	var total int64

	query := r.db.WithContext(ctx).Model(&model.WorkflowRun{}).Where("workflow_id = ?", workflowID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("started_at DESC NULLS LAST, id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *WorkflowRepository) ListByTriggerType(ctx context.Context, triggerType string) ([]model.Workflow, error) {
	var list []model.Workflow
	err := r.db.WithContext(ctx).Where("trigger_type = ? AND enabled = true", triggerType).Find(&list).Error
	return list, err
}

func (r *WorkflowRepository) GetNextRunNumber(ctx context.Context, workflowID string) (int, error) {
	var maxNumber int
	r.db.WithContext(ctx).Model(&model.WorkflowRun{}).
		Where("workflow_id = ?", workflowID).
		Select("COALESCE(MAX(run_number), 0)").
		Scan(&maxNumber)
	return maxNumber + 1, nil
}
