// Package cron
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cron

import (
	"context"
	"devinggo/modules/system/worker/consts"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// TestCronData 测试定时任务的数据结构
type TestCronData struct {
	// TODO: 在这里定义你的参数字段
	// 例如：
	// Name  string `json:"name" v:"required#名称不能为空"`
	// Email string `json:"email" v:"required|email#邮箱不能为空|邮箱格式错误"`
}

func init() {
	// 注册定时任务
	worker.RegisterCronFunc(consts.TEST_CRON_CRON, "测试定时任务", handleTestCron)
}

// handleTestCron 处理定时任务参数
func handleTestCron(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
	if g.IsEmpty(params) {
		glob2.WithWorkLog().Warning(ctx, "[test_cron] 参数为空")
		return
	}
	
	data := new(TestCronData)
	if err := params.Scan(data); err != nil {
		glob2.WithWorkLog().Errorf(ctx, "[test_cron] 参数解析失败: %v", err)
		return
	}
	
	// 可以在这里添加参数验证
	// if err := g.Validator().Data(data).Run(ctx); err != nil {
	//     glob2.WithWorkLog().Errorf(ctx, "[test_cron] 参数验证失败: %v", err)
	//     return
	// }
	
	payload.Data = data
	glob2.WithWorkLog().Debugf(ctx, "[test_cron] 参数设置成功: %+v", data)
}
