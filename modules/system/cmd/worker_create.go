// Package cmd
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	CreateWorker = &gcmd.Command{
		Name:  "worker:create",
		Brief: "创建 Task 或 Cron 任务",
		Description: `
		快速创建 Task 或 Cron 任务，自动生成规范代码
		
		用法: 
		  go run main.go worker:create -name 任务名称 [-module 模块名称] [-type 类型] [-desc 描述]
		
		参数:
		  -name   任务名称（必填），建议使用下划线命名，如: send_email
		  -module 模块名称（可选，默认: system）
		  -type   任务类型（可选，默认: both）
		          task: 仅创建异步任务
		          cron: 仅创建定时任务
		          both: 同时创建任务和定时任务
		  -desc   任务描述（可选），用于生成注释
		
		示例:
		  # 创建发送邮件任务（默认both类型）
		  go run main.go worker:create -name send_email -desc "发送邮件"
		  
		  # 仅创建异步任务
		  go run main.go worker:create -name process_order -type task -desc "处理订单"
		  
		  # 仅创建定时任务
		  go run main.go worker:create -name clean_logs -type cron -desc "清理日志"
		  
		  # 在指定模块中创建
		  go run main.go worker:create -name notify_user -module custom -desc "用户通知"
		`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			CmdInit(ctx, parser)

			// 获取任务名称
			name := gcmd.GetOpt("name").String()
			if name == "" {
				return gerror.New("任务名称不能为空，使用 -name 参数指定")
			}

			// 获取模块名称，默认为 system
			moduleName := gcmd.GetOpt("module", "system").String()

			// 获取类型，默认为 both
			workerType := gcmd.GetOpt("type", "both").String()

			// 获取描述
			description := gcmd.GetOpt("desc").String()
			if description == "" {
				description = name
			}

			// 验证类型
			if workerType != "task" && workerType != "cron" && workerType != "both" {
				return gerror.New("类型必须是 task、cron 或 both")
			}

			// 验证模块是否存在
			modulePath := fmt.Sprintf("./modules/%s", moduleName)
			if !gfile.Exists(modulePath) {
				return gerror.Newf("模块 '%s' 不存在，请先创建模块或使用已存在的模块", moduleName)
			}

			// 打印创建信息
			fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("📦 任务名称: %s\n", name)
			fmt.Printf("📂 所属模块: %s\n", moduleName)
			fmt.Printf("🔖 任务类型: %s\n", workerType)
			fmt.Printf("📝 任务描述: %s\n", description)
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

			// 创建必要的目录
			if err := createWorkerDirectories(ctx, moduleName, workerType); err != nil {
				return err
			}

			// 创建或更新常量文件
			if err := updateConstFile(ctx, moduleName, name, description, workerType); err != nil {
				return err
			}

			// 根据类型创建文件
			if workerType == "cron" || workerType == "both" {
				if err := createCronFile(ctx, moduleName, name, description); err != nil {
					return err
				}
			}

			if workerType == "task" || workerType == "both" {
				if err := createTaskFile(ctx, moduleName, name, description, workerType == "both"); err != nil {
					return err
				}
			}

			// 输出成功信息
			fmt.Println("\n✅ 创建成功！")
			fmt.Println("\n📁 生成的文件:")
			if workerType == "cron" || workerType == "both" {
				fmt.Printf("   • modules/%s/worker/cron/%s_cron.go\n", moduleName, name)
			}
			if workerType == "task" || workerType == "both" {
				fmt.Printf("   • modules/%s/worker/server/%s_worker.go\n", moduleName, name)
			}
			fmt.Printf("   • modules/%s/worker/consts/const.go (已更新)\n", moduleName)

			// 输出下一步提示
			fmt.Println("\n💡 下一步:")
			fmt.Println("   1. 编辑生成的文件，添加业务逻辑")
			fmt.Println("   2. 如果 worker 服务正在运行，需要重启以加载新任务")
			if workerType == "cron" || workerType == "both" {
				fmt.Println("   3. 在后台管理系统中配置定时任务的执行时间")
			}
			fmt.Println()

			return nil
		},
	}
)

// createWorkerDirectories 创建必要的目录
func createWorkerDirectories(ctx context.Context, moduleName, workerType string) error {
	dirs := []string{
		fmt.Sprintf("./modules/%s/worker", moduleName),
		fmt.Sprintf("./modules/%s/worker/consts", moduleName),
	}

	if workerType == "cron" || workerType == "both" {
		dirs = append(dirs, fmt.Sprintf("./modules/%s/worker/cron", moduleName))
	}

	if workerType == "task" || workerType == "both" {
		dirs = append(dirs, fmt.Sprintf("./modules/%s/worker/server", moduleName))
	}

	for _, dir := range dirs {
		if !gfile.Exists(dir) {
			if err := gfile.Mkdir(dir); err != nil {
				return gerror.Wrapf(err, "创建目录 '%s' 失败", dir)
			}
			g.Log().Debugf(ctx, "创建目录: %s", dir)
		}
	}
	return nil
}

