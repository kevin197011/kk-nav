# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

.PHONY: install dev build preview docker-build docker-run docker-dev help

# 安装依赖
install:
	@echo "Installing dependencies..."
	@npm install

# 开发模式
dev:
	@echo "Starting development server..."
	@npm run dev

# 构建生产版本
build:
	@echo "Building for production..."
	@npm run build

# 预览构建结果
preview:
	@echo "Previewing build..."
	@npm run preview

# Docker 构建（生产）
docker-build:
	@echo "Building Docker image..."
	@docker build -t kk-nav-frontend -f Dockerfile .

# Docker 运行（生产）
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 3000:80 --rm kk-nav-frontend

# Docker 开发模式
docker-dev:
	@echo "Starting Docker development container..."
	@docker build -t kk-nav-frontend-dev -f Dockerfile.dev .
	@docker run -p 5173:5173 -v $(PWD):/app --rm kk-nav-frontend-dev

# 帮助
help:
	@echo "Available commands:"
	@echo "  make install      - Install dependencies"
	@echo "  make dev          - Start development server"
	@echo "  make build        - Build for production"
	@echo "  make preview      - Preview production build"
	@echo "  make docker-build - Build Docker image (production)"
	@echo "  make docker-run   - Run Docker container (production)"
	@echo "  make docker-dev   - Run Docker container (development)"
