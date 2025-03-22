package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/routes"
	"github.com/vulnark/vulnark/utils"
)

func init() {
	// 设置配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	// 创建必要的目录
	dirs := []string{
		viper.GetString("log.file_path"),
		viper.GetString("upload.location"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("无法创建目录 %s: %v", dir, err)
		}
	}

	// 设置运行模式
	mode := viper.GetString("server.mode")
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func main() {
	// 初始化数据库连接
	utils.InitDB()
	defer utils.CloseDB()

	// 自动迁移数据库模型
	autoMigrateModels()

	// 创建默认管理员账户
	createDefaultAdmin()

	// 创建Gin路由
	router := gin.Default()

	// 设置路由
	routes.SetupRouter(router)

	port := viper.GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("启动服务器，监听端口: %d\n", port)

	if err := router.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// autoMigrateModels 自动迁移所有数据库模型
func autoMigrateModels() {
	if utils.DBType == "mysql" && utils.DB != nil {
		log.Println("开始迁移数据库模型...")

		// 迁移所有模型
		utils.DB.AutoMigrate(
			&models.User{},
			&models.Vulnerability{},
			&models.Asset{},
			&models.Knowledge{},
			&models.VulnDB{},
			&models.Settings{},
			&models.VulnerabilityAssignment{},
			&models.VulnerabilityAssignmentHistory{},
			&models.ScanTask{},
			&models.ScanResult{},
			&models.CIIntegration{},
			&models.IntegrationHistory{},
		)

		// 使用正确的方式创建关联关系
		// 注意: GORM v2不再支持Related方法，改用关联表来表示多对多关系
		log.Println("数据库迁移完成")
	}
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin() {
	if utils.DBType == "mysql" && utils.DB != nil {
		// 检查是否已存在管理员账户
		var count int64
		utils.DB.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&count)

		// 打印当前管理员账户数量
		log.Printf("当前系统中的管理员账户数量: %d", count)

		// 如果管理员数量为0，创建默认管理员
		if count == 0 {
			// 创建默认管理员账户
			admin := models.User{
				Username:  "admin",
				Password:  "admin123", // 这里会由BeforeCreate钩子自动加密
				Email:     "admin@vulnark.com",
				RealName:  "系统管理员",
				Role:      models.RoleAdmin,
				Active:    true,
				LastLogin: time.Now(), // 设置为当前时间，避免零日期问题
			}

			if err := utils.DB.Create(&admin).Error; err != nil {
				log.Printf("创建默认管理员失败: %v", err)
			} else {
				log.Println("已创建默认管理员账户，用户名: admin, 密码: admin123")
			}
		}
	}
}
