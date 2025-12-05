// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"gorm.io/gorm"
)

// WebHandler Web页面处理器
type WebHandler struct {
	db *gorm.DB
}

// NewWebHandler 创建Web处理器
func NewWebHandler(db *gorm.DB) *WebHandler {
	return &WebHandler{db: db}
}

// Home 首页
func (h *WebHandler) Home(c *gin.Context) {
	// 获取用户信息（如果有token）
	var user *models.User
	if userID, exists := c.Get("user_id"); exists {
		h.db.First(&user, userID)
	}

	// 获取分类
	var categories []models.Category
	h.db.Where("active = ?", true).Order("sort_order").Find(&categories)

	// 获取链接
	var links []models.Link
	query := h.db.Where("status = ?", "active").Preload("Category").Preload("Tags")

	// 搜索
	search := c.Query("search")
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR url ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 标签筛选
	tag := c.Query("tag")
	if tag != "" {
		query = query.Joins("JOIN link_tags ON link_tags.link_id = links.id").
			Joins("JOIN tags ON tags.id = link_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Find(&links)

	// 按分类组织链接
	categoryLinksMap := make(map[uint][]models.Link)
	for _, link := range links {
		categoryLinksMap[link.CategoryID] = append(categoryLinksMap[link.CategoryID], link)
	}

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

	// 渲染模板
	// 渲染布局模板，它会查找 "home-content" block
	c.HTML(http.StatusOK, "layouts/base.html", gin.H{
		"Title":           "首页",
		"User":            user,
		"Categories":      categories,
		"CategoryLinksMap": categoryLinksMap,
		"Tags":            tags,
		"Search":          search,
		"Tag":             tag,
		"Stats":           stats,
	})
}

// Login 登录页面
func (h *WebHandler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "layouts/base.html", gin.H{
		"Title": "登录",
		"User":  nil,
	})
}

// Register 注册页面
func (h *WebHandler) Register(c *gin.Context) {
	c.HTML(http.StatusOK, "layouts/base.html", gin.H{
		"Title": "注册",
		"User":  nil,
	})
}

// AdminDashboard 管理后台首页
func (h *WebHandler) AdminDashboard(c *gin.Context) {
	// 检查用户权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil || !user.IsAdmin() {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// 获取统计数据
	var stats struct {
		TotalLinks      int64
		ActiveLinks     int64
		InactiveLinks   int64
		ErrorLinks      int64
		TotalCategories int64
		TotalTags       int64
		TotalUsers      int64
		TotalClicks     int64
		TodayClicks      int64
		ThisWeekClicks  int64
		ThisMonthClicks int64
	}

	h.db.Model(&models.Link{}).Count(&stats.TotalLinks)
	h.db.Model(&models.Link{}).Where("status = ?", "active").Count(&stats.ActiveLinks)
	h.db.Model(&models.Link{}).Where("status = ?", "inactive").Count(&stats.InactiveLinks)
	h.db.Model(&models.Link{}).Where("status = ?", "error").Count(&stats.ErrorLinks)
	h.db.Model(&models.Category{}).Count(&stats.TotalCategories)
	h.db.Model(&models.Tag{}).Count(&stats.TotalTags)
	h.db.Model(&models.User{}).Count(&stats.TotalUsers)
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

	c.HTML(http.StatusOK, "layouts/admin.html", gin.H{
		"Title":         "仪表盘",
		"User":          user,
		"Stats":         stats,
		"PopularLinks":  popularLinks,
		"RecentClicks":  recentClicks,
	})
}

// AdminCategories 管理后台分类管理
func (h *WebHandler) AdminCategories(c *gin.Context) {
	// 检查用户权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil || !user.IsAdmin() {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "layouts/admin.html", gin.H{
		"Title": "分类管理",
		"User":  user,
	})
}

// AdminLinks 管理后台链接管理
func (h *WebHandler) AdminLinks(c *gin.Context) {
	// 检查用户权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil || !user.IsAdmin() {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "layouts/admin.html", gin.H{
		"Title": "链接管理",
		"User":  user,
	})
}

// AdminTags 管理后台标签管理
func (h *WebHandler) AdminTags(c *gin.Context) {
	// 检查用户权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil || !user.IsAdmin() {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "layouts/admin.html", gin.H{
		"Title": "标签管理",
		"User":  user,
	})
}

// AdminUsers 管理后台用户管理
func (h *WebHandler) AdminUsers(c *gin.Context) {
	// 检查用户权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil || !user.IsAdmin() {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "layouts/admin.html", gin.H{
		"Title": "用户管理",
		"User":  user,
	})
}

// AdminSettings 管理后台系统设置
func (h *WebHandler) AdminSettings(c *gin.Context) {
	// 检查用户权限
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil || !user.IsAdmin() {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "layouts/admin.html", gin.H{
		"Title": "系统设置",
		"User":  user,
	})
}
