package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/artifact/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type ScanHandler struct {
	svc *service.ScanService
}

func NewScanHandler(svc *service.ScanService) *ScanHandler {
	return &ScanHandler{svc: svc}
}

func (h *ScanHandler) TriggerScan(c *gin.Context) {
	var req service.TriggerScanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	scan, err := h.svc.TriggerScan(c.Param("name"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, scan)
}

func (h *ScanHandler) GetScanResults(c *gin.Context) {
	registryID := c.Query("registry_id")
	if registryID == "" {
		response.BadRequest(c, "registry_id is required")
		return
	}
	list, err := h.svc.ListByImage(registryID, c.Param("name"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}
