// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Log      LogConfig
}

// AppConfig 应用配置
type AppConfig struct {
	Name  string
	Env   string
	Port  int
	Debug bool
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string
	Format string
}

var globalConfig *Config

// Load 加载配置
func Load() (*Config, error) {
	// 加载.env文件（如果存在）
	_ = godotenv.Load()

	// 设置viper配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	// 环境变量支持
	viper.AutomaticEnv()

	// 读取配置文件（可选）
	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时使用环境变量
		fmt.Println("Warning: config file not found, using environment variables")
	}

	config := &Config{
		App: AppConfig{
			Name:  getString("APP_NAME", "ops-nav"),
			Env:   getString("APP_ENV", "development"),
			Port:  getInt("APP_PORT", 8080),
			Debug: getBool("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Type:            getString("DB_TYPE", "postgres"),
			Host:            getString("DB_HOST", "localhost"),
			Port:            getInt("DB_PORT", 5432),
			Name:            getString("DB_NAME", "ops_nav"),
			User:            getString("DB_USER", "postgres"),
			Password:        getString("DB_PASSWORD", ""),
			SSLMode:         getString("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: time.Duration(getInt("DB_CONN_MAX_LIFETIME", 300)) * time.Second,
		},
		JWT: JWTConfig{
			Secret:      getString("JWT_SECRET", "change-me-in-production"),
			ExpireHours: getInt("JWT_EXPIRE_HOURS", 24),
		},
		Redis: RedisConfig{
			Host:     getString("REDIS_HOST", "localhost"),
			Port:     getInt("REDIS_PORT", 6379),
			Password: getString("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
		Log: LogConfig{
			Level:  getString("LOG_LEVEL", "info"),
			Format: getString("LOG_FORMAT", "json"),
		},
	}

	globalConfig = config
	return config, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		panic("config not loaded, call Load() first")
	}
	return globalConfig
}

// GetString 获取字符串配置值（支持环境变量）
func (c *Config) GetString(key string) string {
	return getString(key, "")
}

// 辅助函数
func getString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if value := viper.GetInt(key); value != 0 {
		return value
	}
	return defaultValue
}

func getBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	if value := viper.GetBool(key); value {
		return value
	}
	return defaultValue
}
