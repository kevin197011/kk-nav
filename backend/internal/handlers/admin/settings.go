// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package admin

import (
	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// SettingsHandler 管理后台设置处理器
type SettingsHandler struct {
	db *gorm.DB
}

// NewSettingsHandler 创建设置处理器
func NewSettingsHandler(db *gorm.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// Index 获取所有设置
func (h *SettingsHandler) Index(c *gin.Context) {
	var settings []models.Setting
	h.db.Order("key").Find(&settings)

	// 转换为map格式
	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	utils.Success(c, gin.H{
		"settings": settingsMap,
		"list":     settings,
	})
}

// Update 更新设置
func (h *SettingsHandler) Update(c *gin.Context) {
	var req struct {
		Settings map[string]string `json:"settings" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 更新每个设置
	for key, value := range req.Settings {
		var setting models.Setting
		if err := h.db.Where("key = ?", key).First(&setting).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 创建新设置
				setting = models.Setting{
					Key:   key,
					Value: value,
				}
				h.db.Create(&setting)
			} else {
				utils.InternalServerError(c, "Database error")
				return
			}
		} else {
			// 更新现有设置
			setting.Value = value
			h.db.Save(&setting)
		}
	}

	utils.SuccessWithMessage(c, "Settings updated successfully", nil)
}

// GetSetting 获取单个设置
func (h *SettingsHandler) GetSetting(key string) string {
	var setting models.Setting
	if err := h.db.Where("key = ?", key).First(&setting).Error; err != nil {
		// 返回默认值
		if defaultValue, ok := models.DefaultSettings[key]; ok {
			return defaultValue
		}
		return ""
	}
	return setting.Value
}

