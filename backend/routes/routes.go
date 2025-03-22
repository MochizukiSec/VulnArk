package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vulnark/vulnark/controllers"
	"github.com/vulnark/vulnark/middleware"
)

// SetupRouter 配置路由
func SetupRouter(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.CORSMiddleware())

	// 公共路由组
	public := r.Group("/api/v1")
	{
		// 健康检查
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// 用户认证
		userController := new(controllers.UserController)
		public.POST("/auth/login", userController.LoginV2)
	}

	// 需要认证的路由组
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.JWTAuthMiddleware())
	{
		// 用户管理
		userController := new(controllers.UserController)
		authorized.GET("/user/info", userController.GetUserInfo)
		authorized.PUT("/user/update", userController.UpdateUserV2)

		// 仅管理员访问的路由
		admin := authorized.Group("/admin")
		admin.Use(middleware.RequireAdmin())
		{
			admin.GET("/users", userController.ListUsers)
			admin.POST("/users", userController.CreateUser)
			admin.DELETE("/user/:id", userController.DeleteUser)
			admin.PUT("/user/:id/role", userController.ChangeUserRole)
		}

		// 资产管理路由
		assetController := new(controllers.AssetController)
		authorized.GET("/assets", assetController.ListAssets)
		authorized.POST("/assets", assetController.CreateAsset)
		authorized.POST("/assets/import", assetController.BatchImportAssets)
		authorized.POST("/assets/batch-delete", assetController.BatchDeleteAssets)
		authorized.GET("/assets/export", assetController.ExportAssets)
		authorized.GET("/assets/:id", assetController.GetAssetByID)
		authorized.PUT("/assets/:id", assetController.UpdateAsset)
		authorized.DELETE("/assets/:id", assetController.DeleteAsset)
		authorized.GET("/assets/:id/vulnerabilities", assetController.GetAssetVulnerabilities)

		// 漏洞管理路由
		vulnerabilityController := new(controllers.VulnerabilityController)
		authorized.GET("/vulnerabilities", vulnerabilityController.ListVulnerabilities)
		authorized.GET("/vulnerabilities/:id", vulnerabilityController.GetVulnerabilityByID)
		authorized.POST("/vulnerabilities", vulnerabilityController.CreateVulnerability)
		authorized.PUT("/vulnerabilities/:id", vulnerabilityController.UpdateVulnerability)
		authorized.DELETE("/vulnerabilities/:id", vulnerabilityController.DeleteVulnerability)
		authorized.POST("/vulnerabilities/import", vulnerabilityController.BatchImportVulnerabilities)
		authorized.POST("/vulnerabilities/batch-delete", vulnerabilityController.BatchDeleteVulnerabilities)

		// 漏洞分发路由
		assignmentController := new(controllers.VulnerabilityAssignmentController)

		// 添加一个调试处理函数
		authorized.GET("/debug-routes", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "路由调试信息",
				"data": gin.H{
					"path":     c.Request.URL.Path,
					"method":   c.Request.Method,
					"fullPath": c.FullPath(),
				},
			})
		})

		authorized.POST("/vulnerabilities/:id/assign", assignmentController.AssignVulnerability)
		authorized.GET("/vulnerabilities/:id/assignments", assignmentController.GetAssignmentsByVulnerability)
		authorized.GET("/assignments", assignmentController.ListAssignments)
		authorized.GET("/assignments/my", assignmentController.GetMyAssignments)
		authorized.GET("/assignments/:id", assignmentController.GetAssignmentDetails)
		authorized.PUT("/assignments/:id/status", assignmentController.UpdateAssignmentStatus)
		authorized.DELETE("/assignments/:id", assignmentController.DeleteAssignment)

		// 知识库路由
		knowledgeController := new(controllers.KnowledgeController)
		authorized.GET("/knowledge", knowledgeController.ListKnowledgeItems)
		authorized.GET("/knowledge/types", knowledgeController.GetKnowledgeTypes)
		authorized.GET("/knowledge/categories", knowledgeController.GetKnowledgeCategories)
		authorized.GET("/knowledge/:id", knowledgeController.GetKnowledgeByID)
		authorized.POST("/knowledge", knowledgeController.CreateKnowledge)
		authorized.PUT("/knowledge/:id", knowledgeController.UpdateKnowledge)
		authorized.DELETE("/knowledge/:id", knowledgeController.DeleteKnowledge)

		// 漏洞库路由
		vulnDBController := new(controllers.VulnDBController)
		authorized.GET("/vulndb", vulnDBController.ListVulnDBEntries)
		authorized.GET("/vulndb/cve/:cve", vulnDBController.GetVulnDBByCVE)
		authorized.GET("/vulndb/id/:id", vulnDBController.GetVulnDBByID)
		authorized.POST("/vulndb", vulnDBController.CreateVulnDBEntry)
		authorized.PUT("/vulndb/id/:id", vulnDBController.UpdateVulnDBEntry)
		authorized.DELETE("/vulndb/id/:id", vulnDBController.DeleteVulnDBEntry)
		authorized.POST("/vulndb/import", vulnDBController.BatchImportVulnDBEntries)

		// 仪表盘路由
		dashboardController := new(controllers.DashboardController)
		authorized.GET("/dashboard/stats", dashboardController.GetDashboardStats)
		authorized.GET("/dashboard/vuln-trends", dashboardController.GetVulnTrends)
		authorized.GET("/dashboard/severity-distribution", dashboardController.GetSeverityDistribution)
		authorized.GET("/dashboard/asset-vuln-distribution", dashboardController.GetAssetVulnDistribution)
		authorized.GET("/dashboard/priority-vulns", dashboardController.GetPriorityVulns)
		authorized.GET("/dashboard/recent-activities", dashboardController.GetRecentActivities)

		// 系统设置路由 (仅管理员访问)
		settingsRouter := authorized.Group("/settings")
		settingsRouter.Use(middleware.RequireAdmin())
		{
			settingsController := new(controllers.SettingsController)
			settingsRouter.GET("", settingsController.GetSettings)
			settingsRouter.PUT("", settingsController.SaveSettings)
			settingsRouter.POST("/test/jira", settingsController.TestJiraConnection)
			settingsRouter.POST("/test/wechat-login", settingsController.TestWechatLogin)
			settingsRouter.POST("/test/work-wechat", settingsController.TestWorkWechatBot)
			settingsRouter.POST("/test/feishu", settingsController.TestFeishuBot)
			settingsRouter.POST("/test/dingtalk", settingsController.TestDingtalkBot)
			settingsRouter.POST("/test/email", settingsController.TestEmailNotification)
			settingsRouter.POST("/test/ai", settingsController.TestAiService)
			settingsRouter.POST("/test/vulndb", settingsController.TestVulnDBConnection)

			// 添加测试漏洞通知API
			settingsRouter.POST("/test/notification/vulnerability", settingsController.TestVulnerabilityNotification)
		}

		// 自动扫描路由
		scanController := new(controllers.ScanController)
		authorized.GET("/scans", scanController.ListScanTasks)
		authorized.GET("/scans/:id", scanController.GetScanTask)
		authorized.POST("/scans", scanController.CreateScanTask)
		authorized.PUT("/scans/:id", scanController.UpdateScanTask)
		authorized.DELETE("/scans/:id", scanController.DeleteScanTask)
		authorized.POST("/scans/:id/start", scanController.StartScanTask)
		authorized.POST("/scans/:id/cancel", scanController.CancelScanTask)
		authorized.GET("/scans/:id/results", scanController.GetScanResults)
		authorized.POST("/scans/:id/import", scanController.ImportScanResults)

		// AI风险评估路由
		authorized.POST("/ai/risk-assessment", controllers.PerformRiskAssessment)

		// CI/CD集成相关接口
		integrationController := new(controllers.IntegrationController)

		cicdGroup := authorized.Group("/integrations")
		{
			cicdGroup.Use(middleware.JWTAuthMiddleware())
			cicdGroup.Use(middleware.AdminAuthMiddleware())
			cicdGroup.GET("", integrationController.GetIntegrations)
			cicdGroup.POST("", integrationController.CreateIntegration)
			cicdGroup.PUT("/:id", integrationController.UpdateIntegration)
			cicdGroup.PUT("/:id/status", integrationController.UpdateIntegrationStatus)
			cicdGroup.DELETE("/:id", integrationController.DeleteIntegration)
			cicdGroup.POST("/:id/api-key/regenerate", integrationController.RegenerateAPIKey)
			cicdGroup.GET("/:id/history", integrationController.GetIntegrationHistory)
		}

		// 接收CI/CD扫描结果的接口 - 不需要登录，通过API Key认证
		authorized.POST("/webhooks/:type", integrationController.ReceiveScanResult)
	}
}
