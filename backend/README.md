# ShowMeCode Backend

## 技术栈

- **框架**: Gin Web Framework
- **ORM**: GORM + MySQL
- **缓存**: Redis
- **依赖注入**: Google Wire
- **限流**: Alibaba Sentinel
- **对象存储**: 腾讯云 COS
- **AI 服务**: 火山引擎 Doubao / OpenAI

## 目录结构

```
backend/
├── conf/                    # 配置文件
│   ├── config_local.ini    # 本地开发配置
│   ├── config_online.ini   # 线上配置
│   └── acm_template/       # ACM 题目模板
│
├── controller/              # HTTP 请求处理层
│   ├── admin/              # 管理后台接口
│   └── user/               # 用户端接口
│
├── routers/                 # 路由定义
│   ├── admin/              # 管理后台路由
│   └── user/               # 用户端路由
│
├── service/                 # 业务逻辑层
│   ├── common_service/     # 公共服务 (认证、账户、文件)
│   ├── system_service/     # 系统管理 (用户、角色、菜单、API)
│   ├── user_saved_code_service/  # 代码保存
│   ├── visual_document_service/  # 学习文档
│   └── visual_debug_servcie/     # 可视化调试
│       ├── debug_core/     # 调试器核心实现
│       └── ai_analyze_core/  # AI 分析
│
├── dao/                     # 数据访问层
├── models/                  # 数据模型
│   ├── dto/                # 数据传输对象
│   ├── po/                 # 持久化对象
│   └── vo/                 # 视图对象
│
├── interceptor/             # 中间件
│   ├── cors.go             # CORS 跨域
│   ├── logger.go           # 日志记录
│   ├── rate_limiter.go     # 限流
│   └── request.go          # 认证授权
│
├── common/                  # 共享模块
│   ├── config/             # 配置管理
│   ├── ai_provider/        # AI 接口 (OpenAI/火山引擎)
│   ├── error/              # 错误处理
│   ├── file_store/         # 文件存储 (COS)
│   └── logger/             # 日志 (Graylog)
│
├── main.go                  # 入口
├── wire.go                  # 依赖注入配置
└── showmecode.sql             # 数据库脚本
```

## 快速开始

### 环境要求

- Go 1.23+
- MySQL 8.0+
- Redis 6.0+
- Linux (需要 cgroups 支持代码沙箱)
- Docker (调试功能需要)

### 安装依赖工具

### 安装运行

```bash
# 安装依赖
go mod download

# 配置（可选，默认使用线上配置）
cp conf/config_local.ini conf/config.ini
# 修改 config.ini 中的数据库和 Redis 配置

# 编译运行
go build -o fancode .
./fancode
```

### Docker 部署

1. 创建调试器镜像

```bash
docker build -t go-debugger -f Dockerfile-debugger .
```

2. 创建并启动后端容器

```bash
# 构建镜像
docker build -t showmecode -f Dockerfile .

# 启动容器 (需要 Docker-in-Docker 权限)
docker run -d \
  -v /usr/bin/docker:/usr/bin/docker \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/fanCode:/var/fanCode \
  --network=host \
  --name showmecode \
  showmecode
```

### 测试

```bash
go test ./...
go test -v ./service/system_service/...  # 运行指定包测试
```

## API 路由

| 前缀 | 功能 |
|------|------|
| `/api/auth/*` | 认证 (登录、注册、验证码) |
| `/api/account/*` | 账户管理 |
| `/api/saved-code/*` | 代码保存 |
| `/api/debug/*` | 调试接口 |
| `/api/visual/*` | 可视化接口 |
| `/api/visual-document/*` | 学习文档 |
| `/manage/*` | 管理后台 |

## 核心功能

### 代码执行与判题

- 基于 Docker 的沙箱隔离执行
- Linux cgroups 资源限制
- 支持 C、C++、Java、Go 语言

### 可视化调试

- DAP (Debug Adapter Protocol) 协议
- 集成 GDB (C/C++) 和 Delve (Go)
- 实时变量监视和栈帧查看
- 数据结构可视化 (数组、链表、树、图)

### AI 代码分析

- 自动识别数据结构类型
- 代码解释和优化建议
- 支持火山引擎和 OpenAI

## 配置说明

主要配置项 (`conf/config.ini`):

```ini
[server]
port = :8080
mode = release

[mysql]
host = localhost:3306
database = fan_code
username = root
password = xxx

[redis]
host = localhost:6379

[ai]
provider = volcengine  # 或 openai
api_key = xxx
```

## 数据库

数据库 SQL 脚本位于 `./showmecode.sql`，可直接导入创建数据库表结构。
