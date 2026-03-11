---
name: new-module
description: "使用场景：在 DevingGo-Light 项目中创建全新的业务模块。当用户需要添加新的功能模块（如电商模块、CMS模块、工单模块等）而不是在现有 system 模块中添加功能时使用。涵盖完整的模块搭建流程：目录结构、模块注册、路由配置、bootstrap 绑定。"
---

# 新模块创建工作流

```bash
make gen-module name={module}        # 自动创建目录结构（推荐）
make clone-module name={new} source=system  # 克隆已有模块
```

## 手动创建

### 1. `modules/{module}/module.go`

```go
package {module}

func init() {
    m := &{module}Module{}
    m.Name = "{module}"
    modules.Register(m)
}

type {module}Module struct{ modules.BaseModule }

func (m *{module}Module) Start(ctx context.Context, s *ghttp.Server) error {
    s.Group("/", func(group *ghttp.RouterGroup) {
        router.BindController(group)
    })
    return nil
}
```

### 2. `modules/{module}/router/{module}/router.go`

```go
func BindController(group *ghttp.RouterGroup) {
    group.Group("/{module}", func(group *ghttp.RouterGroup) {
        group.Bind(
            // {module}.XxxController,
        ).Middleware(service.Middleware().AdminAuth)
    })
}
```

### 3. Bootstrap 注册

**`modules/bootstrap/modules/modules.go`**：
```go
_ "devinggo/modules/{module}"
```

**`modules/bootstrap/logic/logic.go`**：
```go
_ "devinggo/modules/{module}/logic/{module}"
```

### 4. `hack/config.yaml` 添加 service 生成配置

```yaml
gfcli:
  gen:
    service:
      - srcFolder: "modules/{module}/logic"
        dstFolder: "modules/{module}/service"
        watchFile:  "modules/{module}/logic"
```

### 5. 生成 Service 接口

```bash
make service
```

## 注意事项

- `m.Name` 在所有模块中必须唯一
- 模块间通信通过 `internal/service` 接口，不直接调用对方 Logic
- 新模块的 `service/` 由 `gf gen service` 生成，禁止手动修改

