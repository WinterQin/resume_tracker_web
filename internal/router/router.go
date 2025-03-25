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
		user := authorized.Group("/user")
		{

			//分页查找 带筛选和搜索
			user.GET("/profile", userHandler.GetProfile)
			//新增
			//user.POST("", userHandler.CreateApplication)
			//修改
			user.PUT("/:id", userHandler.UpdateProfile) // 新的更新路由
			//删除
			user.DELETE("/:id", userHandler.DeleteAccount)

		}

		// 申请相关路由
		applications := authorized.Group("/applications")
		{

			//分页查找 带筛选和搜索
			applications.GET("", applicationHandler.GetApplications)
			//新增
			applications.POST("", applicationHandler.CreateApplication)
			//修改
			applications.PUT("/:id", applicationHandler.UpdateApplication) // 新的更新路由
			//删除
			applications.DELETE("/:id", applicationHandler.DeleteApplication)

			//统计信息
			applications.GET("/statistics", applicationHandler.GetStatistics)
			//最近申请
			applications.GET("/recent", applicationHandler.GetRecentApplications) // 获取最近5条申请
			//更新状态
			applications.PATCH("/status", applicationHandler.UpdateStatus)

		}
	}

	return r

}
