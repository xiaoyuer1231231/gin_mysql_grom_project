package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/config"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/util"
)

func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization 头部
		authHeader := c.GetHeader("Authorization")
		fmt.Println("tokenValue", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		// 2. 提取 Bearer Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Bearer", tokenString)

		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Bearer token required",
				"code":  "INVALID_TOKEN_FORMAT",
			})
			c.Abort()
			return
		}

		// 3. 验证 Token
		claims, err := util.ValidateToken(tokenString, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// 4. 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_claims", claims) // 存储完整的 claims

		// 5. 继续处理请求
		c.Next()
	}
}
