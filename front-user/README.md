# ShowMeCode User

ShowMeCode 可视化 OJ 系统的用户端前端 - 在线编程学习平台。

## 技术栈

- **框架**: Vue 3 + TypeScript + Vite
- **UI 组件**: Element Plus
- **状态管理**: Pinia + 持久化插件
- **代码编辑**: Monaco Editor
- **数据可视化**: @antv/g6 + structv2
- **终端模拟**: xterm
- **文档渲染**: v-md-editor + Prism.js

## 目录结构

```
front-user/
├── src/
│   ├── api/                          # API 接口层
│   │   ├── auth/                     # 认证
│   │   ├── account/                  # 账户
│   │   ├── debug/                    # 调试 API (核心)
│   │   ├── visual/                   # 可视化 API
│   │   ├── saved-code/               # 代码保存
│   │   ├── visual-document/          # 学习文档
│   │   └── visual-document-bank/     # 文档库
│   │
│   ├── components/                   # 核心组件
│   │   ├── code-editor/              # 代码编辑器模块
│   │   │   ├── editor/               # Monaco 编辑器
│   │   │   │   ├── hooks/            # useVSCode Hook
│   │   │   │   ├── conf/             # 编辑器配置
│   │   │   │   ├── themes/           # 主题
│   │   │   │   └── utils/            # 断点、调试工具
│   │   │   └── console/              # 调试控制台
│   │   │       ├── coding-button/    # 执行按钮
│   │   │       ├── debug-button/     # 调试控制
│   │   │       └── debug-terminal/   # 调试终端
│   │   │           ├── console.vue   # 程序输出
│   │   │           ├── frames.vue    # 栈帧信息
│   │   │           └── variables.vue # 变量监视
│   │   │
│   │   ├── code-visual/              # 数据结构可视化
│   │   │   ├── visual/               # structv2 渲染
│   │   │   │   ├── layouter/         # 布局器
│   │   │   │   │   ├── array.ts      # 一维数组
│   │   │   │   │   ├── array2d.ts    # 二维数组
│   │   │   │   │   ├── binary-tree.ts # 二叉树
│   │   │   │   │   ├── graph.ts      # 图
│   │   │   │   │   └── link-list.ts  # 链表
│   │   │   │   └── type/             # 类型定义
│   │   │   ├── visual-setting/       # 可视化配置面板
│   │   │   ├── document/             # 文档面板
│   │   │   └── utils/                # 数据结构工具
│   │   │
│   │   ├── Svgicon/                  # SVG 图标
│   │   └── theme-switcher/           # 主题切换
│   │
│   ├── views/                        # 页面组件
│   │   ├── coding/                   # 编码工作区 (核心)
│   │   │   ├── index.vue             # 主组件
│   │   │   ├── visual.vue            # 可视化面板
│   │   │   └── saved-user-code.vue   # 保存代码
│   │   ├── learn/                    # 学习模块
│   │   │   ├── coding-panel/         # 编码面板
│   │   │   ├── document-panel/       # 文档面板
│   │   │   └── directory-tree.vue    # 目录树
│   │   ├── home/                     # 首页
│   │   ├── login/                    # 登录
│   │   ├── register/                 # 注册
│   │   ├── my-profile/               # 个人中心
│   │   └── account-setting/          # 账户设置
│   │
│   ├── store/modules/                # Pinia 状态
│   │   ├── coding.ts                 # 编码状态
│   │   ├── debug.ts                  # 调试状态 (核心)
│   │   ├── visual.ts                 # 可视化配置
│   │   ├── theme.ts                  # 主题状态
│   │   ├── user.ts                   # 用户状态
│   │   └── visual-document.ts        # 文档状态
│   │
│   ├── layout/                       # 布局
│   │   └── layouts/                  # 布局组件
│   │
│   ├── router/                       # 路由配置
│   ├── utils/                        # 工具函数
│   ├── constants/                    # 常量
│   ├── assets/                       # 静态资源
│   └── styles/                       # 全局样式
│
├── public/
├── vite.config.ts
└── package.json
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
pnpm code:check    # lint + format + style 检查
pnpm code:fix      # 自动修复
pnpm lint:fix      # ESLint 修复
pnpm format        # Prettier 格式化
```

## 核心功能

### 代码编辑器

基于 Monaco Editor 的在线代码编辑器:
- 语法高亮
- 代码补全
- 断点设置
- 多语言支持 (C, C++)

### 交互式调试

完整的调试功能:
- 设置/移除断点
- 单步执行 (Step In/Out/Over)
- 继续执行
- 变量监视
- 栈帧查看
- 程序输出

### 数据结构可视化

支持多种数据结构的动态可视化:
- 一维数组 (Array)
- 二维数组 (Array2D)
- 链表 (LinkList)
- 二叉树 (BinaryTree)
- 图 (Graph)

### 学习系统

集成的编程学习模块:
- 结构化知识文档
- 代码示例演示
- 交互式练习

## 页面路由

| 路由 | 页面 | 功能 |
|------|------|------|
| `/home` | 首页 | 平台入口 |
| `/coding` | 编码工作区 | 代码编辑和调试 |
| `/learn/:bankID` | 学习模块 | 文档 + 编程练习 |
| `/login` | 登录 | 用户登录 |
| `/register` | 注册 | 用户注册 |
| `/myprofile` | 个人中心 | 个人信息 |
| `/account/setting` | 账户设置 | 账户配置 |

## 状态管理

### Debug Store (调试核心)

```typescript
// 状态
id              // 调试会话 ID
status          // 'init' | 'compiled' | 'running' | 'stopped' | 'terminated'
breakpoints     // 断点行号列表
lineNum         // 当前调试行号
currentFrameID  // 当前栈帧 ID
outputs         // 调试输出

// 方法
setBreakpoints()  // 设置断点
updateLineNum()   // 更新当前行
```

### Visual Store (可视化配置)

```typescript
// 状态
action          // 可视化开关
isAIEnabled     // AI 自动识别
descriptionType // 当前数据结构类型
descriptions    // 各类型配置
```

### Coding Store

```typescript
// 状态
code      // 当前代码
language  // 编程语言
```

## 调试 API

```typescript
// 创建调试会话
reqCreateDebugSession()

// 启动调试
reqStart({ code, language, breakpoints })

// 单步调试
reqStepIn() / reqStepOut() / reqStepOver()

// 继续执行
reqContinue()

// 获取栈帧和变量
reqGetStackTrace()
reqGetFrameVariables()

// 监听调试事件 (SSE)
reqListenDebugEvent(sessionID)

// 终止调试
reqTerminate()
```

## IDE 推荐配置

- [VS Code](https://code.visualstudio.com/)
- [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (禁用 Vetur)
- [TypeScript Vue Plugin](https://marketplace.visualstudio.com/items?itemName=Vue.vscode-typescript-vue-plugin)
