package handler

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/deploy/service"
	"github.com/zcicd/zcicd-server/pkg/response"
	"gorm.io/gorm"
)

type ApprovalHandler struct {
	svc *service.ApprovalService
}

func NewApprovalHandler(svc *service.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{svc: svc}
}

func (h *ApprovalHandler) Approve(c *gin.Context) {
	var req service.ApproveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		req = service.ApproveReq{}
	}
	approverID := c.GetString("user_id")
	record, err := h.svc.Approve(c.Param("id"), approverID, req)
	if err != nil {
		h.handleNotFoundOrInternal(c, err, "审批记录不存在")
		return
	}
	response.OK(c, record)
}

func (h *ApprovalHandler) Reject(c *gin.Context) {
	var req service.RejectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	approverID := c.GetString("user_id")
	record, err := h.svc.Reject(c.Param("id"), approverID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, record)
}

func (h *ApprovalHandler) Get(c *gin.Context) {
	record, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "审批记录不存在")
		return
	}
	response.OK(c, record)
}

func (h *ApprovalHandler) ListPending(c *gin.Context) {
	approverID := c.GetString("user_id")
	page, pageSize := parsePagination(c)
	list, total, err := h.svc.ListPending(approverID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OKWithPage(c, list, total, page, pageSize)
}

func (h *ApprovalHandler) handleNotFoundOrInternal(c *gin.Context, err error, fallbackNotFound string) {
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
