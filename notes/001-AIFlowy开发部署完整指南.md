# AIFlowy 开发部署完整指南

> 文档编号：001
> 版本：v2.0.2
> 最后更新：2025-12-29

---

## 目录

1. [项目概述](#1-项目概述)
2. [环境要求](#2-环境要求)
3. [快速开始](#3-快速开始)
4. [后端开发](#4-后端开发)
5. [前端开发](#5-前端开发)
6. [数据库配置](#6-数据库配置)
7. [存储配置](#7-存储配置)
8. [缓存配置](#8-缓存配置)
9. [AI 功能配置](#9-ai-功能配置)
10. [Docker 部署](#10-docker-部署)
11. [生产环境部署](#11-生产环境部署)
12. [项目架构](#12-项目架构)
13. [常见问题](#13-常见问题)

---

## 1. 项目概述

AIFlowy 是一个企业级开源 AI 应用开发平台，采用 Java + Vue 3 技术栈，提供以下核心功能：

**AI 功能**
- Bots 智能助手应用
- 业务插件系统
- RAG 知识库
- 工作流智能体编排（AI Workflow）
- 素材中心（AI 自动生成图片、音频、视频）
- 数据中心（可定制化表格、工作流读写）
- 本地模型支持
- 大模型管理

**系统管理**
- 用户/角色/菜单管理
- 部门/岗位管理
- 日志管理
- 数据字典
- 定时任务

---

## 2. 环境要求

### 2.1 后端环境

| 组件 | 版本要求 | 说明 |
|------|----------|------|
| JDK | 8+ | 推荐 JDK 8，项目编译目标版本 |
| Maven | 3.6+ | 构建工具 |
| MySQL | 8.0+ | 主数据库，需支持 utf8mb4 |
| Redis | 6.0+ | 可选，用于分布式缓存 |

### 2.2 前端环境

| 组件 | 版本要求 | 说明 |
|------|----------|------|
| Node.js | ≥20.10.0 | 运行环境 |
| pnpm | ≥9.12.0 | 包管理器（推荐 10.14.0） |

### 2.3 可选组件

| 组件 | 版本 | 用途 |
|------|------|------|
| Ollama | 最新 | 本地大模型运行 |
| Elasticsearch | 8.x | RAG 搜索引擎（可选） |
| Docker | 20.10+ | 容器化部署 |

---

## 3. 快速开始

### 3.1 克隆项目

```bash
git clone https://gitee.com/aiflowy/aiflowy.git
cd aiflowy
```

### 3.2 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE aiflowy CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

# 导入表结构
mysql -u root -p aiflowy < sql/aiflowy-v2.ddl.sql

# 导入初始数据
mysql -u root -p aiflowy < sql/aiflowy-v2.data.sql
```

### 3.3 启动后端

```bash
# 编译项目（首次编译约需 3 分钟）
mvn clean package -DskipTests

# 运行（开发模式）
java -jar aiflowy-starter/aiflowy-starter-all/target/aiflowy-starter-all-2.0.2.jar
```

> **注意**：如遇 `无效的目标发行版: 17` 错误，需确保根 `pom.xml` 中已显式配置 `maven-compiler-plugin`（源码已包含此配置）。

### 3.4 启动前端

```bash
cd aiflowy-ui-admin
pnpm install
pnpm dev
```

### 3.5 访问系统

- 前端地址：http://localhost:5173
- 后端地址：http://localhost:8080
- 默认账号：`admin`
- 默认密码：`123456`

---

## 4. 后端开发

### 4.1 项目结构

```
aiflowy/
├── aiflowy-api/                    # API 层（接口定义与控制器）
│   ├── aiflowy-api-admin/          # 管理端 API
│   ├── aiflowy-api-public/         # 公共 API
│   ├── aiflowy-api-usercenter/     # 用户中心 API
│   └── aiflowy-api-mcp/            # MCP 协议 API
├── aiflowy-commons/                # 通用组件
│   ├── aiflowy-common-ai/          # AI 相关工具
│   ├── aiflowy-common-base/        # 基础工具类
│   ├── aiflowy-common-cache/       # 缓存组件
│   ├── aiflowy-common-web/         # Web 组件
│   ├── aiflowy-common-satoken/     # 认证组件
│   ├── aiflowy-common-file-storage/# 文件存储
│   ├── aiflowy-common-sms/         # 短信服务
│   ├── aiflowy-common-audio/       # 音频处理
│   ├── aiflowy-common-chat-protocol/# 聊天协议
│   └── aiflowy-common-tcaptcha/    # 验证码
├── aiflowy-modules/                # 业务模块
│   ├── aiflowy-module-ai/          # AI 核心模块
│   ├── aiflowy-module-auth/        # 认证模块
│   ├── aiflowy-module-system/      # 系统管理模块
│   ├── aiflowy-module-core/        # 核心业务模块
│   ├── aiflowy-module-datacenter/  # 数据中心模块
│   ├── aiflowy-module-job/         # 定时任务模块
│   ├── aiflowy-module-log/         # 日志模块
│   └── aiflowy-module-autoconfig/  # 自动配置
├── aiflowy-starter/                # 启动器
│   ├── aiflowy-starter-all/        # 全量启动器
│   ├── aiflowy-starter-admin/      # 管理端启动器
│   ├── aiflowy-starter-public/     # 公共端启动器
│   ├── aiflowy-starter-usercenter/ # 用户中心启动器
│   └── aiflowy-starter-codegen/    # 代码生成器
└── sql/                            # 数据库脚本
```

### 4.2 核心依赖版本

| 依赖 | 版本 | 说明 |
|------|------|------|
| Spring Boot | 2.7.18 | Web 框架 |
| MyBatis-Flex | 1.11.3 | ORM 框架 |
| AgentsFlex | 2.0.0-beta.5 | AI Agent 框架 |
| TinyFlow-Java | 2.0.0-beta.6 | 工作流引擎 |
| Sa-Token | 1.40.0 | 认证授权框架 |
| HikariCP | 4.0.3 | 数据库连接池 |
| Hutool | 5.8.36 | 工具类库 |
| FastJSON | 2.0.57 | JSON 处理 |

### 4.3 Maven 常用命令

```bash
# 清理并编译
mvn clean compile

# 打包（跳过测试）
mvn clean package -DskipTests

# 安装到本地仓库
mvn clean install -DskipTests

# 只编译指定模块
mvn clean package -pl aiflowy-starter/aiflowy-starter-all -am

# 运行单元测试
mvn test

# 生成依赖树
mvn dependency:tree
```

### 4.4 启动器说明

| 启动器 | 端口 | 用途 |
|--------|------|------|
| aiflowy-starter-all | 8080 | 全量部署，包含所有模块 |
| aiflowy-starter-admin | 8080 | 仅管理端 |
| aiflowy-starter-public | 8081 | 仅公共端 |
| aiflowy-starter-usercenter | 8082 | 仅用户中心 |
| aiflowy-starter-codegen | 8090 | 代码生成器 |

### 4.5 开发调试

```bash
# 使用 dev 配置启动
java -jar aiflowy-starter-all-2.0.2.jar --spring.profiles.active=dev

# 开启调试端口
java -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=5005 \
     -jar aiflowy-starter-all-2.0.2.jar

# 指定配置文件
java -jar aiflowy-starter-all-2.0.2.jar \
     --spring.config.location=file:./config/application.yml
```

---

## 5. 前端开发

### 5.1 项目结构

AIFlowy 包含多个前端应用：

```
aiflowy-ui-admin/           # 管理端（主应用）
aiflowy-ui-usercenter/      # 用户中心
aiflowy-ui-websdk/          # 嵌入式 Web SDK
docs/                       # VitePress 文档站点
```

### 5.2 管理端结构（aiflowy-ui-admin）

```
aiflowy-ui-admin/
├── app/                        # 主应用
│   ├── src/
│   │   ├── api/               # API 接口定义
│   │   ├── components/        # 业务组件
│   │   ├── views/             # 页面组件
│   │   ├── router/            # 路由配置
│   │   ├── store/             # Pinia 状态管理
│   │   ├── layouts/           # 布局组件
│   │   ├── locales/           # 国际化
│   │   ├── assets/            # 静态资源
│   │   └── utils/             # 工具函数
│   └── vite.config.mts        # Vite 配置
├── packages/                   # 共享包
│   ├── @core/
│   │   ├── base/              # 基础包
│   │   │   ├── design/        # 设计系统
│   │   │   ├── icons/         # 图标
│   │   │   ├── shared/        # 共享工具
│   │   │   └── typings/       # 类型定义
│   │   ├── composables/       # Vue Composables
│   │   ├── preferences/       # 偏好设置
│   │   └── ui-kit/            # UI 组件库
│   │       ├── form-ui/       # 表单组件
│   │       ├── layout-ui/     # 布局组件
│   │       ├── menu-ui/       # 菜单组件
│   │       ├── popup-ui/      # 弹窗组件
│   │       ├── shadcn-ui/     # Shadcn 组件
│   │       └── tabs-ui/       # 标签页组件
│   ├── effects/               # 业务逻辑包
│   │   ├── access/            # 权限控制
│   │   ├── common-ui/         # 通用 UI
│   │   ├── hooks/             # 自定义 Hooks
│   │   ├── layouts/           # 布局管理
│   │   ├── plugins/           # 插件
│   │   └── request/           # HTTP 请求
│   ├── constants/             # 常量定义
│   ├── icons/                 # 图标库
│   ├── locales/               # 国际化
│   ├── preferences/           # 应用偏好
│   ├── stores/                # 全局状态
│   ├── styles/                # 全局样式
│   ├── types/                 # TypeScript 类型
│   └── utils/                 # 工具函数
├── internal/                   # 内部开发工具
│   └── lint-configs/          # Lint 配置
├── scripts/                    # 构建脚本
└── pnpm-workspace.yaml        # 工作区配置
```

### 5.3 常用命令

```bash
# 进入前端目录
cd aiflowy-ui-admin

# 安装依赖
pnpm install

# 开发模式（管理端）
pnpm dev:app

# 开发模式（文档）
pnpm dev:docs

# 构建生产版本
pnpm build

# 仅构建管理端
pnpm build:app

# 代码检查
pnpm lint

# 代码格式化
pnpm format

# 类型检查
pnpm check:type

# 循环依赖检查
pnpm check:circular

# 拼写检查
pnpm check:cspell

# 完整检查（推荐提交前运行）
pnpm check

# 单元测试
pnpm test:unit

# E2E 测试
pnpm test:e2e

# 清理 node_modules 和 dist
pnpm clean

# 重新安装依赖
pnpm reinstall

# 预览构建结果
pnpm preview

# 构建 Docker 镜像
pnpm build:docker
```

### 5.4 代理配置

开发模式下，Vite 会自动代理 API 请求到后端：

```typescript
// app/vite.config.mts
server: {
  proxy: {
    '/api': {
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
      target: 'http://localhost:5320/api',  // 后端地址
      ws: true,                              // 支持 WebSocket
    },
  },
},
```

如需修改后端地址，编辑 `app/vite.config.mts` 中的 `target` 值。

### 5.5 环境变量

创建 `.env.development` 或 `.env.production` 文件：

```bash
# API 基础路径
VITE_API_BASE_URL=/api

# 应用标题
VITE_APP_TITLE=AIFlowy

# 是否启用 Mock
VITE_USE_MOCK=false
```

### 5.6 技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.5.17 | 前端框架 |
| TypeScript | 5.9.3 | 类型系统 |
| Vite | 7.2.2 | 构建工具 |
| Pinia | 3.0.3 | 状态管理 |
| Vue Router | 4.5.1 | 路由管理 |
| Element Plus | 2.10.2 | UI 组件库 |
| Tailwind CSS | 3.4.18 | CSS 框架 |
| VueUse | 13.4.0 | 组合式工具 |
| Axios | 1.10.0 | HTTP 客户端 |

---

## 6. 数据库配置

### 6.1 MySQL 配置

编辑 `aiflowy-starter/aiflowy-starter-all/src/main/resources/application.yml`：

```yaml
spring:
  datasource:
    url: jdbc:mysql://127.0.0.1:3306/aiflowy?useInformationSchema=true&characterEncoding=utf-8
    username: root
    password: 123456
```

### 6.2 连接池配置

项目使用 HikariCP 连接池，可在配置中调整：

```yaml
spring:
  datasource:
    hikari:
      minimum-idle: 5
      maximum-pool-size: 20
      idle-timeout: 30000
      pool-name: AIFlowyHikariCP
      max-lifetime: 2000000
      connection-timeout: 30000
```

### 6.3 数据库脚本

```
sql/
├── aiflowy-v2.ddl.sql     # 表结构定义（约 1086 行）
└── aiflowy-v2.data.sql    # 初始数据
```

**主要数据表前缀：**

| 前缀 | 说明 |
|------|------|
| `tb_bot_*` | Bot 相关表 |
| `tb_document_*` | 文档/知识库相关表 |
| `tb_model_*` | 大模型相关表 |
| `tb_workflow_*` | 工作流相关表 |
| `tb_user_*` | 用户相关表 |
| `tb_role_*` | 角色权限相关表 |
| `tb_menu_*` | 菜单相关表 |
| `tb_dict_*` | 数据字典表 |
| `tb_qrtz_*` | Quartz 定时任务表 |

---

## 7. 存储配置

### 7.1 本地存储

```yaml
aiflowy:
  storage:
    type: local
    local:
      # Windows 示例
      root: 'C:\aiflowy\attachment'
      # Linux 示例
      # root: /www/aiflowy/attachment
      prefix: 'http://localhost:8080/static'

spring:
  web:
    resources:
      # Windows 示例
      static-locations: file:C:\aiflowy\attachment
      # Linux 示例
      # static-locations: file:/www/aiflowy/attachment
  mvc:
    static-path-pattern: /static/**
```

### 7.2 S3 兼容存储

支持阿里云 OSS、AWS S3、MinIO 等：

```yaml
aiflowy:
  storage:
    type: s3
    s3:
      manufacturer: aliyun  # aliyun/aws/minio
      access-key: your-access-key
      secret-key: your-secret-key
      endpoint: oss-cn-hangzhou.aliyuncs.com
      region: cn-hangzhou
      bucket-name: your-bucket-name
      access-policy: 2  # 0:私有 1:公共读 2:公共读写
      prefix: attachment
```

---

## 8. 缓存配置

### 8.1 本地缓存（默认）

```yaml
jetcache:
  cacheType: local
  local:
    default:
      type: linkedhashmap
      keyConvertor: fastjson
```

### 8.2 Redis 缓存

```yaml
jetcache:
  cacheType: remote  # 或 both（本地+远程）
  remote:
    default:
      type: redis
      keyConvertor: fastjson2
      valueEncoder: java
      valueDecoder: java
      broadcastChannel: aiflowy
      poolConfig:
        minIdle: 5
        maxIdle: 20
        maxTotal: 50
      host: 127.0.0.1
      port: 6379
      password: your-password  # 无密码则删除此行
      database: 0
```

---

## 9. AI 功能配置

### 9.1 Ollama 本地模型

```yaml
aiflowy:
  ollama:
    host: http://127.0.0.1:11434
```

### 9.2 RAG 搜索引擎

**Lucene（默认，轻量级）：**

```yaml
rag:
  searcher:
    type: lucene
    lucene:
      indexDirPath: ./luceneKnowledge
```

**Elasticsearch（生产推荐）：**

```yaml
rag:
  searcher:
    type: elasticSearch
    elastic:
      host: https://127.0.0.1:9200
      userName: elastic
      password: your-password
      indexName: aiflowy
```

### 9.3 阿里云语音服务

```yaml
aiflowy:
  audio:
    type: aliAudioService
    ali:
      access-key-id: your-access-key-id
      access-key-secret: your-access-key-secret
      app-key: your-app-key
      voice: siyue  # 发音人
```

### 9.4 搜索引擎节点

```yaml
node:
  reader: 'defaultReader'
  bochaai:
    apiKey: 'your-bocha-api-key'
```

### 9.5 Bot API Key 配置

```yaml
aiflowy:
  aiBot:
    apiKeyMasterKey: your-32-char-master-key  # 用于加密 Bot API Key
```

---

## 10. Docker 部署

### 10.1 docker-compose 部署

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

**docker-compose.yml 配置：**

```yaml
version: '3.8'
services:
  aiflowy:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - SPRING_DATASOURCE_URL=jdbc:mysql://mysql:3306/aiflowy?useInformationSchema=true&characterEncoding=utf-8
      - SPRING_DATASOURCE_USERNAME=root
      - SPRING_DATASOURCE_PASSWORD=123456
    networks:
      - aiflowy-net
    depends_on:
      mysql:
        condition: service_healthy

  mysql:
    image: mysql:8.0
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: aiflowy
    ports:
      - "3306:3306"
    volumes:
      - ./initdb:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10
    networks:
      - aiflowy-net

networks:
  aiflowy-net:
```

### 10.2 单独构建镜像

```bash
# 先编译项目
mvn clean package -DskipTests

# 构建 Docker 镜像
docker build \
  --build-arg VERSION=2.0.2 \
  --build-arg SERVICE_NAME=aiflowy-starter-all \
  -t aiflowy:2.0.2 .

# 运行容器
docker run -d \
  --name aiflowy \
  -p 8080:8080 \
  -e SPRING_DATASOURCE_URL="jdbc:mysql://host.docker.internal:3306/aiflowy?useInformationSchema=true&characterEncoding=utf-8" \
  -e SPRING_DATASOURCE_USERNAME=root \
  -e SPRING_DATASOURCE_PASSWORD=123456 \
  aiflowy:2.0.2
```

### 10.3 Dockerfile 说明

```dockerfile
FROM openjdk:8-jdk-alpine

# 环境变量
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
ENV JAVA_OPTS=""

# 安装字体（验证码需要）
RUN apk --no-cache add ttf-dejavu fontconfig

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 复制 JAR
ADD ./aiflowy-starter/target/aiflowy-starter-all-${VERSION}.jar app.jar

EXPOSE 8080

ENTRYPOINT java ${JAVA_OPTS} -Djava.security.egd=file:/dev/./urandom -jar /app.jar
```

---

## 11. 生产环境部署

### 11.1 配置文件

创建 `application-prod.yml`：

```yaml
spring:
  config:
    activate:
      on-profile: prod
  datasource:
    url: jdbc:mysql://your-db-host:3306/aiflowy?useInformationSchema=true&characterEncoding=utf-8
    username: aiflowy
    password: your-secure-password
    hikari:
      maximum-pool-size: 50
      minimum-idle: 10

aiflowy:
  storage:
    type: s3
    s3:
      # 生产环境使用云存储

jetcache:
  cacheType: remote
  remote:
    default:
      type: redis
      host: your-redis-host
      port: 6379
      password: your-redis-password

logging:
  file:
    path: /var/log/aiflowy/
    name: aiflowy.log
  level:
    root: warn
    tech.aiflowy: info
```

### 11.2 JVM 优化参数

```bash
java \
  -Xms2g \
  -Xmx4g \
  -XX:+UseG1GC \
  -XX:MaxGCPauseMillis=200 \
  -XX:+HeapDumpOnOutOfMemoryError \
  -XX:HeapDumpPath=/var/log/aiflowy/heapdump.hprof \
  -Dfile.encoding=UTF-8 \
  -Duser.timezone=Asia/Shanghai \
  -jar aiflowy-starter-all-2.0.2.jar \
  --spring.profiles.active=prod
```

### 11.3 Systemd 服务

创建 `/etc/systemd/system/aiflowy.service`：

```ini
[Unit]
Description=AIFlowy Application
After=network.target mysql.service redis.service

[Service]
Type=simple
User=aiflowy
Group=aiflowy
WorkingDirectory=/opt/aiflowy
ExecStart=/usr/bin/java \
  -Xms2g -Xmx4g \
  -XX:+UseG1GC \
  -jar /opt/aiflowy/aiflowy-starter-all-2.0.2.jar \
  --spring.profiles.active=prod
ExecStop=/bin/kill -SIGTERM $MAINPID
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
# 启用服务
sudo systemctl daemon-reload
sudo systemctl enable aiflowy
sudo systemctl start aiflowy
sudo systemctl status aiflowy
```

### 11.4 Nginx 反向代理

```nginx
upstream aiflowy_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 前端静态文件
    location / {
        root /var/www/aiflowy;
        try_files $uri $uri/ /index.html;
        gzip on;
        gzip_types text/plain application/json application/javascript text/css;
    }

    # API 代理
    location /api {
        proxy_pass http://aiflowy_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # SSE 支持
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 3600s;
    }

    # WebSocket 代理
    location /ws {
        proxy_pass http://aiflowy_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 3600s;
    }

    # 静态资源
    location /static {
        proxy_pass http://aiflowy_backend;
        proxy_cache_valid 200 1d;
    }
}
```

---

## 12. 项目架构

### 12.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端应用                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Admin UI   │  │ UserCenter  │  │      Web SDK        │  │
│  │  (Vue 3)    │  │  (Vue 3)    │  │   (TypeScript)      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                     HTTP / WebSocket / SSE
                            │
┌─────────────────────────────────────────────────────────────┐
│                        API 层                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Admin API   │  │ Public API  │  │   UserCenter API    │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                       业务模块层                             │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌───────┐  │
│  │   AI    │ │  Auth   │ │ System  │ │  Core   │ │  Job  │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └───────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                       通用组件层                             │
│  ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────────────┐  │
│  │ Base  │ │ Cache │ │  Web  │ │ Auth  │ │ File Storage  │  │
│  └───────┘ └───────┘ └───────┘ └───────┘ └───────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                       基础设施                               │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌───────────────┐   │
│  │  MySQL  │  │  Redis  │  │ Lucene/ │  │   Object      │   │
│  │         │  │         │  │   ES    │  │   Storage     │   │
│  └─────────┘  └─────────┘  └─────────┘  └───────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 12.2 AI 模块架构

```
┌─────────────────────────────────────────────────────────────┐
│                      AI 应用层                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │    Bots     │  │  Workflow   │  │   Material Center   │  │
│  │  智能助手   │  │  工作流编排  │  │     素材生成        │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      AI 能力层                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ AgentsFlex  │  │  TinyFlow   │  │    RAG Engine       │  │
│  │  Agent框架  │  │  工作流引擎  │  │   知识库检索        │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      模型接入层                              │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌───────┐  │
│  │ OpenAI  │ │  文心   │ │  通义   │ │  GLM   │ │Ollama │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └───────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 12.3 聊天协议

AIFlowy 使用自定义 SSE 协议 `aiflowy-chat` v1.1：

**协议域（Domain）：**

| Domain | 说明 |
|--------|------|
| llm | LLM 输出内容 |
| tool | 工具调用 |
| system | 系统消息 |
| business | 业务消息 |
| workflow | 工作流状态 |
| interaction | 用户交互 |
| debug | 调试信息 |

详细规范参见：`/aiflowy-chat-protocol.md`

---

## 13. 常见问题

### 13.1 前端启动报错

**问题：** `pnpm install` 失败

**解决：**
```bash
# 检查 Node 版本
node -v  # 需要 >= 20.10.0

# 检查 pnpm 版本
pnpm -v  # 需要 >= 9.12.0

# 清理缓存重试
pnpm clean
pnpm install
```

### 13.2 后端启动报错

**问题：** 数据库连接失败

**解决：**
1. 确认 MySQL 服务已启动
2. 检查数据库用户权限
3. 确认 `application.yml` 中的连接信息正确
4. 测试连接：`mysql -u root -p -h 127.0.0.1`

### 13.3 验证码不显示

**问题：** 登录验证码显示为空白

**解决：**
Docker 部署时需要安装字体：
```dockerfile
RUN apk --no-cache add ttf-dejavu fontconfig
```

### 13.4 文件上传失败

**问题：** 上传大文件失败

**解决：**
调整配置：
```yaml
spring:
  servlet:
    multipart:
      max-file-size: 100MB
      max-request-size: 100MB
```

Nginx 也需要调整：
```nginx
client_max_body_size 100M;
```

### 13.5 SSE 连接中断

**问题：** AI 对话流式输出中断

**解决：**
Nginx 配置：
```nginx
location /api {
    proxy_buffering off;
    proxy_cache off;
    proxy_read_timeout 3600s;
}
```

### 13.6 跨域问题

**问题：** 前后端分离部署时出现 CORS 错误

**解决：**
后端添加 CORS 配置，或通过 Nginx 统一代理前后端。

### 13.7 内存不足

**问题：** 前端构建时 OOM

**解决：**
项目已配置 8GB 内存限制：
```json
"build": "cross-env NODE_OPTIONS=--max-old-space-size=8192 turbo build"
```

如仍不足，可调整此值或分模块构建。

---

## 附录

### A. 端口列表

| 端口 | 服务 |
|------|------|
| 5173 | 前端开发服务器 |
| 8080 | 后端 API |
| 3306 | MySQL |
| 6379 | Redis |
| 9200 | Elasticsearch |
| 11434 | Ollama |

### B. 相关链接

- 官方文档：https://aiflowy.tech
- GitHub：https://github.com/aiflowy/aiflowy
- Gitee：https://gitee.com/aiflowy/aiflowy

### C. 版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v2.0.2 | 2025-12-26 | 当前版本，新增聊天协议、验证码等 |
| v1.x | - | 稳定版本分支 |

---

*文档编号：001 - AIFlowy 开发部署完整指南*
