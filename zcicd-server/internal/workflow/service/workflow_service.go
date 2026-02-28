package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	appErrors "github.com/zcicd/zcicd-server/pkg/errors"
	"github.com/zcicd/zcicd-server/pkg/mq"

	"github.com/zcicd/zcicd-server/internal/workflow/model"
	"github.com/zcicd/zcicd-server/internal/workflow/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WorkflowService struct {
	repo     *repository.WorkflowRepository
	buildRepo *repository.BuildRepository
	mqClient *mq.Client
}

func NewWorkflowService(repo *repository.WorkflowRepository, buildRepo *repository.BuildRepository, mqClient *mq.Client) *WorkflowService {
	return &WorkflowService{
		repo:      repo,
		buildRepo: buildRepo,
		mqClient:  mqClient,
	}
}

func (s *WorkflowService) Create(ctx context.Context, req *CreateWorkflowRequest) (*model.Workflow, error) {
	wf := &model.Workflow{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		TriggerType: req.TriggerType,
		Enabled:     true,
	}
	if req.TriggerType == "" {
		wf.TriggerType = "manual"
	}
	if req.Enabled != nil {
		wf.Enabled = *req.Enabled
	}
	if req.TriggerConfig != nil {
		data, _ := json.Marshal(req.TriggerConfig)
		wf.TriggerConfig = datatypes.JSON(data)
	}

	for _, stageReq := range req.Stages {
		stage := model.WorkflowStage{
			Name:      stageReq.Name,
			StageType: stageReq.StageType,
			SortOrder: stageReq.SortOrder,
			Timeout:   stageReq.Timeout,
			Enabled:   true,
		}
		if stage.Timeout == 0 {
			stage.Timeout = 3600
		}
		if stageReq.Enabled != nil {
			stage.Enabled = *stageReq.Enabled
		}
		if stageReq.Config != nil {
			data, _ := json.Marshal(stageReq.Config)
			stage.Config = datatypes.JSON(data)
		}

		for _, jobReq := range stageReq.Jobs {
			job := model.StageJob{
				Name:       jobReq.Name,
				JobType:    jobReq.JobType,
				SortOrder:  jobReq.SortOrder,
				TimeoutSec: jobReq.TimeoutSec,
				Enabled:    true,
			}
			if job.TimeoutSec == 0 {
				job.TimeoutSec = 3600
			}
			if jobReq.Enabled != nil {
				job.Enabled = *jobReq.Enabled
			}
			if jobReq.Config != nil {
				data, _ := json.Marshal(jobReq.Config)
				job.Config = datatypes.JSON(data)
			}
			stage.Jobs = append(stage.Jobs, job)
		}
		wf.Stages = append(wf.Stages, stage)
	}

	if err := s.repo.Create(ctx, wf); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建工作流失败", err)
	}
	return wf, nil
}

func (s *WorkflowService) GetByID(ctx context.Context, id string) (*model.Workflow, error) {
	wf, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrWorkflowNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询工作流失败", err)
	}
	return wf, nil
}

func (s *WorkflowService) Update(ctx context.Context, id string, req *UpdateWorkflowRequest) (*model.Workflow, error) {
	wf, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrWorkflowNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询工作流失败", err)
	}

	if req.Name != "" {
		wf.Name = req.Name
	}
	if req.Description != "" {
		wf.Description = req.Description
	}
	if req.TriggerType != "" {
		wf.TriggerType = req.TriggerType
	}
	if req.Enabled != nil {
		wf.Enabled = *req.Enabled
	}
	if req.TriggerConfig != nil {
		data, _ := json.Marshal(req.TriggerConfig)
		wf.TriggerConfig = datatypes.JSON(data)
	}

	if err := s.repo.Update(ctx, wf); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "更新工作流失败", err)
	}
	return wf, nil
}

func (s *WorkflowService) Delete(ctx context.Context, id string) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrWorkflowNotFound
		}
		return appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询工作流失败", err)
	}
	return s.repo.Delete(ctx, id)
}

func (s *WorkflowService) ListByProject(ctx context.Context, projectID string, page, pageSize int) ([]model.Workflow, int64, error) {
	return s.repo.ListByProject(ctx, projectID, page, pageSize)
}

