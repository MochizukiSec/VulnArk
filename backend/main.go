package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vuln-management/config"
	"vuln-management/controllers"
	"vuln-management/middleware"
	"vuln-management/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件，使用环境变量")
	}

	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer config.CloseDB()

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{getEnv("ALLOWED_ORIGINS", "*")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 创建服务
	userService := services.NewUserService(config.Database)
	vulnService := services.NewVulnerabilityService(config.Database)
	reportService := services.NewReportService(config.Database)
	// 创建AI分析服务
	aiAnalysisService := services.NewAIAnalysisService(config.Database)

	// 创建控制器
	userController := controllers.NewUserController(userService)
	vulnController := controllers.NewVulnerabilityController(vulnService)
	reportController := controllers.NewReportController(reportService)
	// 创建AI分析控制器
	aiAnalysisController := controllers.NewAIAnalysisController(aiAnalysisService)

	// 设置API路由
	setupRoutes(r, userController, vulnController, reportController, aiAnalysisController)

	// 获取端口，默认为8000
	port := getEnv("PORT", "8000")

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("服务器运行于 http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器优雅退出")
}

func setupRoutes(r *gin.Engine, userController *controllers.UserController, vulnController *controllers.VulnerabilityController, reportController *controllers.ReportController, aiAnalysisController *controllers.AIAnalysisController) {
	// 公共路由
	public := r.Group("/api")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "服务运行正常"})
		})

		// 身份验证相关路由
		auth := public.Group("/auth")
		{
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		// 漏洞相关路由
		vulnerabilities := protected.Group("/vulnerabilities")
		{
			vulnerabilities.GET("", vulnController.GetAllVulnerabilities)
			vulnerabilities.GET("/:id", vulnController.GetVulnerabilityByID)
			vulnerabilities.POST("", vulnController.CreateVulnerability)
			vulnerabilities.PUT("/:id", vulnController.UpdateVulnerability)
			vulnerabilities.DELETE("/:id", vulnController.DeleteVulnerability)
			vulnerabilities.POST("/import", vulnController.ImportVulnerabilities)
		}

		// 用户相关路由
		users := protected.Group("/users")
		{
			users.GET("/me", userController.GetCurrentUser)
			users.PUT("/me", userController.UpdateCurrentUser)
			users.GET("", middleware.RequireAdmin(), userController.GetAllUsers)
			users.POST("", middleware.RequireAdmin(), userController.CreateUser)
			users.PUT("/:id", middleware.RequireAdmin(), userController.UpdateUser)
			users.DELETE("/:id", middleware.RequireAdmin(), userController.DeleteUser)
		}

		// 仪表盘数据
		protected.GET("/dashboard", userController.GetDashboardData)

		// 报告
		reports := protected.Group("/reports")
		{
			// 生成报告相关
			reports.GET("/summary", reportController.GenerateSummaryReport)
			reports.GET("/detailed", reportController.GenerateDetailedReport)

			// 报告管理
			reports.GET("", reportController.GetAllReports)
			reports.GET("/:id", reportController.GetReportByID)
			reports.POST("", reportController.CreateReport)
			reports.DELETE("/:id", reportController.DeleteReport)
		}

		// AI分析相关路由
		aiAnalysisController.RegisterRoutes(protected)
	}
}

// getEnv 获取环境变量或返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
