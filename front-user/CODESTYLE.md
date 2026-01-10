# 代码风格规范

本项目使用统一的代码风格配置，确保团队成员之间的代码风格一致。

## 工具配置

### 1. ESLint

- 配置文件：`.eslintrc.cjs`
- 用于检查 JavaScript/TypeScript/Vue 代码质量
- 自动修复：`pnpm run lint:fix`

### 2. Prettier

- 配置文件：`.prettierrc.cjs`
- 用于代码格式化
- 格式化命令：`pnpm run format`

### 3. Stylelint

- 配置文件：`.stylelintrc.js`
- 用于检查 CSS/SCSS 代码质量
- 自动修复：`pnpm run lint:style`

### 4. EditorConfig

- 配置文件：`.editorconfig`
- 确保跨编辑器的基本格式一致

## VSCode 配置

### 推荐扩展

- Vue.volar - Vue 3 支持
- Vue.vscode-typescript-vue-plugin - Vue TypeScript 支持
- esbenp.prettier-vscode - Prettier 格式化
- dbaeumer.vscode-eslint - ESLint 支持
- stylelint.vscode-stylelint - Stylelint 支持

### 自动格式化设置

- 保存时自动格式化：已启用
- 默认格式化工具：Prettier
- 保存时自动修复 ESLint 错误：已启用

## 常用命令

```bash
# 检查代码质量
pnpm run code:check

# 自动修复代码格式和质量问题
pnpm run code:fix

# 仅格式化代码
pnpm run format

# 仅检查 ESLint
pnpm run lint

# 自动修复 ESLint 问题
pnpm run lint:fix

# 检查样式代码
pnpm run lint:style:check

# 自动修复样式问题
pnpm run lint:style
```

## 代码风格要点

### TypeScript/JavaScript

- 使用 `const` 和 `let`，避免 `var`
- 使用单引号
- 行尾分号
- 最大行宽：100 字符
- 使用相对路径导入

### Vue

- 组件名使用 PascalCase
- Props 使用 camelCase
- 使用 Vue 3 Composition API
- 模板中的属性使用 kebab-case

### CSS/SCSS

- 使用 2 空格缩进
- 使用单引号
- 属性顺序遵循 Stylelint 标准

## Git Hooks

项目配置了 Git Hooks：

- `pre-commit`：提交前自动格式化和检查代码
- `commit-msg`：检查提交信息格式

## 注意事项

1. 提交代码前会自动运行代码检查和格式化
2. 如果遇到无法自动修复的问题，请手动修复后再次提交
3. 建议在开发时安装推荐的 VSCode 扩展
4. 保持代码风格一致有助于提高代码可读性和维护性
