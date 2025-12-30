# AIFlowy Go 语言重写规划实施任务书

> 文档编号：004
> 版本：v1.2
> 创建日期：2025-12-29
> 最后更新：2025-12-29
> 状态：Stage 1 已完成
> 策略：小步快跑，分阶段验证

---

## 一、项目概述

### 1.1 项目目标

将 AIFlowy 后端从 Java (Spring Boot) 技术栈迁移至 Go 语言技术栈，实现：

- **性能提升**：内存占用降低 70%，启动速度提升 30-50 倍
- **运维简化**：单二进制部署，无 JVM 依赖
- **云原生优化**：更适合 Kubernetes 容器化部署
- **开发效率**：更简洁的代码，更快的编译速度

### 1.2 实施原则

| 原则 | 说明 |
|-----|------|
| **小步快跑** | 每个阶段控制在 1-2 周，快速迭代 |
| **验证优先** | 每个阶段完成后停下验证，通过后再继续 |
| **功能对等** | 优先保证与 Java 版本功能一致 |
| **渐进迁移** | Go 服务与 Java 服务可并行运行 |
| **文档同步** | 每次验证后更新本文档状态 |

### 1.3 技术选型确认

| 层级 | Java 技术 | Go 替代方案 | 确认状态 |
|-----|----------|------------|---------|
| Web 框架 | Spring Boot 2.7.18 | **Echo v4.14** | [x] 已验证 |
| ORM | MyBatis-Flex 1.11.3 | **database/sql** (暂用原生) | [x] 已验证 |
| AI 框架 | AgentsFlex 2.0.0-beta.5 | **Eino** (字节) | [ ] 待确认 |
| 工作流 | TinyFlow 2.0.0-beta.6 | **Temporal** / 自研 | [ ] 待确认 |
| 认证 | Sa-Token 1.40.0 | **Casbin** + jwt | [ ] 待确认 |
| 配置 | Spring Config | **Viper v1.21** | [x] 已验证 |
| 日志 | Logback | **Zap v1.27** | [x] 已验证 |
| 缓存 | JetCache + Caffeine | **go-redis** + **bigcache** | [ ] 待确认 |

---

## 二、现状分析

### 2.1 Java 代码规模

| 指标 | 数量 |
|-----|------|
| Controller | 71 个 |
| Service | 54 个 |
| 数据库表 | 57 个 |
| API 端点 | 260+ 个 |
| Java 文件 | 766 个 |

### 2.2 功能模块划分

#### 核心业务模块

| 模块 | 表数量 | Controller | 复杂度 | 优先级 |
|-----|--------|-----------|--------|--------|
| 认证授权 | 0 | 1 | 低 | P0 |
| 系统管理 | 16 | 15 | 中 | P0 |
| 模型管理 | 2 | 2 | 低 | P1 |
| Bot 管理 | 10 | 10 | 高 | P1 |
| 工作流 | 4 | 5 | 高 | P2 |
| 知识库 RAG | 4 | 4 | 高 | P2 |
| 插件系统 | 4 | 2 | 中 | P2 |
| 数据中枢 | 2 | 2 | 中 | P3 |
| 定时任务 | 13 | 2 | 中 | P3 |

#### 关键 API 接口

**认证相关**：
- `POST /api/v1/auth/login` - 登录
- `POST /api/v1/auth/logout` - 登出
- `GET /api/v1/auth/getPermissions` - 获取权限

**Bot 相关**：
- `POST /api/v1/bot/chat` - 流式聊天 (SSE) ⭐ 核心
- `POST /api/v1/bot/voiceInput` - 语音输入
- `POST /api/v1/bot/prompt/chore/chat` - 提示词优化
- `GET /api/v1/bot/generateConversationId` - 生成会话ID
- CRUD `/api/v1/bot/**`

**工作流相关**：
- `POST /api/v1/workflow/runAsync` - 异步执行 ⭐ 核心
- `POST /api/v1/workflow/getChainStatus` - 状态查询
- `POST /api/v1/workflow/resume` - 恢复执行
- `POST /api/v1/workflow/singleRun` - 单节点测试

---

## 三、分阶段实施计划

### 阶段总览

