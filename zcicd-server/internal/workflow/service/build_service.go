package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	appErrors "github.com/zcicd/zcicd-server/pkg/errors"
	"github.com/zcicd/zcicd-server/pkg/mq"

	"github.com/zcicd/zcicd-server/internal/workflow/engine"
	"github.com/zcicd/zcicd-server/internal/workflow/model"
	"github.com/zcicd/zcicd-server/internal/workflow/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BuildService struct {
	repo         *repository.BuildRepository
	templateRepo *repository.TemplateRepository
	crdManager   *engine.CRDManager
	mqClient     *mq.Client
	namespace    string
}

func NewBuildService(
	repo *repository.BuildRepository,
	templateRepo *repository.TemplateRepository,
	crdManager *engine.CRDManager,
	mqClient *mq.Client,
	namespace string,
) *BuildService {
	return &BuildService{
		repo:         repo,
		templateRepo: templateRepo,
		crdManager:   crdManager,
		mqClient:     mqClient,
		namespace:    namespace,
	}
}

func (s *BuildService) CreateConfig(ctx context.Context, req *CreateBuildConfigRequest) (*model.BuildConfig, error) {
	cfg := &model.BuildConfig{
		ProjectID:      req.ProjectID,
		ServiceID:      req.ServiceID,
		Name:           req.Name,
		TemplateID:     req.TemplateID,
		RepoURL:        req.RepoURL,
		Branch:         req.Branch,
		DockerfilePath: req.DockerfilePath,
		DockerContext:  req.DockerContext,
		ImageRepo:      req.ImageRepo,
		TagStrategy:    req.TagStrategy,
		CacheEnabled:   true,
		BuildScript:    req.BuildScript,
	}
	if cfg.Branch == "" {
		cfg.Branch = "main"
	}
	if cfg.TagStrategy == "" {
		cfg.TagStrategy = "branch-commit"
	}
	if cfg.DockerfilePath == "" {
		cfg.DockerfilePath = "Dockerfile"
	}
	if cfg.DockerContext == "" {
		cfg.DockerContext = "."
	}
	if req.CacheEnabled != nil {
		cfg.CacheEnabled = *req.CacheEnabled
	}
	if req.BuildEnv != nil {
		data, _ := json.Marshal(req.BuildEnv)
		cfg.BuildEnv = datatypes.JSON(data)
	} else {
		cfg.BuildEnv = datatypes.JSON([]byte("{}"))
	}
	if req.Variables != nil {
		data, _ := json.Marshal(req.Variables)
		cfg.Variables = datatypes.JSON(data)
	}

	if err := s.repo.CreateConfig(ctx, cfg); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建构建配置失败", err)
	}
	return cfg, nil
}

func (s *BuildService) GetConfig(ctx context.Context, id string) (*model.BuildConfig, error) {
	cfg, err := s.repo.FindConfigByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrBuildConfigNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建配置失败", err)
	}
	return cfg, nil
}

func (s *BuildService) UpdateConfig(ctx context.Context, id string, req *UpdateBuildConfigRequest) (*model.BuildConfig, error) {
	cfg, err := s.repo.FindConfigByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrBuildConfigNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建配置失败", err)
	}

	if req.Name != "" {
		cfg.Name = req.Name
	}
	if req.TemplateID != nil {
		cfg.TemplateID = req.TemplateID
	}
	if req.RepoURL != "" {
		cfg.RepoURL = req.RepoURL
	}
	if req.Branch != "" {
		cfg.Branch = req.Branch
	}
	if req.DockerfilePath != "" {
		cfg.DockerfilePath = req.DockerfilePath
	}
	if req.DockerContext != "" {
		cfg.DockerContext = req.DockerContext
	}
	if req.ImageRepo != "" {
		cfg.ImageRepo = req.ImageRepo
	}
	if req.TagStrategy != "" {
		cfg.TagStrategy = req.TagStrategy
	}
	if req.CacheEnabled != nil {
		cfg.CacheEnabled = *req.CacheEnabled
	}
	if req.BuildScript != "" {
		cfg.BuildScript = req.BuildScript
	}
	if req.BuildEnv != nil {
		data, _ := json.Marshal(req.BuildEnv)
		cfg.BuildEnv = datatypes.JSON(data)
	}
	if req.Variables != nil {
		data, _ := json.Marshal(req.Variables)
		cfg.Variables = datatypes.JSON(data)
	}

	if err := s.repo.UpdateConfig(ctx, cfg); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "更新构建配置失败", err)
	}
	return cfg, nil
}

func (s *BuildService) DeleteConfig(ctx context.Context, id string) error {
	if _, err := s.repo.FindConfigByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrBuildConfigNotFound
		}
		return appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建配置失败", err)
	}
	return s.repo.DeleteConfig(ctx, id)
}

func (s *BuildService) ListConfigsByProject(ctx context.Context, projectID string, page, pageSize int) ([]model.BuildConfig, int64, error) {
	return s.repo.ListConfigsByProject(ctx, projectID, page, pageSize)
}

func (s *BuildService) ListConfigsByService(ctx context.Context, serviceID string) ([]model.BuildConfig, error) {
	return s.repo.ListConfigsByService(ctx, serviceID)
}

