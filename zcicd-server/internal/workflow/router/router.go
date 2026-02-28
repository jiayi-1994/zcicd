package router

import (
	"github.com/zcicd/zcicd-server/internal/workflow/handler"
	"github.com/zcicd/zcicd-server/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, jwtSecret string, workflowHandler *handler.WorkflowHandler, buildHandler *handler.BuildHandler, templateHandler *handler.TemplateHandler, wsHandler *handler.WSHandler, webhookHandler *handler.WebhookHandler) {
	auth := middleware.JWTAuth(jwtSecret)

	// Workflow routes
	workflows := r.Group("/workflows")
	workflows.Use(auth)
	{
		workflows.GET("", workflowHandler.List)
		workflows.POST("", workflowHandler.Create)
		workflows.GET("/:id", workflowHandler.Get)
		workflows.PUT("/:id", workflowHandler.Update)
		workflows.DELETE("/:id", workflowHandler.Delete)
		workflows.POST("/:id/trigger", workflowHandler.Trigger)
		workflows.GET("/:id/runs", workflowHandler.ListRuns)
		workflows.GET("/:id/runs/:run_id", workflowHandler.GetRun)
		workflows.POST("/:id/runs/:run_id/cancel", workflowHandler.CancelRun)
		workflows.POST("/:id/runs/:run_id/retry", workflowHandler.RetryRun)
	}

	// Webhook routes (no auth - verified by signature/token)
	webhooks := r.Group("/webhooks")
	{
		webhooks.POST("/github", webhookHandler.HandleGitHub)
		webhooks.POST("/gitlab", webhookHandler.HandleGitLab)
	}

	// Build config routes
	buildConfigs := r.Group("/build-configs")
	buildConfigs.Use(auth)
	{
		buildConfigs.GET("", buildHandler.ListConfigs)
		buildConfigs.POST("", buildHandler.CreateConfig)
		buildConfigs.GET("/:id", buildHandler.GetConfig)
		buildConfigs.PUT("/:id", buildHandler.UpdateConfig)
		buildConfigs.DELETE("/:id", buildHandler.DeleteConfig)
		buildConfigs.POST("/:id/trigger", buildHandler.TriggerBuild)
	}

	// Build runs
	buildRuns := r.Group("/build-runs")
	buildRuns.Use(auth)
	{
		buildRuns.GET("", buildHandler.ListRuns)
		buildRuns.GET("/:run_id", buildHandler.GetRun)
		buildRuns.POST("/:run_id/cancel", buildHandler.CancelRun)
		buildRuns.GET("/:run_id/logs", buildHandler.GetRunLogs)
	}

	// WebSocket route (outside auth group, uses token query param)
	r.GET("/build-runs/:run_id/logs/ws", wsHandler.HandleBuildLogs)

	// Build templates
	templates := r.Group("/build-templates")
	templates.Use(auth)
	{
		templates.GET("", templateHandler.List)
		templates.POST("", templateHandler.Create)
		templates.GET("/:id", templateHandler.Get)
		templates.PUT("/:id", templateHandler.Update)
		templates.DELETE("/:id", templateHandler.Delete)
	}
}
