package controllers

import (
	"context"
	"net/http"
	"time"

	"vuln-management/config"
	"vuln-management/middleware"
	"vuln-management/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Register 处理用户注册
func Register(c *gin.Context) {
	var registration models.UserRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式无效", "details": err.Error()})
		return
	}

	// 检查邮箱是否已存在
	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": registration.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该邮箱已被注册"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户时出错", "details": err.Error()})
		return
	}

	// 检查用户名是否已存在
	err = usersCollection.FindOne(ctx, bson.M{"username": registration.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该用户名已被使用"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户时出错", "details": err.Error()})
		return
	}

	// 创建新用户
	now := time.Now()
	newUser := models.User{
		ID:         primitive.NewObjectID(),
		Username:   registration.Username,
		Email:      registration.Email,
		FirstName:  registration.FirstName,
		LastName:   registration.LastName,
		Department: registration.Department,
		Role:       models.RoleUser, // 新用户默认为普通用户角色
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 设置密码哈希
	if err := newUser.SetPassword(registration.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码哈希失败", "details": err.Error()})
		return
	}

	// 保存用户到数据库
	_, err = usersCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败", "details": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(newUser.ID, newUser.Username, newUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败", "details": err.Error()})
		return
	}

	// 返回用户信息和令牌
	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user":    newUser.ToResponse(),
		"token":   token,
	})
}

// Login 处理用户登录
func Login(c *gin.Context) {
	var credentials models.UserCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式无效", "details": err.Error()})
		return
	}

	// 查找用户
	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码不正确"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户时出错", "details": err.Error()})
		}
		return
	}

	// 验证密码
	if !user.CheckPassword(credentials.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码不正确"})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	_, err = usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"last_login": now, "updated_at": now}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新登录时间失败", "details": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败", "details": err.Error()})
		return
	}

	// 返回用户信息和令牌
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    user.ToResponse(),
		"token":   token,
	})
}
