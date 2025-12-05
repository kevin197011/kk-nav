// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"time"
)

// ClickLog 点击日志模型
type ClickLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LinkID    uint      `gorm:"not null;index" json:"link_id" binding:"required"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"`
	IPAddress string    `gorm:"not null;size:45" json:"ip_address" binding:"required"`
	UserAgent string    `gorm:"not null;type:text" json:"user_agent" binding:"required"`
	Referer   string    `gorm:"type:text" json:"referer"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Link Link  `gorm:"foreignKey:LinkID" json:"link,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (ClickLog) TableName() string {
	return "click_logs"
}

