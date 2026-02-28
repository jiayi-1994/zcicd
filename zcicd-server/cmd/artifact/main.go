package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/artifact/handler"
	"github.com/zcicd/zcicd-server/internal/artifact/repository"
	"github.com/zcicd/zcicd-server/internal/artifact/router"
	"github.com/zcicd/zcicd-server/internal/artifact/service"
	"github.com/zcicd/zcicd-server/pkg/config"
	"github.com/zcicd/zcicd-server/pkg/database"
	"github.com/zcicd/zcicd-server/pkg/logger"
	"github.com/zcicd/zcicd-server/pkg/middleware"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger.Init(cfg)

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Repositories
	regRepo := repository.NewRegistryRepository(db)
	scanRepo := repository.NewScanRepository(db)
	chartRepo := repository.NewChartRepository(db)

	// Services
	regSvc := service.NewRegistryService(regRepo)
	scanSvc := service.NewScanService(scanRepo)
	chartSvc := service.NewChartService(chartRepo)

	// Handlers
	regH := handler.NewRegistryHandler(regSvc)
	scanH := handler.NewScanHandler(scanSvc)
	chartH := handler.NewChartHandler(chartSvc)

	// Gin setup
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	router.RegisterRoutes(api, cfg.JWT.Secret, regH, scanH, chartH)

	port := cfg.Server.Port
	if port == 0 {
		port = 8086
	}
	log.Printf("artifact-service starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
