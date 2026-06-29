package custom

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	UserRouter
	EnterpriseRouter
}

var (
	enterpriseApi          = api.ApiGroupApp.CustomApiGroup.EnterpriseApi
	userApi                = api.ApiGroupApp.CustomApiGroup.UserApi
)
