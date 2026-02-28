package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/system/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	stats, err := h.svc.GetOverview()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, stats)
}

func (h *DashboardHandler) Trends(c *gin.Context) {
	days := 30
	if d := c.Query("days"); d != "" {
		if v, err := strconv.Atoi(d); err == nil && v > 0 {
			days = v
		}
	}
	trends, err := h.svc.GetTrends(days)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, trends)
}
