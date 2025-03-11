package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"vuln-management/models"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到.env文件，使用环境变量")
	}

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 设置MongoDB客户端
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("无法连接到MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB连接测试失败: %v", err)
	}
	log.Println("成功连接到MongoDB")

	// 获取数据库和集合
	db := client.Database("vuln_management")
	usersCollection := db.Collection("users")

	// 检查管理员用户是否存在
	var existingAdmin models.User
	err = usersCollection.FindOne(ctx, bson.M{"email": "admin@qq.com"}).Decode(&existingAdmin)
	if err == nil {
		fmt.Println("管理员用户已经存在:", existingAdmin.Username, existingAdmin.Email)
		return
	} else if err != mongo.ErrNoDocuments {
		log.Fatalf("查询管理员用户时出错: %v", err)
	}

	// 创建管理员用户
	now := time.Now()
	adminUser := models.User{
		ID:         primitive.NewObjectID(),
		Username:   "admin",
		Email:      "admin@qq.com",
		FirstName:  "系统",
		LastName:   "管理员",
		Department: "IT部门",
		Role:       "admin",
		Status:     "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 设置密码
	if err := adminUser.SetPassword("Admin123"); err != nil {
		log.Fatalf("设置管理员密码失败: %v", err)
	}

	// 保存管理员用户
	_, err = usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		log.Fatalf("创建管理员用户失败: %v", err)
	}

	fmt.Println("已成功创建管理员用户 admin@qq.com，密码为 Admin123")
}
