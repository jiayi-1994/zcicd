package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/project/handler"
	"github.com/zcicd/zcicd-server/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine, h *handler.ProjectHandler, jwtSecret string) {
	api := r.Group("/api/v1")
	api.Use(middleware.JWTAuth(jwtSecret))

	// Project routes
	projects := api.Group("/projects")
	projects.POST("", h.CreateProject)
	projects.GET("", h.ListProjects)
	projects.GET("/:id", h.GetProject)
	projects.PUT("/:id", h.UpdateProject)
	projects.DELETE("/:id", h.DeleteProject)

	// Service routes (nested under project)
	projects.POST("/:id/services", h.CreateService)
	projects.GET("/:id/services", h.ListServices)

	// Service routes (standalone)
	services := api.Group("/services")
	services.GET("/:id", h.GetService)
	services.PUT("/:id", h.UpdateService)
	services.DELETE("/:id", h.DeleteService)

	// Environment routes (nested under project)
	projects.POST("/:id/environments", h.CreateEnvironment)
	projects.GET("/:id/environments", h.ListEnvironments)

	// Environment routes (standalone)
	envs := api.Group("/environments")
	envs.GET("/:id", h.GetEnvironment)
	envs.PUT("/:id", h.UpdateEnvironment)
	envs.DELETE("/:id", h.DeleteEnvironment)
}
