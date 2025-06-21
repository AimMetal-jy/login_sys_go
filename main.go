package main

import (
	"log"
	"login_sys_go/config"
	"login_sys_go/handlers"
	"login_sys_go/models"
	"login_sys_go/routes"
)

func main() {
	// 初始化数据库配置
	dbConfig := config.NewDatabaseConfig()

	// 连接数据库
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 初始化用户仓库
	userRepo := models.NewUserRepository(db)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(userRepo)

	// 设置路由
	router := routes.SetupRoutes(authHandler)

	// 启动服务器
	log.Println("Starting server on port 8000...")
	log.Println("API endpoints:")
	log.Println("  GET  /api/health        - Health check")
	log.Println("  POST /api/auth/register - User registration")
	log.Println("  POST /api/auth/login    - User login")

	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}