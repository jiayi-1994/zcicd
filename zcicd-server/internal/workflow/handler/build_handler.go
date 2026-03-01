package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/zcicd/zcicd-server/internal/workflow/service"
	appErrors "github.com/zcicd/zcicd-server/pkg/errors"
	"github.com/zcicd/zcicd-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BuildHandler struct {
	svc *service.BuildService
}

func NewBuildHandler(svc *service.BuildService) *BuildHandler {
	return &BuildHandler{svc: svc}
}

func (h *BuildHandler) CreateConfig(c *gin.Context) {
	var req service.CreateBuildConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cfg, err := h.svc.CreateConfig(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 500, 50001, err.Error())
		return
	}
	response.Created(c, cfg)
}

func (h *BuildHandler) GetConfig(c *gin.Context) {
	id := c.Param("id")
	cfg, err := h.svc.GetConfig(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "构建配置不存在")
		return
	}
	response.OK(c, cfg)
}

func (h *BuildHandler) UpdateConfig(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateBuildConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cfg, err := h.svc.UpdateConfig(c.Request.Context(), id, &req)
	if err != nil {
		handleNotFoundOrInternal(c, err, "构建配置不存在")
		return
	}
	response.OK(c, cfg)
}

func (h *BuildHandler) DeleteConfig(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteConfig(c.Request.Context(), id); err != nil {
		handleNotFoundOrInternal(c, err, "构建配置不存在")
		return
	}
	response.OK(c, nil)
}

func (h *BuildHandler) ListConfigs(c *gin.Context) {
	projectID := c.Query("project_id")
	serviceID := c.Query("service_id")

	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	if serviceID != "" {
		list, err := h.svc.ListConfigsByService(c.Request.Context(), serviceID)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		response.OK(c, list)
		return
	}

	if projectID == "" {
		response.BadRequest(c, "project_id or service_id is required")
		return
	}

	list, total, err := h.svc.ListConfigsByProject(c.Request.Context(), projectID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *BuildHandler) TriggerBuild(c *gin.Context) {
	configID := c.Param("id")
	var req service.TriggerBuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = service.TriggerBuildRequest{}
	}

	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	run, err := h.svc.TriggerBuild(c.Request.Context(), configID, userID, &req)
	if err != nil {
		handleNotFoundOrInternal(c, err, "构建配置不存在")
		return
	}
	response.Created(c, run)
}

func (h *BuildHandler) GetRun(c *gin.Context) {
	runID := c.Param("run_id")
	run, err := h.svc.GetRun(c.Request.Context(), runID)
	if err != nil {
		response.NotFound(c, "构建运行不存在")
		return
	}
	response.OK(c, run)
}

func (h *BuildHandler) ListRuns(c *gin.Context) {
	configID := c.Query("config_id")
	projectID := c.Query("project_id")

	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	var list interface{}
	var total int64
	var err error

	if configID != "" {
		list, total, err = h.svc.ListRuns(c.Request.Context(), configID, page, pageSize)
	} else if projectID != "" {
		list, total, err = h.svc.ListAllRuns(c.Request.Context(), projectID, page, pageSize)
	} else {
		response.BadRequest(c, "config_id or project_id is required")
		return
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *BuildHandler) CancelRun(c *gin.Context) {
	runID := c.Param("run_id")
	if err := h.svc.CancelRun(c.Request.Context(), runID); err != nil {
		handleNotFoundOrInternal(c, err, "构建运行不存在")
		return
	}
	response.OK(c, nil)
}

func (h *BuildHandler) GetRunLogs(c *gin.Context) {
	runID := c.Param("run_id")
	response.OK(c, gin.H{
		"run_id":  runID,
		"message": "实时日志请通过 WebSocket 连接获取: /api/v1/build-runs/" + runID + "/logs/ws",
		"logs":    []string{},
	})
}

type TemplateHandler struct {
	svc *service.TemplateService
}

func NewTemplateHandler(svc *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{svc: svc}
}

func (h *TemplateHandler) Create(c *gin.Context) {
	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	tpl, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, 500, 50001, err.Error())
		return
	}
	response.Created(c, tpl)
}

func (h *TemplateHandler) Get(c *gin.Context) {
	id := c.Param("id")
	tpl, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "构建模板不存在")
		return
	}
	response.OK(c, tpl)
}

func (h *TemplateHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req service.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tpl, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		handleNotFoundOrInternal(c, err, "构建模板不存在")
		return
	}
	response.OK(c, tpl)
}

func (h *TemplateHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handleNotFoundOrInternal(c, err, "构建模板不存在")
		return
	}
	response.OK(c, nil)
}

func (h *TemplateHandler) List(c *gin.Context) {
	language := c.Query("language")
	list, err := h.svc.List(c.Request.Context(), language)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func handleNotFoundOrInternal(c *gin.Context, err error, fallbackNotFound string) {
	if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "record not found") {
		response.NotFound(c, fallbackNotFound)
		return
	}
	lowerErr := strings.ToLower(err.Error())
	if strings.Contains(lowerErr, "foreign key") || strings.Contains(lowerErr, "violates foreign key constraint") {
		response.NotFound(c, "关联资源不存在")
		return
	}
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) && strings.Contains(appErr.Message, "不存在") {
		response.NotFound(c, appErr.Message)
		return
	}
	response.Error(c, 500, 50001, err.Error())
}
