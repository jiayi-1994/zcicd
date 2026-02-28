package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zcicd/zcicd-server/internal/system/handler"
	"github.com/zcicd/zcicd-server/pkg/middleware"
)

func RegisterRoutes(r *gin.RouterGroup, jwtSecret string, notifyH *handler.NotifyHandler, clusterH *handler.ClusterHandler, integrationH *handler.IntegrationHandler, auditH *handler.AuditHandler, dashH *handler.DashboardHandler) {
	auth := middleware.JWTAuth(jwtSecret)

	sys := r.Group("/system")
	sys.Use(auth)
	{
		// Dashboard
		dash := sys.Group("/dashboard")
		dash.GET("/overview", dashH.Overview)
		dash.GET("/trends", dashH.Trends)

		clusters := sys.Group("/clusters")
		clusters.GET("", clusterH.List)
		clusters.POST("", clusterH.Create)
		clusters.GET("/:cid", clusterH.Get)
		clusters.PUT("/:cid", clusterH.Update)
		clusters.DELETE("/:cid", clusterH.Delete)

		integrations := sys.Group("/integrations")
		integrations.GET("", integrationH.List)
		integrations.POST("", integrationH.Create)
		integrations.PUT("/:iid", integrationH.Update)
		integrations.DELETE("/:iid", integrationH.Delete)

		sys.GET("/audit-logs", auditH.List)
	}

	notify := r.Group("/notifications")
	notify.Use(auth)
	{
		channels := notify.Group("/channels")
		channels.GET("", notifyH.ListChannels)
		channels.POST("", notifyH.CreateChannel)
		channels.PUT("/:cid", notifyH.UpdateChannel)
		channels.DELETE("/:cid", notifyH.DeleteChannel)

		notify.GET("/rules", notifyH.ListRules)
		notify.POST("/rules", notifyH.CreateRule)
		notify.PUT("/rules/:rid", notifyH.UpdateRule)
	}
}
