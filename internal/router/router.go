package router

import (
	"internship-manager/internal/handler"
	"internship-manager/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // 允许所有源访问
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(config))

	// 创建处理器实例
	userHandler := handler.NewUserHandler()
	applicationHandler := handler.NewApplicationHandler()

	// 公开路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(middleware.JWTAuth())
	{
		// 申请相关路由
		applications := authorized.Group("/applications")
		{
			applications.GET("", applicationHandler.GetApplications)
			applications.POST("", applicationHandler.CreateApplication)

			applications.GET("/recent", applicationHandler.GetRecentApplications) // 获取最近5条申请
			applications.DELETE("/:id", applicationHandler.DeleteApplication)     // 新增的删除路由
			applications.PUT("/:id", applicationHandler.UpdateApplication)        // 新的更新路由
			applications.PUT("/status", applicationHandler.UpdateStatus)
			applications.PUT("/event", applicationHandler.UpdateEvent)

			applications.GET("/statistics", applicationHandler.GetStatistics)
			applications.GET("/upcoming-events", applicationHandler.GetUpcomingEvents)
		}
	}

	return r

}
