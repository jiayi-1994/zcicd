package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/auth/handler"
	"github.com/zcicd/zcicd-server/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine, h *handler.AuthHandler, jwtSecret string) {
	auth := r.Group("/api/v1/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)

	protected := auth.Group("")
	protected.Use(middleware.JWTAuth(jwtSecret))
	protected.GET("/profile", h.GetProfile)
	protected.PUT("/profile", h.UpdateProfile)
	protected.PUT("/password", h.ChangePassword)

	admin := auth.Group("/users")
	admin.Use(middleware.JWTAuth(jwtSecret), middleware.AdminRequired())
	admin.GET("", h.ListUsers)
}
