// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email" binding:"required,email"`
	Username     string    `gorm:"uniqueIndex;not null;size:100" json:"username" binding:"required,min=3,max=100"`
	PasswordHash string    `gorm:"not null;size:255" json:"-"`
	Role         string    `gorm:"not null;default:'user';size:20" json:"role"` // user | admin
	Active       bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	Favorites  []Favorite  `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
	ClickLogs  []ClickLog  `gorm:"foreignKey:UserID" json:"click_logs,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsAdmin 判断是否为管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// DisplayName 获取显示名称
func (u *User) DisplayName() string {
	if u.Username != "" {
		return u.Username
	}
	// 如果没有用户名，返回邮箱前缀
	if u.Email != "" {
		return u.Email[:len(u.Email)-len("@example.com")]
	}
	return "User"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = "user"
	}
	return nil
}

