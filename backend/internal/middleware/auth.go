// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/database"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// AuthMiddleware JWT认证中间件（支持 JWT Token 和 API Token）
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		// 提取Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// 判断是 API Token 还是 JWT Token
		if strings.HasPrefix(token, "kk_") {
			// API Token 认证
			var apiToken models.APIToken
			if err := database.DB.Preload("User").Where("token = ?", token).First(&apiToken).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					utils.Unauthorized(c, "Invalid API token")
					c.Abort()
					return
				}
				utils.InternalServerError(c, "Database error")
				c.Abort()
				return
			}

			// 检查 Token 是否有效
			if !apiToken.IsValid() {
				utils.Unauthorized(c, "Token is inactive or expired")
				c.Abort()
				return
			}

			// 更新最后使用时间
			now := time.Now()
			apiToken.LastUsedAt = &now
			database.DB.Save(&apiToken)

			// 将用户信息存储到上下文
			c.Set("user_id", apiToken.UserID)
			c.Set("username", apiToken.User.Username)
			c.Set("email", apiToken.User.Email)
			c.Set("role", apiToken.User.Role)
		} else {
			// JWT Token 认证
		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		}

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

