package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/system/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type IntegrationHandler struct {
	svc *service.IntegrationService
}

func NewIntegrationHandler(svc *service.IntegrationService) *IntegrationHandler {
	return &IntegrationHandler{svc: svc}
}

func (h *IntegrationHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *IntegrationHandler) Create(c *gin.Context) {
	var req service.CreateIntegrationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	i, err := h.svc.Create(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, i)
}

func (h *IntegrationHandler) Update(c *gin.Context) {
	var req service.UpdateIntegrationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	i, err := h.svc.Update(c.Param("iid"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, i)
}

func (h *IntegrationHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("iid")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}
