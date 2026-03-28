# Admin UI 重构 T0 执行清单（可直接开工）

更新时间：2026-03-28  
对应阶段：`docs/TODO-AdminUI-Refactor-Optimized.md` 中的 `T0：基线收口`

---

## 1. T0 目标

把当前已经落地的 `user/post + 公共 CRUD` 收成“可复制基础版”，避免后续迁移 `role/menu/dept/dict` 时继续复制旧债。

T0 完成后，应满足：

- `system/user` 联调环境冒烟通过。
- `system/post` 与 `system/user` 成为可复用样板。
- 关键 API 和样板页的高频 `any` 明显收敛。
- 列表加载、提交、导入导出等场景不再只 `console.error`。
- 后续页面可以开始复用共享类型，而不是继续页面内临时拼类型。

---

## 2. 执行命令

本阶段至少会反复使用下面两个命令：

```bash
pnpm -F @vben/backend typecheck
pnpm -F @vben/backend run dev
```

建议节奏：

1. 每完成一个小批次文件，跑一次 `typecheck`。
2. `user/post` 页面相关修改合并后，跑一次本地手工冒烟。

---

## 3. 工作包拆解

## T0.1 `system/user` 联调冒烟收口

目标：关闭当前文档里唯一明确遗留的 Week 1 未收口项。

涉及文件：

- `admin-ui/apps/backend/src/views/system/user/index.vue`
- `admin-ui/apps/backend/src/views/system/user/use-user-crud.ts`
- `admin-ui/apps/backend/src/views/system/user/use-user-actions.ts`
- `admin-ui/apps/backend/src/views/system/user/components/user-modal.vue`
- `admin-ui/apps/backend/src/api/system/user.ts`

执行项：

- [ ] 验证查询、重置、分页、刷新。
- [ ] 验证新增、编辑、删除、批量删除。
- [ ] 验证回收站切换、恢复、彻底删除。
- [ ] 验证状态切换、重置密码、更新缓存、设置首页。
- [ ] 验证导入、导出、模板下载。
- [ ] 验证超级管理员保护逻辑。
- [ ] 记录联调中发现的字段不一致或接口异常，优先回填到 API/类型层，不在页面临时打补丁。

完成标准：

- `user` 页核心链路在联调环境可跑通；
- 若有缺陷，缺陷点已明确映射到具体文件。

---

## T0.2 建立共享类型基础

目标：给当前 CRUD 基础和下一批页面提供统一类型入口。

建议新建目录：

- `admin-ui/apps/backend/src/types/`

建议新增文件：

- `admin-ui/apps/backend/src/types/paging.ts`
- `admin-ui/apps/backend/src/types/common.ts`

建议内容：

- [ ] `paging.ts`
  - `PageQuery`
  - `PageInfo`
  - `PageResponse<T>`
- [ ] `common.ts`
  - `IdType`
  - `OptionItem`
  - `TreeOptionItem`
  - `BatchIdsPayload`
  - `StatusValue`

注意：

- 不追求一开始做成“大而全类型库”。
- 只抽当前 `user/post/role/dept/dict` 立刻能用上的最小集合。

完成标准：

- 后续 `api/system/*.ts` 与 `use-crud-page.ts` 可以直接引用共享类型；
- 不再继续在页面里大量写 `Record<string, any>`。

---

## T0.3 给公共 CRUD 基础补类型

目标：先稳住“会被所有页面复制”的基础层。

涉及文件：

- `admin-ui/apps/backend/src/composables/crud/use-crud-page.ts`
- `admin-ui/apps/backend/src/composables/crud/use-dict-options.ts`
- `admin-ui/apps/backend/src/components/crud/crud-toolbar.vue`

执行项：

- [ ] `use-crud-page.ts`
  - 用 `PageResponse<T>` 替代当前 `Promise<any>` / `response: any`
  - 给 `buildParams`、`resolveItems`、`resolveTotal` 明确泛型边界
  - 明确加载失败时的处理策略，避免只打印日志
- [ ] `use-dict-options.ts`
  - 收敛 `DictItem` 的索引签名范围
  - 让 `getDictOptions` / `getMultipleDictOptions` 返回类型更稳定
- [ ] `crud-toolbar.vue`
  - 去掉 `@change="(value: any) => ..."` 这类显式 `any`
  - 保持列显隐、刷新、回收站切换接口稳定

完成标准：

- 公共 CRUD 基础不再成为 `any` 扩散源头；
- 后续页面接入时能吃到更明确的泛型约束。

