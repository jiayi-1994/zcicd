package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zcicd/zcicd-server/internal/deploy/engine"
	"github.com/zcicd/zcicd-server/internal/deploy/model"
	"github.com/zcicd/zcicd-server/internal/deploy/repository"
	"github.com/zcicd/zcicd-server/pkg/mq"
	"gorm.io/datatypes"
)

type DeployService struct {
	deployRepo   *repository.DeployRepository
	approvalRepo *repository.ApprovalRepository
	appManager   *engine.AppManager
	syncCtrl     *engine.SyncController
	rolloutCtrl  *engine.RolloutController
	gitopsWriter *engine.GitOpsWriter
	mqClient     *mq.Client
	argoNS       string
}

func NewDeployService(
	deployRepo *repository.DeployRepository,
	approvalRepo *repository.ApprovalRepository,
	appManager *engine.AppManager,
	syncCtrl *engine.SyncController,
	rolloutCtrl *engine.RolloutController,
	gitopsWriter *engine.GitOpsWriter,
	mqClient *mq.Client,
	argoNS string,
) *DeployService {
	return &DeployService{
		deployRepo:   deployRepo,
		approvalRepo: approvalRepo,
		appManager:   appManager,
		syncCtrl:     syncCtrl,
		rolloutCtrl:  rolloutCtrl,
		gitopsWriter: gitopsWriter,
		mqClient:     mqClient,
		argoNS:       argoNS,
	}
}

// CreateConfig creates a deploy config and corresponding Argo CD Application.
func (s *DeployService) CreateConfig(ctx context.Context, projectID string, req CreateDeployConfigReq) (*model.DeployConfig, error) {
	valuesJSON, _ := json.Marshal(req.ValuesOverride)
	argoAppName := fmt.Sprintf("zcicd-%s-%s", projectID[:8], req.Name)

	config := &model.DeployConfig{
		ProjectID:      projectID,
		ServiceID:      req.ServiceID,
		EnvironmentID:  req.EnvironmentID,
		Name:           req.Name,
		DeployType:     req.DeployType,
		RepoURL:        req.RepoURL,
		TargetRevision: req.TargetRevision,
		ChartPath:      req.ChartPath,
		ValuesOverride: datatypes.JSON(valuesJSON),
		SyncPolicy:     req.SyncPolicy,
		AutoSync:       req.AutoSync,
		SelfHeal:       req.SelfHeal,
		Prune:          req.Prune,
		ArgoAppName:    argoAppName,
		Namespace:      req.Namespace,
	}
	if config.TargetRevision == "" {
		config.TargetRevision = "main"
	}
	if config.SyncPolicy == "" {
		config.SyncPolicy = "manual"
	}

	if err := s.deployRepo.CreateConfig(config); err != nil {
		return nil, err
	}

	// Create Argo CD Application
	if s.appManager != nil {
		argoApp := s.buildArgoApp(config, req.ValuesOverride)
		if err := s.appManager.CreateApp(ctx, argoApp); err != nil {
			// Log but don't fail — Argo CD might not be available
			fmt.Printf("warning: failed to create argo app %s: %v\n", argoAppName, err)
		}
	}

	return config, nil
}

// GetConfig returns a deploy config by ID.
func (s *DeployService) GetConfig(id string) (*model.DeployConfig, error) {
	return s.deployRepo.GetConfig(id)
}

// UpdateConfig updates a deploy config and its Argo CD Application.
func (s *DeployService) UpdateConfig(ctx context.Context, id string, req UpdateDeployConfigReq) (*model.DeployConfig, error) {
	config, err := s.deployRepo.GetConfig(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		config.Name = req.Name
	}
	if req.DeployType != "" {
		config.DeployType = req.DeployType
	}
	if req.RepoURL != "" {
		config.RepoURL = req.RepoURL
	}
	if req.TargetRevision != "" {
		config.TargetRevision = req.TargetRevision
	}
	if req.ChartPath != "" {
		config.ChartPath = req.ChartPath
	}
	if req.ValuesOverride != nil {
		valuesJSON, _ := json.Marshal(req.ValuesOverride)
		config.ValuesOverride = datatypes.JSON(valuesJSON)
	}
	if req.SyncPolicy != "" {
		config.SyncPolicy = req.SyncPolicy
	}
	if req.AutoSync != nil {
		config.AutoSync = *req.AutoSync
	}
	if req.SelfHeal != nil {
		config.SelfHeal = *req.SelfHeal
	}
	if req.Prune != nil {
		config.Prune = *req.Prune
	}
	if req.Namespace != "" {
		config.Namespace = req.Namespace
	}

	if err := s.deployRepo.UpdateConfig(config); err != nil {
		return nil, err
	}

	// Update Argo CD Application
	if s.appManager != nil && config.ArgoAppName != "" {
		var values map[string]interface{}
		if req.ValuesOverride != nil {
			values = req.ValuesOverride
		}
		argoApp := s.buildArgoApp(config, values)
		if err := s.appManager.UpdateApp(ctx, argoApp); err != nil {
			fmt.Printf("warning: failed to update argo app %s: %v\n", config.ArgoAppName, err)
		}
	}

	return config, nil
}

