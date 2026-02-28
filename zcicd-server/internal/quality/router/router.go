package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/quality/handler"
	"github.com/zcicd/zcicd-server/pkg/middleware"
)

func RegisterRoutes(r *gin.RouterGroup, jwtSecret string, testH *handler.TestHandler, scanH *handler.ScanHandler, gateH *handler.QualityGateHandler) {
	auth := middleware.JWTAuth(jwtSecret)

	projects := r.Group("/projects/:project_id")
	projects.Use(auth)
	{
		// Test configs
		tests := projects.Group("/tests")
		tests.GET("", testH.ListConfigs)
		tests.POST("", testH.CreateConfig)
		tests.GET("/:id", testH.GetConfig)
		tests.PUT("/:id", testH.UpdateConfig)
		tests.DELETE("/:id", testH.DeleteConfig)
		tests.POST("/:id/run", testH.TriggerRun)
		tests.GET("/:id/runs", testH.ListRuns)
		tests.GET("/:id/runs/:run_id", testH.GetRun)

		// Scan configs
		scans := projects.Group("/scans")
		scans.GET("", scanH.ListConfigs)
		scans.POST("", scanH.CreateConfig)
		scans.GET("/:sid", scanH.GetConfig)
		scans.PUT("/:sid", scanH.UpdateConfig)
		scans.DELETE("/:sid", scanH.DeleteConfig)
		scans.POST("/:sid/run", scanH.TriggerRun)
		scans.GET("/:sid/runs", scanH.ListRuns)

		// Quality gate
		projects.GET("/quality-gate", gateH.Get)
		projects.PUT("/quality-gate", gateH.Upsert)
	}

	// Scan run detail (cross-project)
	scanRuns := r.Group("/scans/runs")
	scanRuns.Use(auth)
	scanRuns.GET("/:rid", scanH.GetRun)
}
