# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AIFlowy is an enterprise-grade, open-source AI application development platform with:
- **Backend (Go)**: Go 1.21+ + Echo + Eino (aiflowy-go/) - **主要开发**
- **Backend (Java)**: Java 8 + Spring Boot 2.7.18 + Maven - 遗留代码
- **Frontend**: Vue 3 + TypeScript monorepo (pnpm workspaces)
- **Database**: MySQL 8.0

Main branch is development; v1.x branch is stable.

## Common Commands

### Go 后端服务管理 (推荐)

Go 后端运行在端口 8213，使用 Makefile 管理：

```bash
# 服务管理
make start                           # 构建并启动服务 (后台)
make stop                            # 停止服务
make restart                         # 重启服务
make status                          # 查看服务状态 (PID/端口/内存)
make logs                            # 查看实时日志
make dev                             # 前台运行 (开发模式)

# 构建
make go-build                        # 构建 Go 后端
make go-test                         # 运行测试
make go-fmt                          # 格式化代码
make go-tidy                         # 整理依赖

# 版本
make version                         # 显示版本信息
make version-bump-patch              # 1.0.0 -> 1.0.1
```

### API 测试工具 (推荐)

使用 `scripts/api` 脚本测试 API，自动管理 Token：

```bash
# 登录获取 Token (保存到 /tmp/.aiflowy.token，有效期 24 小时)
make login                           # 使用默认账号 admin/123456
./scripts/api login admin 123456     # 指定账号密码

# 常用 API 调用
./scripts/api GET bot/list                              # Bot 列表
./scripts/api GET auth/getUserInfo                      # 当前用户信息
./scripts/api GET bot/getDetail id=263949328253587456   # Bot 详情
./scripts/api POST botCategory/save -d '{"categoryName":"测试分类"}'

# Bot 聊天
./scripts/api POST bot/chat -d '{"botId":"263949328253587456","message":"你好","stream":false}'

# 高级用法
./scripts/api GET bot/list -f data.0.title              # 提取字段
./scripts/api GET bot/list -v                           # 详细模式
./scripts/api GET bot/list -r                           # 原始输出
```

> **Token 文件**: `/tmp/.aiflowy.token`
> **默认账号**: admin / 123456
> **API 基础地址**: http://localhost:8213/api/v1

### Backend (Maven) - 遗留
```bash
mvn clean package                    # Build all modules
mvn clean install                    # Install to local repo
```

### Frontend (aiflowy-ui-admin)
```bash
# 使用 Makefile (推荐)
make ui-install                      # 安装依赖
make ui-dev                          # 运行开发服务器
make ui-build                        # 构建生产版本
make ui-lint                         # 代码检查

# 或直接使用 pnpm
cd aiflowy-ui-admin
pnpm install                         # Install dependencies
pnpm dev                             # Dev server (proxies to localhost:5320)
pnpm build                           # Production build
pnpm lint                            # Run ESLint
pnpm check                           # All checks (circular deps, types, spell)
```

### Database Setup
```bash
# DDL: sql/aiflowy-v2.ddl.sql
# Seed data: sql/aiflowy-v2.data.sql
# Default credentials: admin/123456
```

### Docker
```bash
docker-compose up                    # Full stack (Java + MySQL)
pnpm build:docker                    # Build Docker image
```

## Architecture

### Go 后端结构 (aiflowy-go/) - 主要开发
```
aiflowy-go/
├── cmd/server/           # 入口
├── configs/              # 配置文件
├── internal/
│   ├── config/           # 配置加载
│   ├── dto/              # 数据传输对象
│   ├── entity/           # 实体定义
│   ├── errors/           # 错误处理
│   ├── handler/          # HTTP 处理器
│   │   ├── auth/         # 认证
│   │   ├── bot/          # Bot 管理 + 聊天
│   │   ├── model/        # 模型管理
│   │   └── system/       # 系统管理
│   ├── middleware/       # 中间件 (JWT, CORS, Logger)
│   ├── repository/       # 数据访问层
│   ├── router/           # 路由配置
│   └── service/          # 业务逻辑
│       └── llm/          # LLM 服务 (Eino 框架)
├── pkg/
│   ├── protocol/         # AIFlowy Chat Protocol v1.1
│   ├── response/         # 响应封装
│   └── snowflake/        # ID 生成
└── build/                # 构建产物
```

Key frameworks:
- **Echo**: Web 框架
- **Eino (CloudWeGo)**: AI/LLM 框架
- **JWT**: 认证
- **MySQL**: 数据库

### Java 后端结构 (遗留)
```
aiflowy-api/           # API layer (admin, public, usercenter, mcp)
aiflowy-commons/       # Shared utilities (ai, cache, web, satoken, file-storage)
aiflowy-modules/       # Core business (ai, auth, system, core, datacenter, job, log)
aiflowy-starter/       # Application starters (admin, all, public, usercenter, codegen)
```

Key frameworks:
- **MyBatis-Flex**: ORM layer
- **AgentsFlex 2.0.0-beta.5**: AI agent framework
- **TinyFlow-Java 2.0.0-beta.6**: Workflow orchestration
- **Sa-Token**: Authentication & authorization

### Frontend Monorepo (aiflowy-ui-admin)
```
app/                   # Main admin application
  src/api/             # API client layer
  src/views/           # Page components
  src/store/           # Pinia stores
  src/router/          # Vue Router config
packages/
  @core/ui-kit/        # Shared UI components (shadcn-ui, form-ui, layout-ui)
  @core/composables/   # Reusable Vue composables
  effects/request/     # HTTP client layer
  stores/              # Global Pinia stores
```

UI Libraries: Element Plus + shadcn/ui

### Multiple UIs
- `aiflowy-ui-admin`: Admin dashboard
- `aiflowy-ui-usercenter`: User portal
- `aiflowy-ui-websdk`: Embeddable web SDK

## Chat Protocol

AIFlowy uses a custom SSE-based chat protocol (`aiflowy-chat` v1.1) for streaming AI responses. See `/aiflowy-chat-protocol.md` for full specification.

Protocol domains: `llm`, `tool`, `system`, `business`, `workflow`, `interaction`, `debug`

## Key Patterns

- Backend follows modular monolith pattern with clear separation (API/Commons/Modules/Starters)
- Frontend uses pnpm workspace catalogs for shared dependency versions
- API endpoints: `/api/*` routes proxied to backend in dev mode
- Entity naming convention: `tb_` prefix for database tables

## Development Notes

- Node >= 20.10.0, pnpm 10.14.0 required for frontend
- Java 8 target for backend
- Git hooks configured via lefthook.yml
- Commit format: czg (commitizen)
