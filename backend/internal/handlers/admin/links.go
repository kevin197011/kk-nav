// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package admin

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// LinksHandler 管理后台链接处理器
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
	query := h.db.Preload("Category").Preload("Tags")

	// 搜索
	if search := c.Query("search"); search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR url ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
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

	utils.Success(c, gin.H{
		"link": link,
	})
}

// Create 创建链接
func (h *LinksHandler) Create(c *gin.Context) {
	var req struct {
		Title       string   `json:"title" binding:"required"`
		URL         string   `json:"url" binding:"required,url"`
		Description string   `json:"description"`
		CategoryID  uint     `json:"category_id" binding:"required"`
		SortOrder   int      `json:"sort_order"`
		Status      string   `json:"status"`
		TagNames    []string `json:"tag_names"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查标题是否已存在
	var existingLink models.Link
	if err := h.db.Where("title = ?", req.Title).First(&existingLink).Error; err == nil {
		utils.Error(c, 400, "Link title already exists")
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := h.db.First(&category, req.CategoryID).Error; err != nil {
		utils.NotFound(c, "Category not found")
		return
	}

	// 处理标签
	var tags []models.Tag
	if len(req.TagNames) > 0 {
		for _, tagName := range req.TagNames {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}
			var tag models.Tag
			if err := h.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
				// 创建新标签
				tag = models.Tag{
					Name:  tagName,
					Color: "#007bff", // 默认颜色
				}
				h.db.Create(&tag)
			}
			tags = append(tags, tag)
		}
	}

	link := models.Link{
		Title:       req.Title,
		URL:         req.URL,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		Tags:        tags,
	}

	if link.Status == "" {
		link.Status = "active"
	}

	if err := h.db.Create(&link).Error; err != nil {
		utils.InternalServerError(c, "Failed to create link")
		return
	}

	h.db.Preload("Category").Preload("Tags").First(&link, link.ID)
	utils.SuccessWithMessage(c, "Link created successfully", link)
}

// Update 更新链接
func (h *LinksHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	var link models.Link
	if err := h.db.Preload("Tags").First(&link, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Link not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var req struct {
		Title       string   `json:"title"`
		URL         string   `json:"url"`
		Description string   `json:"description"`
		CategoryID  uint     `json:"category_id"`
		SortOrder   int      `json:"sort_order"`
		Status      string   `json:"status"`
		TagNames    []string `json:"tag_names"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 更新字段
	if req.Title != "" && req.Title != link.Title {
		// 检查新标题是否已被其他链接使用
		var existingLink models.Link
		if err := h.db.Where("title = ? AND id != ?", req.Title, id).First(&existingLink).Error; err == nil {
			utils.Error(c, 400, "Link title already exists")
			return
		}
		link.Title = req.Title
	}
	if req.URL != "" {
		link.URL = req.URL
	}
	if req.Description != "" {
		link.Description = req.Description
	}
	if req.CategoryID > 0 {
		// 检查分类是否存在
		var category models.Category
		if err := h.db.First(&category, req.CategoryID).Error; err != nil {
			utils.NotFound(c, "Category not found")
			return
		}
		link.CategoryID = req.CategoryID
	}
	if req.SortOrder > 0 {
		link.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		link.Status = req.Status
	}

	// 更新标签
	if req.TagNames != nil {
		var tags []models.Tag
		for _, tagName := range req.TagNames {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}
			var tag models.Tag
			if err := h.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
				// 创建新标签
				tag = models.Tag{
					Name:  tagName,
					Color: "#007bff",
				}
				h.db.Create(&tag)
			}
			tags = append(tags, tag)
		}
		link.Tags = tags
	}

	if err := h.db.Save(&link).Error; err != nil {
		utils.InternalServerError(c, "Failed to update link")
		return
	}

	h.db.Preload("Category").Preload("Tags").First(&link, link.ID)
	utils.SuccessWithMessage(c, "Link updated successfully", link)
}

// Delete 删除链接
func (h *LinksHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	if err := h.db.Delete(&models.Link{}, id).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete link")
		return
	}

	utils.SuccessWithMessage(c, "Link deleted successfully", nil)
}

// CheckStatus 检测链接状态
func (h *LinksHandler) CheckStatus(c *gin.Context) {
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

	// 检测链接状态
	status := h.checkLinkStatus(link.URL)
	now := time.Now()
	link.Status = status
	link.LastCheckedAt = &now

	if err := h.db.Save(&link).Error; err != nil {
		utils.InternalServerError(c, "Failed to update link status")
		return
	}

	utils.Success(c, gin.H{
		"link":   link,
		"status": status,
	})
}

// BatchCheckStatus 批量检测链接状态
func (h *LinksHandler) BatchCheckStatus(c *gin.Context) {
	var req struct {
		LinkIDs []uint `json:"link_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if len(req.LinkIDs) == 0 {
		utils.BadRequest(c, "Link IDs required")
		return
	}

	var links []models.Link
	h.db.Where("id IN ?", req.LinkIDs).Find(&links)

	now := time.Now()
	results := make([]gin.H, 0, len(links))

	for _, link := range links {
		status := h.checkLinkStatus(link.URL)
		link.Status = status
		link.LastCheckedAt = &now
		h.db.Save(&link)

		results = append(results, gin.H{
			"id":     link.ID,
			"title":  link.Title,
			"url":    link.URL,
			"status": status,
		})
	}

	utils.SuccessWithMessage(c, "Batch check completed", gin.H{
		"results": results,
		"total":   len(results),
	})
}

// checkLinkStatus 检测链接状态
func (h *LinksHandler) checkLinkStatus(url string) string {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "active"
	}
	return "error"
}

// MoveUp 上移
func (h *LinksHandler) MoveUp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	var link models.Link
	if err := h.db.First(&link, id).Error; err != nil {
		utils.NotFound(c, "Link not found")
		return
	}

	// 查找上一个链接
	var prevLink models.Link
	if err := h.db.Where("category_id = ? AND sort_order < ?", link.CategoryID, link.SortOrder).
		Order("sort_order DESC").First(&prevLink).Error; err != nil {
		utils.Error(c, 400, "Cannot move up")
		return
	}

	// 交换排序
	link.SortOrder, prevLink.SortOrder = prevLink.SortOrder, link.SortOrder
	h.db.Save(&link)
	h.db.Save(&prevLink)

	utils.SuccessWithMessage(c, "Link moved up successfully", nil)
}

// MoveDown 下移
func (h *LinksHandler) MoveDown(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid link ID")
		return
	}

	var link models.Link
	if err := h.db.First(&link, id).Error; err != nil {
		utils.NotFound(c, "Link not found")
		return
	}

	// 查找下一个链接
	var nextLink models.Link
	if err := h.db.Where("category_id = ? AND sort_order > ?", link.CategoryID, link.SortOrder).
		Order("sort_order ASC").First(&nextLink).Error; err != nil {
		utils.Error(c, 400, "Cannot move down")
		return
	}

	// 交换排序
	link.SortOrder, nextLink.SortOrder = nextLink.SortOrder, link.SortOrder
	h.db.Save(&link)
	h.db.Save(&nextLink)

	utils.SuccessWithMessage(c, "Link moved down successfully", nil)
}

