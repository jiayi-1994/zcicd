package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/deploy/engine"
	"github.com/zcicd/zcicd-server/internal/deploy/handler"
	"github.com/zcicd/zcicd-server/internal/deploy/repository"
	"github.com/zcicd/zcicd-server/internal/deploy/router"
	"github.com/zcicd/zcicd-server/internal/deploy/service"
	"github.com/zcicd/zcicd-server/pkg/config"
	"github.com/zcicd/zcicd-server/pkg/database"
	"github.com/zcicd/zcicd-server/pkg/k8s"
	"github.com/zcicd/zcicd-server/pkg/logger"
	"github.com/zcicd/zcicd-server/pkg/middleware"
	"github.com/zcicd/zcicd-server/pkg/mq"
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

	redisClient, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	natsClient, err := mq.NewNATSClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect nats: %v", err)
	}

	k8sClient, err := k8s.NewK8sClient("")
	if err != nil {
		log.Printf("warning: k8s not available: %v (Argo CD features disabled)", err)
	}

	// Argo CD engine
	argoNS := "argocd"
	var appManager *engine.AppManager
	var syncCtrl *engine.SyncController
	var rolloutCtrl *engine.RolloutController
	if k8sClient != nil {
		appManager = engine.NewAppManager(k8sClient.DynamicClient, argoNS)
		syncCtrl = engine.NewSyncController(k8sClient.DynamicClient, argoNS)
		rolloutCtrl = engine.NewRolloutController(k8sClient.DynamicClient, argoNS)
	}
	gitopsWriter := engine.NewGitOpsWriter(redisClient)

	// Start health monitor in background
	var healthMonitor *engine.HealthMonitor
	if k8sClient != nil {
		healthMonitor = engine.NewHealthMonitor(k8sClient.DynamicClient, argoNS)
		healthMonitor.Start(context.Background(), func(appName string, status engine.AppStatus) {
			log.Printf("health change: app=%s sync=%s health=%s", appName, status.SyncStatus, status.HealthStatus)
		})
	}
	_ = healthMonitor

	// Repositories
	deployRepo := repository.NewDeployRepository(db)
	approvalRepo := repository.NewApprovalRepository(db)
	envRepo := repository.NewEnvRepository(db)
	// Services
	deploySvc := service.NewDeployService(deployRepo, approvalRepo, appManager, syncCtrl, rolloutCtrl, gitopsWriter, natsClient, argoNS)
	approvalSvc := service.NewApprovalService(approvalRepo, deployRepo)
	envSvc := service.NewEnvService(envRepo)

	// Handlers
	deployH := handler.NewDeployHandler(deploySvc)
	approvalH := handler.NewApprovalHandler(approvalSvc)
	envH := handler.NewEnvHandler(envSvc)

	// Gin setup
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	router.RegisterRoutes(api, cfg.JWT.Secret, deployH, approvalH, envH)

	port := cfg.Server.Port
	if port == 0 {
		port = 8084
	}
	log.Printf("deploy-service starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
