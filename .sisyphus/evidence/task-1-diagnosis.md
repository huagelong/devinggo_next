# 诊断报告：Vue SFC "Failed to resolve extends base type" 错误

**诊断时间**: 2026-04-19  
**诊断人**: Sisyphus Junior (自动化诊断)

---

## 1. 当前错误状态

### Dev Server 状态
- **启动结果**: 成功（Vite v8.0.0，端口 6000）
- **错误数量**: 3 个 "Failed to resolve extends base type" 错误
- **错误级别**: Pre-transform error（非致命，dev server 可运行）

### 受影响文件
| 文件 | 出错行 | extends 类型 |
|------|--------|-------------|
| `packages/effects/layouts/src/basic/menu/extra-menu.vue` | 12 | `MenuProps` |
| `packages/effects/layouts/src/basic/menu/menu.vue` | 8 | `MenuProps` |
| `packages/effects/layouts/src/basic/menu/mixed-menu.vue` | 13 | `NormalMenuProps` |

### 错误模式
所有错误都遵循同一模式：
```vue
<script lang="ts" setup>
import type { MenuProps } from '@vben-core/menu-ui';
interface Props extends MenuProps { ... }
</script>
```

---

## 2. 根因分析

### 根因结论：**pnpm 依赖安装不完整 + TypeScript tsconfig extends 解析失败**

### 详细机制

**直接原因**: `@vue/compiler-sfc` 在解析 `interface Props extends MenuProps` 时，无法解析 `MenuProps` 的完整类型定义。

**根本原因**（逐层深入）:

1. **TypeScript `moduleResolution` 回退到 Node10 模式**
   - `packages/effects/layouts/tsconfig.json` 使用 `"extends": "@vben/tsconfig/web.json"`
   - 该文件定义了 `"moduleResolution": "bundler"`
   - **但 TypeScript 5.9.3 无法找到 `@vben/tsconfig/web.json`**
   - 导致 extends 静默失败，`moduleResolution` 为 `undefined`（等同于 Node10）

2. **为什么找不到 `@vben/tsconfig/web.json`？**
   - `packages/effects/layouts/package.json` 中没有声明 `"@vben/tsconfig": "workspace:*"` 作为 devDependencies
   - pnpm 严格模式不会为未声明的依赖创建符号链接
   - 验证：`packages/effects/layouts/node_modules/@vben/` 中没有 `tsconfig` 目录
   - 验证：`ts.resolveModuleName('@vben/tsconfig/web.json', ...)` 返回 `NOT FOUND`

3. **Node10 模式下模块解析的差异**
   - `moduleResolution: "bundler"` → 正确解析 `@vben-core/menu-ui` → `src/index.ts`（含类型定义）
   - `moduleResolution: "node"` (Node10) → 错误解析 `@vben-core/menu-ui` → `dist/index.mjs`（编译后 JS）
   - 实测验证：
     ```
     Node10 模式: resolvedFileName = "dist/index.mjs"
     Bundler 模式: resolvedFileName = "src/index.ts"
     ```

4. **compiler-sfc 使用 TS 解析结果**
   - `@vue/compiler-sfc` 调用 `ts.resolveModuleName()` 获取文件路径
   - 拿到 `dist/index.mjs` 后调用 `fileToScope()` 解析类型
   - `.mjs` 文件中没有 TypeScript interface 定义
   - 递归解析失败 → 抛出 "Failed to resolve extends base type"

### 因果链图
```
package.json 缺少 @vben/tsconfig 依赖
    → pnpm 不创建符号链接
    → TS 无法解析 tsconfig extends "@vben/tsconfig/web.json"
    → moduleResolution 回退为 Node10
    → TS 解析 @vben-core/menu-ui 到 dist/index.mjs（而非 src/index.ts）
    → compiler-sfc 无法从 .mjs 文件提取 interface 定义
    → "Failed to resolve extends base type"
```

---

## 3. 版本对比表

### 关键依赖版本对比

