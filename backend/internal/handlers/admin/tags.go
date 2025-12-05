// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// TagsHandler 管理后台标签处理器
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

	// 搜索
	query := h.db
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Model(&models.Tag{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tags)

	// 加载每个标签的链接数量
	for i := range tags {
		tags[i].LinksCount()
	}

	utils.Success(c, gin.H{
		"tags": tags,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (int(total) + pageSize - 1) / pageSize,
		},
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

// Create 创建标签
func (h *TagsHandler) Create(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查名称是否已存在
	var existing models.Tag
	if err := h.db.Where("name = ?", tag.Name).First(&existing).Error; err == nil {
		utils.Error(c, 400, "Tag name already exists")
		return
	}

	if err := h.db.Create(&tag).Error; err != nil {
		utils.InternalServerError(c, "Failed to create tag")
		return
	}

	utils.SuccessWithMessage(c, "Tag created successfully", tag)
}

// Update 更新标签
func (h *TagsHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid tag ID")
		return
	}

	var tag models.Tag
	if err := h.db.First(&tag, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Tag not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	if err := c.ShouldBindJSON(&tag); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查名称是否与其他标签冲突
	var existing models.Tag
	if err := h.db.Where("name = ? AND id != ?", tag.Name, id).First(&existing).Error; err == nil {
		utils.Error(c, 400, "Tag name already exists")
		return
	}

	if err := h.db.Save(&tag).Error; err != nil {
		utils.InternalServerError(c, "Failed to update tag")
		return
	}

	utils.SuccessWithMessage(c, "Tag updated successfully", tag)
}

// Delete 删除标签
func (h *TagsHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid tag ID")
		return
	}

	// 检查是否有链接使用此标签
	var count int64
	h.db.Model(&models.Link{}).
		Joins("JOIN link_tags ON link_tags.link_id = links.id").
		Where("link_tags.tag_id = ?", id).
		Count(&count)

	if count > 0 {
		utils.Error(c, 400, "Cannot delete tag with existing links")
		return
	}

	if err := h.db.Delete(&models.Tag{}, id).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete tag")
		return
	}

	utils.SuccessWithMessage(c, "Tag deleted successfully", nil)
}

