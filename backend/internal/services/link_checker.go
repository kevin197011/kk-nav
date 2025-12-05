// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package services

import (
	"context"
	"net/http"
	"time"

	"kk-nav/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// LinkChecker 链接状态检测服务
type LinkChecker struct {
	db     *gorm.DB
	logger *zap.Logger
	ticker *time.Ticker
	stop   chan struct{}
}

// NewLinkChecker 创建链接检测服务
func NewLinkChecker(db *gorm.DB, logger *zap.Logger) *LinkChecker {
	return &LinkChecker{
		db:     db,
		logger: logger,
		stop:   make(chan struct{}),
	}
}

// Start 启动定时检测任务（每1小时运行一次）
func (lc *LinkChecker) Start(ctx context.Context) {
	// 立即运行一次
	go lc.runCheck()

	// 每1小时运行一次
	lc.ticker = time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-lc.ticker.C:
				lc.runCheck()
			case <-lc.stop:
				lc.logger.Info("Link checker stopped")
				return
			case <-ctx.Done():
				lc.logger.Info("Link checker context cancelled")
				lc.Stop()
				return
			}
		}
	}()

	lc.logger.Info("Link checker started, will run every 1 hour")
}

// Stop 停止定时检测任务
func (lc *LinkChecker) Stop() {
	if lc.ticker != nil {
		lc.ticker.Stop()
	}
	close(lc.stop)
}

// runCheck 执行链接状态检测
func (lc *LinkChecker) runCheck() {
	lc.logger.Info("Starting link status check job")

	// 只检测状态为 active 或 error 的链接，跳过 inactive（手动禁用）
	var links []models.Link
	if err := lc.db.Where("status IN ?", []string{"active", "error"}).Find(&links).Error; err != nil {
		lc.logger.Error("Failed to fetch links for status check", zap.Error(err))
		return
	}

	if len(links) == 0 {
		lc.logger.Info("No links to check")
		return
	}

	lc.logger.Info("Checking links", zap.Int("count", len(links)))

	now := time.Now()
	checkedCount := 0
	activeCount := 0
	errorCount := 0

	for _, link := range links {
		status := lc.checkLinkStatus(link.URL)

		// 只更新状态为 active 或 error，不改变 inactive
		if link.Status != status {
			link.Status = status
			link.LastCheckedAt = &now
			if err := lc.db.Save(&link).Error; err != nil {
				lc.logger.Error("Failed to update link status",
					zap.Uint("link_id", link.ID),
					zap.String("url", link.URL),
					zap.Error(err))
				continue
			}
			checkedCount++
		} else {
			// 即使状态没变，也更新检测时间
			link.LastCheckedAt = &now
			lc.db.Model(&link).Update("last_checked_at", now)
		}

		if status == "active" {
			activeCount++
		} else {
			errorCount++
		}
	}

	lc.logger.Info("Link status check completed",
		zap.Int("total", len(links)),
		zap.Int("checked", checkedCount),
		zap.Int("active", activeCount),
		zap.Int("error", errorCount))
}

// checkLinkStatus 检测单个链接的状态
func (lc *LinkChecker) checkLinkStatus(url string) string {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()

	// HTTP 状态码 200-399 视为正常
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "active"
	}
	return "error"
}

