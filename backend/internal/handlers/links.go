// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// LinksHandler 链接处理器
type LinksHandler struct {
	db *gorm.DB
}

// NewLinksHandler 创建链接处理器
func NewLinksHandler(db *gorm.DB) *LinksHandler {
	return &LinksHandler{db: db}
}

// Index 链接列表
func (h *LinksHandler) Index(c *gin.Context) {
	var links []models.Link
	query := h.db.Where("status = ?", "active").Preload("Category").Preload("Tags")

	// 搜索
	if search := c.Query("search"); search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR url ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 标签筛选
	if tag := c.Query("tag"); tag != "" {
		query = query.Joins("JOIN link_tags ON link_tags.link_id = links.id").
			Joins("JOIN tags ON tags.id = link_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// 分页（前台默认显示所有链接，不分页）
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "1000"))
	offset := (page - 1) * pageSize

	var total int64
	query.Model(&models.Link{}).Count(&total)
	query.Order("sort_order").Offset(offset).Limit(pageSize).Find(&links)

	utils.Success(c, gin.H{
		"links": links,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// Show 链接详情
func (h *LinksHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	var link models.Link
	if err := h.db.Preload("Category").Preload("Tags").First(&link, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Link not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 获取相关链接
	var relatedLinks []models.Link
	h.db.Where("category_id = ? AND id != ? AND status = ?", link.CategoryID, link.ID, "active").
		Order("click_count DESC").
		Limit(6).
		Find(&relatedLinks)

	utils.Success(c, gin.H{
		"link":         link,
		"related_links": relatedLinks,
	})
}

// Click 记录点击
func (h *LinksHandler) Click(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	var link models.Link
	if err := h.db.First(&link, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Link not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 增加点击数
	if err := link.IncrementClickCount(); err != nil {
		utils.InternalServerError(c, "Failed to increment click count")
		return
	}

	// 记录点击日志
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}

	clickLog := models.ClickLog{
		LinkID:    link.ID,
		UserID:    userID,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   c.GetHeader("Referer"),
	}

	if err := h.db.Create(&clickLog).Error; err != nil {
		// 日志记录失败不影响主流程
		_ = err
	}

	// 重定向到链接
	c.Redirect(http.StatusFound, link.URL)
}

// Favorite 添加收藏
func (h *LinksHandler) Favorite(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "Authentication required")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	// 检查链接是否存在
	var link models.Link
	if err := h.db.First(&link, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Link not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 创建收藏
	favorite := models.Favorite{
		UserID: userID.(uint),
		LinkID: link.ID,
	}

	if err := h.db.Create(&favorite).Error; err != nil {
		// 如果已存在，返回成功
		if err == gorm.ErrDuplicatedKey {
			utils.SuccessWithMessage(c, "Already favorited", nil)
			return
		}
		utils.InternalServerError(c, "Failed to create favorite")
		return
	}

	utils.SuccessWithMessage(c, "Favorited successfully", nil)
}

// Unfavorite 取消收藏
func (h *LinksHandler) Unfavorite(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "Authentication required")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	// 删除收藏
	if err := h.db.Where("user_id = ? AND link_id = ?", userID, id).
		Delete(&models.Favorite{}).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete favorite")
		return
	}

	utils.SuccessWithMessage(c, "Unfavorited successfully", nil)
}

// Favorites 我的收藏列表
func (h *LinksHandler) Favorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "Authentication required")
		return
	}

	var favorites []models.Favorite
	h.db.Where("user_id = ?", userID).
		Preload("Link").Preload("Link.Category").Preload("Link.Tags").
		Order("created_at DESC").
		Find(&favorites)

	links := make([]models.Link, 0, len(favorites))
	for _, fav := range favorites {
		if fav.Link.Status == "active" {
			links = append(links, fav.Link)
		}
	}

	utils.Success(c, gin.H{
		"links": links,
		"total": len(links),
	})
}