```
┌──────────────────────────────────────────────────────────────────┐
│  Stage 0: 技术验证 POC (1周)                                     │
│  ├─ 验证 Echo + Ent + MySQL 可行性                               │
│  └─ 验收: Hello World API + 数据库连接                           │
├──────────────────────────────────────────────────────────────────┤
│  Stage 1: 基础框架 (1周)                                         │
│  ├─ 项目结构、配置、日志、错误处理                                │
│  └─ 验收: 健康检查 API + 中间件链                                 │
├──────────────────────────────────────────────────────────────────┤
│  Stage 2: 认证授权 (1周)                                         │
│  ├─ JWT + Casbin + 登录登出                                      │
│  └─ 验收: admin/123456 登录成功                                  │
├──────────────────────────────────────────────────────────────────┤
│  Stage 3: 系统管理 (2周)                                         │
│  ├─ 用户/角色/菜单/部门/字典 CRUD                                 │
│  └─ 验收: 前端系统管理页面完整可用                                │
├──────────────────────────────────────────────────────────────────┤
│  Stage 4: 模型管理 (1周)                                         │
│  ├─ 模型提供商 + 模型配置                                        │
│  └─ 验收: 可添加/管理 LLM 模型                                   │
├──────────────────────────────────────────────────────────────────┤
│  Stage 5: LLM 核心 (2周)                                         │
│  ├─ Eino 集成 + 多模型适配器                                     │
│  └─ 验收: 可调用 OpenAI/DeepSeek/Ollama                          │
├──────────────────────────────────────────────────────────────────┤
│  Stage 6: Bot 基础 (1周)                                         │
│  ├─ Bot CRUD + 分类 + 会话 + 消息                                │
│  └─ 验收: 可创建 Bot 并查看详情                                  │
├──────────────────────────────────────────────────────────────────┤
│  Stage 7: 流式聊天 (2周)                                         │
│  ├─ SSE + AIFlowy Chat Protocol + Memory                         │
│  └─ 验收: 前端 Bot 聊天功能完整可用                              │
├──────────────────────────────────────────────────────────────────┤
│  Stage 8: Tool 系统 (1周)                                        │
│  ├─ Tool 接口 + 执行框架                                         │
│  └─ 验收: AI 可调用工具并返回结果                                │
├──────────────────────────────────────────────────────────────────┤
│  Stage 9: 插件模块 (1周)                                         │
│  ├─ 插件 CRUD + HTTP 工具                                        │
│  └─ 验收: Bot 可使用插件调用外部 API                             │
├──────────────────────────────────────────────────────────────────┤
│  Stage 10: 工作流基础 (2周)                                      │
│  ├─ 工作流 CRUD + DSL 解析 + 节点定义                            │
│  └─ 验收: 可创建/编辑工作流                                      │
├──────────────────────────────────────────────────────────────────┤
│  Stage 11: 工作流执行 (2周)                                      │
│  ├─ 执行引擎 + 状态管理 + 结果记录                               │
│  └─ 验收: 工作流可异步执行并查询状态                             │
├──────────────────────────────────────────────────────────────────┤
│  Stage 12: 知识库管理 (1周)                                      │
│  ├─ 知识库 + 文档 CRUD                                           │
│  └─ 验收: 可上传/管理文档                                        │
├──────────────────────────────────────────────────────────────────┤
│  Stage 13: RAG 向量检索 (2周)                                    │
│  ├─ 文档分块 + 向量化 + 检索                                     │
│  └─ 验收: Bot 可使用知识库回答问题                               │
├──────────────────────────────────────────────────────────────────┤
│  Stage 14: 辅助功能 (1周)                                        │
│  ├─ API 密钥 + 操作日志 + 定时任务                               │
│  └─ 验收: 辅助功能可用                                           │
├──────────────────────────────────────────────────────────────────┤
│  Stage 15: 集成测试 (1周)                                        │
│  ├─ 端到端测试 + 性能测试                                        │
│  └─ 验收: 所有功能正常，性能达标                                 │
├──────────────────────────────────────────────────────────────────┤
│  Stage 16: 生产部署 (1周)                                        │
│  ├─ Docker 镜像 + K8s 配置 + 灰度切换                            │
│  └─ 验收: 生产环境稳定运行                                       │
└──────────────────────────────────────────────────────────────────┘

总计: 约 20 周 (5 个月)
```

---

## 四、详细阶段任务

### Stage 0: 技术验证 POC

**目标**：验证核心技术栈可行性

**时间**：1 周

