package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/system/service"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type ClusterHandler struct {
	svc *service.ClusterService
}

func NewClusterHandler(svc *service.ClusterService) *ClusterHandler {
	return &ClusterHandler{svc: svc}
}

func (h *ClusterHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *ClusterHandler) Create(c *gin.Context) {
	var req service.CreateClusterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cl, err := h.svc.Create(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, cl)
}

func (h *ClusterHandler) Get(c *gin.Context) {
	cl, err := h.svc.Get(c.Param("cid"))
	if err != nil {
		response.NotFound(c, "集群不存在")
		return
	}
	response.OK(c, cl)
}

func (h *ClusterHandler) Update(c *gin.Context) {
	var req service.UpdateClusterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cl, err := h.svc.Update(c.Param("cid"), req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, cl)
}

func (h *ClusterHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("cid")); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}
