// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// TagsHandler 标签处理器
type TagsHandler struct {
	db *gorm.DB
}

// NewTagsHandler 创建标签处理器
func NewTagsHandler(db *gorm.DB) *TagsHandler {
	return &TagsHandler{db: db}
}

// Index 标签列表
func (h *TagsHandler) Index(c *gin.Context) {
	var tags []models.Tag

	// 获取热门标签（按链接数排序）
	query := h.db.Joins("JOIN link_tags ON link_tags.tag_id = tags.id").
		Joins("JOIN links ON links.id = link_tags.link_id").
		Where("links.status = ?", "active").
		Group("tags.id").
		Order("COUNT(links.id) DESC")

	// 限制数量
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 0 {
		query = query.Limit(limit)
	}

	query.Find(&tags)

	// 加载每个标签的链接数量
	for i := range tags {
		tags[i].LinksCount()
	}

	utils.Success(c, gin.H{
		"tags":  tags,
		"total": len(tags),
	})
}

// Show 标签详情
func (h *TagsHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid tag ID")
		return
	}

	var tag models.Tag
	if err := h.db.Preload("Links", "status = ?", "active").First(&tag, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Tag not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	utils.Success(c, gin.H{
		"tag": tag,
	})
}

