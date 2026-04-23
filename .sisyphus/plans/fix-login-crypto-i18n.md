# 修复登录页 crypto.subtle 崩溃和 i18n 缺失

## TL;DR

> **Quick Summary**: 修复两个独立 Bug：(1) 通过 HTTP+IP 访问时 crypto.subtle 为 undefined 导致登录崩溃；(2) authentication.loginFailed 和 common.loginFailed 翻译键缺失。
> 
> **Deliverables**:
> - 添加 webcrypto-liner-shim polyfill，使 crypto.subtle 在非安全上下文（HTTP+IP）下可用
> - 在 4 个 locale 文件中添加缺失的 loginFailed 翻译键
> 
> **Estimated Effort**: Quick
> **Parallel Execution**: YES - 2 waves
> **Critical Path**: Task 1 + Task 2 (parallel) → Task 3 (验证)

---

## Context

### Original Request
用户通过 `http://129.211.23.223:3199/`（HTTP + IP 地址）访问应用，登录时遇到两个问题：
1. `TypeError: Cannot read properties of undefined (reading 'importKey')` — crypto.subtle 在非安全上下文中不可用
2. `[intlify] Not found 'authentication.loginFailed' key in 'zh' locale messages` — 翻译键缺失

### Interview Summary
- 用户选择 **Polyfill 方案**（webcrypto-liner-shim）
- 确认需要同时修复 i18n 缺失问题

### Metis Review
**Identified Gaps** (addressed):
- `common.loginFailed` 也缺失（不只是 `authentication.loginFailed`）→ 已纳入修复范围
- 需要 4 个 locale 文件（不是 2 个）→ 已更新
- polyfill 应在应用入口加载 → 已调整方案
- 不要添加 `loginFailedDesc`（代码中未使用）→ 已排除

---

## Work Objectives

### Core Objective
使登录功能在 HTTP+IP 环境下正常工作，并修复缺失的翻译键。

### Concrete Deliverables
- crypto.subtle polyfill 在应用入口加载
- 4 个 locale 文件添加 loginFailed 翻译键

### Definition of Done
- [ ] 通过 HTTP+IP 访问时 crypto.subtle 可用（不再报 TypeError）
- [ ] 登录失败时显示正确的中文/英文错误消息

### Must Have
- crypto.subtle 在非安全上下文（HTTP+IP）下可用
- `authentication.loginFailed` 在 zh-CN 和 en-US 中有翻译
- `common.loginFailed` 在 zh-CN 和 en-US 中有翻译

### Must NOT Have (Guardrails)
- **禁止**修改 encryptPassword 的加密算法或流程
- **禁止**添加 `loginFailedDesc`（代码中未使用）
- **禁止**修改 `page.json` 中已有的 `page.profile.loginFailed`
- **禁止**修改 `profile/index.vue`
- **禁止**在 `index.html` 中添加 CDN script 标签
- **禁止**重构 encryptPassword 函数的 ECB 逻辑

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** - ALL verification is agent-executed.

### Test Decision
- **Infrastructure exists**: YES (vitest)
- **Automated tests**: None needed (配置和 i18n 修复)
- **Framework**: N/A

### QA Policy
- **Build**: Use Bash - `pnpm build` 验证构建成功
- **i18n**: Use Grep - 验证翻译键存在
- **Dependency**: Use Bash - 验证 polyfill 包安装

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately - 两个独立修复，MAX PARALLEL):
├── Task 1: 添加 crypto.subtle polyfill [quick]
├── Task 2: 修复 4 个 locale 文件的缺失翻译键 [quick]

Wave 2 (After Wave 1 - 验证):
├── Task 3: 验证修复结果 [quick]

