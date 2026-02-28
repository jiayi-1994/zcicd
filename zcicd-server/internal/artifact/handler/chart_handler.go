package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/artifact/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type ChartHandler struct {
	svc *service.ChartService
}

func NewChartHandler(svc *service.ChartService) *ChartHandler {
	return &ChartHandler{svc: svc}
}

func (h *ChartHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *ChartHandler) Create(c *gin.Context) {
	var req service.CreateChartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	chart, err := h.svc.Create(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, chart)
}

func (h *ChartHandler) Get(c *gin.Context) {
	chart, err := h.svc.Get(c.Param("name"))
	if err != nil {
		response.NotFound(c, "Chart 不存在")
		return
	}
	response.OK(c, chart)
}

func (h *ChartHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("name")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}
