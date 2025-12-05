// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handlers

import (
	"kk-nav/internal/models"
	"kk-nav/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingsHandler 设置处理器（公开API）
type SettingsHandler struct {
	db *gorm.DB
}

// NewSettingsHandler 创建设置处理器
func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// GetPublicSettings 获取公开设置（不需要认证）
func (h *SettingsHandler) GetPublicSettings(c *gin.Context) {
	// 只返回公开的设置项
	publicKeys := []string{"site_name", "site_description", "primary_color", "theme"}

	settingsMap := make(map[string]string)
	for _, key := range publicKeys {
		var setting models.Setting
		if err := h.db.Where("key = ?", key).First(&setting).Error; err != nil {
			// 如果数据库中没有，使用默认值
			if defaultValue, ok := models.DefaultSettings[key]; ok {
				settingsMap[key] = defaultValue
			}
		} else {
			settingsMap[key] = setting.Value
		}
	}

	utils.Success(c, gin.H{
		"settings": settingsMap,
	})
}
