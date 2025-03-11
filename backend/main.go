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
	"vuln-management/routes"
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

	// 初始化数据库数据
	if err := config.InitializeDatabase(); err != nil {
		log.Printf("初始化数据失败: %v", err)
		// 仅打印警告，不中止程序
	} else {
		log.Println("数据库初始化完成")
	}

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // 允许所有源
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	r.Use(cors.New(corsConfig))

	// 添加OPTIONS请求处理程序
	r.OPTIONS("/*path", func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if corsConfig.AllowAllOrigins {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			// 检查是否是允许的源
			allowed := false
			for _, allowedOrigin := range corsConfig.AllowOrigins {
				if allowedOrigin == origin {
					allowed = true
					break
				}
			}
			if allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin, Content-Type, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
		return
	})

	// 创建服务
	userService := services.NewUserService(config.Database)
	vulnService := services.NewVulnerabilityService(config.Database)
	reportService := services.NewReportService(config.Database)
	aiAnalysisService := services.NewAIAnalysisService(config.Database)
	assetService := services.NewAssetService(config.Database)

	// 创建控制器
	userController := controllers.NewUserController(userService)
	vulnController := controllers.NewVulnerabilityController(vulnService)
	reportController := controllers.NewReportController(reportService)
	aiAnalysisController := controllers.NewAIAnalysisController(aiAnalysisService)
	assetController := controllers.NewAssetController(assetService)

	// 设置路由
	setupRoutes(r, userController, vulnController, reportController, aiAnalysisController, assetController)

	// 使用routes包中的SetupRouter方法设置漏洞库路由
	routes.SetupRouter(r)

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

func setupRoutes(r *gin.Engine, userController *controllers.UserController, vulnController *controllers.VulnerabilityController, reportController *controllers.ReportController, aiAnalysisController *controllers.AIAnalysisController, assetController *controllers.AssetController) {
	// 处理重复路径问题：添加一个前缀为"/api/api"的路由组，将请求重定向到正确的路径
	apiRedirect := r.Group("/api/api")
	{
		// 处理/api/api/auth/login重复路径
		apiRedirect.POST("/auth/login", userController.Login)

		// 其他可能的重复路径 - 可根据需要添加更多
		apiRedirect.Any("/*path", func(c *gin.Context) {
			// 获取path参数
			path := c.Param("path")
			// 重定向到正确的路径
			c.Request.URL.Path = "/api" + path
			// 继续路由处理
			r.HandleContext(c)
		})
	}

	// 公共路由
	public := r.Group("/api")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "服务运行正常"})
		})

		// 系统初始化路由
		public.POST("/initialize/admin", userController.InitializeAdmin)

		// 身份验证相关路由
		auth := public.Group("/auth")
		{
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

		// 资产管理相关路由
		assets := protected.Group("/assets")
		{
			assets.GET("", assetController.SearchAssets)
			assets.GET("/:id", assetController.GetAssetByID)
			assets.POST("", assetController.CreateAsset)
			assets.PUT("/:id", assetController.UpdateAsset)
			assets.DELETE("/:id", assetController.DeleteAsset)

			// 资产与漏洞关联
			assets.GET("/:id/vulnerabilities", assetController.GetAssetVulnerabilities)
			assets.POST("/:id/vulnerabilities/:vulnId", assetController.AddVulnerabilityToAsset)
			assets.DELETE("/:id/vulnerabilities/:vulnId", assetController.RemoveVulnerabilityFromAsset)

			// 资产备注
			assets.POST("/:id/notes", assetController.AddAssetNote)
		}
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
