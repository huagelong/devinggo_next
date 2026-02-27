# Upload 包使用教程

## 简介

Upload 包是 DevingGo 框架中的文件上传工具包，提供了完整的文件上传解决方案，支持本地存储和云存储（S3兼容），包含普通文件上传、分片上传和网络图片保存等功能。

## 主要功能

- ✅ 普通文件上传
- ✅ 分片上传（支持大文件）
- ✅ 网络图片保存
- ✅ 本地存储和云存储
- ✅ 文件类型校验
- ✅ MD5 校验
- ✅ 自动资源类型识别
- ✅ 文件大小格式化

## 文件说明

| 文件名 | 说明 |
|--------|------|
| upload.go | 核心上传功能，处理普通文件上传 |
| chunk_upload.go | 分片上传功能，用于大文件上传 |
| cloud.go | 云存储配置和操作 |
| save_network_image.go | 网络图片保存功能 |

## 快速开始

### 1. 普通文件上传

```go
package main

import (
    "context"
    "devinggo/modules/system/model/req"
    "devinggo/modules/system/pkg/upload"
)

func uploadFile(ctx context.Context, file *ghttp.UploadFile) {
    // 构建上传参数
    input := &req.FileUploadInput{
        File:        file,
        StorageMode: 1,         // 1: 本地存储, 2: 云存储
        Name:        "",        // 可选：自定义文件名
        RandomName:  true,      // 是否使用随机文件名
    }
    
    // 执行上传
    result, err := upload.Upload(ctx, input)
    if err != nil {
        // 处理错误
        return
    }
    
    // 使用上传结果
    println("文件URL:", result.Url)
    println("文件MD5:", result.Hash)
    println("文件大小:", result.SizeInfo)
}
```

### 2. 分片上传（大文件）

```go
func chunkUpload(ctx context.Context, file *ghttp.UploadFile, index, total int64, hash string) {
    // 构建分片上传参数
    input := &req.ChunkUploadInput{
        File:        file,
        StorageMode: 1,
        Index:       index,     // 当前分片索引（从1开始）
        Total:       total,     // 总分片数
        Hash:        hash,      // 文件唯一标识
        Name:        "large_file.zip",
        Ext:         "zip",
        Type:        "application/zip",
        RandomName:  true,
    }
    
    // 执行分片上传
    result, err := upload.ChunkUpload(ctx, input)
    if err != nil {
        // 处理错误
        return
    }
    
    // 注意：只有最后一个分片上传完成后才会返回结果
    if result != nil {
        println("文件上传完成:", result.Url)
    }
}
```

### 3. 保存网络图片

```go
func saveNetworkImage(ctx context.Context, imageUrl string) {
    // 保存网络图片到本地
    result, err := upload.SaveNetworkImage(
        ctx,
        1,              // 存储模式
        imageUrl,       // 图片URL
        true,           // 是否使用随机文件名
    )
    if err != nil {
        // 处理错误
        return
    }
    
    println("图片已保存:", result.Url)
    println("本地路径:", result.LocalPath)
}
```

## 核心函数说明

### Upload - 普通文件上传

上传单个文件到服务器，支持本地存储和云存储。

**函数签名:**
```go
func Upload(ctx context.Context, in *req.FileUploadInput) (*res.SystemUploadFileRes, error)
```

**参数说明:**
- `ctx`: 上下文对象
- `in.File`: 上传的文件对象
- `in.StorageMode`: 存储模式（1: 本地, 2: 云存储）
- `in.Name`: 自定义文件名（可选）
- `in.RandomName`: 是否使用随机文件名

**返回值:**
```go
type SystemUploadFileRes struct {
    StorageMode int     // 存储模式
    OriginName  string  // 原始文件名
    ObjectName  string  // 对象名称（保存后的文件名）
    Hash        string  // 文件MD5值
    MimeType    string  // 资源类型
    StoragePath string  // 存储路径
    Suffix      string  // 文件后缀
    SizeByte    int64   // 文件大小（字节）
    SizeInfo    string  // 格式化的文件大小
    LocalPath   string  // 本地路径
    Url         string  // 访问URL
}
```

