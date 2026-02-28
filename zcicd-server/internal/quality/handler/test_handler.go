package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/quality/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type TestHandler struct {
	svc *service.TestService
}

func NewTestHandler(svc *service.TestService) *TestHandler {
	return &TestHandler{svc: svc}
}

func (h *TestHandler) CreateConfig(c *gin.Context) {
	var req service.CreateTestConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cfg, err := h.svc.CreateConfig(c.Param("project_id"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, cfg)
}

func (h *TestHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig(c.Param("id"))
	if err != nil {
		response.NotFound(c, "测试配置不存在")
		return
	}
	response.OK(c, cfg)
}

func (h *TestHandler) UpdateConfig(c *gin.Context) {
	var req service.UpdateTestConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cfg, err := h.svc.UpdateConfig(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, cfg)
}

func (h *TestHandler) DeleteConfig(c *gin.Context) {
	if err := h.svc.DeleteConfig(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *TestHandler) ListConfigs(c *gin.Context) {
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

func (h *TestHandler) TriggerRun(c *gin.Context) {
	run, err := h.svc.TriggerRun(c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, run)
}

func (h *TestHandler) GetRun(c *gin.Context) {
	run, err := h.svc.GetRun(c.Param("run_id"))
	if err != nil {
		response.NotFound(c, "测试运行不存在")
		return
	}
	response.OK(c, run)
}

func (h *TestHandler) ListRuns(c *gin.Context) {
	page, pageSize := parsePagination(c)
	list, total, err := h.svc.ListRuns(c.Param("id"), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
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
