// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"fmt"

	"kk-nav/internal/config"
	"kk-nav/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect 连接数据库
func Connect(cfg *config.Config) error {
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.Name + ".db")
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// 配置日志级别
	logLevel := logger.Silent
	if cfg.App.Debug {
		logLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 设置models包的数据库实例
	models.SetDB(DB)

	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not connected")
	}

	// 导入所有模型
	return DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Link{},
		&models.Tag{},
		&models.Favorite{},
		&models.ClickLog{},
		&models.Setting{},
		&models.APIToken{},
	)
}

// InitializeData 初始化默认数据（管理员账号和系统设置）
func InitializeData(cfg *config.Config) error {
	if DB == nil {
		return fmt.Errorf("database not connected")
	}

	// 创建默认管理员账号
	if err := createDefaultAdmin(cfg); err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	// 创建默认系统设置
	if err := createDefaultSettings(); err != nil {
		return fmt.Errorf("failed to create default settings: %w", err)
	}

	return nil
}

// createDefaultAdmin 创建默认管理员账号
func createDefaultAdmin(cfg *config.Config) error {
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

	// 检查管理员是否已存在
	var count int64
	if err := DB.Model(&models.User{}).Where("email = ?", adminEmail).Count(&count).Error; err != nil {
		return err
	}

	// 如果已存在，跳过创建
	if count > 0 {
		return nil
	}

	// 创建管理员账号
	admin := &models.User{
		Email:    adminEmail,
		Username: adminUsername,
		Role:     "admin",
		Active:   true,
	}

	// 哈希密码（需要导入 utils 包）
	// 这里直接使用 bcrypt，避免循环依赖
	passwordHash, err := hashPassword(adminPassword)
	if err != nil {
		return err
	}
	admin.PasswordHash = passwordHash

	// 保存到数据库
	return DB.Create(admin).Error
}

// createDefaultSettings 创建默认系统设置
func createDefaultSettings() error {
	for key, value := range models.DefaultSettings {
		var count int64
		if err := DB.Model(&models.Setting{}).Where("key = ?", key).Count(&count).Error; err != nil {
			return err
		}

		// 如果已存在，跳过创建
		if count > 0 {
			continue
		}

		setting := &models.Setting{
			Key:         key,
			Value:       value,
			Description: "系统设置: " + key,
		}

		if err := DB.Create(setting).Error; err != nil {
			return err
		}
	}

	return nil
}

// hashPassword 哈希密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

