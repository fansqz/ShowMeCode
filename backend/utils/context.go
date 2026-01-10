package utils

import (
	"context"
)

const (
	CtxUserIDKey  = "userID"
	CtxVisitorUID = "visitorUID"
)

type UserInfo struct {
	ID        uint     `json:"id"`
	Avatar    string   `json:"avatar"`
	LoginName string   `json:"loginName"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Roles     []uint   `json:"roles"`
	Menus     []string `json:"menus"`
}

// GetUserIDWithCtx 从ctx中获取userID
func GetUserIDWithCtx(ctx context.Context) uint {
	id := ctx.Value(CtxUserIDKey)
	answer, ok := id.(uint)
	if !ok {
		return 0
	}
	return answer
}

// GetVisitorIDWithCtx 获取访客id
func GetVisitorIDWithCtx(ctx context.Context) string {
	uid, ok := ctx.Value(CtxVisitorUID).(string)
	if !ok {
		return ""
	}
	return uid
}
