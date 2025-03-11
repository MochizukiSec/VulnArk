package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"vuln-management/config"
	"vuln-management/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer config.CloseDB()

	// 检查管理员用户是否已存在
	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": "admin@qq.com"}).Decode(&existingUser)

	if err == nil {
		fmt.Println("管理员账户已存在，无需创建")
		return
	} else if err != mongo.ErrNoDocuments {
		log.Fatalf("查询数据库时出错: %v", err)
	}

	// 创建管理员用户
	adminUser := models.User{
		ID:         primitive.NewObjectID(),
		Username:   "admin",
		Email:      "admin@qq.com",
		FirstName:  "系统",
		LastName:   "管理员",
		Role:       models.RoleAdmin,
		Department: "系统运维",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 设置密码
	if err := adminUser.SetPassword("Admin123"); err != nil {
		log.Fatalf("设置密码时出错: %v", err)
	}

	// 保存到数据库
	result, err := usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		log.Fatalf("创建管理员账户时出错: %v", err)
	}

	fmt.Printf("成功创建管理员账户，ID: %s\n", result.InsertedID)
}
