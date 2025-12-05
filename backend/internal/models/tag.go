// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null;size:100" json:"name" binding:"required,min=1,max=100"`
	Color     string    `gorm:"not null;size:7" json:"color" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Links []Link `gorm:"many2many:link_tags;" json:"links,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// LinksCount 获取激活链接数量
func (t *Tag) LinksCount() int64 {
	var count int64
	db := GetDB()
	if db != nil {
		db.Model(&Link{}).
			Joins("JOIN link_tags ON link_tags.link_id = links.id").
			Where("link_tags.tag_id = ? AND links.status = ?", t.ID, "active").
			Count(&count)
	}
	return count
}

// BeforeCreate 创建前钩子
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	// 如果没有设置颜色，自动生成随机颜色
	if t.Color == "" {
		t.Color = generateRandomColor()
	}
	return nil
}

// generateRandomColor 生成随机颜色
func generateRandomColor() string {
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "#007bff" // 默认颜色
	}
	return "#" + hex.EncodeToString(bytes)
}

