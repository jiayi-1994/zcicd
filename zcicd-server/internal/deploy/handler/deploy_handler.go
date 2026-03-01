package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/deploy/service"
	"github.com/zcicd/zcicd-server/pkg/response"
	"gorm.io/gorm"
)

type DeployHandler struct {
	svc *service.DeployService
}

func NewDeployHandler(svc *service.DeployService) *DeployHandler {
	return &DeployHandler{svc: svc}
}

func (h *DeployHandler) CreateConfig(c *gin.Context) {
	var req service.CreateDeployConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	projectID := c.Param("project_id")
	config, err := h.svc.CreateConfig(c.Request.Context(), projectID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, config)
}

func (h *DeployHandler) GetConfig(c *gin.Context) {
	config, err := h.svc.GetConfig(c.Param("id"))
	if err != nil {
		response.NotFound(c, "部署配置不存在")
		return
	}
	response.OK(c, config)
}

func (h *DeployHandler) UpdateConfig(c *gin.Context) {
	var req service.UpdateDeployConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	config, err := h.svc.UpdateConfig(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.OK(c, config)
}

func (h *DeployHandler) DeleteConfig(c *gin.Context) {
	if err := h.svc.DeleteConfig(c.Request.Context(), c.Param("id")); err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.OK(c, nil)
}

func (h *DeployHandler) ListConfigs(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		response.BadRequest(c, "project_id is required")
		return
	}
	page, pageSize := parsePagination(c)
	list, total, err := h.svc.ListConfigs(projectID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *DeployHandler) ListConfigsByEnv(c *gin.Context) {
	projectID := c.Query("project_id")
	envID := c.Query("environment_id")
	if projectID == "" || envID == "" {
		response.BadRequest(c, "project_id and environment_id are required")
		return
	}
	list, err := h.svc.ListConfigsByEnv(projectID, envID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *DeployHandler) TriggerSync(c *gin.Context) {
	var req service.TriggerSyncReq
	if err := c.ShouldBindJSON(&req); err != nil {
		req = service.TriggerSyncReq{}
	}
	userID := c.GetString("user_id")
	history, err := h.svc.TriggerSync(c.Request.Context(), c.Param("id"), userID, req)
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.Created(c, history)
}

func (h *DeployHandler) Rollback(c *gin.Context) {
	var req service.RollbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetString("user_id")
	history, err := h.svc.Rollback(c.Request.Context(), c.Param("id"), userID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, history)
}

func (h *DeployHandler) GetStatus(c *gin.Context) {
	status, err := h.svc.GetStatus(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.OK(c, status)
}

func (h *DeployHandler) GetResources(c *gin.Context) {
	tree, err := h.svc.GetResources(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.OK(c, tree)
}

func (h *DeployHandler) GetHistory(c *gin.Context) {
	history, err := h.svc.GetHistory(c.Param("history_id"))
	if err != nil {
		response.NotFound(c, "部署历史不存在")
		return
	}
	response.OK(c, history)
}

func (h *DeployHandler) ListHistories(c *gin.Context) {
	configID := c.Param("id")
	page, pageSize := parsePagination(c)
	list, total, err := h.svc.ListHistories(configID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *DeployHandler) GetRolloutStatus(c *gin.Context) {
	status, err := h.svc.GetRolloutStatus(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.OK(c, status)
}

func (h *DeployHandler) PromoteRollout(c *gin.Context) {
	if err := h.svc.PromoteRollout(c.Request.Context(), c.Param("id")); err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.OK(c, nil)
}

func (h *DeployHandler) AbortRollout(c *gin.Context) {
	if err := h.svc.AbortRollout(c.Request.Context(), c.Param("id")); err != nil {
		h.handleNotFoundOrInternal(c, err, "部署配置不存在")
		return
	}
	response.OK(c, nil)
}

func parsePagination(c *gin.Context) (int, int) {
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}
	return page, pageSize
}

func (h *DeployHandler) handleNotFoundOrInternal(c *gin.Context, err error, fallbackNotFound string) {
	if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "record not found") {
		response.NotFound(c, fallbackNotFound)
		return
	}
	lowerErr := strings.ToLower(err.Error())
	if strings.Contains(lowerErr, "foreign key") || strings.Contains(lowerErr, "violates foreign key constraint") {
		response.NotFound(c, "关联资源不存在")
		return
	}
	response.InternalError(c, err.Error())
}
