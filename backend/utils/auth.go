package utils

import (
	"vuln-management/models"

	"github.com/gin-gonic/gin"
)

// GetUserFromContext 从Gin上下文中获取当前用户信息
func GetUserFromContext(c *gin.Context) *models.User {
	// 从上下文中获取用户信息（由JWT中间件设置）
	user, exists := c.Get("user")
	if !exists {
		return nil
	}

	// 类型转换
	currentUser, ok := user.(*models.User)
	if !ok {
		return nil
	}

	return currentUser
}