**完成日期**：2025-12-29

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 0.1 | 创建 Go 项目 `aiflowy-go` | 项目可编译 | [x] 完成 |
| 0.2 | 集成 Echo 框架 | Hello World API 可访问 | [x] 完成 |
| 0.3 | 集成数据库连接 | 使用 database/sql | [x] 完成 |
| 0.4 | 连接 MySQL 数据库 | 可查询 `tb_sys_account` 表 | [x] 完成 |
| 0.5 | 集成 Viper 配置 | 可读取 config.yaml | [x] 完成 |
| 0.6 | 集成 Zap 日志 | 日志可输出到控制台 | [x] 完成 |

**产出**：
- [x] `aiflowy-go/` 项目目录
- [x] `GET /health` API 返回 `{"code":200,"data":{"status":"ok"}}`
- [x] 数据库连接测试通过 (查询到 admin 用户)

**验收结果**：
```bash
# 1. 健康检查
$ curl http://localhost:8211/health
{"code":200,"data":{"status":"ok"}}

# 2. 详细健康检查 (含数据库)
$ curl http://localhost:8211/health/detail
{
  "code": 200,
  "data": {
    "database": {
      "info": {
        "first_login_name": "admin",
        "first_user_id": "1",
        "table": "tb_sys_account",
        "total_count": 1
      },
      "status": "ok"
    },
    "memory": {"alloc_mb": 0, "sys_mb": 8},
    "runtime": {"go_version": "go1.25.5", "num_goroutine": 8},
    "status": "ok",
    "uptime": "1.030528958s",
    "version": "1.0.0"
  }
}
```

**阶段状态**：[x] 已通过 (2025-12-29)

---

### Stage 1: 基础框架

**目标**：搭建完整的项目骨架

**时间**：1 周

**完成日期**：2025-12-29

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 1.1 | 设计项目目录结构 | 符合 Go 标准布局 | [x] 完成 |
| 1.2 | 配置管理模块 | 支持多环境配置 (dev/prod) | [x] 完成 |
| 1.3 | 日志模块 | 支持分级日志、日志轮转 (Zap + Lumberjack) | [x] 完成 |
| 1.4 | 错误处理 | 统一错误响应格式 (BusinessError) | [x] 完成 |
| 1.5 | 中间件链 | Logger + Recovery + CORS + Gzip | [x] 完成 |
| 1.6 | 统一响应格式 | 与 Java 版本 `Result<T>` 一致 | [x] 完成 |
| 1.7 | 雪花 ID 生成器 | 与 Java 版本兼容 | [x] 完成 |

**验收结果**：

```bash
# 1. 健康检查
$ curl http://localhost:8213/health
{"code":200,"data":{"status":"ok"}}

# 2. 错误处理测试
$ curl http://localhost:8213/test/error
{"code":1001,"message":"这是一个测试错误"}

# 3. 雪花ID生成测试
$ curl http://localhost:8213/test/snowflake
{"code":200,"data":{"generated_ids":["263882911898537984",...],
  "parsed_example":{"datacenter_id":1,"worker_id":1}}}

# 4. 配置管理测试
$ curl http://localhost:8213/test/config
{"code":200,"data":{"environment":"development","is_production":false,...}}

# 5. CORS 测试
$ curl -I -H "Origin: http://localhost:8212" http://localhost:8213/health
Access-Control-Allow-Origin: *
```

**阶段状态**：[x] 已通过 (2025-12-29)

**项目目录结构**：

```
aiflowy-go/
├── cmd/
│   └── server/
│       └── main.go              # 入口
├── internal/
│   ├── config/                  # 配置
│   ├── handler/                 # HTTP 处理器 (Controller)
│   │   ├── auth/
│   │   ├── system/
│   │   ├── ai/
│   │   └── common/
│   ├── service/                 # 业务逻辑
│   ├── repository/              # 数据访问
│   ├── entity/                  # Ent Schema
│   ├── middleware/              # 中间件
│   ├── pkg/                     # 内部工具包
│   │   ├── response/            # 统一响应
│   │   ├── errors/              # 错误定义
│   │   └── utils/               # 工具函数
│   └── router/                  # 路由定义
├── pkg/                         # 可导出的公共包
├── configs/                     # 配置文件
│   ├── config.yaml
│   └── config.dev.yaml
├── scripts/                     # 脚本
├── docs/                        # 文档
├── go.mod
├── go.sum
└── Makefile
```

**统一响应格式** (与 Java 版本一致)：

```go
// pkg/response/response.go
type Result[T any] struct {
    Code    int    `json:"code"`
    Message string `json:"message,omitempty"`
    Data    T      `json:"data,omitempty"`
}

func Ok[T any](data T) Result[T] {
    return Result[T]{Code: 200, Data: data}
}

func Fail(code int, message string) Result[any] {
    return Result[any]{Code: code, Message: message}
}
```

