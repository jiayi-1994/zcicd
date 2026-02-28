package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zcicd/zcicd-server/internal/auth/handler"
	"github.com/zcicd/zcicd-server/internal/auth/repository"
	"github.com/zcicd/zcicd-server/internal/auth/router"
	"github.com/zcicd/zcicd-server/internal/auth/service"
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

	rdb, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	// Wire dependencies
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg, rdb)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	router.RegisterRoutes(r, authHandler, cfg.JWT.Secret)

	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8081
	}
	log.Printf("auth-service starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
