// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Link 链接模型
type Link struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"not null;size:255;unique" json:"title" binding:"required,min=1,max=255"`
	URL           string    `gorm:"not null;type:text" json:"url" binding:"required,url"`
	Description   string    `gorm:"type:text" json:"description"`
	CategoryID    uint      `gorm:"not null;index" json:"category_id" binding:"required"`
	SortOrder     int       `gorm:"not null" json:"sort_order"`
	Status        string    `gorm:"not null;default:'active';size:20;index" json:"status"` // active | inactive | error
	ClickCount    int       `gorm:"not null;default:0" json:"click_count"`
	LastCheckedAt *time.Time `gorm:"type:timestamp" json:"last_checked_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	Category   Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags       []Tag       `gorm:"many2many:link_tags;" json:"tags,omitempty"`
	Favorites  []Favorite  `gorm:"foreignKey:LinkID" json:"favorites,omitempty"`
	ClickLogs  []ClickLog  `gorm:"foreignKey:LinkID" json:"click_logs,omitempty"`
}

// TableName 指定表名
func (Link) TableName() string {
	return "links"
}

// IsActive 判断是否激活
func (l *Link) IsActive() bool {
	return l.Status == "active"
}

// IncrementClickCount 增加点击数
func (l *Link) IncrementClickCount() error {
	db := GetDB()
	if db == nil {
		return nil
	}
	return db.Model(l).UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error
}

// BeforeCreate 创建前钩子
func (l *Link) BeforeCreate(tx *gorm.DB) error {
	// 自动设置排序
	if l.SortOrder == 0 {
		var maxOrder int
		tx.Model(&Link{}).Where("category_id = ?", l.CategoryID).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)
		l.SortOrder = maxOrder + 1
	}

	// 规范化URL
	if err := l.normalizeURL(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate 更新前钩子
func (l *Link) BeforeUpdate(tx *gorm.DB) error {
	// 规范化URL
	return l.normalizeURL()
}

// normalizeURL 规范化URL
func (l *Link) normalizeURL() error {
	urlStr := strings.TrimSpace(l.URL)
	if urlStr == "" {
		return nil
	}

	// 如果没有协议，添加https://
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	// 验证URL格式
	_, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	l.URL = urlStr
	return nil
}

