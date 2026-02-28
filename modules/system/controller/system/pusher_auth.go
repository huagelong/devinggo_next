// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"context"

	"devinggo/modules/system/api/system"
	"devinggo/modules/system/controller/base"
	"devinggo/modules/system/myerror"
	"devinggo/modules/system/pkg/websocket"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	PusherAuthController = pusherAuthController{}
)

type pusherAuthController struct {
	base.BaseController
}

// PusherAuth Pusher频道认证端点
// 用于Private/Presence频道的认证签名生成
func (c *pusherAuthController) PusherAuth(ctx context.Context, req *system.PusherAuthReq) (rs *system.PusherAuthRes, err error) {
	rs = &system.PusherAuthRes{}

	// 验证socket_id格式（防止伪造）
	if req.SocketId == "" || req.ChannelName == "" {
		err = myerror.ValidationFailed(ctx, "socket_id and channel_name are required")
		return
	}

	// ⚠️ 验证socket_id是否属于当前用户（防止跨用户攻击）
	// 从Redis获取socket_id对应的ServerName，验证连接存在性
	serverName := websocket.GetServerNameBySocketId4Redis(ctx, req.SocketId)
	if serverName == "" {
		g.Log().Warning(ctx, "Invalid socket_id:", req.SocketId)
		err = myerror.ValidationFailed(ctx, "Invalid socket_id or connection not found")
		return
	}

	// ⚠️ 安全检查：验证频道类型
	if !websocket.RequiresAuth(req.ChannelName) {
		err = myerror.ValidationFailed(ctx, "This channel does not require authentication")
		return
	}

	// TODO: 可选安全增强
	// - 验证用户是否有权访问该频道（根据业务规则）
	// - 记录认证日志（审计）
	// - 限流防止滥用（Rate Limiting）

	// 检查是否为Presence频道
	if websocket.IsPresenceChannel(req.ChannelName) {
		// Phase 3: Presence频道认证（包含channel_data）
		// 生成channel_data（包含用户信息）
		userInfo := map[string]interface{}{
			"name": "User " + gconv.String(c.UserId), // TODO: 从数据库获取真实用户信息
			// 可以添加更多用户信息，如头像、角色等
		}
		channelData, err := websocket.EncodeChannelData(gconv.String(c.UserId), userInfo)
		if err != nil {
			g.Log().Warning(ctx, "EncodeChannelData error:", err)
			err = myerror.ValidationFailed(ctx, "Failed to generate channel_data")
			return rs, err
		}

		// 生成认证签名（包含channel_data）
		auth := websocket.GenerateAuthSignature(req.SocketId, req.ChannelName, channelData)
		rs.Auth = auth
		rs.ChannelData = channelData

		g.Log().Debugf(ctx, "Generated presence auth for user:%d, socket:%s, channel:%s", c.UserId, req.SocketId, req.ChannelName)
	} else {
		// Private频道认证（不包含channel_data）
		auth := websocket.GenerateAuthSignature(req.SocketId, req.ChannelName, "")
		rs.Auth = auth

		g.Log().Debugf(ctx, "Generated private auth for user:%d, socket:%s, channel:%s", c.UserId, req.SocketId, req.ChannelName)
	}

	return
}
