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
// 注意：此接口直接返回Pusher标准格式的JSON，不使用GoFrame响应包装
func (c *pusherAuthController) PusherAuth(ctx context.Context, req *system.PusherAuthReq) (rs *system.PusherAuthRes, err error) {
	r := g.RequestFromCtx(ctx)

	// 验证socket_id格式（防止伪造）
	if req.SocketId == "" || req.ChannelName == "" {
		// 返回Pusher标准错误格式
		r.Response.Status = 403
		r.Response.WriteJson(g.Map{
			"error": "socket_id and channel_name are required",
		})
		r.ExitAll()
		return
	}

	// ⚠️ 验证socket_id是否属于当前用户（防止跨用户攻击）
	// 从Redis获取socket_id对应的ServerName，验证连接存在性
	serverName := websocket.GetServerNameBySocketId4Redis(ctx, req.SocketId)
	if serverName == "" {
		g.Log().Warning(ctx, "Invalid socket_id:", req.SocketId)
		r.Response.Status = 403
		r.Response.WriteJson(g.Map{
			"error": "Invalid socket_id or connection not found",
		})
		r.ExitAll()
		return
	}

	// ⚠️ 安全检查：验证频道类型
	if !websocket.RequiresAuth(req.ChannelName) {
		r.Response.Status = 403
		r.Response.WriteJson(g.Map{
			"error": "This channel does not require authentication",
		})
		r.ExitAll()
		return
	}

	// TODO: 可选安全增强
	// - 验证用户是否有权访问该频道（根据业务规则）
	// - 记录认证日志（审计）
	// - 限流防止滥用（Rate Limiting）

	// ⚠️ 检查是否为 Encrypted Channel（private-encrypted-*）
	if websocket.IsEncryptedChannel(req.ChannelName) {
		// Encrypted Channel 认证：需要生成 shared_secret
		// shared_secret 是一个 32 字节的随机密钥（Base64 编码）
		// Pusher.js 会使用此密钥进行端到端加密
		sharedSecret := websocket.GenerateSharedSecret()

		// 保存 shared_secret 到 Redis（用于服务器端加密推送）
		// Key: pusher:encrypted_secret:{channel_name}
		// TTL: 24小时（可配置）
		saveErr := websocket.SaveSharedSecret(ctx, req.ChannelName, sharedSecret)
		if saveErr != nil {
			g.Log().Warning(ctx, "Failed to save shared_secret:", saveErr)
			// 继续执行，不影响客户端认证
		}

		// 生成认证签名（不包含 channel_data）
		auth := websocket.GenerateAuthSignature(req.SocketId, req.ChannelName, "")

		g.Log().Debugf(ctx, "Generated encrypted channel auth for user:%d, socket:%s, channel:%s", c.UserId, req.SocketId, req.ChannelName)

		// 返回 Pusher 标准格式（包含 shared_secret）
		r.Response.WriteJson(g.Map{
			"auth":          auth,
			"shared_secret": sharedSecret,
		})
		r.ExitAll()
		return
	}

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
			r.Response.Status = 500
			r.Response.WriteJson(g.Map{
				"error": "Failed to generate channel_data",
			})
			r.ExitAll()
			return nil, err
		}

		// 生成认证签名（包含channel_data）
		auth := websocket.GenerateAuthSignature(req.SocketId, req.ChannelName, channelData)

		g.Log().Debugf(ctx, "Generated presence auth for user:%d, socket:%s, channel:%s", c.UserId, req.SocketId, req.ChannelName)

		// 直接返回Pusher标准格式（不包装在GoFrame响应中）
		r.Response.WriteJson(g.Map{
			"auth":         auth,
			"channel_data": channelData,
		})
		r.ExitAll()
	} else {
		// Private频道认证（不包含channel_data）
		auth := websocket.GenerateAuthSignature(req.SocketId, req.ChannelName, "")

		g.Log().Debugf(ctx, "Generated private auth for user:%d, socket:%s, channel:%s", c.UserId, req.SocketId, req.ChannelName)

		// 直接返回Pusher标准格式（不包装在GoFrame响应中）
		r.Response.WriteJson(g.Map{
			"auth": auth,
		})
		r.ExitAll()
	}

	return
}