func (s *WorkflowService) Trigger(ctx context.Context, workflowID string, userID string, req *TriggerWorkflowRequest) (*model.WorkflowRun, error) {
	wf, err := s.repo.FindByID(ctx, workflowID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrWorkflowNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询工作流失败", err)
	}

	if !wf.Enabled {
		return nil, appErrors.NewAppError(40003, "工作流已禁用")
	}

	runNumber, err := s.repo.GetNextRunNumber(ctx, workflowID)
	if err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "获取运行序号失败", err)
	}

	run := &model.WorkflowRun{
		WorkflowID:  workflowID,
		RunNumber:   runNumber,
		Status:      "pending",
		TriggerType: "manual",
		TriggeredBy: &userID,
	}
	if req != nil && req.InputParams != nil {
		data, _ := json.Marshal(req.InputParams)
		run.InputParams = datatypes.JSON(data)
	}

	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建工作流运行失败", err)
	}

	// Publish workflow.started event
	if s.mqClient != nil {
		eventData, _ := json.Marshal(map[string]interface{}{
			"workflow_run_id": run.ID,
			"workflow_id":     workflowID,
			"run_number":      runNumber,
			"triggered_by":    userID,
			"triggered_at":    time.Now().Format(time.RFC3339),
		})
		s.mqClient.Publish(mq.SubjectWorkflowStarted, eventData)
	}

	return run, nil
}

func (s *WorkflowService) GetRun(ctx context.Context, runID string) (*model.WorkflowRun, error) {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrWorkflowRunNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询工作流运行失败", err)
	}
	return run, nil
}

func (s *WorkflowService) UpdateRunStatus(ctx context.Context, runID string, status string, message string) error {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		return err
	}

	run.Status = status
	if message != "" {
		run.ErrorMessage = message
	}

	now := time.Now()
	if status == "running" && run.StartedAt == nil {
		run.StartedAt = &now
	}
	if status == "succeeded" || status == "failed" || status == "cancelled" {
		run.FinishedAt = &now
		if run.StartedAt != nil {
			duration := int(now.Sub(*run.StartedAt).Seconds())
			run.DurationSec = &duration
		}
	}

	return s.repo.UpdateRun(ctx, run)
}

func (s *WorkflowService) ListRuns(ctx context.Context, workflowID string, page, pageSize int) ([]model.WorkflowRun, int64, error) {
	return s.repo.ListRuns(ctx, workflowID, page, pageSize)
}

// CancelRun cancels a pending or running workflow run.
func (s *WorkflowService) CancelRun(ctx context.Context, runID string) error {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		return err
	}
	if run.Status != "pending" && run.Status != "running" {
		return appErrors.NewAppError(40004, "只能取消待执行或运行中的工作流")
	}
	run.Status = "cancelled"
	now := time.Now()
	run.FinishedAt = &now
	if run.StartedAt != nil {
		d := int(now.Sub(*run.StartedAt).Seconds())
		run.DurationSec = &d
	}
	return s.repo.UpdateRun(ctx, run)
}

// RetryRun creates a new run based on a failed run's parameters.
func (s *WorkflowService) RetryRun(ctx context.Context, runID, userID string) (*model.WorkflowRun, error) {
	prev, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		return nil, err
	}
	if prev.Status != "failed" && prev.Status != "cancelled" {
		return nil, appErrors.NewAppError(40004, "只能重试失败或已取消的工作流")
	}
	runNumber, _ := s.repo.GetNextRunNumber(ctx, prev.WorkflowID)
	run := &model.WorkflowRun{
		WorkflowID:  prev.WorkflowID,
		RunNumber:   runNumber,
		Status:      "pending",
		TriggerType: "manual",
		TriggeredBy: &userID,
		InputParams: prev.InputParams,
	}
	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, err
	}
	return run, nil
}

// TriggerByWebhook finds webhook-triggered workflows and triggers matching ones.
func (s *WorkflowService) TriggerByWebhook(ctx context.Context, repoURL, branch, commitSHA string) (int, error) {
	workflows, err := s.repo.ListByTriggerType(ctx, "webhook")
	if err != nil {
		return 0, err
	}
	triggered := 0
	for _, wf := range workflows {
		var cfg map[string]string
		if len(wf.TriggerConfig) > 0 {
			json.Unmarshal(wf.TriggerConfig, &cfg)
		}
		// Match repo URL and optionally branch
		if cfg["repo_url"] != repoURL {
			continue
		}
		if cfgBranch, ok := cfg["branch"]; ok && cfgBranch != "" && cfgBranch != branch {
			continue
		}
		runNumber, _ := s.repo.GetNextRunNumber(ctx, wf.ID)
		run := &model.WorkflowRun{
			WorkflowID:  wf.ID,
			RunNumber:   runNumber,
			Status:      "pending",
			TriggerType: "webhook",
			InputParams: datatypes.JSON(mustMarshal(map[string]string{
				"branch":     branch,
				"commit_sha": commitSHA,
			})),
		}
		if err := s.repo.CreateRun(ctx, run); err != nil {
			continue
		}
		triggered++
	}
	return triggered, nil
}

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