| 依赖 | 本地 catalog | 本地实际安装 | 上游 catalog | 差异级别 |
|------|-------------|-------------|-------------|---------|
| typescript | ^5.9.3 | 5.9.3 / 5.8.2 | ^6.0.2 | **MAJOR** |
| vue | ^3.5.30 | 3.5.30 | ^3.5.32 | patch |
| @vitejs/plugin-vue | ^6.0.5 | 6.0.5 / 5.2.4 | ^6.0.6 | patch |
| @vue/compiler-sfc | (随 vue) | 3.5.30 | (随 vue) | patch |
| vite | ^8.0.0 | 8.0.0 / 5.4.21 | ^8.0.8 | patch |
| vue-tsc | ^3.2.5 | 3.2.5 / 3.2.6 | ^3.2.6 | patch |
| unplugin-vue | **未配置** | **未安装** | ^7.1.1 | **缺失** |

### 版本差异结论

TypeScript 5.9.3 → 6.0.2 是 MAJOR 版本跳跃，但**这不是导致当前错误的直接原因**。即使升级到 TS 6.x，如果 pnpm 依赖安装问题不解决，tsconfig extends 仍然会失败。

上游添加了 `unplugin-vue: ^7.1.1`（用于 rolldown 构建），但 dev server 仍使用 `@vitejs/plugin-vue`。

---

## 4. 配置对比表

### tsconfig/base.json 对比
**结论：完全一致**。本地和上游的 `internal/tsconfig/base.json` 内容完全相同。

### tsconfig/web.json 对比
**结论：完全一致**。本地和上游的 `internal/tsconfig/web.json` 内容完全相同。

### vite 插件配置（internal/vite-config/src/plugins/index.ts）
本地 `viteVue` 配置：
```typescript
viteVue({
  script: {
    defineModel: true,
    // propsDestructure: true,
  },
})
```
与上游一致（代码相同）。

### package.json 依赖差异
**关键差异**：上游在 `packages/effects/layouts/package.json` 中可能声明了 `@vben/tsconfig` 作为 devDependencies（需要进一步验证），或者上游的 pnpm 配置允许隐式访问 workspace 包。

---

## 5. 验证数据

### 实验验证

#### 实验 1：TypeScript 模块解析（Node10 vs Bundler）
```
containingFile: packages/effects/layouts/src/basic/menu/extra-menu.vue
moduleResolution: undefined (Node10)
  → @vben-core/menu-ui → dist/index.mjs ❌

moduleResolution: Bundler
  → @vben-core/menu-ui → src/index.ts ✅
```

#### 实验 2：tsconfig extends 解析
```
ts.findConfigFile → packages/effects/layouts/tsconfig.json
extends: "@vben/tsconfig/web.json"
ts.resolveModuleName("@vben/tsconfig/web.json") → NOT FOUND
```

#### 实验 3：node_modules 符号链接
```
packages/effects/layouts/node_modules/@vben/ 目录内容:
  constants, hooks, icons, locales, preferences, stores, types, utils
  （没有 tsconfig！）
```

---

## 6. 推荐修复路径

### 方案 A：添加缺失的 devDependencies（推荐）

在 `packages/effects/layouts/package.json` 的 `devDependencies` 中添加：
```json
"@vben/tsconfig": "workspace:*"
```

然后运行 `pnpm install` 重新安装依赖。

这会创建符号链接，使 TypeScript 能正确解析 tsconfig extends，从而启用 `moduleResolution: "bundler"`，最终使 compiler-sfc 正确解析 workspace 包中的类型。

**优点**: 最小改动，直接修复根因  
**风险**: 低  
**验证**: 修复后 `pnpm dev` 应不再有 "Failed to resolve extends base type" 错误

### 方案 B：全面同步上游版本（备选）

将 `pnpm-workspace.yaml` 中的 catalog 同步到上游最新版本，特别是：
- `typescript: ^6.0.2`
- `vue: ^3.5.32`
- `@vitejs/plugin-vue: ^6.0.6`
- `vite: ^8.0.8`
- 添加 `unplugin-vue: ^7.1.1`

**优点**: 与上游保持一致  
**风险**: TypeScript 6.x 是大版本升级，可能有破坏性变更  
**注意**: 即使升级版本，如果 devDependencies 问题不解决，错误仍可能存在

### 方案 C：同时执行 A + B（最安全）

先执行方案 A 修复根因，验证错误消失后，再逐步同步上游版本。

---

## 7. 不推荐方案

- **`@vue-ignore` 注释**: 不推荐。会导致基类型属性变为 fallthrough attrs，改变运行时行为
- **`skipExtendsCheck`**: 不推荐。会掩盖所有 extends 解析问题
