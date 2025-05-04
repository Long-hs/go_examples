package middleware

import (
	"double-token-example/pkg/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "未提供token"})
			c.Abort()
			return
		}

		// 解析token
		userID, err := utils.GetUserIDFromToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 将userID存入上下文
		c.Set("user_id", userID)
		c.Next()
	}
}
