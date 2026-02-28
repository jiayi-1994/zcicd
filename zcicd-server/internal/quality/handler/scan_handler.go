package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/quality/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type ScanHandler struct {
	svc *service.ScanService
}

func NewScanHandler(svc *service.ScanService) *ScanHandler {
	return &ScanHandler{svc: svc}
}

func (h *ScanHandler) CreateConfig(c *gin.Context) {
	var req service.CreateScanConfigReq
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

func (h *ScanHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig(c.Param("sid"))
	if err != nil {
		response.NotFound(c, "扫描配置不存在")
		return
	}
	response.OK(c, cfg)
}

func (h *ScanHandler) UpdateConfig(c *gin.Context) {
	var req service.UpdateScanConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cfg, err := h.svc.UpdateConfig(c.Param("sid"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, cfg)
}

func (h *ScanHandler) DeleteConfig(c *gin.Context) {
	if err := h.svc.DeleteConfig(c.Param("sid")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ScanHandler) ListConfigs(c *gin.Context) {
	page, pageSize := parsePagination(c)
	list, total, err := h.svc.ListConfigs(c.Param("project_id"), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *ScanHandler) TriggerRun(c *gin.Context) {
	run, err := h.svc.TriggerRun(c.Param("sid"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, run)
}

func (h *ScanHandler) GetRun(c *gin.Context) {
	run, err := h.svc.GetRun(c.Param("rid"))
	if err != nil {
		response.NotFound(c, "扫描运行不存在")
		return
	}
	response.OK(c, run)
}

func (h *ScanHandler) ListRuns(c *gin.Context) {
	page, pageSize := parsePagination(c)
	list, total, err := h.svc.ListRuns(c.Param("sid"), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}
