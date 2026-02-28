package service

import (
	"context"
	"errors"

	"github.com/zcicd/zcicd-server/internal/project/model"
	"github.com/zcicd/zcicd-server/internal/project/repository"
	appErrors "github.com/zcicd/zcicd-server/pkg/errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	serviceRepo *repository.ServiceRepository
	envRepo     *repository.EnvironmentRepository
}

func NewProjectService(
	projectRepo *repository.ProjectRepository,
	serviceRepo *repository.ServiceRepository,
	envRepo *repository.EnvironmentRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		serviceRepo: serviceRepo,
		envRepo:     envRepo,
	}
}

// ==================== Project CRUD ====================

func (s *ProjectService) CreateProject(ctx context.Context, userID string, req *CreateProjectRequest) (*model.Project, error) {
	existing, _ := s.projectRepo.FindByIdentifier(ctx, req.Identifier)
	if existing != nil {
		return nil, appErrors.ErrProjectExists
	}

	project := &model.Project{
		Name:        req.Name,
		Identifier:  req.Identifier,
		Description: req.Description,
		OwnerID:     userID,
		RepoURL:     req.RepoURL,
	}
	if req.DefaultBranch != "" {
		project.DefaultBranch = req.DefaultBranch
	}
	if req.Visibility != "" {
		project.Visibility = req.Visibility
	}
	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建项目失败", err)
	}
	return project, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrProjectNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询项目失败", err)
	}
	return project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, req *UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrProjectNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询项目失败", err)
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.RepoURL != "" {
		project.RepoURL = req.RepoURL
	}
	if req.DefaultBranch != "" {
		project.DefaultBranch = req.DefaultBranch
	}
	if req.Visibility != "" {
		project.Visibility = req.Visibility
	}
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "更新项目失败", err)
	}
	return project, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	_, err := s.projectRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrProjectNotFound
		}
		return appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询项目失败", err)
	}
	return s.projectRepo.Delete(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context, page, pageSize int, keyword string) ([]model.Project, int64, error) {
	return s.projectRepo.List(ctx, page, pageSize, keyword)
}

// ==================== Service CRUD ====================

func (s *ProjectService) CreateService(ctx context.Context, projectID string, req *CreateServiceRequest) (*model.Service, error) {
	// Verify project exists
	if _, err := s.GetProject(ctx, projectID); err != nil {
		return nil, err
	}

	existing, _ := s.serviceRepo.FindByProjectAndName(ctx, projectID, req.Name)
	if existing != nil {
		return nil, appErrors.ErrServiceExists
	}

	svc := &model.Service{
		ProjectID:       projectID,
		Name:            req.Name,
		ServiceType:     req.ServiceType,
		Language:        req.Language,
		RepoURL:         req.RepoURL,
		HealthCheckPath: req.HealthCheckPath,
		HelmChartPath:   req.HelmChartPath,
	}
	if req.Branch != "" {
		svc.Branch = req.Branch
	}
	if req.DockerfilePath != "" {
		svc.DockerfilePath = req.DockerfilePath
	}
	if req.BuildContext != "" {
		svc.BuildContext = req.BuildContext
	}
	if req.DeployType != "" {
		svc.DeployType = req.DeployType
	}
	if req.HelmValues != nil {
		svc.HelmValues = datatypes.JSON(req.HelmValues)
	}
	if req.Ports != nil {
		svc.Ports = datatypes.JSON(req.Ports)
	}
	if req.EnvVars != nil {
		svc.EnvVars = datatypes.JSON(req.EnvVars)
	}
	if req.Resources != nil {
		svc.Resources = datatypes.JSON(req.Resources)
	}

	if err := s.serviceRepo.Create(ctx, svc); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建服务失败", err)
	}
	return svc, nil
}

func (s *ProjectService) GetService(ctx context.Context, id string) (*model.Service, error) {
	svc, err := s.serviceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrServiceNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询服务失败", err)
	}
	return svc, nil
}

