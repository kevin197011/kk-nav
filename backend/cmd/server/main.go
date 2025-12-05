// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kk-nav/internal/config"
	"kk-nav/internal/database"
	"kk-nav/internal/handlers"
	adminHandlers "kk-nav/internal/handlers/admin"
	"kk-nav/internal/middleware"
	"kk-nav/internal/services"
	"kk-nav/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	var logger *zap.Logger
	if cfg.Log.Format == "json" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// 初始化JWT
	utils.InitJWT(cfg.JWT.Secret)

	// 连接数据库
	if err := database.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}
	defer database.Close()

	// 自动迁移
	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// 初始化默认数据（管理员账号和系统设置）
	if err := database.InitializeData(cfg); err != nil {
		logger.Fatal("Failed to initialize data", zap.Error(err))
	}

	// 设置Gin模式
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin引擎
	r := gin.New()

	// 中间件
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	// 注册路由（将在后续实现）
	registerRoutes(r, logger)

	// 启动链接状态检测服务
	linkChecker := services.NewLinkChecker(database.DB, logger)
	checkerCtx, checkerCancel := context.WithCancel(context.Background())
	defer checkerCancel()
	linkChecker.Start(checkerCtx)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		logger.Info("Server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 停止链接检测服务
	linkChecker.Stop()
	checkerCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// registerRoutes 注册路由
func registerRoutes(r *gin.Engine, logger *zap.Logger) {
	cfg := config.Get()
	db := database.DB

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(db, cfg)
	linksHandler := handlers.NewLinksHandler(db)
	categoriesHandler := handlers.NewCategoriesHandler(db)
	tagsHandler := handlers.NewTagsHandler(db)
	statsHandler := handlers.NewStatsHandler(db)
	settingsHandler := handlers.NewSettingsHandler(db)

	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关（不需要认证）
		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.Me)
		}

		// 前台API（不需要认证）
		apiV1.GET("/categories", categoriesHandler.Index)
		apiV1.GET("/categories/:id", categoriesHandler.Show)
		apiV1.GET("/links", linksHandler.Index)
		apiV1.GET("/links/:id", linksHandler.Show)
		apiV1.POST("/links/:id/click", linksHandler.Click)
		apiV1.GET("/tags", tagsHandler.Index)
		apiV1.GET("/tags/:id", tagsHandler.Show)
		apiV1.GET("/stats", statsHandler.Index)
		apiV1.GET("/settings", settingsHandler.GetPublicSettings)

		// 用户相关（需要认证）
		user := apiV1.Group("", middleware.AuthMiddleware())
		{
			user.POST("/links/:id/favorite", linksHandler.Favorite)
			user.DELETE("/links/:id/unfavorite", linksHandler.Unfavorite)
			user.GET("/favorites", linksHandler.Favorites)
		}
	}

	// 管理后台API（需要管理员权限）
	admin := r.Group("/api/v1/admin", middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 初始化管理后台处理器
		adminDashboardHandler := adminHandlers.NewDashboardHandler(db)
		adminCategoriesHandler := adminHandlers.NewCategoriesHandler(db)
		adminLinksHandler := adminHandlers.NewLinksHandler(db)
		adminTagsHandler := adminHandlers.NewTagsHandler(db)
		adminUsersHandler := adminHandlers.NewUsersHandler(db)
		adminSettingsHandler := adminHandlers.NewSettingsHandler(db)
		adminTokensHandler := adminHandlers.NewTokensHandler(db)

		// 管理后台首页
		admin.GET("/dashboard", adminDashboardHandler.Index)

		// 分类管理
		admin.GET("/categories", adminCategoriesHandler.Index)
		admin.POST("/categories", adminCategoriesHandler.Create)
		admin.GET("/categories/:id", adminCategoriesHandler.Show)
		admin.PUT("/categories/:id", adminCategoriesHandler.Update)
		admin.DELETE("/categories/:id", adminCategoriesHandler.Delete)
		admin.PATCH("/categories/:id/move-up", adminCategoriesHandler.MoveUp)
		admin.PATCH("/categories/:id/move-down", adminCategoriesHandler.MoveDown)

		// 链接管理
		admin.GET("/links", adminLinksHandler.Index)
		admin.POST("/links", adminLinksHandler.Create)
		admin.GET("/links/:id", adminLinksHandler.Show)
		admin.PUT("/links/:id", adminLinksHandler.Update)
		admin.DELETE("/links/:id", adminLinksHandler.Delete)
		admin.POST("/links/:id/check-status", adminLinksHandler.CheckStatus)
		admin.POST("/links/batch-check", adminLinksHandler.BatchCheckStatus)
		admin.PATCH("/links/:id/move-up", adminLinksHandler.MoveUp)
		admin.PATCH("/links/:id/move-down", adminLinksHandler.MoveDown)

		// 标签管理
		admin.GET("/tags", adminTagsHandler.Index)
		admin.POST("/tags", adminTagsHandler.Create)
		admin.GET("/tags/:id", adminTagsHandler.Show)
		admin.PUT("/tags/:id", adminTagsHandler.Update)
		admin.DELETE("/tags/:id", adminTagsHandler.Delete)

		// 用户管理
		admin.GET("/users", adminUsersHandler.Index)
		admin.POST("/users", adminUsersHandler.Create)
		admin.PUT("/users/:id", adminUsersHandler.Update)
		admin.DELETE("/users/:id", adminUsersHandler.Delete)

		// 系统设置
		admin.GET("/settings", adminSettingsHandler.Index)
		admin.PUT("/settings", adminSettingsHandler.Update)

		// Token 管理
		admin.GET("/tokens", adminTokensHandler.Index)
		admin.POST("/tokens", adminTokensHandler.Create)
		admin.GET("/tokens/:id", adminTokensHandler.Show)
		admin.PUT("/tokens/:id", adminTokensHandler.Update)
		admin.DELETE("/tokens/:id", adminTokensHandler.Delete)
	}

	// 静态文件服务（前端资源）
	// 在生产环境中，前端通常由独立的 Nginx 服务提供
	// 这里提供静态文件服务用于开发环境或单容器部署
	// 前端由独立的 Nginx 服务提供，后端只提供 API
	// 所有非 API 路由返回 404
	r.NoRoute(func(c *gin.Context) {
		// 如果是 API 请求，返回 404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "API endpoint not found",
			})
			return
		}
		// 其他请求返回 404（前端由独立服务提供）
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "Not found",
		})
	})
}
