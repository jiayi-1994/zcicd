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

type WorkflowHandler struct {
	svc *service.WorkflowService
}

func NewWorkflowHandler(svc *service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{svc: svc}
}

func (h *WorkflowHandler) Create(c *gin.Context) {
	var req service.CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	wf, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, wf)
}

func (h *WorkflowHandler) Get(c *gin.Context) {
	id := c.Param("id")
	wf, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "工作流不存在")
		return
	}
	response.OK(c, wf)
}

func (h *WorkflowHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	wf, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "工作流不存在")
		return
	}
	response.OK(c, wf)
}

func (h *WorkflowHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.handleNotFoundOrInternal(c, err, "工作流不存在")
		return
	}
	response.OK(c, nil)
}

func (h *WorkflowHandler) List(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		response.BadRequest(c, "project_id is required")
		return
	}

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

	list, total, err := h.svc.ListByProject(c.Request.Context(), projectID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *WorkflowHandler) Trigger(c *gin.Context) {
	id := c.Param("id")
	var req service.TriggerWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body
		req = service.TriggerWorkflowRequest{}
	}

	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "未授权")
		return
	}

	run, err := h.svc.Trigger(c.Request.Context(), id, userID, &req)
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "工作流不存在")
		return
	}
	response.Created(c, run)
}

func (h *WorkflowHandler) GetRun(c *gin.Context) {
	runID := c.Param("run_id")
	run, err := h.svc.GetRun(c.Request.Context(), runID)
	if err != nil {
		response.NotFound(c, "工作流运行不存在")
		return
	}
	response.OK(c, run)
}

func (h *WorkflowHandler) ListRuns(c *gin.Context) {
	workflowID := c.Param("id")

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

	list, total, err := h.svc.ListRuns(c.Request.Context(), workflowID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *WorkflowHandler) CancelRun(c *gin.Context) {
	runID := c.Param("run_id")
	if err := h.svc.CancelRun(c.Request.Context(), runID); err != nil {
		h.handleNotFoundOrInternal(c, err, "工作流运行不存在")
		return
	}
	response.OK(c, nil)
}

func (h *WorkflowHandler) RetryRun(c *gin.Context) {
	runID := c.Param("run_id")
	userID := c.GetString("user_id")
	run, err := h.svc.RetryRun(c.Request.Context(), runID, userID)
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "工作流运行不存在")
		return
	}
	response.Created(c, run)
}

func (h *WorkflowHandler) handleNotFoundOrInternal(c *gin.Context, err error, fallbackNotFound string) {
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
