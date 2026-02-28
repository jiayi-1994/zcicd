package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/deploy/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type EnvHandler struct {
	svc *service.EnvService
}

func NewEnvHandler(svc *service.EnvService) *EnvHandler {
	return &EnvHandler{svc: svc}
}

func (h *EnvHandler) ListVariables(c *gin.Context) {
	list, err := h.svc.ListVariables(c.Param("env_id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *EnvHandler) CreateVariable(c *gin.Context) {
	var req service.EnvVariableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	v, err := h.svc.CreateVariable(c.Param("env_id"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, v)
}

func (h *EnvHandler) UpdateVariable(c *gin.Context) {
	var req service.EnvVariableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	v, err := h.svc.UpdateVariable(c.Param("var_id"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, v)
}

func (h *EnvHandler) DeleteVariable(c *gin.Context) {
	if err := h.svc.DeleteVariable(c.Param("var_id")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *EnvHandler) BatchUpsertVariables(c *gin.Context) {
	var req service.BatchEnvVariablesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.BatchUpsertVariables(c.Param("env_id"), req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *EnvHandler) GetQuota(c *gin.Context) {
	q, err := h.svc.GetQuota(c.Param("env_id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, q)
}

func (h *EnvHandler) UpsertQuota(c *gin.Context) {
	var req service.EnvResourceQuotaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	q, err := h.svc.UpsertQuota(c.Param("env_id"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, q)
}