// updateConstFile 创建或更新常量文件
func updateConstFile(ctx context.Context, moduleName, name, description, workerType string) error {
	constPath := fmt.Sprintf("./modules/%s/worker/consts/const.go", moduleName)

	// 常量名称（转为大写加下划线格式）
	constName := strings.ToUpper(gstr.CaseSnake(name))

	// 检查文件是否存在
	if !gfile.Exists(constPath) {
		// 创建新的常量文件
		content := `// Package consts
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package consts

var (
`

		if workerType == "cron" || workerType == "both" {
			content += fmt.Sprintf("\t%s_CRON = \"%s_cron\" // %s\n", constName, name, description)
		}
		if workerType == "task" || workerType == "both" {
			content += fmt.Sprintf("\t%s_TASK = \"%s_task\" // %s\n", constName, name, description)
		}

		content += ")\n"

		if err := gfile.PutContents(constPath, content); err != nil {
			return gerror.Wrapf(err, "创建常量文件失败")
		}
		g.Log().Debugf(ctx, "创建常量文件: %s", constPath)
	} else {
		// 读取现有文件内容
		content := gfile.GetContents(constPath)

		// 检查常量是否已存在
		if workerType == "cron" || workerType == "both" {
			cronConstName := fmt.Sprintf("%s_CRON", constName)
			if strings.Contains(content, cronConstName) {
				return gerror.Newf("常量 %s 已存在，请检查是否重复创建", cronConstName)
			}
			// 在最后一个 ) 之前插入新常量
			newConst := fmt.Sprintf("\t%s = \"%s_cron\" // %s\n", cronConstName, name, description)
			lastParen := strings.LastIndex(content, ")")
			if lastParen != -1 {
				content = content[:lastParen] + newConst + content[lastParen:]
			}
		}

		if workerType == "task" || workerType == "both" {
			taskConstName := fmt.Sprintf("%s_TASK", constName)
			if strings.Contains(content, taskConstName) {
				return gerror.Newf("常量 %s 已存在，请检查是否重复创建", taskConstName)
			}
			// 在最后一个 ) 之前插入新常量
			newConst := fmt.Sprintf("\t%s = \"%s_task\" // %s\n", taskConstName, name, description)
			lastParen := strings.LastIndex(content, ")")
			if lastParen != -1 {
				content = content[:lastParen] + newConst + content[lastParen:]
			}
		}

		// 写回文件
		if err := gfile.PutContents(constPath, content); err != nil {
			return gerror.Wrapf(err, "更新常量文件失败")
		}
		g.Log().Debugf(ctx, "更新常量文件: %s", constPath)
	}

	return nil
}

// createCronFile 创建 Cron 文件
func createCronFile(ctx context.Context, moduleName, name, description string) error {
	cronPath := fmt.Sprintf("./modules/%s/worker/cron/%s_cron.go", moduleName, name)

	// 检查文件是否已存在
	if gfile.Exists(cronPath) {
		return gerror.Newf("Cron 文件 '%s' 已存在，请检查是否重复创建", cronPath)
	}

	// 常量名称
	constName := strings.ToUpper(gstr.CaseSnake(name))
	// 数据结构名称（驼峰命名）
	structName := gstr.CaseCamel(name) + "Data"
	// 处理函数名
	handlerName := "handle" + gstr.CaseCamel(name)

	content := fmt.Sprintf(`// Package cron
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cron

import (
	"context"
	"devinggo/modules/%s/worker/consts"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// %s %s的数据结构
type %s struct {
	// TODO: 在这里定义你的参数字段
	// 例如：
	// Name  string `+"`json:\"name\" v:\"required#名称不能为空\"`"+`
	// Email string `+"`json:\"email\" v:\"required|email#邮箱不能为空|邮箱格式错误\"`"+`
}

func init() {
	// 注册定时任务
	worker.RegisterCronFunc(consts.%s_CRON, "%s", %s)
}

// %s 处理定时任务参数
func %s(ctx context.Context, payload *glob2.Payload, params *gjson.Json) {
	if g.IsEmpty(params) {
		glob2.WithWorkLog().Warning(ctx, "[%s] 参数为空")
		return
	}
	
	data := new(%s)
	if err := params.Scan(data); err != nil {
		glob2.WithWorkLog().Errorf(ctx, "[%s] 参数解析失败: %%v", err)
		return
	}
	
	// 可以在这里添加参数验证
	// if err := g.Validator().Data(data).Run(ctx); err != nil {
	//     glob2.WithWorkLog().Errorf(ctx, "[%s] 参数验证失败: %%v", err)
	//     return
	// }
	
	payload.Data = data
	glob2.WithWorkLog().Debugf(ctx, "[%s] 参数设置成功: %%+v", data)
}
`, moduleName, structName, description, structName, constName, description, handlerName, handlerName, handlerName, name, structName, name, name, name)

	if err := gfile.PutContents(cronPath, content); err != nil {
		return gerror.Wrapf(err, "创建 Cron 文件失败")
	}

	g.Log().Debugf(ctx, "创建 Cron 文件: %s", cronPath)
	return nil
}

