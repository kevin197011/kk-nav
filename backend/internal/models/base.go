// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package models

import (
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB 设置数据库实例（由database包调用）
func SetDB(database *gorm.DB) {
	db = database
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

