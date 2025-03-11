package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// 获取数据库和集合
	db := client.Database("vuln_management")
	usersCollection := db.Collection("users")

	// 查询admin用户信息
	var adminUser bson.M
	err = usersCollection.FindOne(ctx, bson.M{"username": "admin"}).Decode(&adminUser)
	if err != nil {
		log.Fatalf("查询admin用户失败: %v", err)
	}

	// 输出完整的用户信息
	fmt.Println("\n管理员用户详细信息:")
	fmt.Println("==================================")
	for key, value := range adminUser {
		fmt.Printf("%s: %v\n", key, value)
	}
	fmt.Println("==================================")

	// 提示登录信息
	fmt.Println("\n登录信息:")
	fmt.Println("用户名: admin")
	fmt.Println("电子邮件: admin@qq.com")
	fmt.Println("密码: Admin123 (默认密码)")
	fmt.Println("\n请使用此信息尝试登录。")
}