// createTaskFile 创建 Task 文件
func createTaskFile(ctx context.Context, moduleName, name, description string, hasCron bool) error {
	taskPath := fmt.Sprintf("./modules/%s/worker/server/%s_worker.go", moduleName, name)

	// 检查文件是否已存在
	if gfile.Exists(taskPath) {
		return gerror.Newf("Task 文件 '%s' 已存在，请检查是否重复创建", taskPath)
	}

	// 常量名称
	constName := strings.ToUpper(gstr.CaseSnake(name))
	// 函数名称（驼峰命名）
	funcName := gstr.CaseCamel(name)
	// 数据结构名称
	structName := gstr.CaseCamel(name) + "Data"

	// 如果同时创建了 cron，可以复用数据结构
	var importCron, dataTypeAlias string
	if hasCron {
		importCron = fmt.Sprintf("\t\"devinggo/modules/%s/worker/cron\"\n", moduleName)
		dataTypeAlias = fmt.Sprintf("// 复用 Cron 的数据结构\ntype %s = cron.%s\n", structName, structName)
	} else {
		dataTypeAlias = fmt.Sprintf(`// %s %s的数据结构
type %s struct {
	// TODO: 在这里定义你的参数字段
	// 例如：
	// Name  string `+"`json:\"name\" v:\"required#名称不能为空\"`"+`
	// Email string `+"`json:\"email\" v:\"required|email#邮箱不能为空|邮箱格式错误\"`"+`
}
`, structName, description, structName)
	}

	content := fmt.Sprintf(`// Package server
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package server

import (
	"context"
%s	"devinggo/modules/%s/worker/consts"
	"devinggo/modules/system/pkg/worker"
	glob2 "devinggo/modules/system/pkg/worker/glob"

	"github.com/hibiken/asynq"
)

%s
func init() {
	// 注册任务处理器
	worker.RegisterWorkerFunc(consts.%s_TASK, execute%s)
}

// execute%s 执行%s任务
func execute%s(ctx context.Context, t *asynq.Task) error {
	// 获取任务参数
	data, err := glob2.GetParamters[%s](ctx, t)
	if err != nil {
		glob2.WithWorkLog().Errorf(ctx, "[%%s] 获取参数失败: %%v", consts.%s_TASK, err)
		return err
	}
	
	// 可以在这里添加参数验证
	// if err := g.Validator().Data(data).Run(ctx); err != nil {
	//     glob2.WithWorkLog().Errorf(ctx, "[%%s] 参数验证失败: %%v", consts.%s_TASK, err)
	//     return err
	// }
	
	glob2.WithWorkLog().Infof(ctx, "[%%s] 开始执行任务, 参数: %%+v", consts.%s_TASK, data)
	
	// ========================================
	// TODO: 在这里实现你的业务逻辑
	// ========================================
	// 示例：
	// 1. 数据库操作
	// err = dao.User.Ctx(ctx).Where("id", data.UserId).Update(g.Map{...})
	// 
	// 2. 调用外部服务
	// result, err := service.External().DoSomething(ctx, data)
	// 
	// 3. 发送通知
	// err = service.Notification().Send(ctx, data.Email, "通知内容")
	
	glob2.WithWorkLog().Infof(ctx, "[%%s] 任务执行完成", consts.%s_TASK)
	return nil
}
`, importCron, moduleName, dataTypeAlias, constName, funcName, funcName, description, funcName, structName, constName, constName, constName, constName)

	if err := gfile.PutContents(taskPath, content); err != nil {
		return gerror.Wrapf(err, "创建 Task 文件失败")
	}

	g.Log().Debugf(ctx, "创建 Task 文件: %s", taskPath)
	return nil
}
