package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/artifact/handler"
	"github.com/zcicd/zcicd-server/pkg/middleware"
)

func RegisterRoutes(r *gin.RouterGroup, jwtSecret string, regH *handler.RegistryHandler, scanH *handler.ScanHandler, chartH *handler.ChartHandler) {
	auth := middleware.JWTAuth(jwtSecret)

	artifacts := r.Group("/artifacts")
	artifacts.Use(auth)
	{
		registries := artifacts.Group("/registries")
		registries.GET("", regH.List)
		registries.POST("", regH.Create)
		registries.GET("/:id", regH.Get)
		registries.PUT("/:id", regH.Update)
		registries.DELETE("/:id", regH.Delete)

		images := artifacts.Group("/images")
		images.GET("/:name/scan", scanH.GetScanResults)
		images.POST("/:name/scan", scanH.TriggerScan)

		charts := artifacts.Group("/charts")
		charts.GET("", chartH.List)
		charts.POST("", chartH.Create)
		charts.GET("/:name", chartH.Get)
		charts.DELETE("/:name", chartH.Delete)
	}
}
