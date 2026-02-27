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

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/hibiken/asynq"
)

// ==================== 函数式注册（最新方式） ====================

func init() {
	// 方式1: 使用RegisterWorkerFunc - 最简洁的Worker注册
	worker.RegisterWorkerFunc("example:simple:worker", executeSimpleWorker)

	// 方式2: 使用RegisterCronFunc - 最简洁的Cron注册
	worker.RegisterCronFunc(
		"example:simple:cron",
		"简单的定时任务示例",
		handleSimpleCronParams,
	)

	// 方式3: 使用Builder模式 - 更灵活的注册方式
	worker.NewWorkerBuilder("example:builder:worker").
		WithExecute(executeBuilderWorker).
		Register()

	worker.NewCronBuilder("example:builder:cron", "Builder模式的定时任务").
		WithParamsHandler(handleBuilderCronParams).
		Register()
}

// SimpleWorkerData 简单Worker的数据结构
type SimpleWorkerData struct {
	Message string `json:"message"`
}

// executeSimpleWorker Worker执行函数
func executeSimpleWorker(ctx context.Context, t *asynq.Task) error {
	data, err := glob2.GetParamters[SimpleWorkerData](ctx, t)
	if err != nil {
		return err
	}
	glob2.WithWorkLog().Infof(ctx, "执行简单Worker: %s", data.Message)
	return nil
}

// SimpleCronData 简单Cron的数据结构
type SimpleCronData struct {
	Cmd string `json:"cmd"`
}

// handleSimpleCronParams Cron参数处理函数
func handleSimpleCronParams(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
	if g.IsEmpty(params) {
		return
	}
	data := new(SimpleCronData)
	if err := params.Scan(data); err != nil {
		glob2.WithWorkLog().Errorf(ctx, "解析Cron参数失败: %v", err)
		return
	}
	payload.Data = data
}

// BuilderWorkerData Builder Worker的数据结构
type BuilderWorkerData struct {
	Command string `json:"command"`
}

// executeBuilderWorker Builder Worker执行函数
func executeBuilderWorker(ctx context.Context, t *asynq.Task) error {
	data, err := glob2.GetParamters[BuilderWorkerData](ctx, t)
	if err != nil {
		return err
	}
	glob2.WithWorkLog().Infof(ctx, "执行Builder Worker命令: %s", data.Command)

	// 执行shell命令
	result, err := gproc.ShellExec(ctx, data.Command)
	if err != nil {
		return err
	}
	glob2.WithWorkLog().Infof(ctx, "命令执行结果: %s", result)
	return nil
}

// BuilderCronData Builder Cron的数据结构
type BuilderCronData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// handleBuilderCronParams Builder Cron参数处理函数
func handleBuilderCronParams(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
	if g.IsEmpty(params) {
		return
	}
	data := new(BuilderCronData)
	if err := params.Scan(data); err != nil {
		glob2.WithWorkLog().Errorf(ctx, "解析Builder Cron参数失败: %v", err)
		return
	}
	payload.Data = data
}

// ==================== 发送任务示例 ====================

// SendSimpleTask 发送简单任务
func SendSimpleTask(ctx context.Context, message string) error {
	return worker.NewTaskBuilder(ctx, "example:simple:worker").
		WithData(SimpleWorkerData{Message: message}).
		Send()
}

// SendBuilderTask 发送Builder任务
func SendBuilderTask(ctx context.Context, command string) error {
	return worker.NewTaskBuilder(ctx, "example:builder:worker").
		WithData(BuilderWorkerData{Command: command}).
		WithQueue("critical"). // 设置队列
		Send()
}