Critical Path: Task 1 → Task 3, Task 2 → Task 3
Parallel Speedup: ~50% faster than sequential
```

### Dependency Matrix

| Task | Depends On | Blocks | Wave |
|------|-----------|--------|------|
| 1 | - | 3 | 1 |
| 2 | - | 3 | 1 |
| 3 | 1, 2 | - | 2 |

---

## TODOs

- [ ] 1. 添加 crypto.subtle polyfill

  **What to do**:
  - 在 `apps/backend/package.json` 中安装 `webcrypto-liner-shim` 依赖
  - 在 `apps/backend/src/main.ts`（或等效的应用入口文件）中添加 polyfill 加载逻辑：
    ```typescript
    // 在所有其他代码之前加载 polyfill
    if (typeof crypto !== 'undefined' && !crypto.subtle) {
      await import('webcrypto-liner-shim');
    }
    ```
  - 确保在 `createApp()` 之前完成 polyfill 加载

  **Must NOT do**:
  - 不要修改 encryptPassword 函数
  - 不要在 index.html 中添加 script 标签
  - 不要在 auth.ts 内部添加 polyfill 逻辑

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 2)
  - **Parallel Group**: Wave 1
  - **Blocks**: Task 3
  - **Blocked By**: None

  **References**:
  - `apps/backend/src/store/auth.ts:20-61` - encryptPassword 函数，了解 crypto.subtle 的使用方式
  - `apps/backend/src/main.ts` 或等效入口文件 - polyfill 加载位置
  - `apps/backend/package.json` - 添加依赖

  **Acceptance Criteria**:

  **QA Scenarios:**
  ```
  Scenario: polyfill 依赖已安装
    Tool: Bash
    Steps:
      1. 运行 `cd /home/wang/Documents/devinggo4vb/admin-ui && pnpm ls webcrypto-liner-shim --filter @vben/backend`
    Expected Result: 显示 webcrypto-liner-shim 版本号
    Evidence: .sisyphus/evidence/login-task-1-polyfill.txt

  Scenario: 应用入口加载 polyfill
    Tool: Grep
    Steps:
      1. 搜索 apps/backend/src/ 中包含 'webcrypto-liner-shim' 的文件
    Expected Result: 找到入口文件中的 polyfill 导入
    Evidence: .sisyphus/evidence/login-task-1-polyfill.txt
  ```

  **Commit**: YES
  - Message: `fix(crypto): add webcrypto-liner-shim polyfill for non-secure contexts`
  - Files: `apps/backend/package.json`, `apps/backend/src/main.ts`（或入口文件）, `pnpm-lock.yaml`

- [ ] 2. 修复 4 个 locale 文件的缺失翻译键

  **What to do**:
  - 在 `packages/locales/src/langs/zh-CN/authentication.json` 中添加 `"loginFailed": "登录失败"`
  - 在 `packages/locales/src/langs/en-US/authentication.json` 中添加 `"loginFailed": "Login Failed"`
  - 在 `packages/locales/src/langs/zh-CN/common.json` 中添加 `"loginFailed": "登录失败"`
  - 在 `packages/locales/src/langs/en-US/common.json` 中添加 `"loginFailed": "Login Failed"`
  - **注意**：key 格式参考 `authentication.json` 中的 `"loginSuccess"` 保持对称
  - **注意**：`common.json` 中参考已有的 `"deleteFailed"`, `"uploadFailed"` 等模式

  **Must NOT do**:
  - 不要添加 `loginFailedDesc`（代码中未使用）
  - 不要修改已有的任何 key
  - 不要修改 `page.json` 中的 `page.profile.loginFailed`

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 1)
  - **Parallel Group**: Wave 1
  - **Blocks**: Task 3
  - **Blocked By**: None

  **References**:
  - `packages/locales/src/langs/zh-CN/authentication.json` - 添加 loginFailed
  - `packages/locales/src/langs/en-US/authentication.json` - 添加 loginFailed
  - `packages/locales/src/langs/zh-CN/common.json` - 添加 loginFailed
  - `packages/locales/src/langs/en-US/common.json` - 添加 loginFailed
  - `apps/backend/src/store/auth.ts:129-131` - 查看 $t('common.loginFailed') 和 $t('authentication.loginFailed') 的使用

  **Acceptance Criteria**:

  **QA Scenarios:**
  ```
  Scenario: 4 个 locale 文件都包含 loginFailed
    Tool: Bash
    Steps:
      1. grep -c '"loginFailed"' packages/locales/src/langs/zh-CN/authentication.json
      2. grep -c '"loginFailed"' packages/locales/src/langs/en-US/authentication.json
      3. grep -c '"loginFailed"' packages/locales/src/langs/zh-CN/common.json
      4. grep -c '"loginFailed"' packages/locales/src/langs/en-US/common.json
    Expected Result: 每个文件返回 1
    Evidence: .sisyphus/evidence/login-task-2-i18n.txt
  ```

  **Commit**: YES
  - Message: `fix(i18n): add missing loginFailed translation keys`
  - Files: 4 个 locale JSON 文件

- [ ] 3. 验证修复结果

  **What to do**:
  - 验证 `webcrypto-liner-shim` 依赖已正确安装
  - 验证 4 个 locale 文件都包含 `loginFailed` 键
  - 验证应用入口正确加载 polyfill
  - 运行构建验证无错误

  **Must NOT do**:
  - 不要修改任何代码

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 2
  - **Blocked By**: Tasks 1, 2

  **Acceptance Criteria**:

  **QA Scenarios:**
  ```
  Scenario: 所有修复已正确应用
    Tool: Bash
    Steps:
      1. 验证 polyfill: `grep -r "webcrypto-liner-shim" apps/backend/src/`
      2. 验证 i18n: `grep -r '"loginFailed"' packages/locales/src/langs/`
      3. 验证依赖: `cat apps/backend/package.json | grep webcrypto-liner-shim`
    Expected Result: 所有验证通过
    Evidence: .sisyphus/evidence/login-task-3-verification.txt
  ```

  **Commit**: NO (验证任务)

---

## Commit Strategy

- **Task 1**: `fix(crypto): add webcrypto-liner-shim polyfill for non-secure contexts`
- **Task 2**: `fix(i18n): add missing loginFailed translation keys`

---

## Success Criteria

### Verification Commands
```bash
grep -r "webcrypto-liner-shim" apps/backend/src/  # Expected: polyfill import found
grep -r '"loginFailed"' packages/locales/src/langs/  # Expected: 4 matches
```

### Final Checklist
- [ ] crypto.subtle polyfill 在应用入口加载
- [ ] 4 个 locale 文件包含 loginFailed 翻译键
- [ ] 构建成功无错误
