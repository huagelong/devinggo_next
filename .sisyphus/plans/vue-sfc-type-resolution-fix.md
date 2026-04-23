# 修复 Vue SFC 编译器 "Failed to resolve extends base type" 错误

## TL;DR

> **Quick Summary**: 7 个 .vue 文件使用 `interface Props extends ExternalType` 模式时，Vue SFC 编译器无法解析外部包的类型。需要先诊断根因（版本差异 vs 配置差异），然后根据诊断结果采用正确的修复方案。
> 
> **Deliverables**:
> - 确认根因（依赖版本差异 or 编译配置缺失）
> - 消除 3 个当前报错的 Vite Pre-transform 错误
> - 预防 4 个潜在问题文件的同类错误
> 
> **Estimated Effort**: Medium
> **Parallel Execution**: YES - 3 waves
> **Critical Path**: Task 1 (诊断) → Task 2/3 (修复) → Task 4 (验证)

---

## Context

### Original Request
用户运行 `pnpm dev` 时看到 3 个 Vite Pre-transform 错误，都是 `[@vue/compiler-sfc] Failed to resolve extends base type`，发生在 `extra-menu.vue`、`menu.vue`、`mixed-menu.vue` 三个文件中。

### Interview Summary
**Key Discussions**:
- 用户确认修改代码后出现此问题（但菜单文件本身自 init 以来未修改）
- 用户选择**全面修复**（包括 4 个 shadcn-ui 潜在问题文件）
- 上游 vue-vben-admin 使用完全相同的代码模式且不需要 `@vue-ignore`

**Research Findings**:
- 上游版本差异：TypeScript ^6.0.2（上游）vs 5.9.3（本地），Vue ^3.5.32 vs 3.5.30
- 项目中完全没有 `vueCompilerOptions` 配置（上游也没有）
- `@vue-ignore` 会导致运行时行为变化（基类型属性变为 fallthrough attrs），**不安全**

### Metis Review
**Identified Gaps** (addressed):
- `@vue-ignore` 运行时语义被误解 → 已确认为不安全方案，改为先诊断根因
- 未追问上游为何能正常工作 → 已确认上游代码相同但版本不同
- 7 个组件运行时影响严重性分析 → 已完成，见下文

---

## Work Objectives

### Core Objective
消除 Vue SFC 编译器对跨包 `interface Props extends` 的类型解析错误，同时确保运行时行为不变。

### Concrete Deliverables
- 3 个报错文件的编译错误消除
- 4 个潜在问题文件的同类错误预防
- 运行时行为验证通过

### Definition of Done
- [x] `pnpm dev` 无 Pre-transform 错误
- [x] 涉及的 7 个组件在运行时行为正常
- [x] 4. 全面运行时行为验证
- [x] All "Must Have" present
- [x] All "Must NOT Have" absent (no `/* @vue-ignore */` in codebase)
- [x] 7 affected components render and function correctly
