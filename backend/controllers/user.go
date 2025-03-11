package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"vuln-management/config"
	"vuln-management/middleware"
	"vuln-management/models"
	"vuln-management/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserController 用户相关控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register 注册用户
func (c *UserController) Register(ctx *gin.Context) {
	// 实现注册逻辑
	ctx.JSON(http.StatusOK, gin.H{"message": "注册功能待实现"})
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var credentials models.UserCredentials
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求格式无效", "details": err.Error()})
		return
	}

	// 查找用户
	usersCollection := config.GetCollection(config.UsersCollection)
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := usersCollection.FindOne(dbCtx, bson.M{
		"$or": []bson.M{
			{"email": credentials.Email},
			{"username": credentials.Email}, // 同时支持用户名和邮箱登录
		},
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名/邮箱或密码不正确"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户时出错", "details": err.Error()})
		}
		return
	}

	// 验证密码
	if !user.CheckPassword(credentials.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名/邮箱或密码不正确"})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	_, err = usersCollection.UpdateOne(
		dbCtx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"last_login": now, "updated_at": now}},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新登录时间失败", "details": err.Error()})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败", "details": err.Error()})
		return
	}

	// 返回用户信息和令牌
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// GetCurrentUser 获取当前用户信息
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	// 实现获取当前用户信息逻辑
	ctx.JSON(http.StatusOK, gin.H{"message": "获取当前用户功能待实现"})
}

// UpdateCurrentUser 更新当前用户信息
func (c *UserController) UpdateCurrentUser(ctx *gin.Context) {
	// 实现更新当前用户信息逻辑
	ctx.JSON(http.StatusOK, gin.H{"message": "更新当前用户功能待实现"})
}

// GetAllUsers 获取所有用户列表
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	// 分页参数
	page, _ := ctx.GetQuery("page")
	perPage, _ := ctx.GetQuery("perPage")
	pageNum := 1
	perPageNum := 20

	if page != "" {
		pageNum = parseInt(page, 1)
	}
	if perPage != "" {
		perPageNum = parseInt(perPage, 20)
	}

	// 搜索参数
	search := ctx.Query("search")
	role := ctx.Query("role")
	status := ctx.Query("status")

	usersCollection := config.GetCollection(config.UsersCollection)
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件
	filter := bson.M{}

	// 添加搜索条件
	if search != "" {
		// 创建OR查询，搜索用户名、邮箱或姓名
		filter["$or"] = []bson.M{
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
			{"first_name": bson.M{"$regex": search, "$options": "i"}},
			{"last_name": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// 添加角色筛选
	if role != "" {
		filter["role"] = role
	}

	// 添加状态筛选
	if status != "" {
		filter["status"] = status
	}

	// 查询选项
	opts := options.Find().
		SetSkip(int64((pageNum - 1) * perPageNum)).
		SetLimit(int64(perPageNum)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 执行查询
	cursor, err := usersCollection.Find(dbCtx, filter, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败", "details": err.Error()})
		return
	}
	defer cursor.Close(dbCtx)

	// 统计总数
	total, err := usersCollection.CountDocuments(dbCtx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "统计用户总数失败", "details": err.Error()})
		return
	}

	// 解码结果
	var users []models.User
	if err := cursor.All(dbCtx, &users); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "解析用户数据失败", "details": err.Error()})
		return
	}

	// 转换为响应格式
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": responses,
		"meta": gin.H{
			"total":    total,
			"page":     pageNum,
			"per_page": perPageNum,
			"pages":    (total + int64(perPageNum) - 1) / int64(perPageNum),
		},
	})
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	// 实现创建用户逻辑
	ctx.JSON(http.StatusOK, gin.H{"message": "创建用户功能待实现"})
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 实现更新用户逻辑
	ctx.JSON(http.StatusOK, gin.H{"message": "更新用户功能待实现"})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// 实现删除用户逻辑
	ctx.JSON(http.StatusOK, gin.H{"message": "删除用户功能待实现"})
}

// GetDashboardData 获取仪表盘数据
func (c *UserController) GetDashboardData(ctx *gin.Context) {
	// 调用dashboard.go中的完整实现
	GetDashboardData(ctx)
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证的用户"})
		return
	}

	// 转换为ObjectID
	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID", "details": err.Error()})
		return
	}

	// 查询用户信息
	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败", "details": err.Error()})
		}
		return
	}

	// 返回用户信息（不包含密码）
	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateCurrentUser 更新当前用户信息
