package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"github.com/vulnark/vulnark/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB      *gorm.DB
	MongoDB *mongo.Database
	DBType  string
)

// InitDB 初始化数据库连接
func InitDB() {
	DBType = viper.GetString("database.type")

	switch DBType {
	case "mysql":
		initMysql()
	case "mongodb":
		initMongoDB()
	default:
		initMysql() // 默认使用MySQL
	}
}

// initMysql 初始化MySQL数据库连接
func initMysql() {
	// 优先从环境变量读取，如果没有则从配置文件读取
	host := getConfigString("database.host", "DB_HOST")
	port := getConfigInt("database.port", "DB_PORT")
	username := getConfigString("database.username", "DB_USER")
	password := getConfigString("database.password", "DB_PASSWORD")
	dbname := getConfigString("database.database", "DB_NAME")
	charset := viper.GetString("database.charset")
	if charset == "" {
		charset = "utf8mb4"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, host, port, dbname, charset)

	log.Printf("尝试连接MySQL数据库: %s@%s:%d/%s", username, host, port, dbname)

	var err error
	var retryCount = 5
	var retryInterval = 5 * time.Second

	// 连接重试机制
	for i := 0; i < retryCount; i++ {
		DB, err = gorm.Open("mysql", dsn)
		if err == nil {
			// 测试连接
			if err = DB.DB().Ping(); err == nil {
				break
			}
			DB.Close()
		}
		
		log.Printf("数据库连接失败 (尝试 %d/%d): %v", i+1, retryCount, err)
		if i < retryCount-1 {
			log.Printf("等待 %v 后重试...", retryInterval)
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		log.Fatalf("无法连接到MySQL数据库，已重试%d次: %v", retryCount, err)
	}

	// 设置连接池配置
	maxIdleConns := viper.GetInt("database.max_idle_conns")
	if maxIdleConns == 0 {
		maxIdleConns = 10
	}
	maxOpenConns := viper.GetInt("database.max_open_conns")
	if maxOpenConns == 0 {
		maxOpenConns = 100
	}
	
	DB.DB().SetMaxIdleConns(maxIdleConns)
	DB.DB().SetMaxOpenConns(maxOpenConns)
	DB.DB().SetConnMaxLifetime(time.Hour)

	// 启用日志
	DB.LogMode(viper.GetString("server.mode") == "development")

	log.Printf("成功连接到MySQL数据库 %s@%s:%d/%s", username, host, port, dbname)
}

// initMongoDB 初始化MongoDB数据库连接
func initMongoDB() {
	uri := viper.GetString("database.mongodb.uri")
	dbname := viper.GetString("database.mongodb.database")
	timeout := viper.GetInt("database.mongodb.timeout")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("无法连接到MongoDB: %v", err)
	}

	// 测试连接
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("无法连接到MongoDB: %v", err)
	}

	MongoDB = client.Database(dbname)
	log.Println("成功连接到MongoDB")
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("关闭MySQL连接时出错: %v", err)
		}
	}
}

// AutoMigrate 自动迁移数据库模型
func AutoMigrate() {
	if DBType == "mysql" && DB != nil {
		// 在这里添加所有需要迁移的模型
		// DB.AutoMigrate(&models.User{}, &models.Vulnerability{}, ...)
		log.Println("开始执行数据库模型自动迁移...")

		// 迁移Settings模型
		if err := DB.AutoMigrate(&models.Settings{}).Error; err != nil {
			log.Printf("迁移Settings模型失败: %v", err)
		} else {
			log.Println("Settings模型迁移成功")
		}

		// 检查Settings表结构
		var tableExists int
		if err := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'settings'").Row().Scan(&tableExists); err != nil {
			log.Printf("检查Settings表存在性失败: %v", err)
		} else if tableExists > 0 {
			// 检查JSON列是否正确
			rows, err := DB.Raw("SHOW COLUMNS FROM settings").Rows()
			if err != nil {
				log.Printf("无法获取settings表结构: %v", err)
			} else {
				defer rows.Close()
				var field, fieldType, null, key, defaultValue, extra string
				columnMap := make(map[string]string)

				for rows.Next() {
					err := rows.Scan(&field, &fieldType, &null, &key, &defaultValue, &extra)
					if err != nil {
						log.Printf("无法扫描列信息: %v", err)
						continue
					}
					columnMap[field] = fieldType
					log.Printf("Settings表列: %s, 类型: %s", field, fieldType)
				}

				// 确保JSON列的类型正确
				for _, column := range []string{"integrations", "notifications", "ai"} {
					if columnType, exists := columnMap[column]; exists {
						if columnType != "json" && columnType != "JSON" && !containsSubstring(columnType, "json") {
							log.Printf("警告: settings表的%s列类型为%s，不是JSON类型。尝试修复...", column, columnType)
							if err := DB.Exec(fmt.Sprintf("ALTER TABLE settings MODIFY %s JSON", column)).Error; err != nil {
								log.Printf("修改%s列类型为JSON失败: %v", column, err)
							} else {
								log.Printf("%s列类型已修改为JSON", column)
							}
						}
					}
				}
			}
		}

		log.Println("数据库模型自动迁移完成")
	}
}

// containsSubstring 检查字符串是否包含子字符串（不区分大小写）
func containsSubstring(str, substr string) bool {
	str = strings.ToLower(str)
	substr = strings.ToLower(substr)
	return strings.Contains(str, substr)
}

// getConfigString 优先从环境变量读取字符串配置，如果没有则从配置文件读取
func getConfigString(configKey, envKey string) string {
	// 先尝试从环境变量读取
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	// 再从配置文件读取
	return viper.GetString(configKey)
}

// getConfigInt 优先从环境变量读取整数配置，如果没有则从配置文件读取
func getConfigInt(configKey, envKey string) int {
	// 先尝试从环境变量读取
	if envValue := os.Getenv(envKey); envValue != "" {
		if intValue, err := strconv.Atoi(envValue); err == nil {
			return intValue
		}
	}
	// 再从配置文件读取
	return viper.GetInt(configKey)
}
