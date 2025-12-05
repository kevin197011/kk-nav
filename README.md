# kk-nav - 运维网站智能导航管理站点

[![Backend Build](https://github.com/kevin197011/kk-nav/actions/workflows/docker-publish-backend.yml/badge.svg)](https://github.com/kevin197011/kk-nav/actions/workflows/docker-publish-backend.yml)
[![Frontend Build](https://github.com/kevin197011/kk-nav/actions/workflows/docker-publish-frontend.yml/badge.svg)](https://github.com/kevin197011/kk-nav/actions/workflows/docker-publish-frontend.yml)

一个基于 Golang + React 的高性能运维工具网址导航系统，采用前后端分离架构，提供完整的用户认证、分类管理、链接管理、统计分析等功能。

## ✨ 功能特性

### 🎯 核心功能
- **分类管理**: 支持自定义分类，图标和颜色配置
- **链接管理**: 工具链接的增删改查，支持状态检测
- **标签系统**: 灵活的标签分类和筛选
- **搜索功能**: 全文搜索工具名称、描述和URL
- **收藏功能**: 用户个人收藏夹

### 🔐 用户系统
- **用户认证**: 基于 JWT 的完整认证系统
- **API Token**: 支持创建和管理 API Token，用于程序化访问
- **权限管理**: 管理员和普通用户角色区分
- **个人中心**: 收藏管理和个人设置

### 📊 统计分析
- **访问统计**: 详细的点击统计和分析
- **实时监控**: 链接状态自动检测
- **数据可视化**: 仪表盘展示系统统计

### 🔧 管理后台
- **数据管理**: 分类、链接、标签、用户、Token 管理
- **系统设置**: 站点配置和功能开关
- **统计报表**: 详细的使用统计和分析

## 🛠️ 技术栈

### 后端
- **语言**: Go 1.21+
- **框架**: Gin (RESTful API)
- **数据库**: PostgreSQL / SQLite
- **ORM**: GORM
- **认证**: JWT + API Token
- **日志**: Zap
- **配置**: Viper

### 前端
- **框架**: React 19 + TypeScript
- **构建工具**: Vite
- **UI 组件**: shadcn/ui + Tailwind CSS
- **路由**: React Router
- **状态管理**: Zustand
- **HTTP 客户端**: Axios
- **表单验证**: React Hook Form + Zod

### 部署
- **容器化**: Docker + Docker Compose
- **Web 服务器**: Nginx (前端)

## 🚀 快速开始

### 前置要求
- Docker 20.10+
- Docker Compose 2.0+
- Ruby 2.7+ (用于数据导入脚本，可选)

### 一键启动

```bash
# 1. 克隆项目
git clone <repository-url>
cd kk-nav

# 2. 配置环境变量（可选）
cp .env.docker.example .env
# 编辑 .env 文件，修改管理员账号和密码

# 3. 启动所有服务
docker-compose up -d

# 4. 访问应用
# 前端: http://localhost:3000
# 后端 API (开发环境): http://localhost:8080/api/v1
# API (通过前端代理): http://localhost:3000/api/v1
# 
# 默认管理员账号（使用用户名登录）:
# 用户名: admin
# 密码: admin123
# 
# 注意：首次启动时会自动创建管理员账号和初始化数据
# 可通过环境变量自定义管理员账号（见下方配置说明）

# 5. 导入示例数据（可选）
ruby setup_navigation.rb
# 这将导入 6 个分类和 16 个常用网站链接
# 详见 QUICK_START.md
```

### 生产环境部署

生产环境使用**单端口暴露架构**，只暴露前端端口，后端和数据库端口仅在内部网络使用：

```bash
# 使用生产环境配置
docker compose -f docker-compose.prod.yml up -d

# 访问应用（仅前端端口暴露）
# 前端和 API: http://localhost:3000
# 所有 API 请求通过 Nginx 代理到后端
```

**安全优势**：
- ✅ 只暴露一个端口（3000）
- ✅ 后端 API 不直接暴露
- ✅ 数据库完全隔离
- ✅ 所有请求经过 Nginx 过滤

详细部署指南请参考：[生产环境部署文档](.github/PRODUCTION_DEPLOY.md)

### 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps
```

### 使用预构建镜像

项目提供自动构建的 Docker 镜像，可以直接拉取使用：

```bash
# 拉取最新镜像
docker pull ghcr.io/kevin197011/kk-nav/backend:latest
docker pull ghcr.io/kevin197011/kk-nav/frontend:latest

# 或拉取指定版本
docker pull ghcr.io/kevin197011/kk-nav/backend:v1.0.0
docker pull ghcr.io/kevin197011/kk-nav/frontend:v1.0.0
```

修改 `docker-compose.yml` 使用预构建镜像：

```yaml
services:
  app:
    image: ghcr.io/kevin197011/kk-nav/backend:latest
    # ... 其他配置保持不变

  frontend:
    image: ghcr.io/kevin197011/kk-nav/frontend:latest
    # ... 其他配置保持不变
```

然后启动服务：

```bash
docker-compose up -d
```

详细说明请参考：[Docker 镜像自动构建文档](.github/workflows/README.md)

## 📁 项目结构

```
kk-nav/
├── backend/                 # 后端代码
│   ├── cmd/server/         # 应用入口
│   ├── internal/           # 内部代码
│   │   ├── config/         # 配置管理
│   │   ├── database/       # 数据库连接
│   │   ├── models/         # 数据模型
│   │   ├── handlers/       # HTTP处理器
│   │   │   └── admin/      # 管理后台处理器
│   │   ├── middleware/     # 中间件
│   │   └── utils/          # 工具函数
│   ├── scripts/            # 脚本文件
│   │   ├── migrate/        # 数据库迁移
│   │   └── seed/           # 数据初始化
│   ├── configs/            # 配置文件
│   ├── Dockerfile          # Docker镜像构建
│   └── go.mod              # Go依赖
├── frontend/               # 前端代码（React）
│   ├── src/                # 源代码
│   │   ├── pages/          # 页面组件
│   │   │   └── admin/      # 管理后台页面
│   │   ├── components/     # 通用组件
│   │   ├── stores/         # 状态管理
│   │   ├── lib/            # 工具库
│   │   └── types/          # 类型定义
│   ├── dist/               # 构建产物
│   ├── Dockerfile          # 生产环境 Dockerfile
│   ├── Dockerfile.dev      # 开发环境 Dockerfile
│   ├── nginx.conf          # Nginx 配置
│   └── package.json        # 依赖配置
├── setup_navigation.rb     # Ruby 数据导入脚本（推荐）
├── navigation-data.json    # 导航数据配置文件
├── check_links_status.rb   # 链接状态检查脚本
├── docker-compose.yml      # Docker Compose 开发环境配置
├── docker-compose.prod.yml # Docker Compose 生产环境配置
├── .env.docker.example     # 环境变量示例
└── README.md              # 项目说明
```

## 🔌 API 文档

### 认证相关
```
POST   /api/v1/auth/login          # 用户登录（使用用户名 + 密码）
GET    /api/v1/auth/me             # 获取当前用户信息
POST   /api/v1/auth/logout         # 用户登出
```

**注意**: 系统不开放用户注册，所有用户由管理员在后台创建。

**登录请求示例**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

### 前台 API（公开）
```
GET    /api/v1/categories          # 分类列表
GET    /api/v1/categories/:id      # 分类详情
GET    /api/v1/links               # 链接列表
GET    /api/v1/links/:id           # 链接详情
POST   /api/v1/links/:id/click     # 记录点击
GET    /api/v1/tags                # 标签列表
GET    /api/v1/tags/:id            # 标签详情
GET    /api/v1/stats                # 统计数据
```

### 用户 API（需要认证）
```
POST   /api/v1/links/:id/favorite  # 收藏链接
DELETE /api/v1/links/:id/unfavorite # 取消收藏
GET    /api/v1/favorites           # 获取收藏列表
```

### 管理后台 API（需要管理员权限）
```
# 仪表盘
GET    /api/v1/admin/dashboard     # 仪表盘数据

# 分类管理
GET    /api/v1/admin/categories    # 分类列表
POST   /api/v1/admin/categories    # 创建分类
GET    /api/v1/admin/categories/:id # 分类详情
PUT    /api/v1/admin/categories/:id # 更新分类
DELETE /api/v1/admin/categories/:id # 删除分类
PATCH  /api/v1/admin/categories/:id/move-up   # 上移
PATCH  /api/v1/admin/categories/:id/move-down # 下移

# 链接管理
GET    /api/v1/admin/links         # 链接列表
POST   /api/v1/admin/links         # 创建链接
GET    /api/v1/admin/links/:id     # 链接详情
PUT    /api/v1/admin/links/:id     # 更新链接
DELETE /api/v1/admin/links/:id     # 删除链接
POST   /api/v1/admin/links/:id/check-status    # 检测链接状态
POST   /api/v1/admin/links/batch-check         # 批量检测

# 标签管理
GET    /api/v1/admin/tags          # 标签列表
POST   /api/v1/admin/tags          # 创建标签
GET    /api/v1/admin/tags/:id      # 标签详情
PUT    /api/v1/admin/tags/:id      # 更新标签
DELETE /api/v1/admin/tags/:id      # 删除标签

# 用户管理
GET    /api/v1/admin/users         # 用户列表
POST   /api/v1/admin/users         # 创建用户
PUT    /api/v1/admin/users/:id     # 更新用户
DELETE /api/v1/admin/users/:id     # 删除用户

# Token 管理
GET    /api/v1/admin/tokens        # Token 列表
POST   /api/v1/admin/tokens        # 创建 Token
GET    /api/v1/admin/tokens/:id    # Token 详情
PUT    /api/v1/admin/tokens/:id    # 更新 Token
DELETE /api/v1/admin/tokens/:id    # 删除 Token

# 系统设置
GET    /api/v1/admin/settings      # 获取设置
PUT    /api/v1/admin/settings      # 更新设置
```

### API 响应格式

所有 API 响应都遵循统一格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 具体数据
  }
}
```

- `code: 0` 表示成功
- `code: 非0` 表示错误
- `message` 为响应消息
- `data` 为响应数据

## 🔑 API Token 使用

系统支持两种认证方式：

1. **JWT Token**: 通过登录获取，有过期时间
2. **API Token**: 通过管理后台创建，可以设置过期时间或永不过期

### 使用方式

在 HTTP 请求头中添加 Token：

```
Authorization: Bearer <your-token>
```

### 示例

```bash
# 使用 API Token 调用 API
curl -X GET http://localhost:8080/api/v1/links \
  -H "Authorization: Bearer kk_xxxxxxxxxxxxx"

# 创建链接（需要管理员权限）
curl -X POST http://localhost:8080/api/v1/admin/links \
  -H "Authorization: Bearer kk_xxxxxxxxxxxxx" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "新链接",
    "url": "https://example.com",
    "category_id": 1,
    "status": "active"
  }'
```

### Token 管理

在管理后台的"Token 管理"页面可以创建和管理 API Token：

1. 登录管理后台（需要 admin 用户）
2. 进入"Token 管理"
3. 点击"新建 Token"，填写名称和关联用户
4. Token 创建后只会显示一次，请妥善保管
5. 使用 Token 时在请求头中添加：`Authorization: Bearer <your-token>`

## ⚙️ 配置说明

### 环境变量

创建 `.env` 文件（可从 `.env.docker.example` 复制）：

```bash
# 数据库配置
DB_PASSWORD=your-secure-db-password

# JWT配置（生产环境必须修改）
JWT_SECRET=your-super-secret-jwt-key-at-least-32-characters

# 默认管理员配置（首次启动时创建）
ADMIN_EMAIL=admin@yourcompany.com
ADMIN_USERNAME=superadmin
ADMIN_PASSWORD=YourSecurePassword123!
```

完整环境变量列表：

```bash
# 应用配置
APP_NAME=kk-nav
APP_ENV=production
APP_PORT=8080
APP_DEBUG=false

# 数据库配置
DB_TYPE=postgres
DB_HOST=postgres
DB_PORT=5432
DB_NAME=kk_nav
DB_USER=kk_nav
DB_PASSWORD=kk_nav_password
DB_SSL_MODE=disable

# JWT配置
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRE_HOURS=24

# 默认管理员配置（首次启动时自动创建）
ADMIN_EMAIL=admin@example.com
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
```

### 端口配置

- **前端**: 3000
- **后端 API**: 8080
- **PostgreSQL**: 5432
- **Redis**: 6379

## 🧪 开发

### 本地开发（后端）

```bash
cd backend

# 安装依赖
go mod download

# 运行迁移
go run ./scripts/migrate/main.go

# 初始化数据
go run ./scripts/seed/main.go

# 启动服务
go run ./cmd/server/main.go
```

### 本地开发（前端）

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 📝 默认账号

- **管理员**: 通过环境变量配置（默认用户名 `admin` / 密码 `admin123`）
- 首次运行会自动创建
- **登录方式**: 使用用户名 + 密码登录（不是邮箱）
- **用户注册**: 已禁用，所有用户由管理员在后台创建
- 可通过 `.env` 文件或环境变量自定义：
  ```bash
  ADMIN_EMAIL=admin@yourcompany.com
  ADMIN_USERNAME=youradmin
  ADMIN_PASSWORD=your-secure-password
  ```

## 🛠️ 实用工具

### 数据导入工具

项目提供了便捷的 Ruby 脚本来快速导入导航数据：

```bash
# 使用 Ruby 脚本导入数据（推荐）
ruby setup_navigation.rb

# 检查链接状态
ruby check_links_status.rb
```

**setup_navigation.rb** 特点：
- ✅ 一个文件包含所有配置
- ✅ 自动检测分类是否存在
- ✅ 支持批量导入链接
- ✅ 详细的执行日志
- ✅ 易于自定义和扩展

**配置示例**：
```ruby
CONFIG = {
  api_url: 'https://your-domain.com/api/v1',
  token: 'your-api-token',
  data: [
    {
      category: {
        name: '分类名称',
        icon: '🎯',
        color: '#3b82f6',
        description: '分类描述'
      },
      links: [
        {
          title: '链接标题',
          url: 'https://example.com',
          description: '链接描述'
        }
      ]
    }
  ]
}
```

详细使用说明请参考：
- [QUICK_START.md](QUICK_START.md) - 快速开始指南
- [PRODUCTION_ADD_LINKS.md](PRODUCTION_ADD_LINKS.md) - 生产环境添加链接

### 诊断工具

```bash
# 检查链接状态和数量
ruby check_links_status.rb
```

输出示例：
```
前台链接（用户可见）: 17 个
后台链接（管理员可见）: 17 个
统计数据: 17 个

按分类统计:
  基础架构: 2 个
  开发工具: 3 个
  云服务: 3 个
  ...
```

## 🔒 安全建议

1. **生产环境必须修改**:
   - `JWT_SECRET`: 使用强随机密钥
   - `DB_PASSWORD`: 使用强密码
   - `ADMIN_EMAIL`: 使用真实管理员邮箱
   - `ADMIN_PASSWORD`: 使用强密码（建议 16 位以上）
   - 所有默认密码

2. **使用 HTTPS**: 生产环境建议使用 HTTPS 传输 Token

3. **Token 管理**: 定期轮换 API Token，及时删除不再使用的 Token

4. **环境变量管理**: 
   - 复制 `.env.docker.example` 为 `.env` 并修改敏感信息
   - 不要将 `.env` 文件提交到版本控制系统

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## � 支相关文档

- [QUICK_START.md](QUICK_START.md) - 快速开始和数据导入指南
- [PRODUCTION_ADD_LINKS.md](PRODUCTION_ADD_LINKS.md) - 生产环境添加链接
- [REMOVE_REGISTRATION.md](REMOVE_REGISTRATION.md) - 注册功能移除说明
- [FIX_SORTING.md](FIX_SORTING.md) - 分类排序功能修复
- [FIX_LINKS_PAGINATION.md](FIX_LINKS_PAGINATION.md) - 链接分页问题修复
- [BUILD_FIX.md](BUILD_FIX.md) - 前端构建错误修复
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - 登录方式变更说明

## 🔄 更新日志

### v1.1.0 (2025-12-05)

**新增功能**:
- ✅ 添加 Ruby 数据导入脚本 (`setup_navigation.rb`)
- ✅ 添加链接状态检查工具 (`check_links_status.rb`)
- ✅ 支持通过环境变量配置默认管理员账号
- ✅ 自动初始化管理员账号和系统设置

**功能改进**:
- ✅ 移除用户注册功能，改为管理员创建用户
- ✅ 登录方式从邮箱改为用户名
- ✅ 修复分类排序功能（上移/下移）
- ✅ 修复前台链接显示数量问题（分页大小调整）
- ✅ 优化前端构建流程

**安全增强**:
- ✅ 生产环境单端口暴露架构
- ✅ 后端 API 不直接暴露
- ✅ 数据库完全隔离
- ✅ 支持 API Token 认证

**文档更新**:
- ✅ 完善 README 文档
- ✅ 添加多个使用指南和故障排除文档
- ✅ 添加 Ruby 脚本使用说明

### v1.0.0 (2025-11-01)

- 🎉 初始版本发布
- ✅ 基础功能实现
- ✅ Docker 容器化部署
- ✅ 前后端分离架构

## 📞 支持

如果您在使用过程中遇到问题，请：

1. 查看本文档了解详细说明
2. 查看相关文档（见上方"相关文档"部分）
3. 运行诊断工具检查系统状态
4. 查看 [Issues](../../issues) 中是否有相似问题
5. 创建新的 Issue 描述问题

## 💡 常见问题

### Q: 前台显示的链接数量与后台不一致？

A: 前台只显示 `status=active` 的链接，后台显示所有状态的链接。运行 `ruby check_links_status.rb` 检查链接状态。

### Q: 如何添加新的导航链接？

A: 有三种方式：
1. 在管理后台手动添加
2. 使用 Ruby 脚本批量导入：`ruby setup_navigation.rb`
3. 使用 API 接口添加

### Q: 如何创建新用户？

A: 系统已禁用注册功能，管理员可以在后台"用户管理"中创建新用户。

### Q: 忘记管理员密码怎么办？

A: 可以通过以下方式重置：
1. 修改 `.env` 文件中的 `ADMIN_PASSWORD`
2. 重启服务：`docker compose restart app`
3. 使用新密码登录

### Q: 如何备份数据？

A: 备份 PostgreSQL 数据库：
```bash
docker exec kk-nav-postgres pg_dump -U kk_nav kk_nav > backup.sql
```

恢复数据：
```bash
docker exec -i kk-nav-postgres psql -U kk_nav kk_nav < backup.sql
```