func UpdateCurrentUser(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证的用户"})
		return
	}

	// 转换为ObjectID
	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID", "details": err.Error()})
		return
	}

	// 绑定请求数据
	var updateData struct {
		FirstName       *string `json:"firstName"`
		LastName        *string `json:"lastName"`
		Department      *string `json:"department"`
		ProfilePicture  *string `json:"profilePicture"`
		CurrentPassword *string `json:"currentPassword"`
		NewPassword     *string `json:"newPassword"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式", "details": err.Error()})
		return
	}

	// 设置要更新的字段
	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 首先获取当前用户信息（如果需要更改密码）
	var currentUser models.User
	if updateData.CurrentPassword != nil && updateData.NewPassword != nil {
		err = usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&currentUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败", "details": err.Error()})
			return
		}

		// 验证当前密码
		if !currentUser.CheckPassword(*updateData.CurrentPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "当前密码不正确"})
			return
		}
	}

	// 构建更新文档
	updateDoc := bson.M{"updated_at": time.Now()}

	if updateData.FirstName != nil {
		updateDoc["first_name"] = *updateData.FirstName
	}
	if updateData.LastName != nil {
		updateDoc["last_name"] = *updateData.LastName
	}
	if updateData.Department != nil {
		updateDoc["department"] = *updateData.Department
	}
	if updateData.ProfilePicture != nil {
		updateDoc["profile_picture"] = *updateData.ProfilePicture
	}

	// 更新密码（如果提供）
	if updateData.CurrentPassword != nil && updateData.NewPassword != nil {
		// 已在前面验证了当前密码
		if err := currentUser.SetPassword(*updateData.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码哈希失败", "details": err.Error()})
			return
		}
		updateDoc["password_hash"] = currentUser.PasswordHash
	}

	// 执行更新
	result, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateDoc},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败", "details": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}

// GetAllUsers 获取所有用户（管理员功能）
func GetAllUsers(c *gin.Context) {
	// 分页参数
	page, _ := c.GetQuery("page")
	perPage, _ := c.GetQuery("perPage")
	pageNum := 1
	perPageNum := 20

	if page != "" {
		pageNum = parseInt(page, 1)
	}
	if perPage != "" {
		perPageNum = parseInt(perPage, 20)
	}

	// 搜索参数
	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status")

	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建查询条件
	filter := bson.M{}

	// 添加搜索条件
	if search != "" {
		// 创建OR查询，搜索用户名、邮箱或姓名
		filter["$or"] = []bson.M{
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
			{"first_name": bson.M{"$regex": search, "$options": "i"}},
			{"last_name": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// 添加角色筛选
	if role != "" {
		filter["role"] = role
	}

	// 添加状态筛选
	if status != "" {
		filter["status"] = status
	}

	// 查询选项
	opts := options.Find().
		SetSkip(int64((pageNum - 1) * perPageNum)).
		SetLimit(int64(perPageNum)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	// 执行查询
	cursor, err := usersCollection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败", "details": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// 统计总数
	total, err := usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "统计用户总数失败", "details": err.Error()})
		return
	}

	// 解码结果
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析用户数据失败", "details": err.Error()})
		return
	}

	// 转换为响应格式
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"meta": gin.H{
			"total":    total,
			"page":     pageNum,
			"per_page": perPageNum,
			"pages":    (total + int64(perPageNum) - 1) / int64(perPageNum),
		},
	})
}

// CreateUser 创建新用户（管理员功能）
func CreateUser(c *gin.Context) {
	var registration models.UserRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式", "details": err.Error()})
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

	// 设置角色（如果未提供，使用默认值）
	role := registration.Role
	if role == "" {
		role = models.RoleUser
	}

	// 设置状态（如果未提供，默认为激活状态）
	status := registration.Status
	if status == "" {
		status = models.StatusActive
	}

	// 创建新用户
	now := time.Now()
	newUser := models.User{
		ID:         primitive.NewObjectID(),
		Username:   registration.Username,
		Email:      registration.Email,
		FirstName:  registration.FirstName,
		LastName:   registration.LastName,
		Role:       role,
		Status:     status,
		Department: registration.Department,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// 设置密码
	if err := newUser.SetPassword(registration.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码哈希失败", "details": err.Error()})
		return
	}

	// 保存用户
	result, err := usersCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "用户创建成功",
		"user":    newUser.ToResponse(),
		"id":      result.InsertedID,
	})
}

// UpdateUser 更新用户信息（管理员功能）
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID", "details": err.Error()})
		return
	}

	// 绑定请求数据
	var updateData struct {
		FirstName      *string `json:"first_name"`
		LastName       *string `json:"last_name"`
		Department     *string `json:"department"`
		Role           *string `json:"role"`
		ProfilePicture *string `json:"profile_picture"`
		Status         *string `json:"status"`
		NewPassword    *string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式", "details": err.Error()})
		return
	}

	// 构建更新文档
	updateDoc := bson.M{"updated_at": time.Now()}

	if updateData.FirstName != nil {
		updateDoc["first_name"] = *updateData.FirstName
	}
	if updateData.LastName != nil {
		updateDoc["last_name"] = *updateData.LastName
	}
	if updateData.Department != nil {
		updateDoc["department"] = *updateData.Department
	}
	if updateData.Role != nil {
		updateDoc["role"] = *updateData.Role
	}
	if updateData.ProfilePicture != nil {
		updateDoc["profile_picture"] = *updateData.ProfilePicture
	}
	// 处理状态更新
	if updateData.Status != nil {
		// 验证状态值是否有效
		validStatus := *updateData.Status == models.StatusActive ||
			*updateData.Status == models.StatusInactive ||
			*updateData.Status == models.StatusDisabled

		if !validStatus {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态值"})
			return
		}

		updateDoc["status"] = *updateData.Status
	}

	// 更新密码（如果提供）
	if updateData.NewPassword != nil {
		if len(*updateData.NewPassword) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度必须至少为6个字符"})
			return
		}

		var user models.User
		user.SetPassword(*updateData.NewPassword)
		updateDoc["password_hash"] = user.PasswordHash
	}

	// 执行更新
	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": updateDoc},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败", "details": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取更新后的用户信息
	var updatedUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功，但无法获取更新后的信息"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户信息更新成功",
		"user":    updatedUser.ToResponse(),
	})
}

// DeleteUser 删除用户（管理员功能）
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID", "details": err.Error()})
		return
	}

	// 当前登录用户不能删除自己
	currentUserID, _ := c.Get("userID")
	if currentUserID.(string) == id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除当前登录的用户"})
		return
	}

	usersCollection := config.GetCollection(config.UsersCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := usersCollection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败", "details": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已成功删除"})
}

// parseInt 将字符串转换为整数，失败时返回默认值
func parseInt(str string, defaultValue int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}