**验收检查点**：
```bash
# 1. 错误测试
curl http://localhost:8211/api/v1/test/error
# 期望返回: {"code": 500, "message": "测试错误"}

# 2. 配置加载
# 日志应显示正确的配置值

# 3. CORS 测试
curl -H "Origin: http://localhost:8212" -I http://localhost:8211/health
# 期望: Access-Control-Allow-Origin 头存在
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 2: 认证授权

**目标**：实现登录、JWT、权限校验

**时间**：1 周

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 2.1 | JWT Token 生成/验证 | Token 可生成和验证 | [ ] |
| 2.2 | 登录 API | `/api/v1/auth/login` | [ ] |
| 2.3 | 登出 API | `/api/v1/auth/logout` | [ ] |
| 2.4 | Casbin 权限模型 | RBAC 模型加载 | [ ] |
| 2.5 | 权限校验中间件 | 未授权返回 401 | [ ] |
| 2.6 | 获取权限列表 | `/api/v1/auth/getPermissions` | [ ] |
| 2.7 | 验证码集成 | `/api/v1/public/getCaptcha` | [ ] |

**API 规格**：

```
POST /api/v1/auth/login
Request:
{
    "username": "admin",
    "password": "123456",
    "captchaId": "xxx",
    "captchaCode": "xxx"
}

Response:
{
    "code": 200,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "account": {
            "id": "1",
            "username": "admin",
            "realName": "管理员",
            ...
        }
    }
}
```

**验收检查点**：
```bash
# 1. 获取验证码
curl http://localhost:8211/api/v1/public/getCaptcha
# 期望: 返回验证码图片数据

# 2. 登录
curl -X POST http://localhost:8211/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
# 期望: 返回 token

# 3. 访问受保护接口
curl http://localhost:8211/api/v1/system/account/list \
  -H "Authorization: Bearer <token>"
# 期望: 返回用户列表

# 4. 未授权访问
curl http://localhost:8211/api/v1/system/account/list
# 期望: 返回 401
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 3: 系统管理

**目标**：实现系统管理 CRUD

**时间**：2 周

**涉及数据表**：

| 表名 | 说明 | 优先级 |
|-----|------|--------|
| tb_sys_account | 用户 | P0 |
| tb_sys_role | 角色 | P0 |
| tb_sys_menu | 菜单 | P0 |
| tb_sys_dept | 部门 | P0 |
| tb_sys_dict | 字典 | P1 |
| tb_sys_dict_item | 字典项 | P1 |
| tb_sys_position | 职位 | P1 |
| tb_sys_account_role | 用户-角色关联 | P0 |
| tb_sys_role_menu | 角色-菜单关联 | P0 |
| tb_sys_role_dept | 角色-部门关联 | P1 |
| tb_sys_account_position | 用户-职位关联 | P1 |

**任务清单**：

| # | 任务 | API 端点 | 状态 |
|---|------|---------|------|
| 3.1 | 用户 CRUD | `/api/v1/system/account/**` | [ ] |
| 3.2 | 角色 CRUD | `/api/v1/system/role/**` | [ ] |
| 3.3 | 菜单 CRUD | `/api/v1/system/menu/**` | [ ] |
| 3.4 | 部门 CRUD | `/api/v1/system/dept/**` | [ ] |
| 3.5 | 职位 CRUD | `/api/v1/system/position/**` | [ ] |
| 3.6 | 字典 CRUD | `/api/v1/system/dict/**` | [ ] |
| 3.7 | 字典项 CRUD | `/api/v1/system/dictItem/**` | [ ] |
| 3.8 | 用户角色分配 | 用户编辑时选择角色 | [ ] |
| 3.9 | 角色菜单分配 | 角色编辑时选择菜单 | [ ] |
| 3.10 | 密码加密 | BCrypt 兼容 Java 版本 | [ ] |