### ChunkUpload - 分片上传

用于上传大文件，将文件分割成多个分片依次上传，最后自动合并。

**函数签名:**
```go
func ChunkUpload(ctx context.Context, in *req.ChunkUploadInput) (*res.SystemUploadFileRes, error)
```

**参数说明:**
- `in.Index`: 当前分片索引（从1开始）
- `in.Total`: 总分片数
- `in.Hash`: 文件唯一标识（用于关联各个分片）
- `in.Ext`: 文件扩展名
- `in.Type`: 文件MIME类型

**工作流程:**
1. 客户端将大文件分割成多个分片
2. 按顺序上传每个分片
3. 所有分片上传完成后自动合并
4. 返回完整文件信息

### SaveNetworkImage - 保存网络图片

从网络URL下载图片并保存到服务器。

**函数签名:**
```go
func SaveNetworkImage(ctx context.Context, storageMode int, url string, randomName bool) (*res.SystemUploadFileRes, error)
```

**参数说明:**
- `storageMode`: 存储模式
- `url`: 图片URL地址
- `randomName`: 是否使用随机文件名

## 工具函数

### GetResourceType - 获取资源类型

根据文件的MIME类型自动识别资源类型。

```go
func GetResourceType(mimeType string) string
```

**支持的资源类型:**
- `image`: 图片类型
- `text`: 文本/文档类型（包括PDF、Word等）
- `audio`: 音频类型
- `video`: 视频类型
- `zip`: 压缩包类型
- `other`: 其他类型

**示例:**
```go
resourceType := upload.GetResourceType("image/jpeg")  // 返回 "image"
resourceType := upload.GetResourceType("video/mp4")   // 返回 "video"
```

### CheckFileMineType - 校验文件类型

检查上传的文件类型是否在允许的范围内。

```go
func CheckFileMineType(ctx context.Context, file *ghttp.UploadFile) error
```

**使用示例:**
```go
if err := upload.CheckFileMineType(ctx, file); err != nil {
    // 文件类型不允许上传
    return err
}
```

### CheckImageMineType - 校验图片类型

专门用于检查图片文件类型。

```go
func CheckImageMineType(ctx context.Context, file *ghttp.UploadFile) error
```

### CalcFileMd5 - 计算文件MD5

计算上传文件的MD5哈希值，用于文件唯一性校验。

```go
func CalcFileMd5(file *ghttp.UploadFile) (string, error)
```

**使用示例:**
```go
md5Hash, err := upload.CalcFileMd5(file)
if err != nil {
    return err
}
println("文件MD5:", md5Hash)
```

## 云存储配置

### 配置说明

云存储功能支持所有S3兼容的对象存储服务（如阿里云OSS、腾讯云COS、MinIO等）。

**配置项:**
- `upload_mode`: 上传模式（1: 本地, 2: 云存储）
- `endpoint`: 服务端点
- `access_key`: 访问密钥
- `secret_key`: 密钥
- `region`: 区域
- `bucket`: 存储桶名称
- `use_ssl`: 是否使用SSL
- `hostname_immutable`: 主机名是否不可变
- `del_local`: 上传到云端后是否删除本地文件

### GetCloudUpload - 获取云存储客户端

获取配置好的云存储客户端实例。

```go
func GetCloudUpload(ctx context.Context) (*goss.Goss, error)
```

### PutFromFile - 上传文件到云存储

将本地文件上传到云存储服务。

```go
func PutFromFile(ctx context.Context, filePath string, remotePath string) error
```

**使用示例:**
```go
err := upload.PutFromFile(ctx, "/local/path/file.jpg", "/remote/path/file.jpg")
if err != nil {
    return err
}
```

### IsLocalUpload - 判断是否本地存储

检查当前配置的存储模式是否为本地存储。

```go
func IsLocalUpload(ctx context.Context) bool
```

