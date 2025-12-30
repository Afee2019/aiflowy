# ============================================================================
# AIFlowy - Makefile
# ============================================================================
# AIFlowy 企业级 AI 应用开发平台
# 项目组件:
#   - aiflowy-go/           Go 后端服务
#   - aiflowy-ui-admin/     Vue3 管理后台
#   - aiflowy-ui-usercenter/ Vue3 用户中心
# ============================================================================

.PHONY: all build clean test help
.PHONY: go-build go-run go-test go-lint go-fmt go-vet go-deps go-tidy
.PHONY: start stop restart status logs logs-tail clean-logs dev
.PHONY: ui-install ui-dev ui-build ui-lint ui-check
.PHONY: version version-set version-bump-patch version-bump-minor version-bump-major
.PHONY: login api

# ============================================================================
# 变量定义
# ============================================================================

# 版本信息
VERSION := $(shell cat VERSION 2>/dev/null || echo "0.0.0")
BUILD_TIME := $(shell date '+%Y-%m-%d_%H:%M:%S')
COMMIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

# Go 配置
GO := go
GO_SRC_DIR := aiflowy-go
GO_BUILD_DIR := $(GO_SRC_DIR)/build
GO_CONFIG_FILE := $(GO_SRC_DIR)/configs/config.yaml
CGO_ENABLED := 0
LDFLAGS := -ldflags="-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.CommitSHA=$(COMMIT_SHA)' -s -w"

# 前端配置
UI_ADMIN_DIR := aiflowy-ui-admin

# 端口配置
SERVER_PORT := 8213
UI_PORT := 8212

# 日志配置
LOG_DIR := $(GO_SRC_DIR)/logs
LOG_FILE := $(LOG_DIR)/server.log
PID_FILE := $(LOG_DIR)/server.pid

# API 测试配置
TOKEN_FILE := /tmp/.aiflowy.token
API_USER ?= admin
API_PASS ?= 123456

# ============================================================================
# 默认目标
# ============================================================================

all: go-build

help: ## 显示帮助信息
	@echo "AIFlowy Make 命令"
	@echo "=================================================="
	@echo "版本: $(VERSION)"
	@echo ""
	@echo "🚀 快捷命令："
	@grep -E '^(help|dev|start|stop|restart|status|clean|version):.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🔧 服务管理："
	@grep -E '^(logs|logs-tail|clean-logs):.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🏗️ Go 后端："
	@grep -E '^go-[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[32m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🎨 前端 (Vue)："
	@grep -E '^ui-[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[34m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "🔑 API 测试："
	@grep -E '^(login|api):.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "⚙️ 版本管理："
	@grep -E '^version-[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[35m%-25s\033[0m %s\n", $$1, $$2}'

# ============================================================================
# 版本管理
# ============================================================================

version: ## 显示版本信息
	@echo "Version:    $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Commit:     $(COMMIT_SHA)"
	@echo "Branch:     $(BRANCH)"

version-set: ## 设置版本 (make version-set V=x.y.z)
	@if [ -z "$(V)" ]; then echo "Usage: make version-set V=x.y.z"; exit 1; fi
	@echo "$(V)" > VERSION
	@echo "版本已设置为: $(V)"

version-bump-patch: ## 升级补丁版本 (x.y.Z)
	@current=$$(cat VERSION); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	patch=$$(echo $$current | cut -d. -f3); \
	new="$$major.$$minor.$$((patch + 1))"; \
	echo "$$new" > VERSION; \
	echo "版本已升级: $$current -> $$new"

version-bump-minor: ## 升级次版本 (x.Y.0)
	@current=$$(cat VERSION); \
	major=$$(echo $$current | cut -d. -f1); \
	minor=$$(echo $$current | cut -d. -f2); \
	new="$$major.$$((minor + 1)).0"; \
	echo "$$new" > VERSION; \
	echo "版本已升级: $$current -> $$new"

version-bump-major: ## 升级主版本 (X.0.0)
	@current=$$(cat VERSION); \
	major=$$(echo $$current | cut -d. -f1); \
	new="$$((major + 1)).0.0"; \
	echo "$$new" > VERSION; \
	echo "版本已升级: $$current -> $$new"

# ============================================================================
# Go 后端
# ============================================================================

go-build: ## 构建 Go 后端服务
	@echo "正在构建 Go 后端..."
	@mkdir -p $(GO_BUILD_DIR)
	cd $(GO_SRC_DIR) && CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(LDFLAGS) -o build/aiflowy-go ./cmd/server
	@echo "构建完成: $(GO_BUILD_DIR)/aiflowy-go"

go-run: ## 运行 Go 后端 (前台)
	cd $(GO_SRC_DIR) && $(GO) run ./cmd/server -config configs/config.yaml

