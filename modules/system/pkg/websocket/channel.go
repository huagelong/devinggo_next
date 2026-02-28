// Package websocket
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package websocket

import "strings"

// ChannelType 频道类型
type ChannelType int

const (
	ChannelTypePublic   ChannelType = iota // 公开频道
	ChannelTypePrivate                     // 私有频道 (private-*)
	ChannelTypePresence                    // 在线状态频道 (presence-*)
)

// GetChannelType 根据频道名获取类型
func GetChannelType(channel string) ChannelType {
	if strings.HasPrefix(channel, "presence-") {
		return ChannelTypePresence
	}
	if strings.HasPrefix(channel, "private-") {
		return ChannelTypePrivate
	}
	return ChannelTypePublic
}

// IsPrivateChannel 判断是否为私有频道
func IsPrivateChannel(channel string) bool {
	return strings.HasPrefix(channel, "private-")
}

// IsPresenceChannel 判断是否为在线状态频道
func IsPresenceChannel(channel string) bool {
	return strings.HasPrefix(channel, "presence-")
}

// IsPublicChannel 判断是否为公开频道
func IsPublicChannel(channel string) bool {
	return !IsPrivateChannel(channel) && !IsPresenceChannel(channel)
}

// RequiresAuth 判断频道是否需要认证
func RequiresAuth(channel string) bool {
	return IsPrivateChannel(channel) || IsPresenceChannel(channel)
}
