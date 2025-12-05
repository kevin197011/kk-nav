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

// APIToken API Token 模型
type APIToken struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;size:255" json:"name" binding:"required"`
	Token       string    `gorm:"uniqueIndex;not null;size:255" json:"token"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Active      bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (APIToken) TableName() string {
	return "api_tokens"
}

// BeforeCreate 创建前生成 token
func (t *APIToken) BeforeCreate(tx *gorm.DB) error {
	if t.Token == "" {
		token, err := generateToken()
		if err != nil {
			return err
		}
		t.Token = token
	}
	return nil
}

// generateToken 生成随机 token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "kk_" + hex.EncodeToString(bytes), nil
}

// IsExpired 检查 token 是否过期
func (t *APIToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*t.ExpiresAt)
}

// IsValid 检查 token 是否有效
func (t *APIToken) IsValid() bool {
	return t.Active && !t.IsExpired()
}