go-test: ## 运行 Go 测试
	cd $(GO_SRC_DIR) && $(GO) test -v ./...

go-lint: ## 运行 Go 代码检查
	@which golangci-lint > /dev/null || (echo "golangci-lint 未安装，跳过..." && exit 0)
	cd $(GO_SRC_DIR) && golangci-lint run ./...

go-fmt: ## 格式化 Go 代码
	cd $(GO_SRC_DIR) && $(GO) fmt ./...

go-vet: ## 运行 go vet
	cd $(GO_SRC_DIR) && $(GO) vet ./...

go-tidy: ## 整理 Go 依赖
	cd $(GO_SRC_DIR) && $(GO) mod tidy -v

go-deps: ## 更新 Go 依赖
	cd $(GO_SRC_DIR) && $(GO) mod download
	cd $(GO_SRC_DIR) && $(GO) mod tidy

# ============================================================================
# 服务管理 (前端 + 后端)
# ============================================================================

# 前端日志配置
UI_LOG_DIR := $(UI_ADMIN_DIR)/logs
UI_LOG_FILE := $(UI_LOG_DIR)/ui.log
UI_PID_FILE := $(UI_LOG_DIR)/ui.pid

start: go-build ## 启动服务 (后端 + 前端)
	@START_TIME=$$(date '+%Y-%m-%d %H:%M:%S'); \
	echo "=== 启动 AIFlowy 服务 ==="; \
	echo ""; \
	echo "[1/2] 启动 Go 后端 (端口 $(SERVER_PORT))..."; \
	if lsof -ti:$(SERVER_PORT) >/dev/null 2>&1; then \
		echo "  ⚠ 后端已在运行"; \
	else \
		mkdir -p $(LOG_DIR); \
		nohup $(GO_BUILD_DIR)/aiflowy-go -config $(GO_CONFIG_FILE) > $(LOG_FILE) 2>&1 & echo $$! > $(PID_FILE); \
		sleep 2; \
		if lsof -ti:$(SERVER_PORT) >/dev/null 2>&1; then \
			echo "  ✓ 后端已启动 (PID: $$(cat $(PID_FILE)))"; \
		else \
			echo "  ✗ 后端启动失败"; \
			tail -10 $(LOG_FILE) 2>/dev/null; \
			exit 1; \
		fi; \
	fi; \
	echo ""; \
	echo "[2/2] 启动前端 (端口 $(UI_PORT))..."; \
	if lsof -ti:$(UI_PORT) >/dev/null 2>&1; then \
		echo "  ⚠ 前端已在运行"; \
	else \
		mkdir -p $(UI_LOG_DIR); \
		cd $(UI_ADMIN_DIR) && nohup pnpm dev > logs/ui.log 2>&1 & echo $$! > logs/ui.pid; \
		sleep 3; \
		if lsof -ti:$(UI_PORT) >/dev/null 2>&1; then \
			echo "  ✓ 前端已启动"; \
		else \
			echo "  ✗ 前端启动失败"; \
			tail -10 $(UI_LOG_FILE) 2>/dev/null; \
			exit 1; \
		fi; \
	fi; \
	echo ""; \
	echo "========================================"; \
	echo "  ✓ AIFlowy 启动完成"; \
	echo "  启动时间: $$START_TIME"; \
	echo ""; \
	echo "  后端: http://localhost:$(SERVER_PORT)"; \
	echo "  前端: http://localhost:$(UI_PORT)"; \
	echo "========================================"

stop: ## 停止服务 (前端 + 后端)
	@echo "=== 停止 AIFlowy 服务 ==="; \
	echo ""; \
	echo "[1/2] 停止前端 (端口 $(UI_PORT))..."; \
	PID=$$(lsof -ti:$(UI_PORT) 2>/dev/null); \
	if [ -n "$$PID" ]; then \
		kill $$PID 2>/dev/null || true; \
		sleep 1; \
		if lsof -ti:$(UI_PORT) >/dev/null 2>&1; then \
			kill -9 $$PID 2>/dev/null || true; \
		fi; \
		echo "  ✓ 前端已停止"; \
	else \
		echo "  - 前端未运行"; \
	fi; \
	rm -f $(UI_PID_FILE); \
	echo ""; \
	echo "[2/2] 停止后端 (端口 $(SERVER_PORT))..."; \
	PID=$$(lsof -ti:$(SERVER_PORT) 2>/dev/null); \
	if [ -n "$$PID" ]; then \
		kill $$PID 2>/dev/null || true; \
		sleep 1; \
		if lsof -ti:$(SERVER_PORT) >/dev/null 2>&1; then \
			kill -9 $$PID 2>/dev/null || true; \
		fi; \
		echo "  ✓ 后端已停止"; \
	else \
		echo "  - 后端未运行"; \
	fi; \
	rm -f $(PID_FILE); \
	echo ""; \
	echo "  ✓ 所有服务已停止"

