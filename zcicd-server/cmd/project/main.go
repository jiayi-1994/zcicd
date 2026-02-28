package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/project/handler"
	"github.com/zcicd/zcicd-server/internal/project/repository"
	"github.com/zcicd/zcicd-server/internal/project/router"
	"github.com/zcicd/zcicd-server/internal/project/service"
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

	// Wire dependencies
	projectRepo := repository.NewProjectRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	envRepo := repository.NewEnvironmentRepository(db)

	svc := service.NewProjectService(projectRepo, serviceRepo, envRepo)
	h := handler.NewProjectHandler(svc)

	// Setup Gin
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	router.RegisterRoutes(r, h, cfg.JWT.Secret)

	port := 8082
	log.Printf("project-service starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