---

## T0.4 API 层最小类型化

目标：先把当前已存在的 5 个 `api/system` 文件从“能用”提升到“可继续扩展”。

涉及文件：

- `admin-ui/apps/backend/src/api/system/user.ts`
- `admin-ui/apps/backend/src/api/system/post.ts`
- `admin-ui/apps/backend/src/api/system/role.ts`
- `admin-ui/apps/backend/src/api/system/dept.ts`
- `admin-ui/apps/backend/src/api/system/dict.ts`

执行项：

- [ ] `user.ts`
  - 给列表查询、详情、保存、更新、状态切换、批量删除等补入参类型
  - 给导出参数补类型，不再使用裸 `Record<string, any>`
- [ ] `post.ts`
  - 给列表查询、保存、更新、状态、排序补类型
- [ ] `role.ts`
  - 即使暂时只保留 `getRoleList`，也先补最小返回类型
- [ ] `dept.ts`
  - 给树接口补最小树节点类型
- [ ] `dict.ts`
  - 给字典查询返回值补类型，和 `use-dict-options.ts` 对齐

建议做法：

- 优先引用 `src/types/*` 和当前模块 `model.ts`。
- 不为“未来所有场景”预留超多字段，只定义当前已消费字段。

完成标准：

- 当前 API 文件中高频 `any` 明显减少；
- 下阶段补 `role/menu/dept/dict` 页面时，不需要先推倒这层重写。

---

## T0.5 样板页债务清理：`system/user`

目标：把 `user` 真正整理成模板页，而不是“功能多但类型很松”的特例页。

涉及文件：

- `admin-ui/apps/backend/src/views/system/user/model.ts`
- `admin-ui/apps/backend/src/views/system/user/schemas.ts`
- `admin-ui/apps/backend/src/views/system/user/use-user-crud.ts`
- `admin-ui/apps/backend/src/views/system/user/use-user-actions.ts`
- `admin-ui/apps/backend/src/views/system/user/components/dept-tree.vue`
- `admin-ui/apps/backend/src/views/system/user/components/user-modal.vue`
- `admin-ui/apps/backend/src/views/system/user/index.vue`

执行项：

- [ ] `model.ts`
  - 去掉 `UserListItem` 上过宽的 `[key: string]: any`
  - 补当前页面真实使用字段
- [ ] `use-user-actions.ts`
  - 替换 `row/data/response` 等显式 `any`
  - 抽离文件名解析函数到轻量工具位置
  - 下载逻辑优先复用 `@vben/utils` 的 `downloadFileFromBlob`
  - catch 分支补用户可见错误提示，不再只有 `console.error`
- [ ] `components/user-modal.vue`
  - 补头像上传响应类型
  - 补详情回填数据类型
  - 补 `deptList/roleList/postList` 映射类型
- [ ] `components/dept-tree.vue`
  - 补树节点类型
  - 收敛过滤和遍历中的 `any`
- [ ] `index.vue`
  - 收敛 `roleOptions/postOptions/deptTreeData/statusOptions/...` 的 `any[]`
  - 收敛 `columns`、行回调、下拉项回调的 `any`

完成标准：

- `user` 页仍保持现有功能，但类型边界清晰很多；
- 后续复制 `user` 时不用连带复制一堆 `any` 和内联工具函数。

---

## T0.6 样板页债务清理：`system/post`

目标：把第二个样板页也整理到可复制状态，避免只有 `user` 一页能当模板。

涉及文件：

- `admin-ui/apps/backend/src/views/system/post/model.ts`
- `admin-ui/apps/backend/src/views/system/post/schemas.ts`
- `admin-ui/apps/backend/src/views/system/post/use-post-crud.ts`
- `admin-ui/apps/backend/src/views/system/post/components/post-modal.vue`
- `admin-ui/apps/backend/src/views/system/post/index.vue`

执行项：

- [ ] `model.ts`
  - 去掉 `PostListItem` 上过宽的 `[key: string]: any`
  - 补排序、状态、备注等已使用字段
- [ ] `use-post-crud.ts`
  - 查询参数类型化
- [ ] `post-modal.vue`
  - 给 `open(data?)`、`values` 补类型
  - catch 分支补明确错误提示
- [ ] `index.vue`
  - 收敛 `statusOptions`、`columns`、行回调中的 `any`
  - 如果动作继续增加，评估是否拆出 `use-post-actions.ts`

完成标准：

