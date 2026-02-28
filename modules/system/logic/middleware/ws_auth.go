// Package middleware
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package middleware

import (
	"devinggo/modules/system/model"
	"devinggo/modules/system/myerror"
	websocket2 "devinggo/modules/system/pkg/websocket"
	"devinggo/modules/system/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// ws鉴权中间件
func (s *sMiddleware) WsAuth(r *ghttp.Request) {
	ctx := r.GetCtx()
	sessionId, err := s.parseSessionId(r)
	g.Log().Debug(ctx, "sessionId:", sessionId)
	if err != nil && g.IsEmpty(sessionId) {
		conn, err := websocket2.GetConnection(r)
		if err != nil {
			// Pusher协议连接失败处理
			r.Response.WriteStatus(400)
			r.Response.WriteJson(g.Map{"error": err.Error()})
		} else {
			// Pusher协议错误响应
			errResponse := &websocket2.PusherResponse{
				Event: websocket2.EventError,
				Data:  `{"message":"sessionId miss","code":4009}`,
			}
			conn.WriteJSON(errResponse)
			conn.Close()
		}
		return
	} else {
		r.SetCtxVar(websocket2.SESSION_ID_KEY, sessionId)
	}

	r.Middleware.Next()
}

func (s *sMiddleware) parseSessionId(r *ghttp.Request) (sessionId string, err error) {
	ctx := r.GetCtx()
	sessionIdTmp := r.GetQuery("sessionId")
	token := r.GetQuery("token")

	//权限检查
	if g.IsEmpty(token) {
		return "", myerror.MissingParameter(ctx, "missing sessionId or token")
	}

	claims, err := service.Token().ParseToken(ctx, token.String())
	if err != nil {
		return "", err
	}
	data := claims.Data
	if g.IsEmpty(data) {
		return "", myerror.ValidationFailed(ctx, "claims is empty")
	}

	if g.IsEmpty(sessionIdTmp) {
		var user *model.Identity
		data := claims.Data
		err = gconv.Scan(data, &user)
		if err != nil {
			return "", err
		}
		if g.IsEmpty(user) {
			return "", myerror.ValidationFailed(ctx, "sessionId miss")
		} else {
			return gconv.String(user.Id), nil
		}

	} else {
		return sessionIdTmp.String(), nil
	}
}
