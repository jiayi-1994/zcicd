package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/project/service"
	appErrors "github.com/zcicd/zcicd-server/pkg/errors"
	"github.com/zcicd/zcicd-server/pkg/response"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*appErrors.AppError); ok {
		switch {
		case appErr.Code >= 50000:
			response.InternalError(c, appErr.Message)
		case appErr.Code >= 40400:
			response.NotFound(c, appErr.Message)
		case appErr.Code >= 40300:
			response.Forbidden(c, appErr.Message)
		case appErr.Code >= 40200:
			response.Error(c, http.StatusConflict, appErr.Code, appErr.Message)
		default:
			response.BadRequest(c, appErr.Message)
		}
		return
	}
	response.InternalError(c, "服务器内部错误")
}

// ==================== Project Handlers ====================

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req service.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetString("user_id")
	project, err := h.svc.CreateProject(c.Request.Context(), userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Created(c, project)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	project, err := h.svc.GetProject(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	project, err := h.svc.UpdateProject(c.Request.Context(), id, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteProject(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	projects, total, err := h.svc.ListProjects(c.Request.Context(), page, pageSize, keyword)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OKWithPage(c, projects, total, page, pageSize)
}

// ==================== Service Handlers ====================

func (h *ProjectHandler) CreateService(c *gin.Context) {
	projectID := c.Param("project_id")
	var req service.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	svc, err := h.svc.CreateService(c.Request.Context(), projectID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Created(c, svc)
}

func (h *ProjectHandler) GetService(c *gin.Context) {
	id := c.Param("id")
	svc, err := h.svc.GetService(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, svc)
}

func (h *ProjectHandler) UpdateService(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	svc, err := h.svc.UpdateService(c.Request.Context(), id, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, svc)
}

func (h *ProjectHandler) DeleteService(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteService(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) ListServices(c *gin.Context) {
	projectID := c.Param("project_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	services, total, err := h.svc.ListServices(c.Request.Context(), projectID, page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OKWithPage(c, services, total, page, pageSize)
}

// ==================== Environment Handlers ====================

func (h *ProjectHandler) CreateEnvironment(c *gin.Context) {
	projectID := c.Param("project_id")
	var req service.CreateEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	env, err := h.svc.CreateEnvironment(c.Request.Context(), projectID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.Created(c, env)
}

func (h *ProjectHandler) GetEnvironment(c *gin.Context) {
	id := c.Param("id")
	env, err := h.svc.GetEnvironment(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, env)
}

func (h *ProjectHandler) UpdateEnvironment(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	env, err := h.svc.UpdateEnvironment(c.Request.Context(), id, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, env)
}

func (h *ProjectHandler) DeleteEnvironment(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteEnvironment(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *ProjectHandler) ListEnvironments(c *gin.Context) {
	projectID := c.Param("project_id")
	envs, err := h.svc.ListEnvironments(c.Request.Context(), projectID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, envs)
}