## 文件路径管理

### GetUploadPath - 获取上传基础路径

获取本地上传文件的基础路径。

```go
func GetUploadPath(ctx context.Context) string
```

### GetUploadFilePath - 获取上传文件完整路径

根据资源类型和日期获取文件的完整存储路径。

```go
func GetUploadFilePath(ctx context.Context, resourceType, dateDirName string) string
```

**路径结构:**
```
resource/public/uploads/{resourceType}/{dateDirName}/
```

### GetUploadUrlPath - 获取文件访问URL路径

获取上传文件的URL访问路径。

```go
func GetUploadUrlPath(ctx context.Context, resourceType, dateDirName, fileName string) string
```

**URL结构:**
```
/uploads/{resourceType}/{dateDirName}/{fileName}
```

## 完整使用示例

### 示例1: 带类型校验的文件上传

```go
package controller

import (
    "context"
    "devinggo/modules/system/model/req"
    "devinggo/modules/system/pkg/upload"
    "github.com/gogf/gf/v2/net/ghttp"
)

type FileController struct{}

func (c *FileController) Upload(ctx context.Context, r *ghttp.Request) {
    // 获取上传文件
    file := r.GetUploadFile("file")
    if file == nil {
        r.Response.WriteJson(g.Map{
            "code": 400,
            "msg":  "请选择文件",
        })
        return
    }
    
    // 校验文件类型
    if err := upload.CheckFileMineType(ctx, file); err != nil {
        r.Response.WriteJson(g.Map{
            "code": 400,
            "msg":  err.Error(),
        })
        return
    }
    
    // 执行上传
    input := &req.FileUploadInput{
        File:        file,
        StorageMode: 1,
        RandomName:  true,
    }
    
    result, err := upload.Upload(ctx, input)
    if err != nil {
        r.Response.WriteJson(g.Map{
            "code": 500,
            "msg":  "上传失败",
            "data": err.Error(),
        })
        return
    }
    
    // 返回结果
    r.Response.WriteJson(g.Map{
        "code": 0,
        "msg":  "上传成功",
        "data": result,
    })
}
```

### 示例2: 分片上传完整流程

```go
// 前端分片上传流程
func (c *FileController) ChunkUpload(ctx context.Context, r *ghttp.Request) {
    // 获取分片信息
    file := r.GetUploadFile("file")
    index := r.GetInt64("index")      // 当前分片序号
    total := r.GetInt64("total")      // 总分片数
    hash := r.GetString("hash")       // 文件唯一标识
    ext := r.GetString("ext")         // 文件扩展名
    fileName := r.GetString("name")   // 原始文件名
    mimeType := r.GetString("type")   // MIME类型
    
    // 校验分片文件类型
    if err := upload.CheckChunkFileMineType(ctx, ext); err != nil {
        r.Response.WriteJson(g.Map{
            "code": 400,
            "msg":  err.Error(),
        })
        return
    }
    
    // 执行分片上传
    input := &req.ChunkUploadInput{
        File:        file,
        StorageMode: 1,
        Index:       index,
        Total:       total,
        Hash:        hash,
        Name:        fileName,
        Ext:         ext,
        Type:        mimeType,
        RandomName:  true,
    }
    
    result, err := upload.ChunkUpload(ctx, input)
    if err != nil {
        r.Response.WriteJson(g.Map{
            "code": 500,
            "msg":  "上传失败",
        })
        return
    }
    
    // 判断是否为最后一个分片
    if result != nil {
        // 所有分片上传完成
        r.Response.WriteJson(g.Map{
            "code": 0,
            "msg":  "上传完成",
            "data": result,
        })
    } else {
        // 分片上传成功，等待下一个分片
        r.Response.WriteJson(g.Map{
            "code": 0,
            "msg":  "分片上传成功",
        })
    }
}
```

### 示例3: 保存网络图片

