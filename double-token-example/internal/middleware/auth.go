package middleware

import (
	"double-token-example/internal/config"
	"double-token-example/internal/db"
	"double-token-example/internal/model"
	"double-token-example/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AccessAuth JWT认证中间件
func AccessAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "未提供token"})
			c.Abort()
			return
		}
		// 解析token
		claims, typ, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "解析token：无效的token"})
			c.Abort()
			return
		}
		// 校验token是否在黑名单
		exists, _ := db.GetRedisDB().Exists(c.Request.Context(), "jwt_blacklist:"+claims.ID).Result()
		if exists == 1 {
			c.JSON(401, gin.H{"error": "redis:无效的token"})
			c.Abort()
			return
		}
		// 校验token类型
		if typ != config.Cfg.JWT.AccessTokenType {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}
		// 将userID存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("jti", claims.ID)
		c.Set("exp", claims.ExpiresAt.Time)
		c.Next()
	}
}

// RefreshAuth JWT刷新中间件
func RefreshAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "未提供token"})
			c.Abort()
			return
		}

		// 解析token
		claims, typ, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 校验token类型
		if typ != config.Cfg.JWT.RefreshTokenType {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 校验token是否存在
		var refreshToken *model.RefreshToken
		err = db.GetMySQL().Where("jti = ?", claims.ID).First(&refreshToken).Error
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 将userID存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("jti", claims.ID)
		c.Next()
	}
}
