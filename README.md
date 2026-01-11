# ShowMeCode

ShowMeCode是一个可视化在线判题系统，集成代码执行、调试可视化和编程知识库，为编程学习提供一体化平台。

## 项目架构

```
ShowMeCode/
├── backend/         # Go 后端 API 服务
├── front-admin/     # Vue 3 管理后台
└── front-user/      # Vue 3 用户端
```

### 系统通信流程

```
Front-Admin ──┐
              │ HTTP/REST
              ▼
Front-User ──► Backend (Go) ──► MySQL/Redis
                    │
                    │ Docker + DAP/TCP
                    ▼
              Go-Debugger (GDB wrapper)
```

## 核心功能

- **代码执行与判题**: Docker 沙箱隔离，支持 C/C++/Java/Go
- **可视化调试**: 断点设置、单步执行、变量监视、栈帧查看
- **数据结构可视化**: 数组、链表、二叉树、图的动态展示
- **AI 代码分析**: 自动识别数据结构，代码解释与建议
- **知识库系统**: 结构化编程学习文档

## 技术栈

| 模块 | 技术 |
|------|------|
| 后端 | Go, Gin, GORM, MySQL, Redis, Docker |
| 前端 | Vue 3, TypeScript, Vite, Element Plus, Pinia |
| 编辑器 | Monaco Editor |
| 可视化 | @antv/g6, structv2 |
| 调试 | DAP 协议, GDB, Delve |
| AI | 火山引擎 Doubao, OpenAI |

## 快速开始

### 环境要求

- Go 1.23+
- Node.js 18+
- pnpm 8+
- MySQL 8.0+
- Redis 6.0+
- Docker
- Linux (后端需要 cgroups 支持)

### 启动后端

```bash
cd backend
go mod download
go build -o fancode .
./fancode
```

### 启动管理后台

```bash
cd front-admin
pnpm install
pnpm dev
```

### 启动用户端

```bash
cd front-user
pnpm install
pnpm dev
```

### Docker 部署

```bash
# 构建调试器镜像
cd backend
docker build -t go-debugger -f Dockerfile-debugger .

# 构建并启动后端
docker build -t showmecode -f Dockerfile .
docker run -d \
  -v /usr/bin/docker:/usr/bin/docker \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/fanCode:/var/fanCode \
  --network=host \
  --name showmecode \
  showmecode
```

## 项目模块

### Backend

Go 后端服务，提供 RESTful API。

主要功能:
- 用户认证与权限管理 (RBAC)
- 代码执行与判题引擎
- 可视化调试服务 (DAP 协议)
- 学习文档管理
- AI 代码分析

详见 [backend/README.md](./backend/README.md)

### Front-Admin

Vue 3 管理后台。

主要功能:
- 用户/角色/权限管理
- 菜单配置
- 教程库管理
- 可视化教程编辑

详见 [front-admin/README.md](./front-admin/README.md)

### Front-User

Vue 3 用户端，在线编程学习平台。

主要功能:
- Monaco 代码编辑器
- 交互式调试
- 数据结构可视化
- 编程学习文档

详见 [front-user/README.md](./front-user/README.md)

## API 路由

| 前缀 | 功能 |
|------|------|
| `/api/auth/*` | 认证 |
| `/api/account/*` | 账户 |
| `/api/debug/*` | 调试 |
| `/api/visual/*` | 可视化 |
| `/api/saved-code/*` | 代码保存 |
| `/api/visual-document/*` | 学习文档 |
| `/manage/*` | 管理后台 |

## 开发指南

### 代码规范

前端项目使用 ESLint + Prettier + StyleLint 进行代码检查:

```bash
# front-admin
pnpm lint && pnpm format && pnpm lint:style

# front-user
pnpm code:check  # 或 pnpm code:fix 自动修复
```

后端使用 Go 标准规范:

```bash
go fmt ./...
go vet ./...
```

### 测试

```bash
# 后端测试
cd backend
go test ./...

# 指定包测试
go test -v ./service/system_service/...
```

## 许可证

MIT License