- `post` 页和 `user` 页都可以作为“标准列表 CRUD”样板；
- 两页的结构风格足够接近，适合提炼 v1.1 列表模板。

---

## T0.7 错误反馈与下载逻辑统一

目标：先统一最低限度的交互反馈，不做大封装。

涉及文件：

- `admin-ui/apps/backend/src/adapter/tdesign.ts`
- `admin-ui/apps/backend/src/views/system/user/use-user-actions.ts`
- `admin-ui/apps/backend/src/views/system/user/components/user-modal.vue`
- `admin-ui/apps/backend/src/views/system/post/index.vue`
- `admin-ui/apps/backend/src/views/system/post/components/post-modal.vue`
- `admin-ui/apps/backend/src/composables/crud/use-crud-page.ts`

可复用现有能力：

- `@vben/utils` 已暴露下载工具
- `admin-ui/apps/backend/src/adapter/tdesign.ts` 已统一导出 `message`

执行项：

- [ ] 统一提示入口：页面侧尽量改用 `#/adapter/tdesign` 导出的 `message`
- [ ] 统一 catch 约定：
  - 列表加载失败：给出用户可见错误提示
  - 新增/编辑/删除/恢复/状态切换失败：给出动作级错误提示
  - 保留日志，但不只留日志
- [ ] 统一下载约定：
  - `content-disposition` 文件名解析单独收口
  - Blob 下载优先复用共享下载工具

完成标准：

- 当前样板页里的失败路径，不再是“控制台知道、用户不知道”；
- 导出/模板下载逻辑不再只存在于 `user` 页私有实现中。

---

## T0.8 范围确认与占位页标记

目标：避免进入 T1/T2 时还在讨论“哪些页面算本期范围”。

涉及文件：

- `docs/TODO-AdminUI-Refactor-Optimized.md`
- `admin-ui/apps/backend/src/views/system/logs/loginLog.vue`
- `admin-ui/apps/backend/src/views/system/logs/operLog.vue`
- `admin-ui/apps/backend/src/views/system/logs/apiLog.vue`

执行项：

- [ ] 确认 `logs/*` 是否纳入后续迁移范围
- [ ] 若纳入：
  - 保留现有占位页
  - 在后续任务中明确归属到高定制批次
- [ ] 若不纳入：
  - 明确菜单侧不暴露
  - 文档中写清楚暂不交付，避免后续误判

完成标准：

- 代码目录和任务清单中的页面范围一致；
- 不会出现“代码里有页、菜单里能点、计划里却没人负责”的状态。

---

## 4. 推荐执行顺序

建议严格按下面顺序推进：

1. `T0.1` 联调冒烟  
2. `T0.2` 共享类型  
3. `T0.3` 公共 CRUD 基础补类型  
4. `T0.4` API 最小类型化  
5. `T0.5` `user` 样板债务清理  
6. `T0.6` `post` 样板债务清理  
7. `T0.7` 错误反馈与下载逻辑统一  
8. `T0.8` 范围确认

原因：

- 先把联调问题暴露出来，避免类型整理后又被真实接口推翻。
- 先稳基础层，再清理页面层。
- 先清理 `user/post`，再进入 `role/menu/dept/dict`，返工最少。

---

## 5. T0 完成判定

满足下面条件后，才建议进入 T1/T2：

- [ ] `pnpm -F @vben/backend typecheck` 通过
- [ ] `system/user` 联调冒烟通过
- [ ] `user/post` 样板页高频 `any` 已收敛
- [ ] 公共 CRUD 基础具备共享类型
- [ ] 当前样板页失败路径具备用户可见提示
- [ ] 页面范围边界已确认

---

## 6. 开工优先文件清单

如果要从今天直接开始改，建议先碰这批文件：

- `admin-ui/apps/backend/src/composables/crud/use-crud-page.ts`
- `admin-ui/apps/backend/src/api/system/user.ts`
- `admin-ui/apps/backend/src/api/system/post.ts`
- `admin-ui/apps/backend/src/views/system/user/model.ts`
- `admin-ui/apps/backend/src/views/system/user/use-user-actions.ts`
- `admin-ui/apps/backend/src/views/system/user/components/user-modal.vue`
- `admin-ui/apps/backend/src/views/system/post/model.ts`
- `admin-ui/apps/backend/src/views/system/post/index.vue`
- `admin-ui/apps/backend/src/views/system/post/components/post-modal.vue`

这批文件改完，T0 的主体价值就已经出来了。
