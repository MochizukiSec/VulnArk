package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vulnark/vulnark/middleware"
	"github.com/vulnark/vulnark/models"
	"github.com/vulnark/vulnark/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserController 用户控制器
type UserController struct{}

// Login 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var loginForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 打印登录尝试信息
	log.Printf("用户尝试登录: username=%s, password=%s", loginForm.Username, loginForm.Password)

	var user models.User
	if utils.DBType == "mysql" {
		if err := utils.DB.Where("username = ?", loginForm.Username).First(&user).Error; err != nil {
			log.Printf("用户查询失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
			})
			return
		}
	}

	log.Printf("找到用户: ID=%d, username=%s, role=%s, 存储的密码哈希=%s",
		user.ID, user.Username, user.Role, user.Password)

	// 特殊处理测试用户 - 为了方便测试
	if loginForm.Username == "testadmin999" && loginForm.Password == "testpass123" {
		log.Printf("开发测试账号登录 - 绕过密码验证")

		// 更新最后登录时间
		user.LastLogin = time.Now()
		utils.DB.Save(&user)

		// 生成Token
		token, err := middleware.GenerateToken(&user)
		if err != nil {
			log.Printf("生成Token失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成Token失败",
				"error":   err.Error(),
			})
			return
		}

		log.Printf("测试账号登录成功，生成token: %s", token)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登录成功",
			"data": gin.H{
				"token":     token,
				"user_id":   user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"real_name": user.RealName,
				"role":      user.Role,
				"avatar":    user.Avatar,
			},
		})
		return
	}

	// 验证密码
	if !user.CheckPassword(loginForm.Password) {
		log.Printf("密码验证失败")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 更新最后登录时间
	user.LastLogin = time.Now()
	utils.DB.Save(&user)

	// 生成Token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成Token失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("用户 %s 登录成功", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":     token,
			"user_id":   user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"real_name": user.RealName,
			"role":      user.Role,
			"avatar":    user.Avatar,
		},
	})
}

// LoginV2 改进的用户登录方法
func (uc *UserController) LoginV2(c *gin.Context) {
	var loginForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 打印登录尝试信息（不显示密码）
	log.Printf("用户尝试登录: username=%s", loginForm.Username)

	var user models.User
	if utils.DBType == "mysql" {
		// 确保只查询未删除的用户
		if err := utils.DB.Where("username = ? AND deleted_at IS NULL", loginForm.Username).First(&user).Error; err != nil {
			log.Printf("用户查询失败: %v", err)
			// 延迟响应以防止时序攻击
			time.Sleep(time.Duration(300+time.Now().UnixNano()%200) * time.Millisecond)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
			})
			return
		}
	}

	log.Printf("找到用户: ID=%d, username=%s, role=%s, 密码哈希长度=%d",
		user.ID, user.Username, user.Role, len(user.Password))

	// 特殊处理测试用户 - 为了方便测试
	if loginForm.Username == "testadmin999" && loginForm.Password == "testpass123" {
		log.Printf("开发测试账号登录 - 绕过密码验证")

		// 更新最后登录时间（仅更新指定字段，避免密码再次加密）
		now := time.Now()
		utils.DB.Model(&user).UpdateColumns(map[string]interface{}{
			"last_login": now,
			"updated_at": now,
		})

		// 生成Token
		token, err := middleware.GenerateToken(&user)
		if err != nil {
			log.Printf("生成Token失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "生成Token失败",
				"error":   err.Error(),
			})
			return
		}

		log.Printf("测试账号登录成功")
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登录成功",
			"data": gin.H{
				"token":     token,
				"user_id":   user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"real_name": user.RealName,
				"role":      user.Role,
				"avatar":    user.Avatar,
			},
		})
		return
	}

	// 验证密码
	isValid := user.CheckPassword(loginForm.Password)
	if !isValid {
		log.Printf("密码验证失败")
		// 延迟响应以防止时序攻击
		time.Sleep(time.Duration(300+time.Now().UnixNano()%200) * time.Millisecond)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 更新最后登录时间（仅更新指定字段，避免密码再次加密）
	now := time.Now()
	utils.DB.Model(&user).UpdateColumns(map[string]interface{}{
		"last_login": now,
		"updated_at": now,
	})

	// 生成Token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		log.Printf("生成Token失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成Token失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("用户 %s 登录成功", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":     token,
			"user_id":   user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"real_name": user.RealName,
			"role":      user.Role,
			"avatar":    user.Avatar,
		},
	})
}

// GetUserInfo 获取用户信息
func (uc *UserController) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	// 隐藏密码
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    user,
	})
}

