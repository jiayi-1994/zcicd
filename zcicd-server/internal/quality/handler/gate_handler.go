package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/quality/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type QualityGateHandler struct {
	svc *service.QualityGateService
}

func NewQualityGateHandler(svc *service.QualityGateService) *QualityGateHandler {
	return &QualityGateHandler{svc: svc}
}

func (h *QualityGateHandler) Get(c *gin.Context) {
	g, err := h.svc.Get(c.Param("project_id"))
	if err != nil {
		response.NotFound(c, "质量门禁未配置")
		return
	}
	response.OK(c, g)
}

func (h *QualityGateHandler) Upsert(c *gin.Context) {
	var req service.QualityGateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	g, err := h.svc.Upsert(c.Param("project_id"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, g)
}