restart: stop start ## 重启服务 (前端 + 后端)

status: ## 显示服务状态
	@echo "=== AIFlowy 服务状态 ==="
	@echo ""
	@echo "Go 后端 (端口 $(SERVER_PORT)):"
	@echo "----------------------------------------"
	@PID=$$(lsof -ti:$(SERVER_PORT) 2>/dev/null); \
	if [ -n "$$PID" ]; then \
		echo "  状态: ✓ 运行中"; \
		echo "  PID: $$PID"; \
		echo "  端口: $(SERVER_PORT)"; \
		echo "  运行时间: $$(ps -o etime= -p $$PID 2>/dev/null | xargs)"; \
		echo "  内存占用: $$(ps -o rss= -p $$PID 2>/dev/null | awk '{printf "%.2f MB", $$1/1024}')"; \
		echo "  地址: http://localhost:$(SERVER_PORT)"; \
	else \
		echo "  状态: ✗ 未运行"; \
	fi
	@echo ""
	@echo "前端 (端口 $(UI_PORT)):"
	@echo "----------------------------------------"
	@PID=$$(lsof -ti:$(UI_PORT) 2>/dev/null); \
	if [ -n "$$PID" ]; then \
		echo "  状态: ✓ 运行中"; \
		echo "  PID: $$PID"; \
		echo "  端口: $(UI_PORT)"; \
		echo "  运行时间: $$(ps -o etime= -p $$PID 2>/dev/null | xargs)"; \
		echo "  内存占用: $$(ps -o rss= -p $$PID 2>/dev/null | awk '{printf "%.2f MB", $$1/1024}')"; \
		echo "  地址: http://localhost:$(UI_PORT)"; \
	else \
		echo "  状态: ✗ 未运行"; \
	fi
	@echo "----------------------------------------"

logs: ## 查看 Go 后端实时日志
	@if [ -f $(LOG_FILE) ]; then \
		tail -f $(LOG_FILE); \
	else \
		echo "日志文件不存在: $(LOG_FILE)"; \
	fi

logs-tail: ## 查看最近日志 (最后100行)
	@if [ -f $(LOG_FILE) ]; then \
		tail -100 $(LOG_FILE); \
	else \
		echo "日志文件不存在: $(LOG_FILE)"; \
	fi

clean-logs: ## 清理日志文件
	@echo "正在清理日志..."
	@rm -f $(LOG_DIR)/*.log
	@echo "  ✓ 日志已清理"

dev: go-run ## 启动开发服务器 (前台运行, Go 后端)

# ============================================================================
# 前端 (Vue)
# ============================================================================

ui-install: ## 安装前端依赖
	cd $(UI_ADMIN_DIR) && pnpm install

ui-dev: ## 运行前端开发服务器
	cd $(UI_ADMIN_DIR) && pnpm dev

ui-build: ## 构建前端生产版本
	cd $(UI_ADMIN_DIR) && pnpm build

ui-lint: ## 运行前端代码检查
	cd $(UI_ADMIN_DIR) && pnpm lint

ui-check: ## 运行前端完整检查
	cd $(UI_ADMIN_DIR) && pnpm check

# ============================================================================
# 清理
# ============================================================================

clean: ## 清理所有构建产物
	@echo "正在清理..."
	rm -rf $(GO_BUILD_DIR)
	cd $(GO_SRC_DIR) && $(GO) clean
	@echo "  ✓ 清理完成"

# ============================================================================
# API 测试
# ============================================================================

login: ## 登录获取 API Token 并保存到 /tmp/.aiflowy.token
	@echo "🔐 正在登录..."
	@./scripts/api login $(API_USER) $(API_PASS) && \
		echo "💡 Token 文件: $(TOKEN_FILE)"

api: ## 调用 API (示例: make api ARGS="GET bot/list")
	@if [ -z "$(ARGS)" ]; then \
		echo "用法: make api ARGS=\"<method> <path> [params...]\""; \
		echo ""; \
		echo "示例:"; \
		echo "  make api ARGS=\"GET bot/list\""; \
		echo "  make api ARGS=\"GET auth/getUserInfo\""; \
		echo "  make api ARGS=\"POST bot/chat -d '{\"botId\":\"xxx\",\"message\":\"你好\"}'\""; \
		echo ""; \
		echo "💡 也可以直接使用: ./scripts/api GET bot/list"; \
	else \
		./scripts/api $(ARGS); \
	fi
