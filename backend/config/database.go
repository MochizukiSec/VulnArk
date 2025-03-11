package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
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
	AssetsCollection          = "assets"
	VulnDatabaseCollection    = "vulndatabase" // 漏洞库集合
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

// InitializeDatabase 在系统首次启动时初始化必要的数据
func InitializeDatabase() error {
	// 创建管理员账号
	err := createAdminUser()
	if err != nil {
		log.Printf("初始化管理员账号时出错: %v", err)
		return err
	}

	return nil
}

// createAdminUser 创建默认管理员账号
func createAdminUser() error {
	// 检查管理员账号是否已存在
	usersCollection := GetCollection(UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 查询是否已有管理员账号
	count, err := usersCollection.CountDocuments(ctx, bson.M{"email": "admin@qq.com"})
	if err != nil {
		return err
	}

	// 已存在管理员账号则直接返回
	if count > 0 {
		log.Println("管理员账号已存在，跳过创建")
		return nil
	}

	// 创建管理员用户
	adminUser := bson.M{
		"username":      "admin",
		"email":         "admin@qq.com",
		"password_hash": "", // 将在下面设置密码哈希
		"first_name":    "系统",
		"last_name":     "管理员",
		"role":          "admin",
		"department":    "安全部门",
		"status":        "active",
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}

	// 生成密码哈希
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("Admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	adminUser["password_hash"] = string(passwordHash)

	// 插入管理员用户
	result, err := usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		return err
	}

	log.Printf("成功创建管理员账号，ID: %v", result.InsertedID)
	return nil
}
