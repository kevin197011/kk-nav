.PHONY: build run test clean migrate seed help

# 变量
APP_NAME=kk-nav
BIN_DIR=bin
CMD_DIR=cmd/server

# 构建
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"

# 运行
run:
	@go run ./$(CMD_DIR)

# 测试
test:
	@go test -v ./...

# 测试覆盖率
test-coverage:
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# 清理
clean:
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# 数据库迁移
migrate:
	@go run scripts/migrate/main.go

# 初始化数据
seed:
	@go run scripts/seed/main.go

# 格式化代码
fmt:
	@go fmt ./...

# 代码检查
lint:
	@golangci-lint run

# 安装依赖
deps:
	@go mod download
	@go mod tidy

# Docker 相关
docker-build:
	@echo "Building Docker image..."
	@docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

docker-restart:
	@docker-compose restart

docker-ps:
	@docker-compose ps

docker-exec:
	@docker-compose exec app sh

docker-dev:
	@echo "Starting development environment..."
	@docker-compose -f docker-compose.dev.yml up

docker-init:
	@echo "Initializing database..."
	@docker-compose exec -T app go run ./scripts/migrate/main.go || true
	@docker-compose exec -T app go run ./scripts/seed/main.go || true

docker-start:
	@./scripts/docker-start.sh

# 帮助
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  migrate        - Run database migrations"
	@echo "  seed           - Seed initial data"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  deps           - Download and tidy dependencies"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-logs    - Show container logs"
	@echo "  docker-restart - Restart containers"
	@echo "  docker-ps      - Show running containers"
	@echo "  docker-exec    - Execute command in app container"
	@echo "  docker-dev     - Start development environment"
	@echo "  docker-init    - Initialize database in container"
	@echo "  help           - Show this help message"

