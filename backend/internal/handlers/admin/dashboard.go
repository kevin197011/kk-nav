// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package admin

import (
	"time"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// DashboardHandler 管理后台首页处理器
type DashboardHandler struct {
	db *gorm.DB
}

// NewDashboardHandler 创建首页处理器
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// Index 首页统计
func (h *DashboardHandler) Index(c *gin.Context) {
	var stats struct {
		TotalLinks      int64 `json:"total_links"`
		ActiveLinks     int64 `json:"active_links"`
		InactiveLinks   int64 `json:"inactive_links"`
		ErrorLinks      int64 `json:"error_links"`
		TotalCategories int64 `json:"total_categories"`
		TotalTags       int64 `json:"total_tags"`
		TotalUsers      int64 `json:"total_users"`
		TotalClicks     int64 `json:"total_clicks"`
		TodayClicks     int64 `json:"today_clicks"`
		ThisWeekClicks  int64 `json:"this_week_clicks"`
		ThisMonthClicks int64 `json:"this_month_clicks"`
	}

	// 链接统计
	h.db.Model(&models.Link{}).Count(&stats.TotalLinks)
	h.db.Model(&models.Link{}).Where("status = ?", "active").Count(&stats.ActiveLinks)
	h.db.Model(&models.Link{}).Where("status = ?", "inactive").Count(&stats.InactiveLinks)
	h.db.Model(&models.Link{}).Where("status = ?", "error").Count(&stats.ErrorLinks)

	// 分类和标签统计
	h.db.Model(&models.Category{}).Count(&stats.TotalCategories)
	h.db.Model(&models.Tag{}).Count(&stats.TotalTags)
	h.db.Model(&models.User{}).Count(&stats.TotalUsers)

	// 点击统计
	h.db.Model(&models.ClickLog{}).Count(&stats.TotalClicks)

	// 时间范围统计
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := now.AddDate(0, 0, -7)
	monthStart := now.AddDate(0, -1, 0)

	h.db.Model(&models.ClickLog{}).Where("created_at >= ?", todayStart).Count(&stats.TodayClicks)
	h.db.Model(&models.ClickLog{}).Where("created_at >= ?", weekStart).Count(&stats.ThisWeekClicks)
	h.db.Model(&models.ClickLog{}).Where("created_at >= ?", monthStart).Count(&stats.ThisMonthClicks)

	// 热门链接
	var popularLinks []struct {
		ID         uint   `json:"id"`
		Title      string `json:"title"`
		ClickCount int    `json:"click_count"`
	}
	h.db.Model(&models.Link{}).
		Order("click_count DESC").
		Limit(10).
		Select("id, title, click_count").
		Scan(&popularLinks)

	// 最近点击
	var recentClicks []models.ClickLog
	h.db.Preload("Link").Preload("User").
		Order("created_at DESC").
		Limit(10).
		Find(&recentClicks)

	utils.Success(c, gin.H{
		"stats":        stats,
		"popular_links": popularLinks,
		"recent_clicks": recentClicks,
	})
}

