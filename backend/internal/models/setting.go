// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"time"
)

// Setting 系统设置模型
type Setting struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"uniqueIndex;not null;size:100" json:"key" binding:"required"`
	Value       string    `gorm:"not null;type:text" json:"value" binding:"required"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Setting) TableName() string {
	return "settings"
}

// DefaultSettings 默认设置
var DefaultSettings = map[string]string{
	"site_name":         "运维工具导航",
	"site_description": "专业的运维工具网址导航系统",
	"primary_color":     "#007bff",
	"theme":             "light",
	"enable_registration": "true",
	"enable_link_check":   "true",
	"check_interval_hours": "24",
	"links_per_page":      "12",
	"enable_analytics":    "true",
	"enable_pwa":          "true",
}

