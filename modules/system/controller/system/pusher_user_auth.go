// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"context"
	"fmt"

	"devinggo/modules/system/api/system"
	"devinggo/modules/system/controller/base"
	"devinggo/modules/system/pkg/websocket"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	PusherUserAuthController = cPusherUserAuthController{}
)

type cPusherUserAuthController struct {
	base.BaseController
}

// UserAuth Pusher User Authentication - 用户身份认证
//
// 用途：
// - 识别和认证用户身份（区别于频道授权）
// - 用户级别的事件订阅和触发
// - 跨频道的用户状态追踪
//
// 使用场景：
// - 用户在线状态管理
// - 用户级别的通知
// - 多设备登录检测
// - 用户行为追踪
//
// 客户端用法：
//
//	const pusher = new Pusher('app-key', {
//	    userAuthentication: {
//	        endpoint: '/system/pusher/user-auth'
//	    }
//	});
//	pusher.signin();
//
// 文档：https://pusher.com/docs/channels/server_api/authenticating-users/
func (c *cPusherUserAuthController) UserAuth(ctx context.Context, req *system.PusherUserAuthReq) (res *system.PusherUserAuthRes, err error) {
	r := g.RequestFromCtx(ctx)

	// 验证用户是否已登录
	if c.UserId == 0 {
		writeJSONOrJSONP(r, g.Map{
			"error": "Unauthorized - User not logged in",
		}, 401)
		return
	}

	g.Log().Debugf(ctx, "Pusher User Authentication: user_id=%d, socket_id=%s", c.UserId, req.SocketId)

	// 构建用户数据
	// 必须包含 id 字段（字符串类型）
	// 可以添加任意其他用户信息
	userData := map[string]interface{}{
		"id":   gconv.String(c.UserId), // 必须是字符串
		"name": fmt.Sprintf("User %d", c.UserId),
		// 可以添加更多用户信息，例如：
		// "email": user.Email,
		// "avatar": user.Avatar,
		// "role": user.Role,
	}

	// TODO: 从数据库获取真实用户信息
	// user, err := service.User().GetById(ctx, c.UserId)
	// if err == nil {
	//     userData["name"] = user.Username
	//     userData["email"] = user.Email
	//     userData["avatar"] = user.Avatar
	// }

	// 生成用户认证签名
	auth, err := websocket.GenerateUserAuthSignature(req.SocketId, userData)
	if err != nil {
		g.Log().Error(ctx, "Failed to generate user auth signature:", err)
		writeJSONOrJSONP(r, g.Map{
			"error": "Failed to generate authentication signature",
		}, 500)
		return
	}

	g.Log().Infof(ctx, "User authenticated successfully: user_id=%d, socket_id=%s", c.UserId, req.SocketId)

	// 返回认证响应
	res = &system.PusherUserAuthRes{
		Auth:     auth,
		UserData: userData,
	}

	writeJSONOrJSONP(r, res, 200)

	return
}
