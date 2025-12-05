#!/bin/bash
# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Docker 快速启动脚本

set -e

echo "🚀 启动 kk-nav 服务..."

# 检查 Docker 和 Docker Compose
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! command -v docker compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 构建并启动服务
echo "📦 构建 Docker 镜像..."
docker-compose build

echo "🔧 启动服务..."
docker-compose up -d

echo "⏳ 等待数据库就绪..."
sleep 10

echo "🗄️  执行数据库迁移..."
docker-compose exec -T app go run ./backend/scripts/migrate/main.go || echo "⚠️  迁移可能已执行"

echo "🌱 初始化数据..."
docker-compose exec -T app go run ./backend/scripts/seed/main.go || echo "⚠️  数据可能已初始化"

echo ""
echo "✅ 服务启动完成！"
echo ""
echo "📝 访问地址:"
echo "   前台: http://localhost:8080"
echo "   管理后台: http://localhost:8080/admin"
echo "   默认账号: admin@example.com / admin123"
echo ""
echo "📋 常用命令:"
echo "   查看日志: docker-compose logs -f"
echo "   停止服务: docker-compose down"
echo "   重启服务: docker-compose restart"
echo ""

