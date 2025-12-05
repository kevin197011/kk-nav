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

// StatsHandler 统计处理器
type StatsHandler struct {
	db *gorm.DB
}

// NewStatsHandler 创建统计处理器
func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

// Index 统计数据
func (h *StatsHandler) Index(c *gin.Context) {
	var stats struct {
		TotalLinks      int64 `json:"total_links"`
		TotalCategories int64 `json:"total_categories"`
		TotalTags       int64 `json:"total_tags"`
		TotalClicks     int64 `json:"total_clicks"`
		TodayClicks     int64 `json:"today_clicks"`
		ThisWeekClicks  int64 `json:"this_week_clicks"`
		ThisMonthClicks int64 `json:"this_month_clicks"`
	}

	// 基础统计
	h.db.Model(&models.Link{}).Where("status = ?", "active").Count(&stats.TotalLinks)
	h.db.Model(&models.Category{}).Where("active = ?", true).Count(&stats.TotalCategories)
	h.db.Model(&models.Tag{}).Count(&stats.TotalTags)
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
		Title      string `json:"title"`
		ClickCount int    `json:"click_count"`
	}
	h.db.Model(&models.Link{}).
		Where("status = ?", "active").
		Order("click_count DESC").
		Limit(10).
		Select("title, click_count").
		Scan(&popularLinks)

	utils.Success(c, gin.H{
		"stats":        stats,
		"popular_links": popularLinks,
	})
}