// DeleteConfig deletes a deploy config and its Argo CD Application.
func (s *DeployService) DeleteConfig(ctx context.Context, id string) error {
	config, err := s.deployRepo.GetConfig(id)
	if err != nil {
		return err
	}

	if s.appManager != nil && config.ArgoAppName != "" {
		if err := s.appManager.DeleteApp(ctx, config.ArgoAppName); err != nil {
			fmt.Printf("warning: failed to delete argo app %s: %v\n", config.ArgoAppName, err)
		}
	}

	return s.deployRepo.DeleteConfig(id)
}

// ListConfigs returns deploy configs for a project.
func (s *DeployService) ListConfigs(projectID string, page, pageSize int) ([]model.DeployConfig, int64, error) {
	return s.deployRepo.ListConfigs(projectID, page, pageSize)
}

// ListConfigsByEnv returns deploy configs for a project+environment.
func (s *DeployService) ListConfigsByEnv(projectID, envID string) ([]model.DeployConfig, error) {
	return s.deployRepo.ListConfigsByEnv(projectID, envID)
}

// TriggerSync triggers a sync for a deploy config.
func (s *DeployService) TriggerSync(ctx context.Context, configID, userID string, req TriggerSyncReq) (*model.DeployHistory, error) {
	config, err := s.deployRepo.GetConfig(configID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	history := &model.DeployHistory{
		DeployConfigID: configID,
		Revision:       req.Revision,
		Status:         "syncing",
		SyncStatus:     "OutOfSync",
		TriggeredBy:    userID,
		StartedAt:      &now,
	}
	if err := s.deployRepo.CreateHistory(history); err != nil {
		return nil, err
	}

	// Write values to GitOps repo if override provided
	if s.gitopsWriter != nil && req.Revision != "" {
		var values map[string]interface{}
		if len(config.ValuesOverride) > 0 {
			json.Unmarshal(config.ValuesOverride, &values)
		}
		if len(values) > 0 {
			commitSHA, gitErr := s.gitopsWriter.UpdateValues(ctx, config.RepoURL, config.TargetRevision, config.ChartPath+"/values.yaml", values)
			if gitErr != nil {
				fmt.Printf("warning: gitops write failed: %v\n", gitErr)
			} else {
				history.GitopsCommit = commitSHA
				s.publishEvent(mq.SubjectGitOpsUpdate, config.ProjectID, userID, history)
			}
		}
	}

	// Trigger Argo CD sync
	if s.syncCtrl != nil && config.ArgoAppName != "" {
		result, err := s.syncCtrl.TriggerSync(ctx, config.ArgoAppName, req.Revision)
		if err != nil {
			history.Status = "failed"
			history.ErrorMessage = err.Error()
			finished := time.Now()
			history.FinishedAt = &finished
			history.Duration = int(finished.Sub(now).Seconds())
			s.deployRepo.UpdateHistory(history)
			s.publishEvent(mq.SubjectDeployFailed, config.ProjectID, userID, history)
			return history, err
		}
		history.SyncStatus = result.Status
		history.HealthStatus = result.Health
		history.Revision = result.Revision
	}

	// Update status based on sync result
	if history.SyncStatus == "Synced" {
		history.Status = "succeeded"
		finished := time.Now()
		history.FinishedAt = &finished
		history.Duration = int(finished.Sub(now).Seconds())
		s.publishEvent(mq.SubjectDeploySucceeded, config.ProjectID, userID, history)
	}
	s.deployRepo.UpdateHistory(history)

	return history, nil
}

// Rollback rolls back to a previous deployment.
func (s *DeployService) Rollback(ctx context.Context, configID, userID string, req RollbackReq) (*model.DeployHistory, error) {
	prevHistory, err := s.deployRepo.GetHistory(req.HistoryID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	history := &model.DeployHistory{
		DeployConfigID: configID,
		Revision:       prevHistory.Revision,
		Status:         "syncing",
		TriggeredBy:    userID,
		RollbackFrom:   &req.HistoryID,
		StartedAt:      &now,
	}
	if err := s.deployRepo.CreateHistory(history); err != nil {
		return nil, err
	}

	config, err := s.deployRepo.GetConfig(configID)
	if err != nil {
		return nil, err
	}

	if s.syncCtrl != nil && config.ArgoAppName != "" {
		result, syncErr := s.syncCtrl.TriggerSync(ctx, config.ArgoAppName, prevHistory.Revision)
		if syncErr != nil {
			history.Status = "failed"
			history.ErrorMessage = syncErr.Error()
		} else {
			history.SyncStatus = result.Status
			history.HealthStatus = result.Health
			history.Status = "succeeded"
		}
	} else {
		history.Status = "succeeded"
	}

	finished := time.Now()
	history.FinishedAt = &finished
	history.Duration = int(finished.Sub(now).Seconds())
	s.deployRepo.UpdateHistory(history)
	s.publishEvent(mq.SubjectDeployRollback, config.ProjectID, userID, history)

	return history, nil
}

// GetStatus returns the current deploy status from Argo CD.
func (s *DeployService) GetStatus(ctx context.Context, configID string) (*engine.AppStatus, error) {
	config, err := s.deployRepo.GetConfig(configID)
	if err != nil {
		return nil, err
	}
	if s.appManager == nil || config.ArgoAppName == "" {
		return &engine.AppStatus{SyncStatus: "Unknown", HealthStatus: "Unknown"}, nil
	}
	return s.appManager.GetApp(ctx, config.ArgoAppName)
}

// GetResources returns the resource tree from Argo CD.
func (s *DeployService) GetResources(ctx context.Context, configID string) (*engine.ResourceTree, error) {
	config, err := s.deployRepo.GetConfig(configID)
	if err != nil {
		return nil, err
	}
	if s.appManager == nil || config.ArgoAppName == "" {
		return &engine.ResourceTree{}, nil
	}
	return s.appManager.GetResourceTree(ctx, config.ArgoAppName)
}

// GetHistory returns a deploy history by ID.
func (s *DeployService) GetHistory(id string) (*model.DeployHistory, error) {
	return s.deployRepo.GetHistory(id)
}

// ListHistories returns deploy histories for a config.
func (s *DeployService) ListHistories(configID string, page, pageSize int) ([]model.DeployHistory, int64, error) {
	return s.deployRepo.ListHistories(configID, page, pageSize)
}

// buildArgoApp converts a DeployConfig to an engine.ArgoApp.
func (s *DeployService) buildArgoApp(config *model.DeployConfig, values map[string]interface{}) engine.ArgoApp {
	return engine.ArgoApp{
		Name:           config.ArgoAppName,
		Namespace:      s.argoNS,
		RepoURL:        config.RepoURL,
		TargetRevision: config.TargetRevision,
		Path:           config.ChartPath,
		DestNamespace:  config.Namespace,
		ValuesOverride: values,
		SyncPolicy:     config.SyncPolicy,
		AutoSync:       config.AutoSync,
		SelfHeal:       config.SelfHeal,
		Prune:          config.Prune,
	}
}

// GetRolloutStatus returns the Argo Rollout status for a deploy config.
func (s *DeployService) GetRolloutStatus(ctx context.Context, configID string) (*engine.RolloutStatus, error) {
	config, err := s.deployRepo.GetConfig(configID)
	if err != nil {
		return nil, err
	}
	if s.rolloutCtrl == nil || config.ArgoAppName == "" {
		return &engine.RolloutStatus{Phase: "Unknown"}, nil
	}
	return s.rolloutCtrl.GetStatus(ctx, config.ArgoAppName)
}

// PromoteRollout promotes a canary/bluegreen rollout.
func (s *DeployService) PromoteRollout(ctx context.Context, configID string) error {
	config, err := s.deployRepo.GetConfig(configID)
	if err != nil {
		return err
	}
	if s.rolloutCtrl == nil || config.ArgoAppName == "" {
		return fmt.Errorf("rollout controller not available")
	}
	return s.rolloutCtrl.Promote(ctx, config.ArgoAppName)
}

// AbortRollout aborts a canary/bluegreen rollout.
func (s *DeployService) AbortRollout(ctx context.Context, configID string) error {
	config, err := s.deployRepo.GetConfig(configID)
	if err != nil {
		return err
	}
	if s.rolloutCtrl == nil || config.ArgoAppName == "" {
		return fmt.Errorf("rollout controller not available")
	}
	return s.rolloutCtrl.Abort(ctx, config.ArgoAppName)
}

// publishEvent publishes a NATS event.
func (s *DeployService) publishEvent(subject, projectID, userID string, history *model.DeployHistory) {
	if s.mqClient == nil {
		return
	}
	event := mq.Event{
		EventType:   subject,
		Timestamp:   time.Now().Format(time.RFC3339),
		ProjectID:   projectID,
		TriggeredBy: userID,
		Payload:     history,
	}
	data, _ := json.Marshal(event)
	s.mqClient.Publish(subject, data)
}
