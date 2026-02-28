// Package websocket
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package websocket

import "context"

// 日志相关方法（暂时保留空实现）
func saveLog(ctx context.Context, msg *PusherRequest, id int64) {
	// TODO: 实现Pusher协议日志保存
}

func (c *Client) updateMsgLog(ctx context.Context, msg *PusherResponse) {
	// TODO: 实现Pusher协议日志更新
}
