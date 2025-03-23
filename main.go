package main

import (
	"internship-manager/internal/middleware"
	"internship-manager/internal/router"
	"internship-manager/pkg/database"
	"log"
	"os"
	"strconv"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 从环境变量获取数据库配置
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))

	// 初始化MySQL连接
	err := database.InitMySQL(&database.MySQLConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         dbPort,
		Username:     getEnv("DB_USER", "resume_winter"),
		Password:     getEnv("DB_PASSWORD", "resume_qwt123456"),
		DBName:       getEnv("DB_NAME", "internship_manager"),
		Charset:      getEnv("DB_CHARSET", "utf8mb4"),
		MaxIdleConns: maxIdleConns,
		MaxOpenConns: maxOpenConns,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 从环境变量获取JWT密钥
	jwtKey := getEnv("JWT_KEY", "winter-key")
	middleware.InitJWT(jwtKey)

	// 设置路由
	r := router.SetupRouter()

	// 从环境变量获取服务器端口
	port := getEnv("SERVER_PORT", "8080")

	// 启动服务器
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
