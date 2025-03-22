package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CORSMiddleware 跨域资源共享中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取允许的源
		allowedOrigins := viper.GetStringSlice("security.cors.allowed_origins")
		allowedOriginsStr := strings.Join(allowedOrigins, ", ")
		if allowedOriginsStr == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origin := c.Request.Header.Get("Origin")
			if origin != "" {
				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == origin || allowedOrigin == "*" {
						c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		}

		// 获取允许的方法
		allowedMethods := viper.GetStringSlice("security.cors.allowed_methods")
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

		// 获取允许的头
		allowedHeaders := viper.GetStringSlice("security.cors.allowed_headers")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

		// 设置允许凭证
		if viper.GetBool("security.cors.allow_credentials") {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
