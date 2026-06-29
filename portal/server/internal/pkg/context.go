package pkg

import (
	"context"
	"portal/internal/model"
)

const (
	UserInfoKey = "user_info" // 用户信息
)

// SetUserToCtx 将用户信息保存到上下文中
func SetUserToCtx(ctx context.Context, user *model.XUser) context.Context {
	ctx = context.WithValue(ctx, UserInfoKey, user)
	return ctx
}

// DetailUserFromCtx 从上下文中获取用户信息
func DetailUserFromCtx(ctx context.Context) *model.XUser {
	userInfo := ctx.Value(UserInfoKey)
	return userInfo.(*model.XUser)
}
