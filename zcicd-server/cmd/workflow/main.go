package main

import (
	"fmt"
	"log"

	"github.com/zcicd/zcicd-server/internal/workflow/engine"
	"github.com/zcicd/zcicd-server/internal/workflow/handler"
	"github.com/zcicd/zcicd-server/internal/workflow/repository"
	"github.com/zcicd/zcicd-server/internal/workflow/router"
	"github.com/zcicd/zcicd-server/internal/workflow/service"
	"github.com/zcicd/zcicd-server/pkg/config"
	"github.com/zcicd/zcicd-server/pkg/database"
	"github.com/zcicd/zcicd-server/pkg/k8s"
	"github.com/zcicd/zcicd-server/pkg/logger"
	"github.com/zcicd/zcicd-server/pkg/middleware"
	"github.com/zcicd/zcicd-server/pkg/mq"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger.Init(cfg)

	// Initialize database
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize Redis
	redisClient, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	// Initialize NATS
	natsClient, err := mq.NewNATSClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect nats: %v", err)
	}

	// Initialize K8s client (for Tekton CRD management)
	k8sClient, err := k8s.NewK8sClient("")
	if err != nil {
		log.Printf("warning: failed to connect k8s: %v (Tekton features disabled)", err)
	}

	// Initialize Tekton adapter
	var crdManager *engine.CRDManager
	namespace := "zcicd"
	if k8sClient != nil {
		crdManager = engine.NewCRDManager(k8sClient.DynamicClient)
	}

	// Initialize repositories
	workflowRepo := repository.NewWorkflowRepository(db)
	buildRepo := repository.NewBuildRepository(db)
	templateRepo := repository.NewTemplateRepository(db)

	// Initialize services
	workflowSvc := service.NewWorkflowService(workflowRepo, buildRepo, natsClient)
	buildSvc := service.NewBuildService(buildRepo, templateRepo, crdManager, natsClient, namespace)
	templateSvc := service.NewTemplateService(templateRepo)

	// Initialize handlers
	workflowHandler := handler.NewWorkflowHandler(workflowSvc)
	buildHandler := handler.NewBuildHandler(buildSvc)
	templateHandler := handler.NewTemplateHandler(templateSvc)
	wsHandler := handler.NewWSHandler(redisClient)
	webhookHandler := handler.NewWebhookHandler(workflowSvc, buildSvc)

	// Setup Gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// Health check
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	router.RegisterRoutes(api, cfg.JWT.Secret, workflowHandler, buildHandler, templateHandler, wsHandler, webhookHandler)

	port := cfg.Server.Port
	if port == 0 {
		port = 8083
	}
	log.Printf("workflow-service starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
