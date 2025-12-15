# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安装构建依赖（包括 Node.js 用于前端构建）
RUN apk add --no-cache git nodejs npm

# 复制后端代码
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制后端源代码
COPY backend/ .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/kk-nav ./cmd/server

# 注意：前端现在有独立的 Dockerfile，可以在 docker-compose 中单独构建
# 如果需要在这里构建，取消下面的注释：
# WORKDIR /app/frontend
# COPY frontend/package*.json ./
# RUN npm ci
# COPY frontend/ .
# RUN npm run build
# WORKDIR /app

# 运行阶段
FROM alpine:3.19

RUN apk --no-cache --no-scripts add ca-certificates tzdata && \
    update-ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/kk-nav .
# 复制配置文件
COPY --from=builder /app/configs ./configs
# 注意：前端构建产物由前端服务提供，或通过 volume 挂载

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./kk-nav"]

