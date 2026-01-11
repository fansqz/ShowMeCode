# ShowMeCode Admin

ShowMeCode 可视化 OJ 系统的管理后台前端。

## 技术栈

- **框架**: Vue 3 + TypeScript + Vite
- **UI 组件**: Element Plus
- **状态管理**: Pinia + 持久化插件
- **路由**: Vue Router 4
- **代码编辑**: Monaco Editor
- **文档编辑**: Vue Markdown Editor
- **HTTP 客户端**: Axios

## 目录结构

```
front-admin/
├── src/
│   ├── api/                      # API 接口层
│   │   ├── auth/                 # 认证接口
│   │   ├── account/              # 账户接口
│   │   ├── user/                 # 用户管理
│   │   ├── role/                 # 角色管理
│   │   ├── menu/                 # 菜单管理
│   │   ├── api/                  # API 权限管理
│   │   ├── visual-document/      # 可视化文档
│   │   ├── visual-document-bank/ # 文档库
│   │   └── common/               # 通用接口
│   │
│   ├── views/                    # 页面组件
│   │   ├── login/                # 登录页
│   │   ├── home/                 # 首页仪表板
│   │   ├── permissions/          # 权限管理模块
│   │   │   ├── user/             # 用户管理
│   │   │   ├── role/             # 角色管理
│   │   │   ├── menu/             # 菜单管理
│   │   │   └── api/              # API 管理
│   │   ├── visual-document/      # 可视化教程
│   │   │   ├── bank/             # 教程库
│   │   │   └── document/         # 教程编辑
│   │   └── 404/                  # 404 页面
│   │
│   ├── components/               # 可复用组件
│   │   ├── Svgicon/              # SVG 图标
│   │   ├── logo/                 # Logo
│   │   └── text-button/          # 文本按钮
│   │
│   ├── layout/                   # 布局组件
│   │   ├── menu/                 # 侧边菜单
│   │   ├── header/               # 顶部栏
│   │   │   ├── breadcrumb/       # 面包屑
│   │   │   └── setting/          # 设置菜单
│   │   └── main/                 # 主内容区
│   │
│   ├── store/modules/            # Pinia 状态
│   │   ├── user.ts               # 用户状态
│   │   └── visual-document.ts    # 文档状态
│   │
│   ├── router/                   # 路由配置
│   │   ├── routers.ts            # 路由定义
│   │   └── index.ts              # 路由实例
│   │
│   ├── utils/                    # 工具函数
│   │   ├── request.ts            # Axios 封装
│   │   ├── format.ts             # 格式化工具
│   │   └── time.ts               # 时间工具
│   │
│   ├── assets/                   # 静态资源
│   ├── styles/                   # 全局样式
│   ├── constants/                # 常量定义
│   ├── main.ts                   # 入口
│   ├── App.vue                   # 根组件
│   └── premisstion.ts            # 权限守卫
│
├── public/                       # 公共资源
├── vite.config.ts                # Vite 配置
├── tsconfig.json                 # TypeScript 配置
├── .env.development              # 开发环境变量
└── .env.production               # 生产环境变量
```

## 快速开始

### 环境要求

- Node.js 18+
- pnpm 8+

### 安装运行

```bash
# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 生产构建
pnpm build:pro
```

### 代码检查

```bash
pnpm lint          # ESLint 检查
pnpm format        # Prettier 格式化
pnpm lint:style    # StyleLint 检查
```

## 主要功能

### 权限管理 (RBAC)

- **用户管理**: 创建、编辑、删除用户，分配角色
- **角色管理**: 角色增删改查，权限分配
- **菜单管理**: 动态菜单配置
- **API 管理**: 接口权限控制

### 可视化教程管理

- **教程库**: 教程分类和组织
- **教程编辑**:
  - Monaco Editor 代码编辑
  - Markdown 文档编辑
  - 断点配置

### 动态路由

根据用户权限动态生成菜单和路由，实现权限控制。

## 页面路由

| 路由 | 页面 | 功能 |
|------|------|------|
| `/login` | 登录 | 用户登录、注册 |
| `/home` | 首页 | 仪表板 |
| `/manage/permissions/user` | 用户管理 | 用户增删改查 |
| `/manage/permissions/role` | 角色管理 | 角色权限配置 |
| `/manage/permissions/menu` | 菜单管理 | 菜单配置 |
| `/manage/permissions/api` | API 管理 | 接口权限 |
| `/manage/visual-document/bank` | 教程库 | 教程分类管理 |
| `/manage/visual-document/:bankID` | 教程编辑 | 编辑教程内容 |

## 状态管理

### User Store

```typescript
// 状态
token        // 登录令牌
menuRoutes   // 菜单路由
username     // 用户名
avatar       // 头像

// 方法
userLogin()   // 登录
userInfo()    // 获取用户信息
userLogout()  // 登出
```

### Visual Document Store

```typescript
// 状态
bankID       // 教程库 ID
id           // 文档 ID
title        // 标题
content      // 内容
codeList     // 代码块列表
```

## IDE 推荐配置

- [VS Code](https://code.visualstudio.com/)
- [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (禁用 Vetur)
- [TypeScript Vue Plugin](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin)
