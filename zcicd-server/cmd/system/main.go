package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/system/handler"
	"github.com/zcicd/zcicd-server/internal/system/repository"
	"github.com/zcicd/zcicd-server/internal/system/router"
	"github.com/zcicd/zcicd-server/internal/system/service"
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
	notifyRepo := repository.NewNotifyRepository(db)
	ruleRepo := repository.NewRuleRepository(db)
	clusterRepo := repository.NewClusterRepository(db)
	integrationRepo := repository.NewIntegrationRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	dashRepo := repository.NewDashboardRepository(db)

	// Services
	notifySvc := service.NewNotifyService(notifyRepo, ruleRepo)
	clusterSvc := service.NewClusterService(clusterRepo)
	integrationSvc := service.NewIntegrationService(integrationRepo)
	auditSvc := service.NewAuditService(auditRepo)
	dashSvc := service.NewDashboardService(dashRepo)

	// Handlers
	notifyH := handler.NewNotifyHandler(notifySvc)
	clusterH := handler.NewClusterHandler(clusterSvc)
	integrationH := handler.NewIntegrationHandler(integrationSvc)
	auditH := handler.NewAuditHandler(auditSvc)
	dashH := handler.NewDashboardHandler(dashSvc)

	// Gin setup
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	router.RegisterRoutes(api, cfg.JWT.Secret, notifyH, clusterH, integrationH, auditH, dashH)

	port := cfg.Server.Port
	if port == 0 {
		port = 8087
	}
	log.Printf("system-service starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
