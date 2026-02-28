package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/system/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type NotifyHandler struct {
	svc *service.NotifyService
}

func NewNotifyHandler(svc *service.NotifyService) *NotifyHandler {
	return &NotifyHandler{svc: svc}
}

func (h *NotifyHandler) ListChannels(c *gin.Context) {
	list, err := h.svc.ListChannels()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *NotifyHandler) CreateChannel(c *gin.Context) {
	var req service.CreateChannelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	ch, err := h.svc.CreateChannel(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, ch)
}

func (h *NotifyHandler) UpdateChannel(c *gin.Context) {
	var req service.UpdateChannelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	ch, err := h.svc.UpdateChannel(c.Param("cid"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, ch)
}

func (h *NotifyHandler) DeleteChannel(c *gin.Context) {
	if err := h.svc.DeleteChannel(c.Param("cid")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *NotifyHandler) ListRules(c *gin.Context) {
	list, err := h.svc.ListRules()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *NotifyHandler) CreateRule(c *gin.Context) {
	var req service.CreateRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	r, err := h.svc.CreateRule(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, r)
}

func (h *NotifyHandler) UpdateRule(c *gin.Context) {
	var req service.UpdateRuleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	r, err := h.svc.UpdateRule(c.Param("rid"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, r)
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
