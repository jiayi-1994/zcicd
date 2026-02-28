package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/system/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type AuditHandler struct {
	svc *service.AuditService
}

func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	filters := map[string]string{
		"user_id":       c.Query("user_id"),
		"action":        c.Query("action"),
		"resource_type": c.Query("resource_type"),
		"project_id":    c.Query("project_id"),
	}
	list, total, err := h.svc.List(page, pageSize, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}