func (s *BuildService) TriggerBuild(ctx context.Context, configID string, userID string, req *TriggerBuildRequest) (*model.BuildRun, error) {
	cfg, err := s.repo.FindConfigByID(ctx, configID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrBuildConfigNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建配置失败", err)
	}

	runNumber, err := s.repo.GetNextRunNumber(ctx, configID)
	if err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "获取运行序号失败", err)
	}

	branch := cfg.Branch
	if req != nil && req.Branch != nil && *req.Branch != "" {
		branch = *req.Branch
	}

	var commitSHA string
	if req != nil && req.CommitSHA != nil {
		commitSHA = *req.CommitSHA
	}

	imageTag := s.generateImageTag(branch, commitSHA, runNumber, cfg.TagStrategy)

	run := &model.BuildRun{
		BuildConfigID: configID,
		RunNumber:     runNumber,
		Status:        "pending",
		Branch:        branch,
		CommitSHA:     commitSHA,
		ImageTag:      imageTag,
		TriggeredBy:   &userID,
	}

	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建构建运行失败", err)
	}

	// Generate and submit Tekton TaskRun
	if s.crdManager != nil {
		buildModel := &engine.BuildModel{
			BuildConfigID:  configID,
			RunID:          run.ID,
			RunNumber:      runNumber,
			ProjectID:      cfg.ProjectID,
			ServiceName:    cfg.Name,
			Namespace:      s.namespace,
			RepoURL:        cfg.RepoURL,
			Branch:         branch,
			CommitSHA:      commitSHA,
			BuildScript:    cfg.BuildScript,
			DockerfilePath: cfg.DockerfilePath,
			DockerContext:  cfg.DockerContext,
			ImageRepo:      cfg.ImageRepo,
			ImageTag:       imageTag,
			CacheEnabled:   cfg.CacheEnabled,
		}

		templateEngine := engine.NewTemplateEngine("")
		taskRunYAML, err := templateEngine.RenderTaskRun(buildModel)
		if err != nil {
			// Log error but don't fail the run creation
			fmt.Printf("Failed to render TaskRun YAML: %v\n", err)
		} else {
			_, err = s.crdManager.CreateTaskRun(ctx, s.namespace, taskRunYAML)
			if err != nil {
				fmt.Printf("Failed to create TaskRun: %v\n", err)
			} else {
				tektonRef := fmt.Sprintf("build-%s-run-%d", configID[:8], runNumber)
				run.TektonRef = tektonRef
				run.Status = "running"
				now := time.Now()
				run.StartedAt = &now
				s.repo.UpdateRun(ctx, run)
			}
		}
	}

	// Publish build.started event
	if s.mqClient != nil {
		eventData, _ := json.Marshal(map[string]interface{}{
			"build_run_id":    run.ID,
			"build_config_id": configID,
			"run_number":      runNumber,
			"project_id":      cfg.ProjectID,
			"service_id":      cfg.ServiceID,
			"image_tag":       imageTag,
			"triggered_by":    userID,
			"triggered_at":    time.Now().Format(time.RFC3339),
		})
		s.mqClient.Publish(mq.SubjectBuildStarted, eventData)
	}

	return run, nil
}

func (s *BuildService) generateImageTag(branch, commitSHA string, runNumber int, strategy string) string {
	switch strategy {
	case "timestamp":
		return fmt.Sprintf("%s-%d", branch, time.Now().Unix())
	case "latest":
		return "latest"
	case "branch-commit":
		fallthrough
	default:
		if commitSHA != "" && len(commitSHA) >= 8 {
			return fmt.Sprintf("%s-%s", branch, commitSHA[:8])
		}
		return fmt.Sprintf("%s-%d", branch, runNumber)
	}
}

func (s *BuildService) GetRun(ctx context.Context, runID string) (*model.BuildRun, error) {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrBuildRunNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建运行失败", err)
	}
	return run, nil
}

func (s *BuildService) UpdateRunStatus(ctx context.Context, runID string, status string, message string) error {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		return err
	}

	run.Status = status
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
		if message != "" {
			run.CommitMessage = message
		}
	}

	if err := s.repo.UpdateRun(ctx, run); err != nil {
		return err
	}

	// Publish completion event
	if s.mqClient != nil && (status == "succeeded" || status == "failed") {
		subject := mq.SubjectBuildCompleted
		if status == "failed" {
			subject = mq.SubjectBuildFailed
		}
		eventData, _ := json.Marshal(map[string]interface{}{
			"build_run_id":    run.ID,
			"build_config_id": run.BuildConfigID,
			"status":          status,
			"image_tag":       run.ImageTag,
			"image_digest":    run.ImageDigest,
			"duration_sec":    run.DurationSec,
			"finished_at":     now.Format(time.RFC3339),
		})
		s.mqClient.Publish(subject, eventData)
	}

	return nil
}

func (s *BuildService) ListRuns(ctx context.Context, buildConfigID string, page, pageSize int) ([]model.BuildRun, int64, error) {
	return s.repo.ListRuns(ctx, buildConfigID, page, pageSize)
}

func (s *BuildService) ListAllRuns(ctx context.Context, projectID string, page, pageSize int) ([]model.BuildRun, int64, error) {
	return s.repo.ListAllRuns(ctx, projectID, page, pageSize)
}

func (s *BuildService) CancelRun(ctx context.Context, runID string) error {
	run, err := s.repo.FindRunByID(ctx, runID)
	if err != nil {
		return err
	}

	if run.Status != "pending" && run.Status != "running" {
		return appErrors.NewAppError(40004, "只能取消待执行或运行中的构建")
	}

	// Cancel Tekton TaskRun if exists
	if s.crdManager != nil && run.TektonRef != "" {
		if err := s.crdManager.CancelTaskRun(ctx, s.namespace, run.TektonRef); err != nil {
			fmt.Printf("Failed to cancel TaskRun: %v\n", err)
		}
	}

	run.Status = "cancelled"
	now := time.Now()
	run.FinishedAt = &now
	if run.StartedAt != nil {
		duration := int(now.Sub(*run.StartedAt).Seconds())
		run.DurationSec = &duration
	}

	return s.repo.UpdateRun(ctx, run)
}
