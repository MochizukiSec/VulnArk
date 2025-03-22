package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/vulnark/vulnark/models"
)

// 自定义JWT声明结构
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("JWT认证失败: 请求中缺少Authorization头部")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证token",
			})
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			log.Printf("JWT认证失败: Authorization头格式错误: %s", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证头格式错误，应为 'Bearer {token}'",
			})
			c.Abort()
			return
		}

		// 解析token
		token := parts[1]
		if len(token) > 10 {
			log.Printf("尝试解析JWT Token: %s...", token[:10])
		} else {
			log.Printf("尝试解析JWT Token: %s (令牌过短)", token)
		}

		claims, err := ParseToken(token)
		if err != nil {
			log.Printf("JWT Token解析失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 确保UserID有效
		if claims.UserID == 0 {
			log.Printf("JWT认证失败: Token中的UserID为0")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的用户ID",
			})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		log.Printf("JWT认证成功: UserID=%d, Role=%s", claims.UserID, claims.Role)

		// 同时设置userID和user_id以确保兼容性
		c.Set("userID", claims.UserID)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// GenerateToken 生成JWT token
func GenerateToken(user *models.User) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(time.Duration(viper.GetInt("auth.token_expire")) * time.Hour)

	// 创建JWT声明
	claims := &Claims{
		UserID: user.ID,
		Role:   string(user.Role),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "vulnark",
		},
	}

	// 使用声明创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥对token进行签名
	tokenString, err := token.SignedString([]byte(viper.GetString("auth.jwt_secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("意外的签名方法")
		}
		return []byte(viper.GetString("auth.jwt_secret")), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// RequireAdmin 仅允许管理员访问
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != string(models.RoleAdmin) {
			log.Printf("权限检查失败: 需要管理员权限，当前角色=%v", role)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireManager 仅允许管理者及以上访问
func RequireManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != string(models.RoleAdmin) && role != string(models.RoleManager)) {
			log.Printf("权限检查失败: 需要管理者或管理员权限，当前角色=%v", role)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "需要管理者或管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
