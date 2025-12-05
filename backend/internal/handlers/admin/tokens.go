// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package admin

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// TokensHandler API Token 处理器
type TokensHandler struct {
	db *gorm.DB
}

// NewTokensHandler 创建 Token 处理器
func NewTokensHandler(db *gorm.DB) *TokensHandler {
	return &TokensHandler{db: db}
}

// Index Token 列表
func (h *TokensHandler) Index(c *gin.Context) {
	var tokens []models.APIToken
	query := h.db.Preload("User")

	// 搜索
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ? OR token ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 用户筛选
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 激活状态筛选
	if active := c.Query("active"); active != "" {
		if active == "true" {
			query = query.Where("active = ?", true)
		} else if active == "false" {
			query = query.Where("active = ?", false)
		}
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Model(&models.APIToken{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tokens)

	utils.Success(c, gin.H{
		"tokens": tokens,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// Show Token 详情
func (h *TokensHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid token ID")
		return
	}

	var token models.APIToken
	if err := h.db.Preload("User").First(&token, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Token not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
	})
}

// Create 创建 Token
func (h *TokensHandler) Create(c *gin.Context) {
	var req struct {
		Name      string     `json:"name" binding:"required"`
		UserID    uint       `json:"user_id" binding:"required"`
		ExpiresAt *time.Time `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查用户是否存在
	var user models.User
	if err := h.db.First(&user, req.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "User not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 创建 Token
	token := models.APIToken{
		Name:      req.Name,
		UserID:    req.UserID,
		ExpiresAt: req.ExpiresAt,
		Active:    true,
	}

	if err := h.db.Create(&token).Error; err != nil {
		utils.InternalServerError(c, "Failed to create token")
		return
	}

	// 重新加载以获取生成的 token
	h.db.Preload("User").First(&token, token.ID)

	utils.SuccessWithMessage(c, "Token created successfully", token)
}

// Update 更新 Token
func (h *TokensHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid token ID")
		return
	}

	var req struct {
		Name      string     `json:"name"`
		Active    *bool      `json:"active"`
		ExpiresAt *time.Time `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var token models.APIToken
	if err := h.db.First(&token, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Token not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 更新字段
	if req.Name != "" {
		token.Name = req.Name
	}
	if req.Active != nil {
		token.Active = *req.Active
	}
	if req.ExpiresAt != nil {
		token.ExpiresAt = req.ExpiresAt
	}

	if err := h.db.Save(&token).Error; err != nil {
		utils.InternalServerError(c, "Failed to update token")
		return
	}

	h.db.Preload("User").First(&token, token.ID)

	utils.SuccessWithMessage(c, "Token updated successfully", token)
}

// Delete 删除 Token
func (h *TokensHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid token ID")
		return
	}

	if err := h.db.Delete(&models.APIToken{}, id).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete token")
		return
	}

	utils.SuccessWithMessage(c, "Token deleted successfully", nil)
}