// UpdateUser 更新用户信息
func (uc *UserController) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	var updateForm struct {
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&updateForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	// 更新字段
	if updateForm.Email != "" {
		user.Email = updateForm.Email
	}
	if updateForm.RealName != "" {
		user.RealName = updateForm.RealName
	}
	if updateForm.Phone != "" {
		user.Phone = updateForm.Phone
	}
	if updateForm.Avatar != "" {
		user.Avatar = updateForm.Avatar
	}
	if updateForm.Password != "" {
		user.Password = updateForm.Password
	}

	user.UpdatedAt = time.Now()

	if err := utils.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户失败",
			"error":   err.Error(),
		})
		return
	}

	// 隐藏密码
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    user,
	})
}

// UpdateUserV2 改进版用户信息更新
func (uc *UserController) UpdateUserV2(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	var updateForm struct {
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&updateForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := utils.DB.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	log.Printf("更新用户信息: ID=%d, username=%s", user.ID, user.Username)

	// 使用map来更新字段，避免零值覆盖和密码加密问题
	updateMap := make(map[string]interface{})
	updateMap["updated_at"] = time.Now()

	// 只更新提供的非空字段
	if updateForm.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var count int64
		utils.DB.Model(&models.User{}).Where("email = ? AND id != ? AND deleted_at IS NULL",
			updateForm.Email, userID).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "邮箱已被其他用户使用",
			})
			return
		}
		updateMap["email"] = updateForm.Email
	}
	if updateForm.RealName != "" {
		updateMap["real_name"] = updateForm.RealName
	}
	if updateForm.Phone != "" {
		updateMap["phone"] = updateForm.Phone
	}
	if updateForm.Avatar != "" {
		updateMap["avatar"] = updateForm.Avatar
	}

	// 密码单独处理
	if updateForm.Password != "" {
		// 使用bcrypt直接生成密码哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateForm.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("密码加密失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "密码处理失败",
				"error":   err.Error(),
			})
			return
		}
		updateMap["password"] = string(hashedPassword)
		log.Printf("用户密码已更新: ID=%d", user.ID)
	}

	// 只有当有字段需要更新时才执行更新
	if len(updateMap) > 1 { // 大于1是因为至少有updated_at字段
		if err := utils.DB.Model(&user).Updates(updateMap).Error; err != nil {
			log.Printf("更新用户失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新用户失败",
				"error":   err.Error(),
			})
			return
		}

		// 重新查询用户以获取更新后的信息
		utils.DB.Where("id = ?", userID).First(&user)
		log.Printf("用户信息更新成功: ID=%d", user.ID)
	}

	// 隐藏密码
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data":    user,
	})
}

// ListUsers 获取用户列表
func (uc *UserController) ListUsers(c *gin.Context) {
	var users []models.User
	if err := utils.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 隐藏所有用户密码
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    users,
	})
}

// CreateUser 管理员创建用户（替代注册功能）
func (uc *UserController) CreateUser(c *gin.Context) {
	// 验证当前用户是否为管理员
	role, exists := c.Get("role")
	if !exists || role != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "需要管理员权限",
		})
		return
	}

	var createForm struct {
		Username string `json:"username" binding:"required,min=4,max=20"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email" binding:"required,email"`
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		Active   bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&createForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var count int64
	utils.DB.Model(&models.User{}).Where("username = ? AND deleted_at IS NULL", createForm.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	utils.DB.Model(&models.User{}).Where("email = ? AND deleted_at IS NULL", createForm.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "邮箱已存在",
		})
		return
	}

	// 设置角色，默认为浏览者
	userRole := models.RoleViewer
	if createForm.Role != "" {
		userRole = models.Role(createForm.Role)
	}

	// 设置激活状态，默认为激活
	active := true
	if createForm.Active == false {
		active = false
	}

	// 创建用户
	now := time.Now()
	user := models.User{
		Username:  createForm.Username,
		Password:  createForm.Password, // 密码将在BeforeCreate钩子中加密
		Email:     createForm.Email,
		RealName:  createForm.RealName,
		Phone:     createForm.Phone,
		Role:      userRole,
		Active:    active,
		LastLogin: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	log.Printf("管理员创建新用户: username=%s, email=%s, role=%s", user.Username, user.Email, user.Role)

	if err := utils.DB.Create(&user).Error; err != nil {
		log.Printf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建用户失败",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("用户创建成功: ID=%d, username=%s", user.ID, user.Username)

	// 隐藏密码
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建成功",
		"data":    user,
	})
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	if err := utils.DB.Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ChangeUserRole 更改用户角色
func (uc *UserController) ChangeUserRole(c *gin.Context) {
	userID := c.Param("id")

	var roleForm struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&roleForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	user.Role = models.Role(roleForm.Role)
	user.UpdatedAt = time.Now()

	if err := utils.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更改用户角色失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更改角色成功",
		"data": gin.H{
			"user_id": user.ID,
			"role":    user.Role,
		},
	})
}
