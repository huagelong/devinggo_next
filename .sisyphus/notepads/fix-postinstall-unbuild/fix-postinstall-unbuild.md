# 修复 pnpm postinstall unbuild 找不到错误

## 问题描述
pnpm postinstall 脚本运行 `pnpm -r run stub --if-present` 时，多个子包的 stub 脚本执行 `pnpm unbuild --stub` 时出现 "Command not found: unbuild" 错误。

## 解决方案
为 8 个工作区包添加 `"unbuild": "catalog:"` 到 devDependencies

## 修改的文件

### 1. packages/@core/base/shared/package.json
- **状态**: 已有 devDependencies
- **修改**: 在现有 devDependencies 中追加 `"unbuild": "catalog:"`

### 2. internal/lint-configs/oxlint-config/package.json
- **状态**: 无 devDependencies
- **修改**: 创建 devDependencies 字段并添加 `"unbuild": "catalog:"`

### 3. internal/lint-configs/oxfmt-config/package.json
- **状态**: 无 devDependencies
- **修改**: 创建 devDependencies 字段并添加 `"unbuild": "catalog:"`

### 4. internal/node-utils/package.json
- **状态**: 无 devDependencies
- **修改**: 创建 devDependencies 字段并添加 `"unbuild": "catalog:"`

### 5. internal/vite-config/package.json
- **状态**: 已有 devDependencies
- **修改**: 在现有 devDependencies 中追加 `"unbuild": "catalog:"`

### 6. internal/lint-configs/eslint-config/package.json
- **状态**: 已有 devDependencies
- **修改**: 在现有 devDependencies 中追加 `"unbuild": "catalog:"`

### 7. scripts/turbo-run/package.json
- **状态**: 无 devDependencies
- **修改**: 创建 devDependencies 字段并添加 `"unbuild": "catalog:"`

### 8. scripts/vsh/package.json
- **状态**: 无 devDependencies
- **修改**: 创建 devDependencies 字段并添加 `"unbuild": "catalog:"`

## 验证结果
运行 `pnpm install` 成功：
- postinstall 脚本成功执行所有 stub 命令
- 无 "Command not found: unbuild" 错误
- 所有 8 个包的 stub 操作都完成

## 其他发现
- prepare 阶段有 lefthook 未找到的警告，但这不是本次修复的重点
- 所有子包正确引用了 catalog 中的 unbuild 版本（^3.6.1）
