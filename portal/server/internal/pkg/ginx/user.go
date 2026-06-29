package ginx

import (
	"portal/internal/model"
	"portal/internal/pkg/errorx"

	"github.com/gin-gonic/gin"
)

const UserInfoKey = "user_info"

func SetUserToCtx(ctx *gin.Context, user *model.XUser) {
	ctx.Set(UserInfoKey, user)
}

// DetailUserFromCtx 从上下文中获取用户信息
func DetailUserFromCtx(ctx *gin.Context) (*model.XUser, error) {
	data, ok := ctx.Get(UserInfoKey)
	if ok {
		return data.(*model.XUser), nil
	}
	return nil, errorx.UserLoginStatusException
}
