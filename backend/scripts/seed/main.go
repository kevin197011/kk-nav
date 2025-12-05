// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"log"

	"kk-nav/internal/config"
	"kk-nav/internal/database"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
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

	db := database.DB

	// 从环境变量读取管理员配置
	adminEmail := cfg.GetString("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@example.com"
	}
	adminUsername := cfg.GetString("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := cfg.GetString("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	// 创建管理员用户
	passwordHash, _ := utils.HashPassword(adminPassword)
	admin := models.User{
		Email:        adminEmail,
		Username:     adminUsername,
		PasswordHash: passwordHash,
		Role:         "admin",
		Active:       true,
	}
	if err := db.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}
	log.Printf("Created admin user: %s", admin.Email)

	// 创建默认设置
	for key, value := range models.DefaultSettings {
		setting := models.Setting{
			Key:         key,
			Value:       value,
			Description: "系统设置: " + key,
		}
		if err := db.FirstOrCreate(&setting, models.Setting{Key: key}).Error; err != nil {
			log.Printf("Warning: Failed to create setting %s: %v", key, err)
		}
	}
	log.Println("Created default settings")

	// 创建分类
	categoriesData := []models.Category{
		{Name: "监控工具", Icon: "📊", Color: "#007bff", Description: "系统监控和性能分析工具", SortOrder: 1, Active: true},
		{Name: "日志分析", Icon: "📝", Color: "#28a745", Description: "日志收集、分析和可视化工具", SortOrder: 2, Active: true},
		{Name: "容器管理", Icon: "🐳", Color: "#17a2b8", Description: "容器化和编排管理工具", SortOrder: 3, Active: true},
		{Name: "云服务", Icon: "☁️", Color: "#ffc107", Description: "云计算平台和服务", SortOrder: 4, Active: true},
		{Name: "开发工具", Icon: "🛠️", Color: "#6f42c1", Description: "开发和调试工具", SortOrder: 5, Active: true},
		{Name: "网络工具", Icon: "🌐", Color: "#fd7e14", Description: "网络诊断和管理工具", SortOrder: 6, Active: true},
	}

	var categories []models.Category
	for _, cat := range categoriesData {
		var existing models.Category
		if err := db.Where("name = ?", cat.Name).First(&existing).Error; err != nil {
			if err := db.Create(&cat).Error; err != nil {
				log.Printf("Warning: Failed to create category %s: %v", cat.Name, err)
				continue
			}
			categories = append(categories, cat)
		} else {
			categories = append(categories, existing)
		}
	}
	log.Printf("Created %d categories", len(categories))

	// 创建标签
	tagsData := []models.Tag{
		{Name: "监控", Color: "#007bff"},
		{Name: "可视化", Color: "#28a745"},
		{Name: "告警", Color: "#dc3545"},
		{Name: "日志", Color: "#6c757d"},
		{Name: "容器", Color: "#17a2b8"},
		{Name: "Docker", Color: "#0db7ed"},
		{Name: "Kubernetes", Color: "#326ce5"},
		{Name: "AWS", Color: "#ff9900"},
		{Name: "开源", Color: "#28a745"},
		{Name: "商业", Color: "#ffc107"},
	}

	var tags []models.Tag
	for _, tag := range tagsData {
		var existing models.Tag
		if err := db.Where("name = ?", tag.Name).First(&existing).Error; err != nil {
			if err := db.Create(&tag).Error; err != nil {
				log.Printf("Warning: Failed to create tag %s: %v", tag.Name, err)
				continue
			}
			tags = append(tags, tag)
		} else {
			tags = append(tags, existing)
		}
	}
	log.Printf("Created %d tags", len(tags))

	// 创建链接
	linksData := []struct {
		Title       string
		URL         string
		Description string
		Category    string
		TagNames    []string
	}{
		{"Grafana", "https://grafana.com", "开源的数据可视化和监控平台，支持多种数据源", "监控工具", []string{"监控", "可视化", "开源"}},
		{"Prometheus", "https://prometheus.io", "开源的监控和告警系统，专为云原生环境设计", "监控工具", []string{"监控", "告警", "开源"}},
		{"Zabbix", "https://www.zabbix.com", "企业级开源监控解决方案", "监控工具", []string{"监控", "企业级", "开源"}},
		{"Elasticsearch", "https://www.elastic.co/elasticsearch", "分布式搜索和分析引擎", "日志分析", []string{"日志", "搜索", "开源"}},
		{"Kibana", "https://www.elastic.co/kibana", "Elasticsearch的数据可视化平台", "日志分析", []string{"日志", "可视化", "开源"}},
		{"Docker", "https://www.docker.com", "容器化平台，简化应用部署和管理", "容器管理", []string{"容器", "Docker", "开源"}},
		{"Kubernetes", "https://kubernetes.io", "容器编排和管理平台", "容器管理", []string{"容器", "Kubernetes", "编排", "开源"}},
		{"AWS Console", "https://console.aws.amazon.com", "Amazon Web Services管理控制台", "云服务", []string{"AWS", "云计算", "商业"}},
	}

	for i, linkData := range linksData {
		// 查找分类
		var category models.Category
		for _, cat := range categories {
			if cat.Name == linkData.Category {
				category = cat
				break
			}
		}

		// 查找标签
		var linkTags []models.Tag
		for _, tagName := range linkData.TagNames {
			for _, tag := range tags {
				if tag.Name == tagName {
					linkTags = append(linkTags, tag)
					break
				}
			}
		}

		link := models.Link{
			Title:       linkData.Title,
			URL:         linkData.URL,
			Description: linkData.Description,
			CategoryID:  category.ID,
			SortOrder:   i + 1,
			Status:      "active",
			Tags:        linkTags,
		}

		var existing models.Link
		if err := db.Where("title = ?", link.Title).First(&existing).Error; err != nil {
			if err := db.Create(&link).Error; err != nil {
				log.Printf("Warning: Failed to create link %s: %v", link.Title, err)
			}
		}
	}
	log.Println("Created sample links")

	log.Println("Database seeding completed successfully!")
	log.Printf("Admin account: %s / %s", adminEmail, adminPassword)
}

