// Package upload 使用示例
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package upload_test

import (
	"context"
	"devinggo/modules/system/pkg/upload"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ExampleSimpleUpload 简单文件上传示例
func ExampleSimpleUpload(ctx context.Context, file *ghttp.UploadFile) {
	// 最简单的上传方式
	result, err := upload.NewUploader(ctx).UploadFile(file)
	if err != nil {
		// 处理错误
		return
	}

	// 使用结果
	_ = result.Url
}

// ExampleCustomUpload 自定义上传配置示例
func ExampleCustomUpload(ctx context.Context, file *ghttp.UploadFile) {
	// 链式调用配置上传参数
	result, err := upload.NewUploader(ctx).
		UseCloudStorage().     // 使用云存储
		SetRandomName(false).  // 保留原始文件名
		SetValidateType(true). // 启用类型验证
		UploadFile(file)

	if err != nil {
		return
	}

	_ = result
}

// ExampleImageUpload 图片上传示例
func ExampleImageUpload(ctx context.Context, file *ghttp.UploadFile) {
	// 上传图片（自动验证图片类型）
	result, err := upload.NewUploader(ctx).
		SetStorageMode(upload.StorageModeLocal). // 使用本地存储
		SetRandomName(true).                     // 使用随机文件名
		UploadImage(file)

	if err != nil {
		return
	}

	_ = result
}

// ExampleSaveFromURL 从URL保存文件示例
func ExampleSaveFromURL(ctx context.Context, imageUrl string) {
	// 从URL保存图片
	result, err := upload.NewUploader(ctx).
		UseLocalStorage().   // 使用本地存储
		SetRandomName(true). // 使用随机文件名
		SaveFromURL(imageUrl)

	if err != nil {
		return
	}

	_ = result
}

// ExampleValidateFile 文件验证示例
func ExampleValidateFile(ctx context.Context, file *ghttp.UploadFile) {
	// 创建验证器
	validator := upload.NewValidator(ctx)

	// 验证文件类型
	if err := validator.ValidateFile(file); err != nil {
		// 文件类型不允许
		return
	}

	// 验证图片类型
	if err := validator.ValidateImage(file); err != nil {
		// 图片类型不允许
		return
	}

	// 验证扩展名
	if err := validator.ValidateExtension("jpg", true); err != nil {
		// 扩展名不允许
		return
	}
}

// ExampleGetResourceType 资源类型识别示例
func ExampleGetResourceType() {
	// 根据MIME类型获取资源类型
	imageType := upload.GetResourceType("image/jpeg")     // 返回 "image"
	videoType := upload.GetResourceType("video/mp4")      // 返回 "video"
	textType := upload.GetResourceType("application/pdf") // 返回 "text"

	_, _, _ = imageType, videoType, textType
}

// ExampleFormatSize 格式化文件大小示例
func ExampleFormatSize() {
	// 格式化文件大小
	size1 := upload.FormatSize(1024)               // 1.00 KB
	size2 := upload.FormatSize(1024 * 1024)        // 1.00 MB
	size3 := upload.FormatSize(1024 * 1024 * 1024) // 1.00 GB

	_, _, _ = size1, size2, size3
}
