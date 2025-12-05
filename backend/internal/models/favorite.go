// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"time"

	"gorm.io/gorm"
)

// Favorite 收藏模型
type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_user_link" json:"user_id" binding:"required"`
	LinkID    uint      `gorm:"not null;index;uniqueIndex:idx_user_link" json:"link_id" binding:"required"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Link Link `gorm:"foreignKey:LinkID" json:"link,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// BeforeCreate 创建前钩子
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	// 检查是否已存在
	var count int64
	tx.Model(&Favorite{}).
		Where("user_id = ? AND link_id = ?", f.UserID, f.LinkID).
		Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

