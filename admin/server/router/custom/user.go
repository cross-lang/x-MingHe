package custom

import (
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

// InitUserRouter 注册用户路由
func (p *UserRouter) InitUserRouter(router *gin.RouterGroup) {
	UserGroup := router.Group("users")
	{
		// 查询用户列表
		UserGroup.GET("list", userApi.ListUser)
		// 查询用户详情
		UserGroup.GET(":user_id/detail", userApi.DetailUser)
		// 编辑用户
		UserGroup.POST("update", userApi.UpdateUser)
		// 禁用用户
		UserGroup.POST(":user_id/disable", userApi.DisableUser)
		// 启用用户
		UserGroup.POST(":user_id/enable", userApi.EnableUser)
	}
}
