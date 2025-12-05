// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handlers

import (
	"github.com/gin-gonic/gin"
	"kk-nav/internal/config"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthHandler 认证处理器
type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:  db,
		cfg: cfg,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.Error(c, 400, "Email already exists")
		return
	}

	// 检查用户名是否已存在
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.Error(c, 400, "Username already exists")
		return
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password")
		return
	}

	// 创建用户
	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         "user",
		Active:       true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role, h.cfg.JWT.ExpireHours)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 查找用户（使用用户名）
	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Unauthorized(c, "Invalid username or password")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 检查用户是否激活
	if !user.Active {
		utils.Unauthorized(c, "User account is inactive")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role, h.cfg.JWT.ExpireHours)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Me 获取当前用户信息
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"role":     user.Role,
		"active":   user.Active,
	})
}

// Logout 用户登出（客户端删除Token即可）
func (h *AuthHandler) Logout(c *gin.Context) {
	utils.SuccessWithMessage(c, "Logged out successfully", nil)
}
