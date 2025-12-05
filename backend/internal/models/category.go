// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null;size:100" json:"name" binding:"required,min=1,max=100"`
	Icon        string    `gorm:"not null;default:'📁';size:50" json:"icon" binding:"required"`
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"not null;default:'#007bff';size:7" json:"color" binding:"required"`
	SortOrder   int       `gorm:"uniqueIndex;not null" json:"sort_order"`
	Active      bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Links []Link `gorm:"foreignKey:CategoryID" json:"links,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// LinksCount 获取激活链接数量
func (c *Category) LinksCount() int64 {
	var count int64
	db := GetDB()
	if db != nil {
		db.Model(&Link{}).Where("category_id = ? AND status = ?", c.ID, "active").Count(&count)
	}
	return count
}

// TotalClicks 获取总点击数
func (c *Category) TotalClicks() int64 {
	var total int64
	db := GetDB()
	if db != nil {
		db.Model(&Link{}).Where("category_id = ?", c.ID).Select("COALESCE(SUM(click_count), 0)").Scan(&total)
	}
	return total
}

// BeforeCreate 创建前钩子
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.SortOrder == 0 {
		// 自动设置排序
		var maxOrder int
		tx.Model(&Category{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)
		c.SortOrder = maxOrder + 1
	}
	return nil
}

