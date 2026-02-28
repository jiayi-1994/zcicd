package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/deploy/handler"
	"github.com/zcicd/zcicd-server/pkg/middleware"
)

func RegisterRoutes(r *gin.RouterGroup, jwtSecret string, deployH *handler.DeployHandler, approvalH *handler.ApprovalHandler, envH *handler.EnvHandler) {
	auth := middleware.JWTAuth(jwtSecret)

	deploys := r.Group("/deploys")
	deploys.Use(auth)
	{
		deploys.GET("", deployH.ListConfigs)
		deploys.GET("/by-env", deployH.ListConfigsByEnv)
		deploys.POST("", deployH.CreateConfig)
		deploys.GET("/:id", deployH.GetConfig)
		deploys.PUT("/:id", deployH.UpdateConfig)
		deploys.DELETE("/:id", deployH.DeleteConfig)
		deploys.POST("/:id/sync", deployH.TriggerSync)
		deploys.POST("/:id/rollback", deployH.Rollback)
		deploys.GET("/:id/status", deployH.GetStatus)
		deploys.GET("/:id/resources", deployH.GetResources)
		deploys.GET("/:id/history", deployH.ListHistories)
		deploys.GET("/:id/history/:history_id", deployH.GetHistory)
		deploys.GET("/:id/rollout", deployH.GetRolloutStatus)
		deploys.POST("/:id/rollout/promote", deployH.PromoteRollout)
		deploys.POST("/:id/rollout/abort", deployH.AbortRollout)
	}

	approvals := r.Group("/approvals")
	approvals.Use(auth)
	{
		approvals.GET("/pending", approvalH.ListPending)
		approvals.GET("/:id", approvalH.Get)
		approvals.POST("/:id/approve", approvalH.Approve)
		approvals.POST("/:id/reject", approvalH.Reject)
	}

	envVars := r.Group("/environments")
	envVars.Use(auth)
	{
		envVars.GET("/:env_id/variables", envH.ListVariables)
		envVars.POST("/:env_id/variables", envH.CreateVariable)
		envVars.PUT("/:env_id/variables/:var_id", envH.UpdateVariable)
		envVars.DELETE("/:env_id/variables/:var_id", envH.DeleteVariable)
		envVars.PUT("/:env_id/variables/batch", envH.BatchUpsertVariables)
		envVars.GET("/:env_id/quota", envH.GetQuota)
		envVars.PUT("/:env_id/quota", envH.UpsertQuota)
	}
}
