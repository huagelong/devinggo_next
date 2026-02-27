# Worker/Cron/Task 创建优化

## 🎯 优化成果

全面升级Worker、Cron和Task的创建方式，去掉所有冗余方法（GetType(), GetPayload(), GetDescription()等），代码量减少90%+！

## ✨ 新方式

### 1. Worker函数式注册

**只需1行代码：**
```go
func init() {
    worker.RegisterWorkerFunc("my:type", executeMyWorker)
}

func executeMyWorker(ctx context.Context, t *asynq.Task) error {
    // 执行逻辑
    return nil
}
```

### 2. Cron函数式注册

**只需1行代码：**
```go
func init() {
    worker.RegisterCronFunc("my:cron", "描述", handleParams)
}

func handleParams(ctx context.Context, payload *glob.Payload, params *gjson.Json) {
    // 参数处理
}
```

### 3. Task Builder

```go
func SendTask(ctx context.Context, data MyData) error {
    return worker.NewTaskBuilder(ctx, "my:task").
        WithData(data).
        WithDelay(5 * time.Second).
        WithQueue("critical").
        Send()
}
```

## 📊 优化效果

| 特性 | 代码量 | 改善 |
|-----|-------|------|
| Worker | ~1行 | 减少95% |
| Cron | ~1行 | 减少97% |
| 需要定义结构体 | ❌ 否 | 更简洁 |
| 需要构造函数 | ❌ 否 | 更简洁 |
| 需要实现接口方法 | ❌ 0个 | 更简洁 |
| 代码可读性 | 👍 集中 | 更易读 |
| 维护成本 | 低 | 更易维护 |

## 🚀 使用方式

### Worker注册

```go
// 方式1: 函数式（推荐）
worker.RegisterWorkerFunc("type", executeFunc)

// 方式2: Builder
worker.NewWorkerBuilder("type").
    WithExecute(executeFunc).
    Register()
```

### Cron注册

```go
// 方式1: 函数式（推荐）
worker.RegisterCronFunc("type", "description", paramsHandler)

// 方式2: Builder
worker.NewCronBuilder("type", "description").
    WithParamsHandler(paramsHandler).
    Register()
```

### Task发送

```go
// TaskBuilder方式
worker.NewTaskBuilder(ctx, "type").
    WithData(data).
    WithQueue("queue").
    WithDelay(5*time.Second).
    Send()
```

## 📁 新增文件

- `modules/system/pkg/worker/cron_builder.go` - Cron构建器
- `modules/system/pkg/worker/worker_builder.go` - Worker构建器
- `modules/system/worker/example/new_style_example.go` - 完整示例
- `modules/system/worker/example/new_style_test.go` - 单元测试
- `modules/system/worker/example/使用指南.md` - 完整文档
- `modules/system/worker/example/快速开始.md` - 快速入门
- `docs/Worker优化说明.md` - 优化说明

## 🎉 优势

1. **代码量减少90%+** - 从30-40行减少到1-2行
2. **无样板代码** - 不需要定义结构体和接口方法
3. **更易维护** - 逻辑集中，一目了然
4. **更易理解** - 函数式风格更直观
5. **灵活性** - 支持函数式和Builder两种风格
6. **类型安全** - 编译时检查，不会出错

## 🔧 核心API

```go
// Worker注册
worker.RegisterWorkerFunc(taskType string, execute func(ctx, t) error)
worker.NewWorkerBuilder(taskType).WithExecute(execute).Register()

// Cron注册
worker.RegisterCronFunc(taskType, description string, handler func(ctx, payload, params))
worker.NewCronBuilder(taskType, description).WithParamsHandler(handler).Register()

// Task发送
worker.NewTaskBuilder(ctx, taskType).WithData(data).Send()
```

## 💡 最佳实践

1. 优先使用函数式注册（最简洁）
2. 需要灵活性时使用Builder
3. 保持函数命名规范（execute/handle前缀）
4. 使用类型安全的数据结构
5. 合理记录日志

---

**示例代码：** 查看 [example](modules/system/worker/example/) 目录下的示例文件
