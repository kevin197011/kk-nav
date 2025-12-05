// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// HomeHandler 首页处理器
type HomeHandler struct {
	db *gorm.DB
}

// NewHomeHandler 创建首页处理器
func NewHomeHandler(db *gorm.DB) *HomeHandler {
	return &HomeHandler{db: db}
}

// Index 首页
func (h *HomeHandler) Index(c *gin.Context) {
	// 获取分类
	var categories []models.Category
	h.db.Where("active = ?", true).Order("sort_order").Find(&categories)

	// 获取链接
	var links []models.Link
	query := h.db.Where("status = ?", "active").Preload("Category").Preload("Tags")

	// 搜索
	if search := c.Query("search"); search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR url ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 标签筛选
	if tag := c.Query("tag"); tag != "" {
		query = query.Joins("JOIN link_tags ON link_tags.link_id = links.id").
			Joins("JOIN tags ON tags.id = link_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Find(&links)

	// 获取热门标签
	var tags []models.Tag
	h.db.Joins("JOIN link_tags ON link_tags.tag_id = tags.id").
		Joins("JOIN links ON links.id = link_tags.link_id").
		Where("links.status = ?", "active").
		Group("tags.id").
		Order("COUNT(links.id) DESC").
		Limit(20).
		Find(&tags)

	// 统计数据
	var stats struct {
		TotalLinks      int64
		TotalCategories int64
		TotalClicks     int64
		TodayClicks     int64
	}
	h.db.Model(&models.Link{}).Where("status = ?", "active").Count(&stats.TotalLinks)
	h.db.Model(&models.Category{}).Where("active = ?", true).Count(&stats.TotalCategories)
	h.db.Model(&models.ClickLog{}).Count(&stats.TotalClicks)

	// 今日点击统计
	todayStart := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	h.db.Model(&models.ClickLog{}).Where("created_at >= ?", todayStart).Count(&stats.TodayClicks)

	utils.Success(c, gin.H{
		"categories": categories,
		"links":      links,
		"tags":       tags,
		"stats":      stats,
	})
}