**验收检查点**：
```bash
# 使用前端进行验收
# 1. 打开 http://localhost:8212
# 2. 登录 admin/123456
# 3. 验证以下页面功能:
#    - 系统管理 > 用户管理 (增删改查)
#    - 系统管理 > 角色管理 (增删改查 + 分配权限)
#    - 系统管理 > 菜单管理 (增删改查)
#    - 系统管理 > 部门管理 (增删改查)
#    - 系统管理 > 字典管理 (增删改查)
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 4: 模型管理

**目标**：模型提供商和模型配置管理

**时间**：1 周

**涉及数据表**：
- `tb_model_provider` - 模型提供商
- `tb_model` - 模型配置

**任务清单**：

| # | 任务 | API 端点 | 状态 |
|---|------|---------|------|
| 4.1 | 模型提供商 CRUD | `/api/v1/modelProvider/**` | [ ] |
| 4.2 | 模型 CRUD | `/api/v1/model/**` | [ ] |
| 4.3 | 模型列表按提供商分组 | `GET /api/v1/model/getList` | [ ] |

**验收检查点**：
```bash
# 使用前端验收
# 1. AI 模型 > 模型配置
# 2. 添加模型提供商 (OpenAI, DeepSeek, Ollama)
# 3. 添加模型 (gpt-4, deepseek-chat, llama3)
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 5: LLM 核心

**目标**：集成 Eino 框架，实现多模型调用

**时间**：2 周

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 5.1 | 集成 Eino 框架 | 依赖安装成功 | [ ] |
| 5.2 | OpenAI 适配器 | 可调用 OpenAI API | [ ] |
| 5.3 | DeepSeek 适配器 | 可调用 DeepSeek API | [ ] |
| 5.4 | Ollama 适配器 | 可调用本地 Ollama | [ ] |
| 5.5 | 模型工厂 | 根据配置动态创建模型实例 | [ ] |
| 5.6 | 测试 API | `/api/v1/ai/test` 简单问答 | [ ] |

**模型工厂设计**：

```go
// internal/service/model_factory.go
type ModelFactory interface {
    CreateChatModel(model *entity.Model) (llm.ChatModel, error)
    CreateEmbeddingModel(model *entity.Model) (llm.EmbeddingModel, error)
}

// 根据 provider 类型创建对应的模型实例
func (f *modelFactory) CreateChatModel(model *entity.Model) (llm.ChatModel, error) {
    switch model.Provider.ProviderType {
    case "openai":
        return f.createOpenAIChatModel(model)
    case "deepseek":
        return f.createDeepSeekChatModel(model)
    case "ollama":
        return f.createOllamaChatModel(model)
    default:
        return nil, fmt.Errorf("unsupported provider: %s", model.Provider.ProviderType)
    }
}
```

**验收检查点**：
```bash
# 1. 测试 OpenAI
curl -X POST http://localhost:8211/api/v1/ai/test \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"modelId": "<openai-model-id>", "prompt": "Hello"}'
# 期望: 返回 AI 回复

# 2. 测试 DeepSeek (类似)
# 3. 测试 Ollama (类似)
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 6: Bot 基础

**目标**：Bot CRUD 和基础数据管理

**时间**：1 周

**涉及数据表**：
- `tb_bot` - 机器人
- `tb_bot_category` - 机器人分类
- `tb_bot_conversation` - 对话会话
- `tb_bot_message` - 对话消息
- `tb_bot_model` - 机器人模型关联
- `tb_bot_recently_used` - 最近使用

**任务清单**：

| # | 任务 | API 端点 | 状态 |
|---|------|---------|------|
| 6.1 | Bot CRUD | `/api/v1/bot/**` | [ ] |
| 6.2 | Bot 分类 CRUD | `/api/v1/botCategory/**` | [ ] |
| 6.3 | 会话管理 | `/api/v1/botConversation/**` | [ ] |
| 6.4 | 消息管理 | `/api/v1/botMessage/**` | [ ] |
| 6.5 | 生成会话 ID | `GET /api/v1/bot/generateConversationId` | [ ] |
| 6.6 | Bot 详情 | `GET /api/v1/bot/getDetail` | [ ] |
| 6.7 | 更新 LLM 配置 | `POST /api/v1/bot/updateLlmOptions` | [ ] |

**验收检查点**：
```bash
# 使用前端验收
# 1. AI 应用 > 机器人管理
# 2. 创建 Bot
# 3. 编辑 Bot 配置
# 4. 选择模型
# 5. 配置系统提示词
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 7: 流式聊天

**目标**：实现 SSE 流式聊天，完整的 AIFlowy Chat Protocol

**时间**：2 周

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 7.1 | SSE 基础设施 | Echo SSE 支持 | [ ] |
| 7.2 | AIFlowy Chat Protocol | 实现 v1.1 协议 | [ ] |
| 7.3 | 聊天 API | `POST /api/v1/bot/chat` | [ ] |
| 7.4 | Memory 管理 | 历史消息上下文 | [ ] |
| 7.5 | 消息持久化 | 保存用户消息和 AI 回复 | [ ] |
| 7.6 | 会话管理 | 创建/更新会话 | [ ] |
| 7.7 | ChatOptions | temperature, topK, topP 等 | [ ] |
| 7.8 | Thinking 模式 | 支持 thinking 过程输出 | [ ] |

**AIFlowy Chat Protocol v1.1 格式**：

```
event: message
data: {"domain":"llm","type":"text","content":"Hello"}

event: message
data: {"domain":"llm","type":"thinking","content":"Let me think..."}

event: message
data: {"domain":"tool","type":"call","content":{"name":"search","args":{}}}

event: message
data: {"domain":"system","type":"done","content":""}
```

**验收检查点**：
```bash
# 使用前端验收
# 1. 打开一个 Bot 聊天界面
# 2. 发送消息 "你好"
# 3. 验证:
#    - 流式输出 (字符逐个显示)
#    - 消息保存 (刷新后消息还在)
#    - 上下文 (可进行多轮对话)
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 8: Tool 系统

**目标**：实现 Tool 调用框架

**时间**：1 周

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 8.1 | Tool 接口定义 | Tool interface | [ ] |
| 8.2 | Tool 注册机制 | Tool registry | [ ] |
| 8.3 | Tool 执行器 | Tool executor | [ ] |
| 8.4 | LLM Tool 调用 | Function calling 支持 | [ ] |
| 8.5 | Tool 结果返回 | 结果返回给 LLM | [ ] |

**Tool 接口设计**：

```go
// internal/service/tool/tool.go
type Tool interface {
    Name() string
    Description() string
    Parameters() map[string]interface{}
    Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

type ToolRegistry struct {
    tools map[string]Tool
}

func (r *ToolRegistry) Register(tool Tool) {
    r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Execute(ctx context.Context, name string, args map[string]interface{}) (interface{}, error) {
    tool, ok := r.tools[name]
    if !ok {
        return nil, fmt.Errorf("tool not found: %s", name)
    }
    return tool.Execute(ctx, args)
}
```

**验收检查点**：
```bash
# 验证 Tool 调用
# 1. 创建一个简单的内置 Tool (如: 获取当前时间)
# 2. 在聊天中问 "现在几点了"
# 3. 验证 AI 调用 Tool 并返回结果
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 9: 插件模块

**目标**：插件管理和 HTTP 工具

**时间**：1 周

**涉及数据表**：
- `tb_plugin` - 插件
- `tb_plugin_item` - 插件项 (具体工具)
- `tb_plugin_category` - 插件分类
- `tb_bot_plugin` - Bot-插件关联

**任务清单**：

| # | 任务 | API 端点 | 状态 |
|---|------|---------|------|
| 9.1 | 插件 CRUD | `/api/v1/plugin/**` | [ ] |
| 9.2 | 插件项 CRUD | `/api/v1/pluginItem/**` | [ ] |
| 9.3 | HTTP Tool 执行器 | 调用外部 API | [ ] |
| 9.4 | Bot-插件关联 | `/api/v1/botPlugin/**` | [ ] |
| 9.5 | 插件转 Tool | PluginItem.toTool() | [ ] |

**验收检查点**：
```bash
# 使用前端验收
# 1. 创建一个 HTTP 插件 (如: 天气 API)
# 2. 将插件关联到 Bot
# 3. 聊天中问 "北京天气怎么样"
# 4. 验证 Bot 调用插件并返回结果
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 10: 工作流基础

**目标**：工作流 CRUD 和 DSL 解析

**时间**：2 周

**涉及数据表**：
- `tb_workflow` - 工作流
- `tb_workflow_category` - 工作流分类
- `tb_bot_workflow` - Bot-工作流关联

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 10.1 | 工作流 CRUD | `/api/v1/workflow/**` | [ ] |
| 10.2 | 工作流分类 | `/api/v1/workflowCategory/**` | [ ] |
| 10.3 | DSL 解析器 | 解析工作流 JSON | [ ] |
| 10.4 | 节点类型定义 | LLM/Tool/Condition/... | [ ] |
| 10.5 | 获取运行参数 | `GET /api/v1/workflow/getRunningParameters` | [ ] |
| 10.6 | 导入/导出 | import/export API | [ ] |
| 10.7 | Bot-工作流关联 | `/api/v1/botWorkflow/**` | [ ] |
| 10.8 | 工作流转 Tool | Workflow.toTool() | [ ] |

**验收检查点**：
```bash
# 使用前端验收
# 1. 工作流 > 工作流管理
# 2. 创建工作流
# 3. 编辑节点
# 4. 保存工作流
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 11: 工作流执行

**目标**：工作流执行引擎

**时间**：2 周

**涉及数据表**：
- `tb_workflow_exec_result` - 执行结果
- `tb_workflow_exec_step` - 执行步骤

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 11.1 | 执行引擎核心 | ChainExecutor | [ ] |
| 11.2 | 异步执行 | `POST /api/v1/workflow/runAsync` | [ ] |
| 11.3 | 状态查询 | `POST /api/v1/workflow/getChainStatus` | [ ] |
| 11.4 | 恢复执行 | `POST /api/v1/workflow/resume` | [ ] |
| 11.5 | 单节点运行 | `POST /api/v1/workflow/singleRun` | [ ] |
| 11.6 | 结果记录 | 保存执行结果和步骤 | [ ] |
| 11.7 | 节点执行器 | LLM/Tool/Condition 节点 | [ ] |
| 11.8 | 人工确认节点 | 支持暂停和恢复 | [ ] |

**验收检查点**：
```bash
# 使用前端验收
# 1. 创建一个简单工作流 (开始 -> LLM节点 -> 结束)
# 2. 运行工作流
# 3. 查看执行状态
# 4. 查看执行结果
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 12: 知识库管理

**目标**：知识库和文档管理

**时间**：1 周

**涉及数据表**：
- `tb_document_collection` - 知识库
- `tb_document` - 文档
- `tb_bot_document_collection` - Bot-知识库关联

**任务清单**：

| # | 任务 | API 端点 | 状态 |
|---|------|---------|------|
| 12.1 | 知识库 CRUD | `/api/v1/documentCollection/**` | [ ] |
| 12.2 | 文档 CRUD | `/api/v1/document/**` | [ ] |
| 12.3 | 文件上传 | `/api/v1/upload/**` | [ ] |
| 12.4 | Bot-知识库关联 | `/api/v1/botDocumentCollection/**` | [ ] |

**验收检查点**：
```bash
# 使用前端验收
# 1. 知识库 > 知识库管理
# 2. 创建知识库
# 3. 上传文档 (PDF/TXT)
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 13: RAG 向量检索

**目标**：文档向量化和检索

**时间**：2 周

**涉及数据表**：
- `tb_document_chunk` - 文档分块
- `tb_document_history` - 文档历史

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 13.1 | 文档解析 | PDF/Word/Excel/TXT | [ ] |
| 13.2 | 文本分块 | DocumentSplitter | [ ] |
| 13.3 | 向量化 | Embedding Model 调用 | [ ] |
| 13.4 | 向量存储 | 本地存储 / Milvus | [ ] |
| 13.5 | 相似度检索 | 向量检索 API | [ ] |
| 13.6 | Rerank | 重排序 | [ ] |
| 13.7 | 知识库 Tool | DocumentCollection.toTool() | [ ] |
| 13.8 | RAG 集成到聊天 | Bot 使用知识库回答 | [ ] |

**验收检查点**：
```bash
# 使用前端验收
# 1. 上传一个文档到知识库
# 2. 将知识库关联到 Bot
# 3. 聊天中问文档相关问题
# 4. 验证 Bot 基于文档回答
```

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 14: 辅助功能

**目标**：API 密钥、日志、定时任务等

**时间**：1 周

**任务清单**：

| # | 任务 | API 端点 | 状态 |
|---|------|---------|------|
| 14.1 | 系统 API 密钥 | `/api/v1/system/apiKey/**` | [ ] |
| 14.2 | Bot API 密钥 | `/api/v1/botApiKey/**` | [ ] |
| 14.3 | 操作日志 | `/api/v1/system/log/**` | [ ] |
| 14.4 | 系统配置 | `/api/v1/system/option/**` | [ ] |
| 14.5 | 定时任务 | `/api/v1/system/job/**` | [ ] |
| 14.6 | 语音输入 | `POST /api/v1/bot/voiceInput` | [ ] |
| 14.7 | 提示词优化 | `POST /api/v1/bot/prompt/chore/chat` | [ ] |

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 15: 集成测试

**目标**：完整功能测试和性能测试

**时间**：1 周

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 15.1 | 单元测试 | 核心模块覆盖率 > 70% | [ ] |
| 15.2 | 集成测试 | 主要流程自动化测试 | [ ] |
| 15.3 | 性能测试 | 压测报告 | [ ] |
| 15.4 | 内存测试 | 内存泄漏检查 | [ ] |
| 15.5 | 并发测试 | 100 并发无错误 | [ ] |
| 15.6 | Bug 修复 | 修复发现的问题 | [ ] |

**性能指标目标**：

| 指标 | Java 基准 | Go 目标 |
|-----|----------|---------|
| 启动时间 | 5s | < 500ms |
| 内存占用 | 500MB | < 150MB |
| QPS (简单 API) | 5000 | > 10000 |
| P99 延迟 | 100ms | < 50ms |

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

### Stage 16: 生产部署

**目标**：Docker 镜像、K8s 配置、灰度切换

**时间**：1 周

**任务清单**：

| # | 任务 | 验收标准 | 状态 |
|---|------|---------|------|
| 16.1 | Dockerfile | 多阶段构建 | [ ] |
| 16.2 | Docker 镜像 | 镜像大小 < 50MB | [ ] |
| 16.3 | K8s 配置 | Deployment + Service | [ ] |
| 16.4 | 健康检查 | Liveness + Readiness | [ ] |
| 16.5 | 配置管理 | ConfigMap / Secret | [ ] |
| 16.6 | 灰度部署 | 10% -> 50% -> 100% | [ ] |
| 16.7 | 监控告警 | Prometheus + Grafana | [ ] |
| 16.8 | 文档完善 | 运维手册 | [ ] |

**阶段状态**：[ ] 未开始 / [ ] 进行中 / [ ] 待验证 / [ ] 已通过

---

## 五、验证流程

### 5.1 每阶段验证步骤

```
1. 开发完成
   ↓
2. 自测通过
   ↓
3. 更新本文档 (标记任务完成)
   ↓
4. 通知用户验证
   ↓
5. 用户验收
   ├─ 通过 → 更新阶段状态为 [已通过]，继续下一阶段
   └─ 不通过 → 记录问题，修复后重新验证
```

### 5.2 验证记录模板

```markdown
## 阶段 X 验证记录

**验证日期**: YYYY-MM-DD
**验证人**: xxx

### 验证项目

| # | 验证项 | 预期结果 | 实际结果 | 状态 |
|---|-------|---------|---------|------|
| 1 | xxx | xxx | xxx | ✅/❌ |

### 发现问题

| # | 问题描述 | 严重程度 | 状态 |
|---|---------|---------|------|
| 1 | xxx | 高/中/低 | 待修复/已修复 |

### 验证结论

[ ] 通过 - 可以进入下一阶段
[ ] 不通过 - 需要修复后重新验证
```

---

## 六、风险控制

### 6.1 技术风险

| 风险 | 概率 | 影响 | 缓解措施 |
|-----|------|------|---------|
| Eino 框架功能不足 | 中 | 高 | 备选 LangChainGo |
| 工作流引擎复杂度高 | 高 | 高 | 优先使用 Temporal |
| RAG 性能不达标 | 中 | 中 | 使用 Milvus 替代本地存储 |
| 前端兼容性问题 | 低 | 中 | 保持 API 响应格式一致 |

### 6.2 进度风险

| 风险 | 概率 | 影响 | 缓解措施 |
|-----|------|------|---------|
| 某阶段超时 | 中 | 中 | 及时调整计划，必要时简化功能 |
| 依赖库版本问题 | 低 | 低 | 锁定依赖版本 |
| 人员变动 | 低 | 高 | 文档完善，代码规范 |

---

## 七、更新历史

| 版本 | 日期 | 更新内容 | 更新人 |
|-----|------|---------|--------|
| v1.0 | 2025-12-29 | 初始版本 | Claude |

---

## 八、附录

### A. 相关文档

- [003-后端功能分析与Go重写可行性评估](003-后端功能分析与Go重写可行性评估.md)
- [001-AIFlowy开发部署完整指南](001-AIFlowy开发部署完整指南.md)

### B. 参考资源

- [Echo 官方文档](https://echo.labstack.com/)
- [Ent 官方文档](https://entgo.io/)
- [Eino 项目](https://github.com/cloudwego/eino)
- [Temporal Go SDK](https://github.com/temporalio/sdk-go)
- [Casbin](https://casbin.org/)

### C. 联系方式

如有问题，请联系项目负责人。

---

*本文档将随项目进展持续更新*