```go
func (c *FileController) SaveNetworkImage(ctx context.Context, r *ghttp.Request) {
    // 获取图片URL
    imageUrl := r.GetString("url")
    if imageUrl == "" {
        r.Response.WriteJson(g.Map{
            "code": 400,
            "msg":  "请提供图片URL",
        })
        return
    }
    
    // 保存网络图片
    result, err := upload.SaveNetworkImage(ctx, 1, imageUrl, true)
    if err != nil {
        r.Response.WriteJson(g.Map{
            "code": 500,
            "msg":  "保存失败",
            "data": err.Error(),
        })
        return
    }
    
    // 返回结果
    r.Response.WriteJson(g.Map{
        "code": 0,
        "msg":  "保存成功",
        "data": result,
    })
}
```

### 示例4: 同时使用云存储

```go
func (c *FileController) UploadToCloud(ctx context.Context, r *ghttp.Request) {
    file := r.GetUploadFile("file")
    
    // 上传到云存储
    input := &req.FileUploadInput{
        File:        file,
        StorageMode: 2,  // 云存储模式
        RandomName:  true,
    }
    
    result, err := upload.Upload(ctx, input)
    if err != nil {
        r.Response.WriteJson(g.Map{
            "code": 500,
            "msg":  "上传失败",
        })
        return
    }
    
    r.Response.WriteJson(g.Map{
        "code": 0,
        "msg":  "上传成功",
        "data": result,
    })
}
```

## 配置文件示例

在系统配置表中需要配置以下上传相关参数：

```yaml
upload_config:
  # 上传模式: 1-本地存储, 2-云存储
  upload_mode: "1"
  
  # 允许上传的文件类型
  upload_allow_file: "doc,docx,xls,xlsx,ppt,pptx,pdf,zip,rar,txt"
  
  # 允许上传的图片类型
  upload_allow_image: "jpg,jpeg,png,gif,bmp,webp"
  
  # 云存储配置（当upload_mode为2时生效）
  endpoint: "oss-cn-hangzhou.aliyuncs.com"
  access_key: "your_access_key"
  secret_key: "your_secret_key"
  region: "cn-hangzhou"
  bucket: "your-bucket-name"
  use_ssl: "true"
  hostname_immutable: "false"
  
  # 上传到云端后是否删除本地文件
  del_local: "false"
```

## 注意事项

1. **文件大小限制**: 需要在服务器和框架配置中设置合适的上传大小限制
2. **分片大小**: 建议分片大小为 5MB - 10MB，根据网络环境调整
3. **文件类型校验**: 建议在上传前进行文件类型校验，防止上传危险文件
4. **MD5校验**: 可以使用MD5值进行文件去重和完整性校验
5. **云存储**: 使用云存储时，需要确保配置正确且有相应的权限
6. **路径安全**: 所有文件路径都经过处理，防止路径遍历攻击
7. **并发上传**: 分片上传支持并发，但需要注意服务器负载

## 常见问题

### Q: 如何支持更大的文件？
A: 使用分片上传功能，将大文件分割成多个小分片上传。

### Q: 如何防止重复上传？
A: 可以通过文件的MD5值进行判断，上传前先检查系统中是否已存在相同MD5的文件。

### Q: 云存储上传失败如何处理？
A: 文件会先保存到本地，如果云存储上传失败，可以根据 `del_local` 配置决定是否保留本地文件。

### Q: 如何自定义文件存储路径？
A: 可以修改配置文件中的 `upload.dir` 配置项，或直接调用相关路径函数。

### Q: 支持哪些云存储服务？
A: 支持所有S3兼容的对象存储服务，包括阿里云OSS、腾讯云COS、AWS S3、MinIO等。

## 更新日志

- v1.0.0: 初始版本，支持基本上传功能
- v1.1.0: 增加分片上传功能
- v1.2.0: 增加云存储支持
- v1.3.0: 增加网络图片保存功能

## 技术支持

如有问题，请访问：
- GitHub: https://github.com/huagelong/devinggo
- 文档: 查看项目文档获取更多信息

## 许可证

本项目遵循 MIT 许可证。详见 LICENSE 文件。