func (s *ProjectService) UpdateService(ctx context.Context, id string, req *UpdateServiceRequest) (*model.Service, error) {
	svc, err := s.serviceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrServiceNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询服务失败", err)
	}

	if req.Name != "" {
		svc.Name = req.Name
	}
	if req.ServiceType != "" {
		svc.ServiceType = req.ServiceType
	}
	if req.Language != "" {
		svc.Language = req.Language
	}
	if req.RepoURL != "" {
		svc.RepoURL = req.RepoURL
	}
	if req.Branch != "" {
		svc.Branch = req.Branch
	}
	if req.DockerfilePath != "" {
		svc.DockerfilePath = req.DockerfilePath
	}
	if req.BuildContext != "" {
		svc.BuildContext = req.BuildContext
	}
	if req.DeployType != "" {
		svc.DeployType = req.DeployType
	}
	if req.HelmChartPath != "" {
		svc.HelmChartPath = req.HelmChartPath
	}
	if req.HealthCheckPath != "" {
		svc.HealthCheckPath = req.HealthCheckPath
	}
	if req.HelmValues != nil {
		svc.HelmValues = datatypes.JSON(req.HelmValues)
	}
	if req.Ports != nil {
		svc.Ports = datatypes.JSON(req.Ports)
	}
	if req.EnvVars != nil {
		svc.EnvVars = datatypes.JSON(req.EnvVars)
	}
	if req.Resources != nil {
		svc.Resources = datatypes.JSON(req.Resources)
	}
	if req.Status != "" {
		svc.Status = req.Status
	}

	if err := s.serviceRepo.Update(ctx, svc); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "更新服务失败", err)
	}
	return svc, nil
}

func (s *ProjectService) DeleteService(ctx context.Context, id string) error {
	_, err := s.serviceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrServiceNotFound
		}
		return appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询服务失败", err)
	}
	return s.serviceRepo.Delete(ctx, id)
}

func (s *ProjectService) ListServices(ctx context.Context, projectID string, page, pageSize int) ([]model.Service, int64, error) {
	return s.serviceRepo.ListByProject(ctx, projectID, page, pageSize)
}

// ==================== Environment CRUD ====================

func (s *ProjectService) CreateEnvironment(ctx context.Context, projectID string, req *CreateEnvRequest) (*model.Environment, error) {
	if _, err := s.GetProject(ctx, projectID); err != nil {
		return nil, err
	}

	existing, _ := s.envRepo.FindByProjectAndName(ctx, projectID, req.Name)
	if existing != nil {
		return nil, appErrors.ErrEnvExists
	}

	env := &model.Environment{
		ProjectID:    projectID,
		Name:         req.Name,
		EnvType:      req.EnvType,
		Namespace:    req.Namespace,
		ClusterID:    req.ClusterID,
		IsProduction: req.IsProduction,
		AutoDeploy:   req.AutoDeploy,
	}
	if req.DeployStrategy != nil {
		env.DeployStrategy = datatypes.JSON(req.DeployStrategy)
	}
	if req.GlobalEnvVars != nil {
		env.GlobalEnvVars = datatypes.JSON(req.GlobalEnvVars)
	}

	if err := s.envRepo.Create(ctx, env); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建环境失败", err)
	}
	return env, nil
}

func (s *ProjectService) GetEnvironment(ctx context.Context, id string) (*model.Environment, error) {
	env, err := s.envRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrEnvNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询环境失败", err)
	}
	return env, nil
}

func (s *ProjectService) UpdateEnvironment(ctx context.Context, id string, req *UpdateEnvRequest) (*model.Environment, error) {
	env, err := s.envRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrEnvNotFound
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询环境失败", err)
	}

	if req.Name != "" {
		env.Name = req.Name
	}
	if req.EnvType != "" {
		env.EnvType = req.EnvType
	}
	if req.Namespace != "" {
		env.Namespace = req.Namespace
	}
	if req.ClusterID != "" {
		env.ClusterID = req.ClusterID
	}
	if req.IsProduction != nil {
		env.IsProduction = *req.IsProduction
	}
	if req.AutoDeploy != nil {
		env.AutoDeploy = *req.AutoDeploy
	}
	if req.DeployStrategy != nil {
		env.DeployStrategy = datatypes.JSON(req.DeployStrategy)
	}
	if req.GlobalEnvVars != nil {
		env.GlobalEnvVars = datatypes.JSON(req.GlobalEnvVars)
	}
	if req.Status != "" {
		env.Status = req.Status
	}

	if err := s.envRepo.Update(ctx, env); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "更新环境失败", err)
	}
	return env, nil
}

func (s *ProjectService) DeleteEnvironment(ctx context.Context, id string) error {
	_, err := s.envRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.ErrEnvNotFound
		}
		return appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询环境失败", err)
	}
	return s.envRepo.Delete(ctx, id)
}

func (s *ProjectService) ListEnvironments(ctx context.Context, projectID string) ([]model.Environment, error) {
	return s.envRepo.ListByProject(ctx, projectID)
}
