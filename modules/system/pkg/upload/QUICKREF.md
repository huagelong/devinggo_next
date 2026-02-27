# Upload 包快速参考

## 🚀 新特性速览

### 1. 链式调用 API
```go
// 最简单的方式
result, err := upload.NewUploader(ctx).UploadFile(file)

// 完整配置
result, err := upload.NewUploader(ctx).
    UseCloudStorage().      // 云存储
    SetRandomName(true).    // 随机文件名
    SetValidateType(true).  // 类型验证
    UploadFile(file)
```

### 2. 统一验证器
```go
validator := upload.NewValidator(ctx)
err := validator.ValidateFile(file)      // 验证文件
err := validator.ValidateImage(file)     // 验证图片
err := validator.ValidateExtension("jpg", true)  // 验证扩展名
```

### 3. 常量定义
```go
upload.StorageModeLocal     // 本地存储 = 1
upload.StorageModeCloud     // 云存储 = 2

upload.ResourceTypeImage    // 图片类型
upload.ResourceTypeVideo    // 视频类型
upload.ResourceTypeAudio    // 音频类型
```

## 📝 常用示例

### 上传文件
```go
// 本地存储
result, err := upload.NewUploader(ctx).
    UseLocalStorage().
    UploadFile(file)

// 云存储
result, err := upload.NewUploader(ctx).
    UseCloudStorage().
    UploadFile(file)
```

### 上传图片
```go
result, err := upload.NewUploader(ctx).
    SetRandomName(true).
    UploadImage(file)
```

### 从URL保存
```go
result, err := upload.NewUploader(ctx).
    UseLocalStorage().
    SaveFromURL("https://example.com/image.jpg")
```

### 自定义文件名
```go
result, err := upload.NewUploader(ctx).
    SetCustomName("myfile.jpg").
    SetRandomName(false).
    UploadFile(file)
```

### 禁用类型验证
```go
result, err := upload.NewUploader(ctx).
    SetValidateType(false).
    UploadFile(file)
```

## 🔄 迁移指南

### 文件上传
```go
// 旧方式
input := &req.FileUploadInput{File: file, StorageMode: 1, RandomName: true}
result, err := upload.Upload(ctx, input)

// 新方式
result, err := upload.NewUploader(ctx).UseLocalStorage().UploadFile(file)
```

### 文件验证
```go
// 旧方式（已废弃）
err := upload.CheckFileMineType(ctx, file)  // 注意拼写错误

// 新方式
err := upload.NewValidator(ctx).ValidateFile(file)

// 或使用修正后的兼容函数
err := upload.CheckFileMimeType(ctx, file)  // 已标记废弃
```

### 网络图片
```go
// 旧方式
result, err := upload.SaveNetworkImage(ctx, 1, url, true)

// 新方式
result, err := upload.NewUploader(ctx).UseLocalStorage().SaveFromURL(url)
```

## ⚡ 优化亮点

### 代码质量
- ✅ 消除所有硬编码字符串和数字
- ✅ 修正拼写错误（Mine → Mime）
- ✅ 统一错误处理和错误上下文
- ✅ 改进资源管理（使用 defer）
- ✅ 提取公共函数减少重复

### API 设计
- ✅ 支持链式调用，代码更流畅
- ✅ 提供合理默认值，简化使用
- ✅ 完全向后兼容，平滑升级
- ✅ 类型安全的常量定义
- ✅ 职责分离，易于扩展

### 用户体验
- ✅ API 更简洁直观
- ✅ 错误信息更详细
- ✅ 文档更完善
- ✅ 示例代码丰富
- ✅ 提供迁移指南

## 📚 文档

- [README.md](readme.md) - 完整使用教程
- [OPTIMIZATION.md](OPTIMIZATION.md) - 详细优化说明
- [examples_test.go](examples_test.go) - 代码示例

## 🎯 推荐做法

1. **新项目**: 直接使用 `NewUploader()` API
2. **旧项目**: 逐步迁移，旧API仍可用
3. **验证**: 使用 `NewValidator()` 替代旧的 Check 函数
4. **常量**: 使用预定义常量替代魔法数字
5. **错误**: 检查并处理所有错误返回值

## ⚠️ 注意

- 旧的 `CheckFileMineType` 更名为 `CheckFileMimeType`
- 旧的 `CheckImageMineType` 更名为 `CheckImageMimeType`
- 旧函数已标记废弃但仍可用
- 推荐尽快迁移到新API
