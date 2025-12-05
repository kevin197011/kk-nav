#!/bin/bash
# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Docker 初始化脚本
# 用于在容器中执行数据库迁移和初始化数据

set -e

echo "等待数据库连接..."
sleep 5

echo "执行数据库迁移..."
/app/kk-nav migrate || go run ./scripts/migrate/main.go

echo "初始化数据..."
/app/kk-nav seed || go run ./scripts/seed/main.go

echo "初始化完成！"

