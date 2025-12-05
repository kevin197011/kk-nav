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

// CategoriesHandler 管理后台分类处理器
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
	h.db.Order("sort_order").Find(&categories)

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
	if err := h.db.First(&category, id).Error; err != nil {
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

// Create 创建分类
func (h *CategoriesHandler) Create(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查名称是否已存在
	var existing models.Category
	if err := h.db.Where("name = ?", category.Name).First(&existing).Error; err == nil {
		utils.Error(c, 400, "Category name already exists")
		return
	}

	if err := h.db.Create(&category).Error; err != nil {
		utils.InternalServerError(c, "Failed to create category")
		return
	}

	utils.SuccessWithMessage(c, "Category created successfully", category)
}

// Update 更新分类
func (h *CategoriesHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	var category models.Category
	if err := h.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Category not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查名称是否与其他分类冲突
	var existing models.Category
	if err := h.db.Where("name = ? AND id != ?", category.Name, id).First(&existing).Error; err == nil {
		utils.Error(c, 400, "Category name already exists")
		return
	}

	if err := h.db.Save(&category).Error; err != nil {
		utils.InternalServerError(c, "Failed to update category")
		return
	}

	utils.SuccessWithMessage(c, "Category updated successfully", category)
}

// Delete 删除分类
func (h *CategoriesHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	// 检查是否有链接使用此分类
	var count int64
	h.db.Model(&models.Link{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		utils.Error(c, 400, "Cannot delete category with existing links")
		return
	}

	if err := h.db.Delete(&models.Category{}, id).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete category")
		return
	}

	utils.SuccessWithMessage(c, "Category deleted successfully", nil)
}

// MoveUp 上移
func (h *CategoriesHandler) MoveUp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	var category models.Category
	if err := h.db.First(&category, id).Error; err != nil {
		utils.NotFound(c, "Category not found")
		return
	}

	// 查找上一个分类
	var prevCategory models.Category
	if err := h.db.Where("sort_order < ?", category.SortOrder).
		Order("sort_order DESC").First(&prevCategory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, 400, "Already at the top")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 交换排序
	category.SortOrder, prevCategory.SortOrder = prevCategory.SortOrder, category.SortOrder
	h.db.Save(&category)
	h.db.Save(&prevCategory)

	utils.SuccessWithMessage(c, "Category moved up successfully", nil)
}

// MoveDown 下移
func (h *CategoriesHandler) MoveDown(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid category ID")
		return
	}

	var category models.Category
	if err := h.db.First(&category, id).Error; err != nil {
		utils.NotFound(c, "Category not found")
		return
	}

	// 查找下一个分类
	var nextCategory models.Category
	if err := h.db.Where("sort_order > ?", category.SortOrder).
		Order("sort_order ASC").First(&nextCategory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(c, 400, "Already at the bottom")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	// 交换排序
	category.SortOrder, nextCategory.SortOrder = nextCategory.SortOrder, category.SortOrder
	h.db.Save(&category)
	h.db.Save(&nextCategory)

	utils.SuccessWithMessage(c, "Category moved down successfully", nil)
}

