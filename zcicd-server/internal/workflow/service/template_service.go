package service

import (
	"context"
	"encoding/json"
	"errors"

	appErrors "github.com/zcicd/zcicd-server/pkg/errors"

	"github.com/zcicd/zcicd-server/internal/workflow/model"
	"github.com/zcicd/zcicd-server/internal/workflow/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TemplateService struct {
	repo *repository.TemplateRepository
}

func NewTemplateService(repo *repository.TemplateRepository) *TemplateService {
	return &TemplateService{repo: repo}
}

func (s *TemplateService) Create(ctx context.Context, userID string, req *CreateTemplateRequest) (*model.BuildTemplate, error) {
	tpl := &model.BuildTemplate{
		Name:        req.Name,
		Language:    req.Language,
		Framework:   req.Framework,
		Description: req.Description,
		BuildScript: req.BuildScript,
		DockerfileTpl: req.DockerfileTpl,
		TektonTaskTpl: req.TektonTaskTpl,
		IsSystem:    false,
		CreatedBy:   &userID,
	}

	if req.BuildEnv != nil {
		data, _ := json.Marshal(req.BuildEnv)
		tpl.BuildEnv = datatypes.JSON(data)
	}

	if err := s.repo.Create(ctx, tpl); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "创建构建模板失败", err)
	}
	return tpl, nil
}

func (s *TemplateService) GetByID(ctx context.Context, id string) (*model.BuildTemplate, error) {
	tpl, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewAppError(40404, "构建模板不存在")
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建模板失败", err)
	}
	return tpl, nil
}

func (s *TemplateService) Update(ctx context.Context, id string, req *CreateTemplateRequest) (*model.BuildTemplate, error) {
	tpl, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewAppError(40404, "构建模板不存在")
		}
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建模板失败", err)
	}

	if tpl.IsSystem {
		return nil, appErrors.NewAppError(40301, "系统模板不可修改")
	}

	if req.Name != "" {
		tpl.Name = req.Name
	}
	if req.Language != "" {
		tpl.Language = req.Language
	}
	if req.Framework != "" {
		tpl.Framework = req.Framework
	}
	if req.Description != "" {
		tpl.Description = req.Description
	}
	if req.BuildScript != "" {
		tpl.BuildScript = req.BuildScript
	}
	if req.DockerfileTpl != "" {
		tpl.DockerfileTpl = req.DockerfileTpl
	}
	if req.TektonTaskTpl != "" {
		tpl.TektonTaskTpl = req.TektonTaskTpl
	}
	if req.BuildEnv != nil {
		data, _ := json.Marshal(req.BuildEnv)
		tpl.BuildEnv = datatypes.JSON(data)
	}

	if err := s.repo.Update(ctx, tpl); err != nil {
		return nil, appErrors.Wrap(appErrors.ErrDatabaseError.Code, "更新构建模板失败", err)
	}
	return tpl, nil
}

func (s *TemplateService) Delete(ctx context.Context, id string) error {
	tpl, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.NewAppError(40404, "构建模板不存在")
		}
		return appErrors.Wrap(appErrors.ErrDatabaseError.Code, "查询构建模板失败", err)
	}

	if tpl.IsSystem {
		return appErrors.NewAppError(40301, "系统模板不可删除")
	}

	return s.repo.Delete(ctx, id)
}

func (s *TemplateService) List(ctx context.Context, language string) ([]model.BuildTemplate, error) {
	return s.repo.List(ctx, language)
}

func (s *TemplateService) ListSystem(ctx context.Context) ([]model.BuildTemplate, error) {
	return s.repo.ListSystem(ctx)
}
