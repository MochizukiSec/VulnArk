package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// MongoClient 是MongoDB客户端实例
	MongoClient *mongo.Client
	// Database 是主数据库
	Database *mongo.Database
)

// 集合名称常量
const (
	UsersCollection           = "users"
	VulnerabilitiesCollection = "vulnerabilities"
	ReportsCollection         = "reports"
)

// InitDB 初始化MongoDB连接
func InitDB() error {
	// 获取MongoDB连接字符串
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// 设置MongoDB客户端选项
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second)

	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 检查连接
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	log.Println("成功连接到MongoDB!")

	// 获取数据库名称
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "vuln_management"
	}

	// 设置全局变量
	MongoClient = client
	Database = client.Database(dbName)

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Printf("关闭MongoDB连接出错: %v", err)
		} else {
			log.Println("MongoDB连接已关闭")
		}
	}
}

// GetCollection 获取指定名称的集合
func GetCollection(name string) *mongo.Collection {
	return Database.Collection(name)
}
