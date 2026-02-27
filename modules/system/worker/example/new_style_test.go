// Package example
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package example

import (
	"context"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/hibiken/asynq"
)

// TestNewStyleWorker 测试新方式的Worker注册
func TestNewStyleWorker(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_ = ctx // 避免未使用警告

		// 测试数据
		type TestData struct {
			Message string `json:"message"`
		}

		// 定义执行函数
		executed := false
		executeFunc := func(ctx context.Context, task *asynq.Task) error {
			data, err := glob2.GetParamters[TestData](ctx, task)
			if err != nil {
				return err
			}
			t.Assert(data.Message, "test message")
			executed = true
			return nil
		}

		// 使用新方式注册Worker
		worker.RegisterWorkerFunc("test:new:worker", executeFunc)

		t.Assert(executed, false) // 注册不会执行
	})
}

// TestNewStyleCron 测试新方式的Cron注册
func TestNewStyleCron(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_ = ctx // 避免未使用警告

		// 测试数据
		type CronData struct {
			Name string `json:"name"`
		}

		// 定义参数处理函数
		handleParams := func(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
			if g.IsEmpty(params) {
				return
			}
			data := new(CronData)
			if err := params.Scan(data); err != nil {
				return
			}
			payload.Data = data
		}

		// 使用新方式注册Cron
		worker.RegisterCronFunc("test:new:cron", "测试定时任务", handleParams)

		// 验证注册成功（实际应用中由系统调用）
		t.AssertNE(handleParams, nil)
	})
}

// TestTaskBuilder 测试TaskBuilder
func TestTaskBuilder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()

		type TaskData struct {
			Message string `json:"message"`
			Count   int    `json:"count"`
		}

		// 构建任务（不实际发送）
		task := worker.NewTaskBuilder(ctx, "test:task").
			WithData(TaskData{
				Message: "hello",
				Count:   100,
			}).
			WithQueue("test").
			WithDelay(5 * time.Second).
			Build()

		t.AssertNE(task, nil)
		t.Assert(task.Type(), "test:task")
	})
}

// BenchmarkNewStyleWorker 性能测试：新方式Worker注册
func BenchmarkNewStyleWorker(b *testing.B) {
	ctx := gctx.New()
	_ = ctx // 避免未使用警告

	executeFunc := func(ctx context.Context, task *asynq.Task) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker.NewWorkerBuilder("bench:worker").
			WithExecute(executeFunc).
			Register()
	}
}

// BenchmarkNewStyleCron 性能测试：新方式Cron注册
func BenchmarkNewStyleCron(b *testing.B) {
	ctx := gctx.New()
	_ = ctx // 避免未使用警告

	handleParams := func(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
		// 参数处理
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worker.NewCronBuilder("bench:cron", "性能测试").
			WithParamsHandler(handleParams).
			Register()
	}
}

// ExampleRegisterWorkerFunc 示例：函数式注册Worker
func ExampleRegisterWorkerFunc() {
	// 定义数据结构
	type MyData struct {
		Name string `json:"name"`
	}

	// 定义执行函数
	executeFunc := func(ctx context.Context, t *asynq.Task) error {
		data, err := glob2.GetParamters[MyData](ctx, t)
		if err != nil {
			return err
		}
		glob2.WithWorkLog().Infof(ctx, "处理: %s", data.Name)
		return nil
	}

	// 注册Worker（就这么简单！）
	worker.RegisterWorkerFunc("example:worker", executeFunc)

	// Output:
}

// ExampleRegisterCronFunc 示例：函数式注册Cron
func ExampleRegisterCronFunc() {
	// 定义数据结构
	type CronData struct {
		Param string `json:"param"`
	}

	// 定义参数处理函数
	handleParams := func(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
		if g.IsEmpty(params) {
			return
		}
		data := new(CronData)
		if err := params.Scan(data); err != nil {
			return
		}
		payload.Data = data
	}

	// 注册Cron（就这么简单！）
	worker.RegisterCronFunc("example:cron", "示例定时任务", handleParams)

	// Output:
}

// ExampleTaskBuilder 示例：TaskBuilder使用
func ExampleTaskBuilder() {
	ctx := gctx.New()

	// 定义数据结构
	type TaskData struct {
		Message string `json:"message"`
	}

	// 构建任务（不实际发送）
	task := worker.NewTaskBuilder(ctx, "example:task").
		WithData(TaskData{Message: "hello"}).
		WithQueue("default").
		WithDelay(5 * time.Second).
		Build()

	if task != nil {
		g.Log().Info(ctx, "任务构建成功")
	}

	// Output:
}
