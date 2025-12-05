// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"log"

	"kk-nav/internal/config"
	"kk-nav/internal/database"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer database.Close()

	// 执行迁移
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully!")
}

