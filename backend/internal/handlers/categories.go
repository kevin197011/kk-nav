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

// CategoriesHandler 分类处理器
type CategoriesHandler struct {
	db *gorm.DB
}

// NewCategoriesHandler 创建分类处理器
func NewCategoriesHandler(db *gorm.DB) *CategoriesHandler {
	return &CategoriesHandler{db: db}
}

// Index 分类列表
func (h *CategoriesHandler) Index(c *gin.Context) {
	var categories []models.Category
	query := h.db

	// 只显示激活的分类
	if c.Query("active") != "false" {
		query = query.Where("active = ?", true)
	}

	query.Order("sort_order").Find(&categories)

	// 加载每个分类的链接数量
	for i := range categories {
		categories[i].LinksCount()
		categories[i].TotalClicks()
	}

	utils.Success(c, gin.H{
		"categories": categories,
		"total":      len(categories),
	})
}

// Show 分类详情
func (h *CategoriesHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	var category models.Category
	if err := h.db.Preload("Links", "status = ?", "active").First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Category not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	utils.Success(c, gin.H{
		"category": category,
	})
}
