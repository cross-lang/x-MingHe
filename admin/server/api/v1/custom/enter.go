package custom

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	UserApi
	EnterpriseApi
}

var (
	userService                = service.ServiceGroupApp.CustomServiceGroup.UserService
	enterpriseService          = service.ServiceGroupApp.CustomServiceGroup.EnterpriseService
)
