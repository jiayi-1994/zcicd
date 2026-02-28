package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/artifact/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type RegistryHandler struct {
	svc *service.RegistryService
}

func NewRegistryHandler(svc *service.RegistryService) *RegistryHandler {
	return &RegistryHandler{svc: svc}
}

func (h *RegistryHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *RegistryHandler) Create(c *gin.Context) {
	var req service.CreateRegistryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	reg, err := h.svc.Create(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, reg)
}

func (h *RegistryHandler) Get(c *gin.Context) {
	reg, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "镜像仓库不存在")
		return
	}
	response.OK(c, reg)
}

func (h *RegistryHandler) Update(c *gin.Context) {
	var req service.UpdateRegistryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	reg, err := h.svc.Update(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, reg)
}

func (h *RegistryHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}
