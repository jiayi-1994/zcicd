package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/quality/handler"
	"github.com/zcicd/zcicd-server/internal/quality/repository"
	"github.com/zcicd/zcicd-server/internal/quality/router"
	"github.com/zcicd/zcicd-server/internal/quality/service"
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
	testRepo := repository.NewTestRepository(db)
	scanRepo := repository.NewScanRepository(db)
	gateRepo := repository.NewQualityGateRepository(db)

	// Services
	testSvc := service.NewTestService(testRepo)
	scanSvc := service.NewScanService(scanRepo)
	gateSvc := service.NewQualityGateService(gateRepo)

	// Handlers
	testH := handler.NewTestHandler(testSvc)
	scanH := handler.NewScanHandler(scanSvc)
	gateH := handler.NewQualityGateHandler(gateSvc)

	// Gin setup
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	router.RegisterRoutes(api, cfg.JWT.Secret, testH, scanH, gateH)

	port := cfg.Server.Port
	if port == 0 {
		port = 8085
	}
	log.Printf("quality-service starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
